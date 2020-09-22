-- +migrate Up
CREATE TABLE IF NOT EXISTS "organisation"."organisation"
(
 "organisation_id" varchar(36) NOT NULL,
 "display_name"    varchar(36) NOT NULL,
 "locality"        varchar(50) NOT NULL,
 "registration_no" varchar(50) NOT NULL,
 "contact_email"   varchar(50) NOT NULL,
 "contact_phone"   varchar(50) NOT NULL,
 "description"     varchar(50) NOT NULL,
 "achievements"    varchar(50) NOT NULL,
 "owner"           varchar(36) NOT NULL,
 "type_of_org"     integer NOT NULL,
 CONSTRAINT "PK_organisation" PRIMARY KEY ( "organisation_id" ),
 CONSTRAINT "FK_26" FOREIGN KEY ( "owner" ) REFERENCES "volunteer"."volunteer" ( "volunteer_id" )
);

CREATE INDEX "fkIdx_26" ON "organisation"."organisation"
(
 "owner"
);

-- +migrate Down
DROP INDEX IF EXISTS "fkIdx_26";
DROP TABLE IF EXISTS organisation.organisation;
