package db

import (
	"strconv"
	"time"

	"github.com/jackc/pgx"

	"github.com/hackallcode/db-homework/internal/pkg/models"
	"github.com/hackallcode/db-homework/internal/pkg/verifier"
)

const (
	InsertThreadQuery          = `INSERT INTO threads (forum_id, user_id, created, slug, title, message) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
	IncThreadsCount            = `UPDATE forums SET threads = threads + 1 WHERE id = $1`
	GetThreadIdByIdOrSlugQuery = `SELECT id, slug FROM threads WHERE id = $1 OR LOWER(slug) = LOWER($2)`
	GetThreadIdBySlugQuery     = `SELECT id, slug FROM threads WHERE LOWER(slug) = LOWER($1)`
	GetThreadByIdQuery   = `SELECT t.id, f.slug, u.nickname, t.created, t.slug, t.title, t.message, t.votes
	FROM threads AS t JOIN forums AS f ON t.forum_id = f.id JOIN users AS u ON t.user_id = u.id WHERE t.id = $1`
	GetThreadByIdOrSlugQuery   = `SELECT t.id, f.slug, u.nickname, t.created, t.slug, t.title, t.message, t.votes
	FROM threads AS t JOIN forums AS f ON t.forum_id = f.id JOIN users AS u ON t.user_id = u.id WHERE t.id = $1 OR LOWER(t.slug) = LOWER($2)`
	GetThreadBySlugQuery = `SELECT t.id, f.slug, u.nickname, t.created, t.slug, t.title, t.message, t.votes
	FROM threads AS t JOIN forums AS f ON t.forum_id = f.id JOIN users AS u ON t.user_id = u.id WHERE LOWER(t.slug) = LOWER($1)`
	GetThreadsByForumIdQuery = `SELECT t.id, u.nickname, t.created, t.slug, t.title, t.message, t.votes
	FROM threads AS t JOIN users AS u ON t.user_id = u.id WHERE t.forum_id = $1`
	UpdateThreadQuery    = `UPDATE threads SET title = $2, message = $3 WHERE id = $1`
	GetVoteQuery         = `SELECT id, vote FROM votes WHERE thread_id = $1 AND user_id = $2`
	InsertVoteQuery      = `INSERT INTO votes (thread_id, user_id, vote) VALUES ($1, $2, $3)`
	UpdateVoteQuery      = `UPDATE  votes SET vote = $2 WHERE id = $1`
	VoteThreadQuery      = `UPDATE threads SET votes = $2 WHERE id = $1`
	TruncateThreadsQuery = `TRUNCATE TABLE threads CASCADE`
)

func CreateThread(tx *pgx.Tx, thread *models.Thread) error {
	if thread.Slug != "" {
		existingThread, err := GetThread(tx, thread.Slug)
		if err == nil {
			*thread = *existingThread
			return models.AlreadyExistsError
		}
		if err != models.NotFoundError {
			return err
		}
	}

	forumId, forumSlug, err := GetForumId(tx, thread.Forum)
	if err != nil {
		return err
	}
	thread.Forum = forumSlug

	userData, err := GetUser(tx, thread.Author)
	if err != nil {
		return err
	}
	thread.Author = userData.Nickname

	if verifier.IsEmpty(thread.Created) {
		thread.Created = time.Now().Format(DateFormat)
	}
	row := tx.QueryRow(InsertThreadQuery, forumId, userData.Id, thread.Created, thread.Slug, thread.Title, thread.Message)
	err = row.Scan(&thread.Id)
	if err != nil {
		if pgxErr, ok := err.(pgx.PgError); ok {
			// Unique violation
			if pgxErr.Code == UniqueErrorCode {
				return models.AlreadyExistsError
			}
		}
		return err
	}

	_, err = tx.Exec(IncThreadsCount, forumId)
	if err != nil {
		return err
	}
	return nil
}

func GetThreadId(tx *pgx.Tx, slugOrId string) (int64, string, error) {
	id, err := strconv.ParseInt(slugOrId, 10, 64)
	slug := slugOrId
	if err == nil {
		row := tx.QueryRow(GetThreadIdByIdOrSlugQuery, id, slug)
		err = row.Scan(&id, &slug)
	} else {
		row := tx.QueryRow(GetThreadIdBySlugQuery, slug)
		err = row.Scan(&id, &slug)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return id, slug, models.NotFoundError
		}
		return id, slug, err
	}
	return id, slug, nil
}

func GetThread(tx *pgx.Tx, slugOrId string) (*models.Thread, error) {
	var created time.Time
	thread := &models.Thread{}

	id, err := strconv.ParseInt(slugOrId, 10, 64)
	slug := slugOrId
	if err == nil {
		row := tx.QueryRow(GetThreadByIdOrSlugQuery, id, slug)
		err = row.Scan(&thread.Id, &thread.Forum, &thread.Author, &created, &thread.Slug, &thread.Title, &thread.Message, &thread.Votes)
	} else {
		row := tx.QueryRow(GetThreadBySlugQuery, slug)
		err = row.Scan(&thread.Id, &thread.Forum, &thread.Author, &created, &thread.Slug, &thread.Title, &thread.Message, &thread.Votes)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.NotFoundError
		}
		return nil, err
	}
	thread.Created = created.Format(DateFormat)
	return thread, nil
}

func GetThreadById(tx *pgx.Tx, id int64) (*models.Thread, error) {
	var created time.Time
	thread := &models.Thread{}

	row := tx.QueryRow(GetThreadByIdQuery, id)
	err := row.Scan(&thread.Id, &thread.Forum, &thread.Author, &created, &thread.Slug, &thread.Title, &thread.Message, &thread.Votes)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.NotFoundError
		}
		return nil, err
	}
	thread.Created = created.Format(DateFormat)
	return thread, nil
}

func UpdateThread(tx *pgx.Tx, updateData *models.ThreadUpdate) (*models.Thread, error) {
	threadData, err := GetThread(tx, updateData.SlugOrId)
	if err != nil {
		return threadData, err
	}

	if !verifier.IsEmpty(updateData.Title) || !verifier.IsEmpty(updateData.Message) {
		if verifier.IsEmpty(updateData.Title) {
			updateData.Title = threadData.Title
		}
		if verifier.IsEmpty(updateData.Message) {
			updateData.Message = threadData.Message
		}

		_, err = tx.Exec(UpdateThreadQuery, threadData.Id, updateData.Title, updateData.Message)
		if err != nil {
			return threadData, err
		}
		threadData.Title = updateData.Title
		threadData.Message = updateData.Message
	}

	return threadData, nil
}

func VoteThread(tx *pgx.Tx, updateData *models.Vote) (*models.Thread, error) {
	threadData, err := GetThread(tx, updateData.SlugOrId)
	if err != nil {
		return nil, err
	}

	userData, err := GetUser(tx, updateData.Nickname)
	if err != nil {
		return nil, err
	}

	hasVote := true
	voteId := int64(0)
	voteValue := int64(0)
	row := tx.QueryRow(GetVoteQuery, threadData.Id, userData.Id)
	err = row.Scan(&voteId, &voteValue)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		hasVote = false
	}

	updateVote := updateData.Vote - voteValue
	if updateVote != 0 {
		threadData.Votes += updateVote
		_, err = tx.Exec(VoteThreadQuery, threadData.Id, threadData.Votes)
		if err != nil {
			return nil, err
		}
		if hasVote {
			_, err = tx.Exec(UpdateVoteQuery, voteId, updateData.Vote)
			if err != nil {
				return nil, err
			}
		} else {
			_, err = tx.Exec(InsertVoteQuery, threadData.Id, userData.Id, updateData.Vote)
			if err != nil {
				return nil, err
			}
		}
	}

	return threadData, nil
}

func GetForumThreads(tx *pgx.Tx, slug, limit, since, desc string) (models.Threads, error) {
	threads := models.Threads{}

	forumId, forumSlug, err := GetForumId(tx, slug)
	if err != nil {
		return threads, err
	}

	where := ""
	if since != "" {
		if desc == "true" {
			where = " AND t.created <= TIMESTAMPTZ '" + since + "'"
		} else {
			where = " AND t.created >= TIMESTAMPTZ '" + since + "'"
		}
	}

	orderBy := ""
	if desc == "true" {
		orderBy = " ORDER BY t.created DESC, t.id DESC"
	} else { // if desc == "false"
		orderBy = " ORDER BY t.created ASC, t.id ASC"
	}

	limitStr := ""
	if limit != "" {
		limitStr = " LIMIT " + limit
	}

	rows, err := tx.Query(GetThreadsByForumIdQuery+where+orderBy+limitStr, forumId)
	if err != nil {
		return threads, err
	}
	defer rows.Close()

	for rows.Next() {
		thread := models.Thread{}
		var created time.Time
		err = rows.Scan(&thread.Id, &thread.Author, &created, &thread.Slug, &thread.Title, &thread.Message, &thread.Votes)
		if err != nil {
			return threads, err
		}
		thread.Forum = forumSlug
		thread.Created = created.Format(DateFormat)
		threads = append(threads, thread)
	}
	return threads, nil
}

func TruncateThreads(tx *pgx.Tx) (err error) {
	_, err = tx.Exec(TruncateThreadsQuery)
	return err
}
