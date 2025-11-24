package migration

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUserRole, downUserRole)
}

func upUserRole(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS user_role (
		id BIGSERIAL PRIMARY KEY,
		role_id BIGINT NOT NULL,
		user_id BIGINT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP,
		deleted_at TIMESTAMP,

		CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_user_role_user_id ON user_role(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_role_role_id ON user_role(role_id);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downUserRole(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS user_role;`)
	if err != nil {
		return err
	}
	return nil
}
