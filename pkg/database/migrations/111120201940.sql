-- +migrate Up

-- Diff code generated with pgModeler (PostgreSQL Database Modeler)
-- pgModeler version: 0.9.3-beta1
-- Diff date: 2020-11-11 19:40:56
-- Source model: sayitsocial
-- Database: sayitsocial
-- PostgreSQL version: 12.0

-- [ Diff summary ]
-- Dropped objects: 0
-- Created objects: 2
-- Changed objects: 1
-- Truncated tables: 0

SET search_path=public,pg_catalog,auth,organisation,volunteer;
-- ddl-end --


-- [ Created objects ] --
-- object: organisation.follower_bridge | type: TABLE --
-- DROP TABLE IF EXISTS organisation.follower_bridge CASCADE;
CREATE TABLE organisation.follower_bridge (
	generated_id varchar(36),
	organisation_id varchar(36),
	volunteer_id varchar(36)
);

-- [ Created foreign keys ] --
-- object: "FK_FOL_VOL" | type: CONSTRAINT --
-- ALTER TABLE organisation.follower_bridge DROP CONSTRAINT IF EXISTS "FK_FOL_VOL" CASCADE;
ALTER TABLE organisation.follower_bridge ADD CONSTRAINT "FK_FOL_VOL" FOREIGN KEY (volunteer_id)
REFERENCES volunteer.volunteer (volunteer_id) MATCH FULL
ON DELETE CASCADE ON UPDATE NO ACTION;
-- ddl-end --

