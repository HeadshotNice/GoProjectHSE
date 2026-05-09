package authjwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func New(secret, issuer string, ttl time.Duration) *Manager {
	return &Manager{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}
}

func (m *Manager) Issue(userID int64) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": m.issuer,
		"sub": fmt.Sprintf("%d", userID),
		"iat": now.Unix(),
		"exp": now.Add(m.ttl).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(m.secret)
}

func (m *Manager) ParseUserID(tokenString string) (int64, error) {
	tok, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	}, jwt.WithIssuer(m.issuer))
	if err != nil {
		return 0, err
	}

	if !tok.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return 0, fmt.Errorf("missing sub")
	}

	var userID int64
	_, err = fmt.Sscanf(sub, "%d", &userID)
	if err != nil || userID <= 0 {
		return 0, fmt.Errorf("bad sub")
	}
	return userID, nil
}

