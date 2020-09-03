-- +migrate Up
CREATE SCHEMA auth;
CREATE SCHEMA organisation;
CREATE SCHEMA volunteer;
CREATE SCHEMA bridge;

-- +migrate Down
DROP SCHEMA auth;
DROP SCHEMA organisation;
DROP SCHEMA volunteer;
DROP SCHEMA bridge;
