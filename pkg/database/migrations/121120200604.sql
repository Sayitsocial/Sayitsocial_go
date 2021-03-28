-- +migrate Up

-- Diff code generated with pgModeler (PostgreSQL Database Modeler)
-- pgModeler version: 0.9.3-beta1
-- Diff date: 2020-11-12 06:04:54
-- Source model: sayitsocial
-- Database: sayitsocial
-- PostgreSQL version: 12.0

-- [ Diff summary ]
-- Dropped objects: 0
-- Created objects: 4
-- Changed objects: 0
-- Truncated tables: 0

SET check_function_bodies = false;
-- ddl-end --

SET search_path=public,pg_catalog,auth,organisation,volunteer,events;
-- ddl-end --


-- [ Created objects ] --
-- object: followers | type: COLUMN --
-- ALTER TABLE organisation.organisation DROP COLUMN IF EXISTS followers CASCADE;
ALTER TABLE organisation.organisation ADD COLUMN followers bigint NOT NULL DEFAULT 0;
-- ddl-end --


-- object: organisation.update_followers | type: FUNCTION --
-- DROP FUNCTION IF EXISTS organisation.update_followers() CASCADE;
-- +migrate StatementBegin
CREATE FUNCTION organisation.update_followers ()
	RETURNS trigger
	LANGUAGE plpgsql
	VOLATILE 
	CALLED ON NULL INPUT
	SECURITY INVOKER
	COST 1
	AS $$
BEGIN
	IF NEW IS NOT NULL THEN
		UPDATE organisation.organisation SET followers = (SELECT COUNT("generated_id") FROM organisation.follower_bridge) WHERE organisation_id = NEW.organisation_id;
	ELSE 
		UPDATE organisation.organisation SET followers = (SELECT COUNT("generated_id") FROM organisation.follower_bridge) WHERE organisation_id = OLD.organisation_id;
	END IF;
END;

	$$;
-- +migrate StatementEnd
-- ddl-end --
ALTER FUNCTION organisation.update_followers() OWNER TO postgres;
-- ddl-end --

-- object: update_followers | type: TRIGGER --
-- DROP TRIGGER IF EXISTS update_followers ON organisation.follower_bridge CASCADE;
CREATE TRIGGER update_followers
	AFTER INSERT OR DELETE
	ON organisation.follower_bridge
	FOR EACH ROW
	EXECUTE PROCEDURE organisation.update_followers();
-- ddl-end --



-- [ Created foreign keys ] --
-- object: "FK_FOL_ORG" | type: CONSTRAINT --
-- ALTER TABLE organisation.follower_bridge DROP CONSTRAINT IF EXISTS "FK_FOL_ORG" CASCADE;
ALTER TABLE organisation.follower_bridge ADD CONSTRAINT "FK_FOL_ORG" FOREIGN KEY (organisation_id)
REFERENCES organisation.organisation (organisation_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

