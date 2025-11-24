package migration

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upVerificationTokens, downVerificationTokens)
}

func upVerificationTokens(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS verification_tokens (
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		token VARCHAR(255) NOT NULL,
		token_type VARCHAR(20) NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP,
		deleted_at TIMESTAMP,

		CONSTRAINT fk_verification_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_verification_tokens_user_id ON verification_tokens(user_id);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downVerificationTokens(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS verification_tokens;`)
	if err != nil {
		return err
	}
	return nil
}
