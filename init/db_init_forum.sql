CREATE TABLE IF NOT EXISTS users
(
  id       BIGSERIAL    NOT NULL
    CONSTRAINT users_pk PRIMARY KEY,
  nickname VARCHAR(32)  NOT NULL,
  email    VARCHAR(255) NOT NULL,
  fullname TEXT         NOT NULL,
  about    TEXT
);

CREATE INDEX IF NOT EXISTS users_email_uindex
  ON users (LOWER(email));

CREATE INDEX IF NOT EXISTS users_nickname_index
  ON users (LOWER(nickname));

-- CREATE INDEX IF NOT EXISTS users_nickname_email_index
--   ON users (LOWER(nickname), LOWER(email));

CREATE TABLE IF NOT EXISTS forums
(
  id      BIGSERIAL    NOT NULL
    CONSTRAINT forums_pk PRIMARY KEY,
  slug    VARCHAR(128) NOT NULL,
  title   VARCHAR(128) NOT NULL,
  user_id BIGINT       NOT NULL,
  posts   BIGINT       NOT NULL DEFAULT 0,
  threads INT          NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS forums_slug_index
  ON forums (LOWER(slug));

CREATE TABLE IF NOT EXISTS threads
(
  id       BIGSERIAL    NOT NULL
    CONSTRAINT threads_pk PRIMARY KEY,
  forum_id BIGINT       NOT NULL,
  user_id  BIGINT       NOT NULL,
  created  TIMESTAMPTZ  NOT NULL,
  slug     VARCHAR(128) NOT NULL,
  title    VARCHAR(128) NOT NULL,
  message  TEXT         NOT NULL,
  votes    INT          NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_threads_id ON threads (id);
CREATE INDEX IF NOT EXISTS idx_threads_slug ON threads (LOWER(slug));
CREATE INDEX IF NOT EXISTS threads_slug_index
  ON threads (forum_id);

CREATE TABLE IF NOT EXISTS posts
(
  id        BIGSERIAL   NOT NULL
    CONSTRAINT posts_pk PRIMARY KEY,
  thread_id BIGINT      NOT NULL,
  user_id   BIGINT      NOT NULL,
  created   TIMESTAMPTZ NOT NULL,
  parent_id BIGINT      NOT NULL DEFAULT 0,
  message   TEXT        NOT NULL,
  is_edited BOOLEAN     NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_posts_id ON posts (id);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id ON posts (thread_id, id);
-- CREATE INDEX IF NOT EXISTS idx_posts_thread_id0 ON posts (thread_id, id) WHERE parent_id = 0;
CREATE INDEX IF NOT EXISTS idx_posts_thread_id_created ON posts (id, created, thread_id);

CREATE TABLE IF NOT EXISTS votes
(
  id        BIGSERIAL   NOT NULL
    CONSTRAINT votes_pk PRIMARY KEY,
  thread_id BIGINT      NOT NULL,
  user_id   BIGINT      NOT NULL,
  vote      BIGINT      NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS votes_thread_user_uindex
  ON votes (thread_id, user_id);


-- CREATE FUNCTION IF NOT EXISTS user_update(in_nickname VARCHAR(32), in_email VARCHAR(255), in_fullname TEXT, in_about TEXT)
--   RETURNS users
--   LANGUAGE plpgsql
-- AS
-- $BODY$
-- DECLARE
--   user_data users;
-- BEGIN
--   SELECT * INTO STRICT user_data FROM users WHERE LOWER(nickname) = LOWER(in_nickname);
--
--   IF in_email != '' THEN
--     user_data.email = in_email;
--   END IF;
--   IF in_fullname != '' THEN
--     user_data.fullname = in_fullname;
--   END IF;
--   IF in_about != '' THEN
--     user_data.about = in_about;
--   END IF;
--
--   UPDATE users
--   SET email    = user_data.email,
--       fullname = user_data.fullname,
--       about    = user_data.about
--   WHERE id = user_data.id;
--
--   RETURN user_data;
-- END
-- $BODY$;
