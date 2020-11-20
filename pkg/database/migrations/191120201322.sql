-- +migrate Up

-- Diff code generated with pgModeler (PostgreSQL Database Modeler)
-- pgModeler version: 0.9.3-beta1
-- Diff date: 2020-11-19 13:22:24
-- Source model: sayitsocial
-- Database: sayitsocial
-- PostgreSQL version: 12.0

-- [ Diff summary ]
-- Dropped objects: 1
-- Created objects: 0
-- Changed objects: 0
-- Truncated tables: 0

SET check_function_bodies = false;
-- ddl-end --

SET search_path=public,pg_catalog,auth,organisation,volunteer,events;
-- ddl-end --


-- [ Dropped objects ] --
ALTER TABLE organisation.organisation DROP CONSTRAINT IF EXISTS "FK_26" CASCADE;
-- ddl-end --

