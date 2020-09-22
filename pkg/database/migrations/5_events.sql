-- +migrate Up

CREATE TABLE IF NOT EXISTS "public"."events"
(
 "event_id"    varchar(36) NOT NULL,
 "name"        varchar(50) NOT NULL,
 "description" varchar(50) NOT NULL,
 "category"    integer NOT NULL,
 CONSTRAINT "PK_events" PRIMARY KEY ( "event_id" ),
 CONSTRAINT "FK_53" FOREIGN KEY ( "category" ) REFERENCES "event_category" ( "generated_id" )
);

CREATE INDEX "fkIdx_53" ON "public"."events"
(
 "category"
);

-- +migrate Down
DROP INDEX IF EXISTS "fkIdx_53";
DROP TABLE IF EXISTS "public"."events";