package db

import (
	"github.com/hackallcode/db-homework/internal/pkg/models"
	"github.com/jackc/pgx"
)

const (
	TruncateQuery = `TRUNCATE TABLE users; TRUNCATE TABLE forums; TRUNCATE TABLE threads; TRUNCATE TABLE posts; TRUNCATE TABLE votes`
	CountRowsQuery = `SELECT (SELECT COUNT(*) FROM users), (SELECT COUNT(*) FROM forums), (SELECT COUNT(*) FROM threads), (SELECT COUNT(*) FROM posts)`
)

func TruncateAll(tx *pgx.Tx) error {
	_, err := tx.Exec(TruncateQuery)
	return err
}

func Status(tx *pgx.Tx) (*models.Status, error) {
	statusData := &models.Status{}

	row := tx.QueryRow(CountRowsQuery)
	err := row.Scan(&statusData.User, &statusData.Forum, &statusData.Thread, &statusData.Post)
	if err != nil {
		return nil, err
	}

	return statusData, nil
}
