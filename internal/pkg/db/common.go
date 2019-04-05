package db

import (
	"errors"
	"io/ioutil"

	"github.com/jackc/pgx"
)

/********************/
/*      ERRORS      */
/********************/

var (
	AlreadyInitError = errors.New("db already initialized")
	NotInitError     = errors.New("db wasn't initialized")
)

/********************/
/*  BASE FUNCTIONS  */
/********************/

var conn *pgx.ConnPool

func Open() (err error) {
	if conn != nil {
		return AlreadyInitError
	}
	conn, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Database: database,
			User:     user,
			Password: password,
		},
		MaxConnections: maxConnections,
	})
	if err != nil {
		return err
	}

	if query, err := ioutil.ReadFile("init/init.sql"); err != nil {
		return err
	} else {
		if _, err := conn.Exec(string(query)); err != nil {
			return err
		}
	}
	return
}

func Close() error {
	if conn == nil {
		return NotInitError
	}
	conn.Close()
	return nil
}

func Begin() (tx *pgx.Tx, err error) {
	if conn == nil {
		return tx, NotInitError
	}

	return conn.Begin()
}

func QueryRow(query string, args ...interface{}) (row *pgx.Row, err error) {
	tx, err := Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row = tx.QueryRow(query, args...)

	return row, tx.Commit()
}

func Query(query string, args ...interface{}) (rows *pgx.Rows, err error) {
	tx, err := Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	rows, err = tx.Query(query, args...)
	if err != nil {
		return
	}

	return rows, tx.Commit()
}

func Exec(query string, args ...interface{}) (tag pgx.CommandTag, err error) {
	tx, err := Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	tag, err = tx.Exec(query, args...)
	if err != nil {
		return
	}

	return tag, tx.Commit()
}
