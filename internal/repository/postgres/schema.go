package postgres

import (
	"context"
	"database/sql"
)

func InitSchema(ctx context.Context, db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS dbtest_logs (
			id BIGSERIAL PRIMARY KEY,
			line TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS users_email_lower_unique ON users (lower(email))`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS pass_hash TEXT`,
		`CREATE TABLE IF NOT EXISTS orders (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			status TEXT NOT NULL DEFAULT 'created',
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			completed_at TIMESTAMPTZ
		)`,
		`ALTER TABLE orders ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT now()`,
		`ALTER TABLE orders ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ`,
		`UPDATE orders SET updated_at = COALESCE(updated_at, created_at, now())`,
		`CREATE TABLE IF NOT EXISTS documents (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending_review',
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`,
		`ALTER TABLE documents ADD COLUMN IF NOT EXISTS title TEXT`,
		`ALTER TABLE documents ADD COLUMN IF NOT EXISTS content TEXT`,
		`ALTER TABLE documents ADD COLUMN IF NOT EXISTS status TEXT`,
		`ALTER TABLE documents ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT now()`,
		`ALTER TABLE documents ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT now()`,
		`ALTER TABLE documents ALTER COLUMN status SET DEFAULT 'pending_review'`,
		`CREATE SEQUENCE IF NOT EXISTS documents_id_seq`,
		`ALTER TABLE documents ALTER COLUMN id SET DEFAULT nextval('documents_id_seq')`,
		`SELECT setval('documents_id_seq', COALESCE((SELECT MAX(id) FROM documents), 1), (SELECT COUNT(*) > 0 FROM documents))`,
		`UPDATE documents SET title = COALESCE(title, 'Документ')`,
		`UPDATE documents SET content = COALESCE(content, '')`,
		`UPDATE documents SET status = COALESCE(NULLIF(status, ''), 'pending_review')`,
		`UPDATE documents SET created_at = COALESCE(created_at, now())`,
		`UPDATE documents SET updated_at = COALESCE(updated_at, created_at, now())`,
		`DO $$
		DECLARE col record;
		BEGIN
			FOR col IN
				SELECT column_name
				FROM information_schema.columns
				WHERE table_name = 'documents'
				  AND is_nullable = 'NO'
				  AND column_name NOT IN ('id', 'user_id', 'title', 'content', 'status', 'created_at')
			LOOP
				EXECUTE format('ALTER TABLE documents ALTER COLUMN %I DROP NOT NULL', col.column_name);
			END LOOP;
		END $$`,
		`DO $$
		DECLARE pass_hash_udt text;
		BEGIN
			IF EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_name = 'users' AND column_name = 'password'
			) THEN
				EXECUTE 'UPDATE users SET password_hash = password WHERE password_hash IS NULL';
			END IF;

			SELECT c.udt_name INTO pass_hash_udt
			FROM information_schema.columns c
			WHERE c.table_name = 'users' AND c.column_name = 'pass_hash'
			LIMIT 1;

			IF pass_hash_udt = 'bytea' THEN
				EXECUTE 'UPDATE users SET pass_hash = convert_to(password_hash, ''UTF8'') WHERE pass_hash IS NULL AND password_hash IS NOT NULL';
				EXECUTE 'UPDATE users SET password_hash = convert_from(pass_hash, ''UTF8'') WHERE password_hash IS NULL AND pass_hash IS NOT NULL';
				EXECUTE 'ALTER TABLE users ALTER COLUMN pass_hash DROP NOT NULL';
			ELSE
				EXECUTE 'UPDATE users SET pass_hash = password_hash WHERE pass_hash IS NULL AND password_hash IS NOT NULL';
				EXECUTE 'UPDATE users SET password_hash = pass_hash WHERE password_hash IS NULL AND pass_hash IS NOT NULL';
				EXECUTE 'ALTER TABLE users ALTER COLUMN pass_hash DROP NOT NULL';
			END IF;
		END $$`,
	}

	for _, s := range stmts {
		if _, err := db.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}
