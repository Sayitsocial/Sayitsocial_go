-- +migrate Up
CREATE TABLE IF NOT EXISTS "public"."event_host_bridge"
(
 "generated_id"    varchar(36) NOT NULL,
 "event_id"        varchar(36) NOT NULL,
 "volunteer_id"    varchar(36) NULL,
 "organisation_id" varchar(36) NULL,
 CONSTRAINT "PK_event_host_bridge" PRIMARY KEY ( "generated_id" ),
 CONSTRAINT "FK_38" FOREIGN KEY ( "event_id" ) REFERENCES "public"."events" ( "event_id" ),
 CONSTRAINT "FK_41" FOREIGN KEY ( "volunteer_id" ) REFERENCES "volunteer"."volunteer" ( "volunteer_id" ),
 CONSTRAINT "FK_44" FOREIGN KEY ( "organisation_id" ) REFERENCES "organisation"."organisation" ( "organisation_id" )
);

CREATE INDEX "fkIdx_38" ON "public"."event_host_bridge"
(
 "event_id"
);

CREATE INDEX "fkIdx_41" ON "public"."event_host_bridge"
(
 "volunteer_id"
);

CREATE INDEX "fkIdx_44" ON "public"."event_host_bridge"
(
 "organisation_id"
);

-- +migrate Down
DROP TABLE IF EXISTS "public"."event_host_bridge";