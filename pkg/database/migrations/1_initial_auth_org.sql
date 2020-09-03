-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.auth (
    uid varchar primary key,
    username varchar unique NOT NULL ,
    password varchar NOT NULL ,
    typeOfUser varchar NOT NULL

);


-- +migrate Down
DROP TABLE auth.auth;

