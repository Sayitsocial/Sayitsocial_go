-- +migrate Up

-- Diff code generated with pgModeler (PostgreSQL Database Modeler)
-- pgModeler version: 0.9.3-beta1
-- Diff date: 2020-11-11 20:03:05
-- Source model: sayitsocial
-- Database: sayitsocial
-- PostgreSQL version: 12.0

-- [ Diff summary ]
-- Dropped objects: 17
-- Created objects: 18
-- Changed objects: 0
-- Truncated tables: 0

SET search_path=public,pg_catalog,auth,organisation,volunteer,events;
-- ddl-end --


-- [ Dropped objects ] --
ALTER TABLE public.event_attendee_bridge DROP CONSTRAINT IF EXISTS "FK_65" CASCADE;
-- ddl-end --
ALTER TABLE public.event_attendee_bridge DROP CONSTRAINT IF EXISTS "FK_62" CASCADE;
-- ddl-end --
ALTER TABLE public.event_host_bridge DROP CONSTRAINT IF EXISTS "FK_44" CASCADE;
-- ddl-end --
ALTER TABLE public.event_host_bridge DROP CONSTRAINT IF EXISTS "FK_41" CASCADE;
-- ddl-end --
ALTER TABLE public.event_host_bridge DROP CONSTRAINT IF EXISTS "FK_38" CASCADE;
-- ddl-end --
ALTER TABLE public.events DROP CONSTRAINT IF EXISTS "FK_53" CASCADE;
-- ddl-end --
DROP INDEX IF EXISTS public."fkIdx_65" CASCADE;
-- ddl-end --
DROP INDEX IF EXISTS public."fkIdx_62" CASCADE;
-- ddl-end --
DROP TABLE IF EXISTS public.event_attendee_bridge CASCADE;
-- ddl-end --
DROP INDEX IF EXISTS public."fkIdx_44" CASCADE;
-- ddl-end --
DROP INDEX IF EXISTS public."fkIdx_41" CASCADE;
-- ddl-end --
DROP INDEX IF EXISTS public."fkIdx_38" CASCADE;
-- ddl-end --
DROP TABLE IF EXISTS public.event_host_bridge CASCADE;
-- ddl-end --
DROP INDEX IF EXISTS public.tbl_geog_gist CASCADE;
-- ddl-end --
DROP INDEX IF EXISTS public."fkIdx_53" CASCADE;
-- ddl-end --
DROP TABLE IF EXISTS public.events CASCADE;
-- ddl-end --
DROP TABLE IF EXISTS public.event_category CASCADE;
-- ddl-end --


-- [ Created objects ] --
-- object: events | type: SCHEMA --
-- DROP SCHEMA IF EXISTS events CASCADE;
CREATE SCHEMA events;
-- ddl-end --

-- object: events.events | type: TABLE --
-- DROP TABLE IF EXISTS events.events CASCADE;
CREATE TABLE events.events (
	event_id character varying(36) NOT NULL,
	name character varying(50) NOT NULL,
	description character varying(50) NOT NULL,
	category integer NOT NULL,
	type_of_event integer NOT NULL,
	start_time bigint NOT NULL,
	host_time bigint NOT NULL,
	location geography NOT NULL,
	trending_index smallint DEFAULT 0,
	CONSTRAINT "PK_events" PRIMARY KEY (event_id)

);
-- ddl-end --

-- object: "fkIdx_53" | type: INDEX --
-- DROP INDEX IF EXISTS events."fkIdx_53" CASCADE;
CREATE INDEX "fkIdx_53" ON events.events
	USING btree
	(
	  category
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: tbl_geog_gist | type: INDEX --
-- DROP INDEX IF EXISTS events.tbl_geog_gist CASCADE;
CREATE INDEX tbl_geog_gist ON events.events
	USING gist
	(
	  location
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: events.event_host_bridge | type: TABLE --
-- DROP TABLE IF EXISTS events.event_host_bridge CASCADE;
CREATE TABLE events.event_host_bridge (
	generated_id character varying(36) NOT NULL,
	event_id character varying(36) NOT NULL,
	volunteer_id character varying(36),
	organisation_id character varying(36),
	CONSTRAINT "PK_event_host_bridge" PRIMARY KEY (generated_id)

);
-- ddl-end --

-- object: "fkIdx_38" | type: INDEX --
-- DROP INDEX IF EXISTS events."fkIdx_38" CASCADE;
CREATE INDEX "fkIdx_38" ON events.event_host_bridge
	USING btree
	(
	  event_id
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: "fkIdx_41" | type: INDEX --
-- DROP INDEX IF EXISTS events."fkIdx_41" CASCADE;
CREATE INDEX "fkIdx_41" ON events.event_host_bridge
	USING btree
	(
	  volunteer_id
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: "fkIdx_44" | type: INDEX --
-- DROP INDEX IF EXISTS events."fkIdx_44" CASCADE;
CREATE INDEX "fkIdx_44" ON events.event_host_bridge
	USING btree
	(
	  organisation_id
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: events.event_attendee_bridge | type: TABLE --
-- DROP TABLE IF EXISTS events.event_attendee_bridge CASCADE;
CREATE TABLE events.event_attendee_bridge (
	generated_id character varying(36) NOT NULL,
	volunteer_id character varying(36) NOT NULL,
	event_id character varying(36) NOT NULL,
	CONSTRAINT "PK_event_attendee_bridge" PRIMARY KEY (generated_id)

);
-- ddl-end --

-- object: "fkIdx_62" | type: INDEX --
-- DROP INDEX IF EXISTS events."fkIdx_62" CASCADE;
CREATE INDEX "fkIdx_62" ON events.event_attendee_bridge
	USING btree
	(
	  volunteer_id
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: "fkIdx_65" | type: INDEX --
-- DROP INDEX IF EXISTS events."fkIdx_65" CASCADE;
CREATE INDEX "fkIdx_65" ON events.event_attendee_bridge
	USING btree
	(
	  event_id
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: events.event_category | type: TABLE --
-- DROP TABLE IF EXISTS events.event_category CASCADE;
CREATE TABLE events.event_category (
	generated_id integer NOT NULL DEFAULT nextval('public.event_category_generated_id_seq'::regclass),
	name character varying(50) NOT NULL,
	CONSTRAINT "PK_event_category" PRIMARY KEY (generated_id)

);
-- ddl-end --

-- [ Created foreign keys ] --
-- object: "FK_53" | type: CONSTRAINT --
-- ALTER TABLE events.events DROP CONSTRAINT IF EXISTS "FK_53" CASCADE;
ALTER TABLE events.events ADD CONSTRAINT "FK_53" FOREIGN KEY (category)
REFERENCES events.event_category (generated_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: "FK_38" | type: CONSTRAINT --
-- ALTER TABLE events.event_host_bridge DROP CONSTRAINT IF EXISTS "FK_38" CASCADE;
ALTER TABLE events.event_host_bridge ADD CONSTRAINT "FK_38" FOREIGN KEY (event_id)
REFERENCES events.events (event_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: "FK_41" | type: CONSTRAINT --
-- ALTER TABLE events.event_host_bridge DROP CONSTRAINT IF EXISTS "FK_41" CASCADE;
ALTER TABLE events.event_host_bridge ADD CONSTRAINT "FK_41" FOREIGN KEY (volunteer_id)
REFERENCES volunteer.volunteer (volunteer_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: "FK_44" | type: CONSTRAINT --
-- ALTER TABLE events.event_host_bridge DROP CONSTRAINT IF EXISTS "FK_44" CASCADE;
ALTER TABLE events.event_host_bridge ADD CONSTRAINT "FK_44" FOREIGN KEY (organisation_id)
REFERENCES organisation.organisation (organisation_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: "FK_62" | type: CONSTRAINT --
-- ALTER TABLE events.event_attendee_bridge DROP CONSTRAINT IF EXISTS "FK_62" CASCADE;
ALTER TABLE events.event_attendee_bridge ADD CONSTRAINT "FK_62" FOREIGN KEY (volunteer_id)
REFERENCES volunteer.volunteer (volunteer_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: "FK_65" | type: CONSTRAINT --
-- ALTER TABLE events.event_attendee_bridge DROP CONSTRAINT IF EXISTS "FK_65" CASCADE;
ALTER TABLE events.event_attendee_bridge ADD CONSTRAINT "FK_65" FOREIGN KEY (event_id)
REFERENCES events.events (event_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

