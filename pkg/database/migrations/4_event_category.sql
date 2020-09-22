-- +migrate Up

CREATE TABLE IF NOT EXISTS "public"."event_category"
(
 "generated_id" serial NOT NULL,
 "name"         varchar(50) NOT NULL,
 CONSTRAINT "PK_event_category" PRIMARY KEY ( "generated_id" )
);

-- +migrate Down
DROP TABLE IF EXISTS "public"."event_category";
