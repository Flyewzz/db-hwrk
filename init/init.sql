-- DROP USER IF EXISTS forum;
-- CREATE USER forum WITH
--   LOGIN
--   NOSUPERUSER
--   INHERIT
--   NOCREATEDB
--   NOCREATEROLE
--   NOREPLICATION
--   CONNECTION LIMIT -1
--   PASSWORD 'forum';

-- DROP DATABASE IF EXISTS forum;
-- CREATE DATABASE forum
--   WITH
--   OWNER = forum
--   ENCODING = 'UTF8'
--   LC_COLLATE = 'Russian_Russia.1251'
--   LC_CTYPE = 'Russian_Russia.1251'
--   TABLESPACE = pg_default
--   CONNECTION LIMIT = -1;

DROP TABLE IF EXISTS users;
CREATE TABLE users
(
  id       BIGSERIAL    NOT NULL
    CONSTRAINT users_pk PRIMARY KEY,
  nickname VARCHAR(32)  NOT NULL,
  email    VARCHAR(255) NOT NULL,
  fullname TEXT         NOT NULL,
  about    TEXT
);

CREATE UNIQUE INDEX users_nickname_uindex
  ON users (LOWER(nickname));

CREATE UNIQUE INDEX users_email_uindex
  ON users (LOWER(email));

CREATE INDEX users_nickname_index
  ON users (LOWER(nickname));

CREATE INDEX users_nickname_email_index
  ON users (LOWER(nickname), LOWER(email));

DROP TABLE IF EXISTS forums;
CREATE TABLE forums
(
  id      BIGSERIAL    NOT NULL
    CONSTRAINT forums_pk PRIMARY KEY,
  slug    VARCHAR(128) NOT NULL,
  title   VARCHAR(128) NOT NULL,
  user_id BIGINT       NOT NULL,
  posts   BIGINT       NOT NULL DEFAULT 0,
  threads INT          NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX forums_slug_uindex
  ON forums (LOWER(slug));

CREATE INDEX forums_slug_index
  ON forums (LOWER(slug));

DROP TABLE IF EXISTS threads;
CREATE TABLE threads
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

CREATE INDEX threads_slug_index
  ON threads (forum_id);

DROP TABLE IF EXISTS posts;
CREATE TABLE posts
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

CREATE INDEX posts_slug_index
  ON posts (thread_id);

DROP TABLE IF EXISTS votes;
CREATE TABLE votes
(
  id        BIGSERIAL   NOT NULL
    CONSTRAINT votes_pk PRIMARY KEY,
  thread_id BIGINT      NOT NULL,
  user_id   BIGINT      NOT NULL,
  vote      BIGINT      NOT NULL
);

CREATE INDEX votes_thread_user_index
  ON votes (thread_id, user_id);

CREATE UNIQUE INDEX votes_thread_user_uindex
  ON votes (thread_id, user_id);
