-- TABLES

-- ---------------
-- users
-- ---------------

CREATE TABLE IF NOT EXISTS
users (
  id            SERIAL    PRIMARY KEY,
  name          TEXT      NOT NULL UNIQUE,
  email         TEXT      NOT NULL UNIQUE,
  password_hash TEXT      NOT NULL,
  is_activated  BOOLEAN   NOT NULL,
  extra         JSONB     NOT NULL,
  created_at    TIMESTAMP NOT NULL,
  updated_at    TIMESTAMP NOT NULL
);

-- ---------------
-- tokens
-- ---------------

DO $$
  BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'token_type') THEN
      CREATE TYPE token_type AS ENUM('activate', 'reset');
    END IF;
  END
$$;

CREATE TABLE IF NOT EXISTS
tokens (
  id         SERIAL     PRIMARY KEY,
  user_id    INTEGER    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  token      TEXT       NOT NULL,
  type       token_type NOT NULL,
  expires_at TIMESTAMP  NOT NULL,
  created_at TIMESTAMP  NOT NULL
);

-- ---------------
-- challenges
-- ---------------

DO $$
  BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'difficulty') THEN
      CREATE TYPE difficulty AS ENUM('easy', 'medium', 'hard');
    END IF;
  END
$$;

CREATE TABLE IF NOT EXISTS
challenges (
  id          SERIAL          PRIMARY KEY,
  title       TEXT            NOT NULL,
  categories  TEXT[]          NOT NULL,
  points      INTEGER         NOT NULL,
  description TEXT            NOT NULL,
  difficulty  difficulty      NOT NULL,
  flag_hash   TEXT            NOT NULL,
  is_locked   BOOLEAN         NOT NULL,
  created_at  TIMESTAMP       NOT NULL,
  updated_at  TIMESTAMP       NOT NULL
);

-- ---------------
-- solutions
-- ---------------

CREATE TABLE IF NOT EXISTS
solutions (
  user_id      INTEGER   NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  challenge_id INTEGER   NOT NULL REFERENCES challenges (id) ON DELETE CASCADE,
  created_at   TIMESTAMP NOT NULL,
  PRIMARY KEY (user_id, challenge_id)
);

-- ---------------
-- announcements
-- ---------------

CREATE TABLE IF NOT EXISTS
announcements (
  id           SERIAL    PRIMARY KEY,
  title        TEXT      NOT NULL,
  body         TEXT      NOT NULL,
  challenge_id INTEGER   REFERENCES challenges (id) ON DELETE CASCADE,
  created_at   TIMESTAMP NOT NULL,
  updated_at   TIMESTAMP NOT NULL
);

-- ---------------
-- likes
-- ---------------

CREATE TABLE IF NOT EXISTS
likes (
  user_id      INTEGER   NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  challenge_id INTEGER   NOT NULL REFERENCES challenges (id) ON DELETE CASCADE,
  created_at   TIMESTAMP NOT NULL,
  PRIMARY KEY (user_id, challenge_id)
);

-- ---------------
-- settings
-- ---------------

DO $$
  BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'scoring_type') THEN
      CREATE TYPE scoring_type AS ENUM('classic', 'dynamic');
    END IF;
  END
$$;

CREATE TABLE IF NOT EXISTS
settings (
    starts_at      TIMESTAMP    NOT NULL,
    ends_at        TIMESTAMP    NOT NULL,
    categories     TEXT[]       NOT NULL,
    scoring_type   scoring_type NOT NULL,
    scoring_params JSONB        NOT NULL
);

-- VIEWS

-- ----------------------
-- challenges_solutions
-- ----------------------

CREATE OR REPLACE VIEW challenges_solutions AS
SELECT
  challenge_id,
  COUNT(*) AS count,
  ARRAY_AGG(user_id) AS users
FROM solutions
GROUP BY challenge_id;

-- ----------------------
-- challenges_likes
-- ----------------------

CREATE OR REPLACE VIEW challenges_likes AS
SELECT
  challenge_id,
  COUNT(*) AS count,
  ARRAY_AGG(user_id) AS users
FROM likes
GROUP BY challenge_id;

-- ----------------------
-- challenges_hints
-- ----------------------

CREATE OR REPLACE VIEW challenges_hints AS
SELECT
  challenge_id,
  COUNT(*) AS count,
  ARRAY_AGG(id) AS announcements
FROM announcements
WHERE challenge_id IS NOT NULL
GROUP BY challenge_id;
