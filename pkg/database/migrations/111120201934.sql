-- +migrate Up

-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler  version: 0.9.3-beta1
-- PostgreSQL version: 13.0
-- Project Site: pgmodeler.io
-- Model Author: ---

-- Database creation must be performed outside a multi lined SQL file. 
-- These commands were put in this file only as a convenience.
-- 
-- object: sayitsocial | type: DATABASE --
-- DROP DATABASE IF EXISTS sayitsocial;
-- CREATE DATABASE sayitsocial
-- 	ENCODING = 'UTF8'
-- 	LC_COLLATE = 'C.UTF-8'
-- 	LC_CTYPE = 'C.UTF-8'
-- 	TABLESPACE = pg_default
-- 	OWNER = postgres;
-- ddl-end --


-- object: auth | type: SCHEMA --
-- DROP SCHEMA IF EXISTS auth CASCADE;
CREATE SCHEMA auth;
-- ddl-end --

-- object: organisation | type: SCHEMA --
-- DROP SCHEMA IF EXISTS organisation CASCADE;
CREATE SCHEMA organisation;
-- ddl-end --

-- object: volunteer | type: SCHEMA --
-- DROP SCHEMA IF EXISTS volunteer CASCADE;
CREATE SCHEMA volunteer;
-- ddl-end --

SET search_path TO pg_catalog,public,auth,organisation,volunteer;
-- ddl-end --

-- -- object: public.geography | type: TYPE --
-- -- DROP TYPE IF EXISTS public.geography CASCADE;
-- CREATE TYPE public.geography;
-- -- ddl-end --

-- object: auth.auth | type: TABLE --
-- DROP TABLE IF EXISTS auth.auth CASCADE;
CREATE TABLE auth.auth (
	uid character varying NOT NULL,
	username character varying NOT NULL,
	password character varying NOT NULL,
	typeofuser character varying NOT NULL,
	CONSTRAINT auth_pkey PRIMARY KEY (uid),
	CONSTRAINT auth_username_key UNIQUE (username)

);
-- ddl-end --

-- object: volunteer.volunteer | type: TABLE --
-- DROP TABLE IF EXISTS volunteer.volunteer CASCADE;
CREATE TABLE volunteer.volunteer (
	volunteer_id character varying(36) NOT NULL,
	display_name character varying(50) NOT NULL,
	contact_email character varying(50) NOT NULL,
	contact_phone character varying(50) NOT NULL,
	bio character varying(256) NOT NULL,
	joined bigint NOT NULL,
	CONSTRAINT "PK_volunteer" PRIMARY KEY (volunteer_id)

);
-- ddl-end --

-- -- object: public.geography | type: TYPE --
-- -- DROP TYPE IF EXISTS public.geography CASCADE;
-- CREATE TYPE public.geography (
-- 	, INTERNALLENGTH = VARIABLE
-- 	, ALIGNMENT = double precision
-- 	, STORAGE = main
-- 	, DELIMITER = ':'
-- 	, CATEGORY = 'U'
-- );
-- -- ddl-end --
-- COMMENT ON TYPE public.geography IS E'postgis type: The type representing spatial features with geodetic (ellipsoidal) coordinate systems.';
-- -- ddl-end --
-- 
-- object: organisation.organisation | type: TABLE --
-- DROP TABLE IF EXISTS organisation.organisation CASCADE;
CREATE TABLE organisation.organisation (
	organisation_id character varying(36) NOT NULL,
	display_name character varying(36) NOT NULL,
	locality character varying(50) NOT NULL,
	registration_no character varying(50) NOT NULL,
	contact_email character varying(50) NOT NULL,
	contact_phone character varying(50) NOT NULL,
	description character varying(50) NOT NULL,
	achievements character varying(50) NOT NULL,
	owner character varying(36) NOT NULL,
	type_of_org integer NOT NULL,
	location public.geography NOT NULL,
	CONSTRAINT "PK_organisation" PRIMARY KEY (organisation_id)

);
-- ddl-end --

-- object: "fkIdx_26" | type: INDEX --
-- DROP INDEX IF EXISTS organisation."fkIdx_26" CASCADE;
CREATE INDEX "fkIdx_26" ON organisation.organisation
	USING btree
	(
	  owner
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: tbl_geog_gist1 | type: INDEX --
-- DROP INDEX IF EXISTS organisation.tbl_geog_gist1 CASCADE;
CREATE INDEX tbl_geog_gist1 ON organisation.organisation
	USING gist
	(
	  location
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: public.event_category_generated_id_seq | type: SEQUENCE --
-- DROP SEQUENCE IF EXISTS public.event_category_generated_id_seq CASCADE;
CREATE SEQUENCE public.event_category_generated_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

-- ddl-end --

-- object: public.event_category | type: TABLE --
-- DROP TABLE IF EXISTS public.event_category CASCADE;
CREATE TABLE public.event_category (
	generated_id integer NOT NULL DEFAULT nextval('public.event_category_generated_id_seq'::regclass),
	name character varying(50) NOT NULL,
	CONSTRAINT "PK_event_category" PRIMARY KEY (generated_id)

);
-- ddl-end --

-- object: public.events | type: TABLE --
-- DROP TABLE IF EXISTS public.events CASCADE;
CREATE TABLE public.events (
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
-- DROP INDEX IF EXISTS public."fkIdx_53" CASCADE;
CREATE INDEX "fkIdx_53" ON public.events
	USING btree
	(
	  category
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: tbl_geog_gist | type: INDEX --
-- DROP INDEX IF EXISTS public.tbl_geog_gist CASCADE;
CREATE INDEX tbl_geog_gist ON public.events
	USING gist
	(
	  location
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: public.event_host_bridge | type: TABLE --
-- DROP TABLE IF EXISTS public.event_host_bridge CASCADE;
CREATE TABLE public.event_host_bridge (
	generated_id character varying(36) NOT NULL,
	event_id character varying(36) NOT NULL,
	volunteer_id character varying(36),
	organisation_id character varying(36),
	CONSTRAINT "PK_event_host_bridge" PRIMARY KEY (generated_id)

);
-- ddl-end --

-- object: "fkIdx_38" | type: INDEX --
-- DROP INDEX IF EXISTS public."fkIdx_38" CASCADE;
CREATE INDEX "fkIdx_38" ON public.event_host_bridge
	USING btree
	(
	  event_id
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: "fkIdx_41" | type: INDEX --
-- DROP INDEX IF EXISTS public."fkIdx_41" CASCADE;
CREATE INDEX "fkIdx_41" ON public.event_host_bridge
	USING btree
	(
	  volunteer_id
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: "fkIdx_44" | type: INDEX --
-- DROP INDEX IF EXISTS public."fkIdx_44" CASCADE;
CREATE INDEX "fkIdx_44" ON public.event_host_bridge
	USING btree
	(
	  organisation_id
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: public.event_attendee_bridge | type: TABLE --
-- DROP TABLE IF EXISTS public.event_attendee_bridge CASCADE;
CREATE TABLE public.event_attendee_bridge (
	generated_id character varying(36) NOT NULL,
	volunteer_id character varying(36) NOT NULL,
	event_id character varying(36) NOT NULL,
	CONSTRAINT "PK_event_attendee_bridge" PRIMARY KEY (generated_id)

);
-- ddl-end --

-- object: "fkIdx_62" | type: INDEX --
-- DROP INDEX IF EXISTS public."fkIdx_62" CASCADE;
CREATE INDEX "fkIdx_62" ON public.event_attendee_bridge
	USING btree
	(
	  volunteer_id
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: "fkIdx_65" | type: INDEX --
-- DROP INDEX IF EXISTS public."fkIdx_65" CASCADE;
CREATE INDEX "fkIdx_65" ON public.event_attendee_bridge
	USING btree
	(
	  event_id
	)
	WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: "FK_26" | type: CONSTRAINT --
-- ALTER TABLE organisation.organisation DROP CONSTRAINT IF EXISTS "FK_26" CASCADE;
ALTER TABLE organisation.organisation ADD CONSTRAINT "FK_26" FOREIGN KEY (owner)
REFERENCES volunteer.volunteer (volunteer_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: "FK_53" | type: CONSTRAINT --
-- ALTER TABLE public.events DROP CONSTRAINT IF EXISTS "FK_53" CASCADE;
ALTER TABLE public.events ADD CONSTRAINT "FK_53" FOREIGN KEY (category)
REFERENCES public.event_category (generated_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: "FK_38" | type: CONSTRAINT --
-- ALTER TABLE public.event_host_bridge DROP CONSTRAINT IF EXISTS "FK_38" CASCADE;
ALTER TABLE public.event_host_bridge ADD CONSTRAINT "FK_38" FOREIGN KEY (event_id)
REFERENCES public.events (event_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: "FK_41" | type: CONSTRAINT --
-- ALTER TABLE public.event_host_bridge DROP CONSTRAINT IF EXISTS "FK_41" CASCADE;
ALTER TABLE public.event_host_bridge ADD CONSTRAINT "FK_41" FOREIGN KEY (volunteer_id)
REFERENCES volunteer.volunteer (volunteer_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: "FK_44" | type: CONSTRAINT --
-- ALTER TABLE public.event_host_bridge DROP CONSTRAINT IF EXISTS "FK_44" CASCADE;
ALTER TABLE public.event_host_bridge ADD CONSTRAINT "FK_44" FOREIGN KEY (organisation_id)
REFERENCES organisation.organisation (organisation_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: "FK_62" | type: CONSTRAINT --
-- ALTER TABLE public.event_attendee_bridge DROP CONSTRAINT IF EXISTS "FK_62" CASCADE;
ALTER TABLE public.event_attendee_bridge ADD CONSTRAINT "FK_62" FOREIGN KEY (volunteer_id)
REFERENCES volunteer.volunteer (volunteer_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: "FK_65" | type: CONSTRAINT --
-- ALTER TABLE public.event_attendee_bridge DROP CONSTRAINT IF EXISTS "FK_65" CASCADE;
ALTER TABLE public.event_attendee_bridge ADD CONSTRAINT "FK_65" FOREIGN KEY (event_id)
REFERENCES public.events (event_id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- -- object: public.geography | type: TYPE --
-- -- DROP TYPE IF EXISTS public.geography CASCADE;
-- CREATE TYPE public.geography (
-- 	, INTERNALLENGTH = VARIABLE
-- 	, ALIGNMENT = double precision
-- 	, STORAGE = main
-- 	, DELIMITER = ':'
-- 	, CATEGORY = 'U'
-- );
-- -- ddl-end --
-- COMMENT ON TYPE public.geography IS E'postgis type: The type representing spatial features with geodetic (ellipsoidal) coordinate systems.';
-- -- ddl-end --
-- 
-- object: "grant_CU_eb94f049ac" | type: PERMISSION --
GRANT CREATE,USAGE
   ON SCHEMA public
   TO postgres;
-- ddl-end --

-- object: "grant_CU_cd8e46e7b6" | type: PERMISSION --
GRANT CREATE,USAGE
   ON SCHEMA public
   TO PUBLIC;
-- ddl-end --


