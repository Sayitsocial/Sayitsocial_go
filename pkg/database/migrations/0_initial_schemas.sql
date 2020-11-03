-- +migrate Up
CREATE EXTENSION postgis;

CREATE SCHEMA auth;
CREATE SCHEMA organisation;
CREATE SCHEMA volunteer;

-- +migrate Down
DROP SCHEMA auth;
DROP SCHEMA organisation;
DROP SCHEMA volunteer;
