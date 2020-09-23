-- +migrate Up
CREATE TABLE IF NOT EXISTS "public"."event_attendee_bridge"
(
 "generated_id" varchar(36) NOT NULL,
 "volunteer_id" varchar(36) NOT NULL,
 "event_id"     varchar(36) NOT NULL,
 CONSTRAINT "PK_event_attendee_bridge" PRIMARY KEY ( "generated_id" ),
 CONSTRAINT "FK_62" FOREIGN KEY ( "volunteer_id" ) REFERENCES "volunteer"."volunteer" ( "volunteer_id" ),
 CONSTRAINT "FK_65" FOREIGN KEY ( "event_id" ) REFERENCES "public"."events" ( "event_id" )
);

CREATE INDEX "fkIdx_62" ON "event_attendee_bridge"
(
 "volunteer_id"
);

CREATE INDEX "fkIdx_65" ON "event_attendee_bridge"
(
 "event_id"
);

-- +migrate Down
DROP TABLE IF EXISTS "public"."event_attendee_bridge";