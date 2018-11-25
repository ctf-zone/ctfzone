-- VIEWS

-- ----------------------
-- challenges_hints
-- ----------------------

DROP VIEW IF EXISTS challenges_hints;

-- ----------------------
-- challenges_likes
-- ----------------------

DROP VIEW IF EXISTS challenges_likes;

-- ----------------------
-- challenges_solutions
-- ----------------------

DROP VIEW IF EXISTS challenges_solutions;

-- TABLES

-- ---------------
-- settings
-- ---------------

DROP TABLE IF EXISTS settings;
DROP TYPE IF EXISTS scoring_type;

-- ---------------
-- likes
-- ---------------

DROP TABLE IF EXISTS likes;

-- ---------------
-- announcements
-- ---------------

DROP TABLE IF EXISTS announcements;

-- ---------------
-- solutions
-- ---------------

DROP TABLE IF EXISTS solutions;

-- ---------------
-- challenges
-- ---------------

DROP TABLE IF EXISTS challenges;
DROP TYPE IF EXISTS difficulty;

-- ---------------
-- tokens
-- ---------------

DROP TABLE IF EXISTS tokens;
DROP TYPE IF EXISTS token_type;

-- ---------------
-- users
-- ---------------

DROP TABLE IF EXISTS users;
