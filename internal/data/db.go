package data

import (
	"database/sql"
	apierror "finstar/internal/error"
)

type DB struct {
	db *sql.DB
}

func NewDbRepository(db *sql.DB) *DB {
	return &DB{db: db}
}

// TODO
func (r *DB) FindUser(userId int) (bool, error) {
	var id int
	err := r.db.QueryRow("SELECT id from users where id = $1", userId).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *DB) Deposited(userId int, total float32) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", total, userId)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil {

		return err
	}
	return nil
}

func (r *DB) Transfer(userIdFrom int, userIdTo int, total float32) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var balance float32
	err = r.db.QueryRow("SELECT balance from users where id = $1", userIdFrom).Scan(&balance)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if balance < total {
		_ = tx.Rollback()
		return apierror.LowBalance
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", total, userIdFrom)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", total, userIdTo)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil {
		return err
	}
	return nil
}
