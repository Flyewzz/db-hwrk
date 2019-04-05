package db

import (
	"strconv"
	"time"

	"github.com/jackc/pgx"

	"github.com/hackallcode/db-homework/internal/pkg/verifier"

	"github.com/hackallcode/db-homework/internal/pkg/models"
)

const (
	CheckPostQuery            = `SELECT id, thread_id FROM posts WHERE id = $1 AND thread_id = $2`
	GetForumByThreadQuery     = `SELECT f.id, f.slug, t.id, t.slug FROM threads AS t JOIN forums AS f ON (t.id = $1 OR LOWER(t.slug) = LOWER($2)) AND t.forum_id = f.id`
	GetForumByThreadSlugQuery = `SELECT f.id, f.slug, t.id, t.slug FROM threads AS t JOIN forums AS f ON LOWER(t.slug) = LOWER($1) AND t.forum_id = f.id`
	InsertPostQuery           = `INSERT INTO posts (thread_id, user_id, created, parent_id, message) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	IncPostsCount             = `UPDATE forums SET posts = posts + 1 WHERE id = $1`
	GetPostByIdQuery          = `SELECT p.id, f.slug, p.thread_id, u.nickname, p.created, p.parent_id, p.message, p.is_edited
	FROM posts AS p 
	JOIN threads AS t ON p.thread_id = t.id 
	JOIN forums AS f ON t.forum_id = f.id 
	JOIN users AS u ON p.user_id = u.id 
	WHERE p.id = $1`
	GetPostsByThreadIdQuery = `SELECT p.id, f.slug, p.thread_id, u.nickname, p.created, p.parent_id, p.message, p.is_edited
	FROM posts AS p 
	JOIN threads AS t ON p.thread_id = t.id 
	JOIN forums AS f ON t.forum_id = f.id 
	JOIN users AS u ON p.user_id = u.id 
	WHERE p.thread_id = $1`

	ParentTreeSinceQuery1 = `
	SELECT id, slug, thread_id, nickname, created, parent_id, message, is_edited FROM (
  	SELECT *, SUM(is_parent) OVER(PARTITION BY is_greater ORDER BY path ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS parent_number FROM (
    SELECT *, SUM(more) OVER(ORDER BY path ROWS BETWEEN UNBOUNDED PRECEDING AND 1 PRECEDING) AS is_greater FROM (
      WITH RECURSIVE id_tree (id, path, more, is_parent) AS (
        (SELECT p.id, array[p.id] AS path, CASE WHEN p.id = `
	ParentTreeSinceQuery2 = ` THEN 1 ELSE 0 END AS more, 1 AS is_parent
          FROM posts AS p WHERE p.parent_id = 0 AND p.thread_id = $1
          ORDER BY p.created, p.id
        ) UNION ALL
        (
          SELECT p.id, array_append(path, p.id), CASE WHEN p.id = `
	ParentTreeSinceQueryDesc1 = `
	SELECT id, slug, thread_id, nickname, created, parent_id, message, is_edited FROM (
  	SELECT *, SUM(is_parent) OVER(PARTITION BY is_greater ORDER BY path ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS parent_number FROM (
    SELECT *, SUM(more) OVER(ORDER BY path ROWS BETWEEN UNBOUNDED PRECEDING AND 1 PRECEDING) AS is_greater FROM (
      WITH RECURSIVE id_tree (id, path, more, is_parent) AS (
        (SELECT p.id, array[-p.id] AS path, CASE WHEN p.id = `
	ParentTreeSinceQueryDesc2 = ` THEN 1 ELSE 0 END AS more, 1 AS is_parent
          FROM posts AS p WHERE p.parent_id = 0 AND p.thread_id = $1
          ORDER BY p.created DESC, path
        ) UNION ALL
        (
          SELECT p.id, array_append(path, p.id), CASE WHEN p.id = `
	ParentTreeSinceQuery3 = ` THEN 1 ELSE 0 END AS more, 0 AS is_parent
          FROM posts AS p JOIN id_tree AS rt ON rt.id = p.parent_id
        )
      )
      SELECT p.id, f.slug, p.thread_id, u.nickname, p.created, p.parent_id, p.message, p.is_edited, it.path, it.more, it.is_parent
      FROM id_tree AS it JOIN posts AS p ON it.id = p.id
      JOIN threads AS t ON p.thread_id = t.id
      JOIN forums AS f ON t.forum_id = f.id
      JOIN users AS u ON p.user_id = u.id
    ) AS r
  	) AS r
  	WHERE is_greater > 0
	) as r `
	ParentTreeSinceQuery4 = ` ORDER BY path `

	SinceBefore = `SELECT id, slug, thread_id, nickname, created, parent_id, message, is_edited FROM (
  SELECT *, SUM(r.more) OVER(ROWS BETWEEN UNBOUNDED PRECEDING AND 1 PRECEDING) AS moreSum FROM (`
	GetTreePostsByThreadIdQuery1 = ` WITH RECURSIVE id_tree (id, path) AS (
    	(SELECT p.id, array[p.id] FROM posts AS p WHERE p.parent_id = 0 AND p.thread_id = $1 ORDER BY p.created, p.id `
	GetTreePostsByThreadIdQuery1Desc = ` WITH RECURSIVE id_tree (id, path) AS (
    	(SELECT p.id, array[-p.id] as path FROM posts AS p WHERE p.parent_id = 0 AND p.thread_id = $1 ORDER BY p.created DESC, path `
	GetTreePostsByThreadIdQuery2 = ` ) UNION ALL
    	SELECT p.id, array_append(path, p.id) FROM posts AS p JOIN id_tree AS rt ON rt.id = p.parent_id
  	)
	SELECT p.id, f.slug, p.thread_id, u.nickname, p.created, p.parent_id, p.message, p.is_edited `
	SinceMiddle1 = `, CASE WHEN `
	SinceMiddle2 = ` THEN 1 ELSE 0 END AS more`
	GetTreePostsByThreadIdQuery3 = ` FROM id_tree AS it JOIN posts AS p ON it.id = p.id 
	JOIN threads AS t ON p.thread_id = t.id 
	JOIN forums AS f ON t.forum_id = f.id 
	JOIN users AS u ON p.user_id = u.id 
	ORDER BY it.path `
	SinceAfter                   = `) AS r) AS r WHERE r.moreSum > 0`
	UpdatePostQuery              = `UPDATE posts SET message = $2, is_edited = true WHERE id = $1`
	TruncatePostsQuery           = `TRUNCATE TABLE posts CASCADE`
)

func CreatePosts(tx *pgx.Tx, slugOrId string, posts *models.Posts) error {
	forumData := models.Forum{}
	threadData := models.Thread{}
	created := time.Now().Format(DateFormat)

	threadId, err := strconv.ParseInt(slugOrId, 10, 64)
	if err == nil {
		row := tx.QueryRow(GetForumByThreadQuery, threadId, slugOrId)
		err = row.Scan(&forumData.Id, &forumData.Slug, &threadData.Id, &threadData.Slug)
	} else {
		row := tx.QueryRow(GetForumByThreadSlugQuery, slugOrId)
		err = row.Scan(&forumData.Id, &forumData.Slug, &threadData.Id, &threadData.Slug)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.NotFoundError
		}
		return err
	}

	for i := range *posts {
		post := &(*posts)[i]
		post.Forum = forumData.Slug
		post.ThreadId = threadData.Id

		userData, err := GetUser(tx, post.Author)
		if err != nil {
			return err
		}
		post.Author = userData.Nickname

		if post.Parent != 0 {
			row := tx.QueryRow(CheckPostQuery, post.Parent, post.ThreadId)
			err = row.Scan(&post.Parent, &post.ThreadId)
			if err != nil {
				if err == pgx.ErrNoRows {
					return models.ConflictDataError
				}
				return err
			}
		}

		if verifier.IsEmpty(post.Created) {
			post.Created = created
		}

		row := tx.QueryRow(InsertPostQuery, post.ThreadId, userData.Id, post.Created, post.Parent, post.Message)
		err = row.Scan(&post.Id)
		if err != nil {
			if pgxErr, ok := err.(pgx.PgError); ok {
				// Unique violation
				if pgxErr.Code == UniqueErrorCode {
					return models.AlreadyExistsError
				}
			}
			return err
		}

		_, err = tx.Exec(IncPostsCount, forumData.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetPost(tx *pgx.Tx, id int64) (*models.Post, error) {
	var created time.Time
	post := &models.Post{}
	row := tx.QueryRow(GetPostByIdQuery, id)
	err := row.Scan(&post.Id, &post.Forum, &post.ThreadId, &post.Author, &created, &post.Parent, &post.Message, &post.IsEdited)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.NotFoundError
		}
		return nil, err
	}
	post.Created = created.Format(DateFormat)
	return post, nil
}

func GetFullPost(tx *pgx.Tx, id int64, related []string) (*models.PostFull, error) {
	fullPostData := &models.PostFull{}

	postData, err := GetPost(tx, id)
	if err != nil {
		return nil, err
	}
	fullPostData.Post = postData

	for _, v := range related {
		switch v {
		case "user":
			userData, err := GetUser(tx, postData.Author)
			if err != nil {
				return nil, err
			}
			fullPostData.Author = userData
		case "thread":
			threadData, err := GetThreadById(tx, postData.ThreadId)
			if err != nil {
				return nil, err
			}
			fullPostData.Thread = threadData
		case "forum":
			forumData, err := GetForum(tx, postData.Forum)
			if err != nil {
				return nil, err
			}
			fullPostData.Forum = forumData
		}
	}

	return fullPostData, nil
}

func UpdatePost(tx *pgx.Tx, updateData *models.PostUpdate) (*models.Post, error) {
	postData, err := GetPost(tx, updateData.Id)
	if err != nil {
		return postData, err
	}

	if !verifier.IsEmpty(updateData.Message) && postData.Message != updateData.Message {
		_, err = tx.Exec(UpdatePostQuery, updateData.Id, updateData.Message)
		if err != nil {
			return postData, err
		}
		postData.Message = updateData.Message
		postData.IsEdited = true
	}

	return postData, nil
}

func GetThreadPostsQuery(limit, since, sort, desc string) string {
	limitStr := ""
	if limit != "" {
		limitStr = " LIMIT " + limit
	}

	switch sort {
	case "tree":
		pathSort := ""
		if desc == "true" {
			pathSort = " DESC "
		}
		if since != "" {
			return SinceBefore + GetTreePostsByThreadIdQuery1 +
				GetTreePostsByThreadIdQuery2 + SinceMiddle1 + " p.id = " + since + SinceMiddle2 + GetTreePostsByThreadIdQuery3 + pathSort +
				SinceAfter + limitStr
		}
		return GetTreePostsByThreadIdQuery1 + GetTreePostsByThreadIdQuery2 + GetTreePostsByThreadIdQuery3 + pathSort + limitStr
	case "parent_tree":
		if since != "" {
			query1 := ""
			if desc == "true" {
				query1 = ParentTreeSinceQueryDesc1
			} else { // if desc == "false"
				query1 = ParentTreeSinceQuery1
			}

			query2 := ""
			if desc == "true" {
				query2 = ParentTreeSinceQuery2
			} else { // if desc == "false"
				query2 = ParentTreeSinceQueryDesc2
			}

			where := ""
			if limit != "" {
				where = " WHERE parent_number <= " + limit
			}

			return query1 + since + query2 + since + ParentTreeSinceQuery3 + where + ParentTreeSinceQuery4
		}
		query1 := ""
		if desc == "true" {
			query1 = GetTreePostsByThreadIdQuery1Desc
		} else { // if desc == "false"
			query1 = GetTreePostsByThreadIdQuery1
		}

		return query1 + limitStr + GetTreePostsByThreadIdQuery2 + GetTreePostsByThreadIdQuery3
	default:
		where := ""
		if since != "" {
			if desc == "true" {
				where = "p.id < " + since
			} else {
				where = "p.id > " + since
			}
		}

		orderBy := ""
		if desc == "true" {
			orderBy = " ORDER BY p.created DESC, p.id DESC "
		} else { // if desc == "false"
			orderBy = " ORDER BY p.created ASC, p.id ASC "
		}

		if since != "" {
			return GetPostsByThreadIdQuery + " AND " + where + orderBy + limitStr
		}
		return GetPostsByThreadIdQuery + orderBy + limitStr
	}
}

func GetThreadPosts(tx *pgx.Tx, slugOrId, limit, since, sort, desc string) (models.Posts, error) {
	posts := models.Posts{}

	forumId, _, err := GetThreadId(tx, slugOrId)
	if err != nil {
		return posts, err
	}

	rows, err := tx.Query(GetThreadPostsQuery(limit, since, sort, desc), forumId)
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		post := models.Post{}
		var created time.Time
		err = rows.Scan(&post.Id, &post.Forum, &post.ThreadId, &post.Author, &created, &post.Parent, &post.Message, &post.IsEdited)
		if err != nil {
			return posts, err
		}
		post.Created = created.Format(DateFormat)
		posts = append(posts, post)
	}
	return posts, nil
}

func TruncatePosts(tx *pgx.Tx) (err error) {
	_, err = tx.Exec(TruncatePostsQuery)
	return err
}
