package data

import (
	"context"
	"database/sql"
	"errors"
	"finstar/config"
	apierror "finstar/internal/error"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

func RunMigrations() error {

	if config.Get().PostgresConn == "" {
		return errors.New("PostgresConn blank")
	}
	m, err := migrate.New(
		"file://migrations",
		config.Get().PostgresConn,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

type DB struct {
	db *pgx.Conn
}

func NewDbRepository(db *pgx.Conn) *DB {
	return &DB{db: db}
}

func (r *DB) FindUser(ctx context.Context, userId int) (bool, error) {
	var id int

	err := r.db.QueryRow(ctx, "SELECT id from users where id = $1", userId).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *DB) Deposited(ctx context.Context, userId int, total float32) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", total, userId)

	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {

		return err
	}
	return nil
}

func (r *DB) Transfer(ctx context.Context, userIdFrom int, userIdTo int, total float32) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	var balance float32
	err = r.db.QueryRow(ctx, "SELECT balance from users where id = $1", userIdFrom).Scan(&balance)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if balance < total {
		_ = tx.Rollback(ctx)
		return apierror.LowBalance
	}

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance - $1 WHERE id = $2", total, userIdFrom)

	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", total, userIdTo)

	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return err
	}
	return nil
}
