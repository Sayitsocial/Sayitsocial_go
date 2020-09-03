-- +migrate Up

ALTER TABLE organisation.organisation
ADD COLUMN type_of_org integer NOT NULL default 0;
