package migration

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upRoles, downRoles)
}

func upRoles(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS roles (
		id BIGSERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP,
	    deleted_at TIMESTAMP
	);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downRoles(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS roles;`)
	if err != nil {
		return err
	}
	return nil
}
