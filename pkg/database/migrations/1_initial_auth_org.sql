-- +migrate Up
CREATE TABLE IF NOT EXISTS auth.auth (
                                        username varchar primary key,
                                        password varchar,
                                        typeOfUser varchar

);


-- +migrate Down
DROP TABLE auth.auth;

