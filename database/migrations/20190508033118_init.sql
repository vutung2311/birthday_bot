-- +goose Up
-- +goose StatementBegin
CREATE TABLE birthdays
(
    id          INTEGER      NOT NULL PRIMARY KEY,
    person_name VARCHAR(255) NOT NULL,
    birthday    DATE
);

CREATE TABLE users
(
    id       INTEGER      NOT NULL PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role     INTEGER DEFAULT 2 NOT NULL
);

CREATE UNIQUE INDEX users_id_unique
    ON users (id);

CREATE UNIQUE INDEX birthdays_id_unique
    ON birthdays (id);

CREATE INDEX birthdays_birthday ON birthdays(birthday);

INSERT INTO users (username, password, role)
VALUES ('admin',
        '$2a$10$zggA2NziFJ7F2Hu0o.T/gOREbU0ec70wtOPG681ybU.EFDcVPQUwy',
        1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS birthdays;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
