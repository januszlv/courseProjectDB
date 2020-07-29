DROP TRIGGER IF EXISTS user_update_tg ON "user";
DROP TRIGGER IF EXISTS user_delete_tg ON "user";
DROP TRIGGER IF EXISTS game_update_tg ON game;
DROP TRIGGER IF EXISTS game_delete_tg ON game;
DROP TRIGGER IF EXISTS match_process_tg ON match;
DROP TRIGGER IF EXISTS move_process_tg ON move;


DROP FUNCTION IF EXISTS check_update_id;
DROP FUNCTION IF EXISTS change_deleted;
DROP FUNCTION IF EXISTS match_process;
DROP FUNCTION IF EXISTS move_process;

DROP TABLE IF EXISTS leaderboard;
DROP TABLE IF EXISTS friendship;
DROP TABLE IF EXISTS move;
DROP TABLE IF EXISTS match;
DROP TABLE IF EXISTS game;
DROP TABLE IF EXISTS "user";

CREATE TABLE "user"
(
  id        INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
  user_name varchar(50) unique                           NOT NULL
    CONSTRAINT empty_name_check CHECK (length(user_name) > 0 ),
  full_name varchar(50)                                  NOT NULL,
  password  varchar(50)                                  NOT NULL
    CONSTRAINT valid_password CHECK ( length(password) > 6 ),
  motto     varchar(100),
  deleted   boolean DEFAULT false
);

CREATE TABLE game
(
  id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
  name        varchar(50) unique                           NOT NULL
    CONSTRAINT empty_game_name_check CHECK ( length(name) > 0 ),
  description varchar(100) unique                          NOT NULL
    CONSTRAINT empty_game_description_check CHECK ( length(description) > 0 ),
  deleted     boolean DEFAULT false
);

CREATE TABLE match
(
  id       INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
  user_id  INT                                          NOT NULL REFERENCES "user" (id),
  game_id  INT                                          NOT NULL REFERENCES game (id),
  score    INT,
  datetime timestamp with time zone                     NOT NULL DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE move
(
  id       INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
  match_id INT                                          NOT NULL REFERENCES match (id),
  datetime timestamp with time zone                     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  info     text                                         NOT NULL

);

CREATE TABLE friendship
(
  id       INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
  user_id  INT                                          NOT NULL REFERENCES "user" (id),
  user2_id INT                                          NOT NULL REFERENCES "user" (id),

  unique (user_id, user2_id)
);

CREATE TABLE leaderboard
(
  id      INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
  user_id INT                                          NOT NULL REFERENCES "user" (id),
  game_id INT                                          NOT NULL REFERENCES game (id),
  score   INT                                          NOT NULL,

  unique (user_id, game_id)
);

-- общая функция для проверки, не меняется ли id
CREATE FUNCTION check_update_id() RETURNS trigger AS
$$
BEGIN
  IF NEW.id <> OLD.id THEN
    RAISE EXCEPTION 'ID cannot be changed';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- общая функция для смены флага deleted
CREATE FUNCTION change_deleted() RETURNS trigger AS
$$
BEGIN
  EXECUTE format('UPDATE %I.%I SET deleted=true WHERE id=%s;', TG_TABLE_SCHEMA, TG_TABLE_NAME, OLD.id);
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- функция для удаления и обновления match
CREATE FUNCTION match_process() RETURNS trigger AS
$$
BEGIN
  IF (TG_OP = 'DELETE') THEN
    RAISE EXCEPTION 'Match cannot be deleted';
  ELSEIF (TG_OP = 'UPDATE') THEN
    IF NEW.id <> OLD.id OR
       NEW.user_id <> OLD.user_id OR
       NEW.game_id <> OLD.game_id OR
       NEW.datetime <> OLD.datetime THEN
      RAISE EXCEPTION 'ID cannot be changed';
    END IF;
    RETURN NEW;
  END IF;
END;
$$ LANGUAGE plpgsql;

-- функция для удаления, обновления и вставки move
CREATE FUNCTION move_process() RETURNS trigger AS
$$
DECLARE
  cur_score INT;
BEGIN
  IF (TG_OP = 'DELETE') THEN
    RAISE EXCEPTION 'Move cannot be deleted';
  ELSEIF (TG_OP = 'UPDATE') THEN
    RAISE EXCEPTION 'Move cannot be updated';
  ELSEIF (TG_OP = 'INSERT') THEN
    SELECT INTO cur_score score FROM match WHERE id = NEW.match_id;
    IF cur_score != 0 THEN
      RAISE EXCEPTION 'cannot add move to ended match';
    END IF;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- триггеры
CREATE TRIGGER user_update_tg
  BEFORE UPDATE
  ON "user"
  FOR EACH ROW
EXECUTE PROCEDURE check_update_id();

CREATE TRIGGER user_delete_tg
  BEFORE DELETE
  ON "user"
  FOR EACH ROW
EXECUTE PROCEDURE change_deleted();

CREATE TRIGGER game_update_tg
  BEFORE UPDATE
  ON game
  FOR EACH ROW
EXECUTE PROCEDURE check_update_id();

CREATE TRIGGER game_delete_tg
  BEFORE DELETE
  ON game
  FOR EACH ROW
EXECUTE PROCEDURE change_deleted();

CREATE TRIGGER match_process_tg
  BEFORE DELETE OR UPDATE
  ON match
  FOR EACH ROW
EXECUTE PROCEDURE match_process();

CREATE TRIGGER move_process_tg
  BEFORE UPDATE OR INSERT OR DELETE
  ON move
  FOR EACH ROW
EXECUTE PROCEDURE move_process();

-- triggers check
-- INSERT INTO "user" (user_name, full_name, password)
-- VALUES ('user name', 'full name', 'password');
SELECT *
FROM "user";
SELECT *
FROM game;
SELECT *
FROM match;
SELECT *
FROM move;
SELECT *
FROM friendship;
SELECT *
FROM leaderboard;
--
-- INSERT INTO game (name, description)
-- VALUES ('game name', 'description of game');
-- SELECT *
-- FROM game;
--
-- INSERT INTO match (user_id, game_id, score)
-- VALUES (1, 1, NULL);
-- SELECT *
-- FROM match;
--
-- INSERT INTO move (match_id, info)
-- VALUES (1, 'info');
-- SELECT *
-- FROM move;
--
-- -- проверка триггеров для юзера
-- UPDATE "user"
-- SET user_name = 'Januszzz', motto = 'play harder...'
-- WHERE id = 1; -- (+)
-- UPDATE "user"
-- SET user_name = 'Janusz'
-- WHERE id = 1; -- (+)
--
-- DELETE
-- FROM "user"
-- WHERE id = 1; -- (+)
--
-- -- проверка триггеров для game
-- UPDATE game
-- SET id = 2
-- WHERE id = 1; -- (+)
-- UPDATE game
-- SET name = 'OTHER game name'
-- WHERE id = 1; -- (+)
--
-- DELETE
-- FROM game
-- WHERE id = 1; -- (+)
--
-- -- проверка триггеров для match
-- UPDATE match
-- SET user_id = 3
-- WHERE user_id = 1; -- (+)
-- UPDATE match
-- SET score = NULL
-- WHERE id = 1; -- (+)
--
-- DELETE
-- FROM match
-- WHERE id = 1; -- (+)
--
-- -- проверка триггеров для move
-- INSERT INTO move (match_id, info)
-- VALUES (1, 'interesting game'); -- (+)
--
-- UPDATE move
-- SET info = 'NEW INFO'
-- WHERE id = 1; -- (+)
--
-- DELETE
-- FROM move
-- WHERE id = 1; -- (+)

-- SELECT *
-- FROM leaderboard
-- WHERE id = 11
--   AND deleted IS NOT TRUE
