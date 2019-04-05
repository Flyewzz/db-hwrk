package db

import (
	"github.com/jackc/pgx"

	"github.com/hackallcode/db-homework/internal/pkg/models"
	"github.com/hackallcode/db-homework/internal/pkg/verifier"
)

const (
	InsertUserQuery               = `INSERT INTO users (nickname, email, fullname, about) VALUES ($1, $2, $3, $4)`
	GetUserByNicknameQuery        = `SELECT id, nickname, email, fullname, about FROM users WHERE LOWER(nickname) = LOWER($1)`
	GetUserByNicknameOrEmailQuery = `SELECT id, nickname, email, fullname, about FROM users 
	WHERE LOWER(nickname) = LOWER($1) OR LOWER(email) = LOWER($2)`
	GetUsersByForumIdQuery = `SELECT DISTINCT u.id, u.nickname, u.email, u.fullname, u.about, LOWER(u.nickname) AS lower_nickname
	FROM threads AS t LEFT JOIN posts AS p ON t.id = p.thread_id
 	JOIN users AS u ON t.user_id = u.id OR p.user_id = u.id
	WHERE t.forum_id = $1`
	UpdateUserQuery    = `UPDATE users SET email = $2, fullname = $3, about = $4 WHERE LOWER(nickname) = LOWER($1)`
	TruncateUsersQuery = `TRUNCATE TABLE users CASCADE`
)

func CreateUser(tx *pgx.Tx, user models.User) (models.Users, error) {
	users := models.Users{}

	rows, err := tx.Query(GetUserByNicknameOrEmailQuery, user.Nickname, user.Email)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		dbUser := models.User{}
		err = rows.Scan(&dbUser.Id, &dbUser.Nickname, &dbUser.Email, &dbUser.FullName, &dbUser.About)
		if err != nil {
			return users, err
		}
		users = append(users, dbUser)
	}
	// Author already exists
	if len(users) > 0 {
		return users, models.AlreadyExistsError
	}

	// Create user
	_, err = tx.Exec(InsertUserQuery, user.Nickname, user.Email, user.FullName, user.About)
	if err != nil {
		return users, err
	}
	return users, nil
}

func GetUser(tx *pgx.Tx, nickname string) (*models.User, error) {
	user := &models.User{}
	row := tx.QueryRow(GetUserByNicknameQuery, nickname)
	err := row.Scan(&user.Id, &user.Nickname, &user.Email, &user.FullName, &user.About)
	if err == pgx.ErrNoRows {
		return nil, models.NotFoundError
	}
	return user, err
}

func UpdateUser(tx *pgx.Tx, updateData *models.UserUpdate) error {
	if verifier.IsEmpty(updateData.Email) || verifier.IsEmpty(updateData.FullName) || verifier.IsEmpty(updateData.About) {
		userData := &models.User{}
		row := tx.QueryRow(GetUserByNicknameQuery, updateData.Nickname)
		err := row.Scan(&userData.Id, &userData.Nickname, &userData.Email, &userData.FullName, &userData.About)
		if err != nil {
			if err == pgx.ErrNoRows {
				return models.NotFoundError
			}
			return err
		}
		if verifier.IsEmpty(updateData.Email) {
			updateData.Email = userData.Email
		}
		if verifier.IsEmpty(updateData.FullName) {
			updateData.FullName = userData.FullName
		}
		if verifier.IsEmpty(updateData.About) {
			updateData.About = userData.About
		}
	}

	tag, err := tx.Exec(UpdateUserQuery, updateData.Nickname, updateData.Email, updateData.FullName, updateData.About)
	if err != nil {
		if pgxErr, ok := err.(pgx.PgError); ok {
			// Unique violation
			if pgxErr.Code == UniqueErrorCode {
				return models.AlreadyExistsError
			}
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return models.NotFoundError
	}

	return nil
}

func GetForumUsers(tx *pgx.Tx, slug, limit, since, desc string) (models.Users, error) {
	users := models.Users{}

	forumId, _, err := GetForumId(tx, slug)
	if err != nil {
		return users, err
	}

	where := ""
	if since != "" {
		if desc == "true" {
			where = " AND LOWER(u.nickname) < LOWER('" + since + "')"
		} else {
			where = " AND LOWER(u.nickname) > LOWER('" + since + "')"
		}
	}

	orderBy := " ORDER BY lower_nickname ASC"
	if desc == "true" {
		orderBy = " ORDER BY lower_nickname DESC"
	}

	limitStr := ""
	if limit != "" {
		limitStr = " LIMIT " + limit
	}

	rows, err := tx.Query(GetUsersByForumIdQuery+where+orderBy+limitStr, forumId)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		lowerNickname := ""
		err = rows.Scan(&user.Id, &user.Nickname, &user.Email, &user.FullName, &user.About, &lowerNickname)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func TruncateUsers(tx *pgx.Tx) error {
	_, err := tx.Exec(TruncateUsersQuery)
	return err
}
