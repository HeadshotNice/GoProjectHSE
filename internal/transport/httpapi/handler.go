package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"HSE/internal/usecase"
	"HSE/internal/usecase/authjwt"
)

type Handler struct {
	uc *usecase.Usecase

	jwt *authjwt.Manager
}

func NewHandler(uc *usecase.Usecase, jwtSecret, jwtIssuer string) *Handler {
	return &Handler{
		uc:  uc,
		jwt: authjwt.New(jwtSecret, jwtIssuer, 0),
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /metrics", metricsHandler())
	mux.HandleFunc("GET /", h.handleIndex)
	mux.HandleFunc("GET /test", h.handleTest)
	mux.HandleFunc("POST /dbtest", h.handleDBTest)

	mux.HandleFunc("POST /auth/register", h.handleRegister)
	mux.HandleFunc("POST /auth/login", h.handleLogin)

	mux.Handle("POST /documents", h.authMiddleware(http.HandlerFunc(h.handleCreateDocument)))
	mux.Handle("GET /documents", h.authMiddleware(http.HandlerFunc(h.handleListDocuments)))

	return metricsMiddleware(mux)
}

func (h *Handler) handleIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(indexHTML()))
}

func (h *Handler) handleTest(w http.ResponseWriter, r *http.Request) {
	msg, err := h.uc.TestHello(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(msg))
}

func (h *Handler) handleDBTest(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad body")
		return
	}
	line := strings.TrimSpace(string(b))
	if err := h.uc.DBTestInsert(r.Context(), line); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrBadRequest) {
			code = http.StatusBadRequest
		}
		writeError(w, code, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}
	id, err := h.uc.Register(r.Context(), strings.TrimSpace(req.Email), req.Password)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrBadRequest) {
			code = http.StatusBadRequest
		}
		if errors.Is(err, usecase.ErrAlreadyExists) {
			code = http.StatusConflict
		}
		writeError(w, code, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user_id": id})
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}
	token, err := h.uc.Login(r.Context(), strings.TrimSpace(req.Email), req.Password)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrBadRequest) {
			code = http.StatusBadRequest
		}
		if errors.Is(err, usecase.ErrUnauthorized) {
			code = http.StatusUnauthorized
		}
		writeError(w, code, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token})
}

type ctxKey int

const userIDKey ctxKey = 1

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		userID, err := h.jwt.ParseUserID(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "bad token")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func userIDFromCtx(ctx context.Context) int64 {
	v := ctx.Value(userIDKey)
	if v == nil {
		return 0
	}
	id, _ := v.(int64)
	return id
}

func (h *Handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromCtx(r.Context())
	orderID, err := h.uc.CreateOrder(r.Context(), userID)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrUnauthorized) {
			code = http.StatusUnauthorized
		}
		writeError(w, code, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"order_id": orderID})
}

func (h *Handler) handleListOrders(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromCtx(r.Context())
	activeOnly := r.URL.Query().Get("active") == "true"
	orders, err := h.uc.ListOrders(r.Context(), userID, activeOnly)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrUnauthorized) {
			code = http.StatusUnauthorized
		}
		writeError(w, code, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"orders": orders})
}

type createDocumentReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *Handler) handleCreateDocument(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromCtx(r.Context())

	var req createDocumentReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}

	documentID, err := h.uc.CreateDocument(
		r.Context(),
		userID,
		strings.TrimSpace(req.Title),
		strings.TrimSpace(req.Content),
	)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrBadRequest) {
			code = http.StatusBadRequest
		}
		if errors.Is(err, usecase.ErrUnauthorized) {
			code = http.StatusUnauthorized
		}
		writeError(w, code, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"document_id": documentID})
}

func (h *Handler) handleListDocuments(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromCtx(r.Context())
	documents, err := h.uc.ListDocuments(r.Context(), userID)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrUnauthorized) {
			code = http.StatusUnauthorized
		}
		writeError(w, code, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"documents": documents})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]any{"error": msg})
}
