package db

import (
	"github.com/jackc/pgx"

	"github.com/hackallcode/db-homework/internal/pkg/models"
)

const (
	InsertForumQuery    = `INSERT INTO forums (user_id, slug, title) VALUES ($1, $2, $3)`
	GetForumIdBySlugQuery = `SELECT id, slug FROM forums WHERE LOWER(slug) = LOWER($1)`
	GetForumBySlugQuery = `SELECT f.id, u.nickname, f.slug, f.title, f.threads, f.posts
	FROM forums AS f JOIN users AS u ON f.user_id = u.id
	WHERE LOWER(f.slug) = LOWER($1)`
	TruncateForumsQuery = `TRUNCATE TABLE forums CASCADE`
)

func CreateForum(tx *pgx.Tx, forum *models.Forum) error {
	existingForum, err := GetForum(tx, forum.Slug)
	if err == nil {
		*forum = *existingForum
		return models.AlreadyExistsError
	}
	if err != models.NotFoundError {
		return err
	}

	userData, err := GetUser(tx, forum.User)
	if err != nil {
		return err
	}
	forum.User = userData.Nickname

	_, err = tx.Exec(InsertForumQuery, userData.Id, forum.Slug, forum.Title)
	if err != nil {
		if pgxErr, ok := err.(pgx.PgError); ok {
			// Unique violation
			if pgxErr.Code == UniqueErrorCode {
				return models.AlreadyExistsError
			}
		}
		return err
	}
	return nil
}

func GetForumId(tx *pgx.Tx, slug string) (int64, string, error) {
	row := tx.QueryRow(GetForumIdBySlugQuery, slug)
	id := int64(0)
	err := row.Scan(&id, &slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return id, slug, models.NotFoundError
		}
		return id, slug, err
	}
	return id, slug, nil
}

func GetForum(tx *pgx.Tx, slug string) (*models.Forum, error) {
	forum := &models.Forum{}
	row := tx.QueryRow(GetForumBySlugQuery, slug)
	err := row.Scan(&forum.Id, &forum.User, &forum.Slug, &forum.Title, &forum.Threads, &forum.Posts)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.NotFoundError
		}
		return nil, err
	}
	return forum, nil
}

func TruncateForums(tx *pgx.Tx) (err error) {
	_, err = tx.Exec(TruncateForumsQuery)
	return err
}
