-- +migrate Up
CREATE TABLE IF NOT EXISTS "volunteer"."volunteer"
(
 "volunteer_id"  varchar(36) NOT NULL,
 "display_name"  varchar(50) NOT NULL,
 "contact_email" varchar(50) NOT NULL,
 "contact_phone" varchar(50) NOT NULL,
 "bio"           varchar(256) NOT NULL,
 "joined"        date NOT NULL,
 CONSTRAINT "PK_volunteer" PRIMARY KEY ( "volunteer_id" )
);

-- +migrate Down
DROP TABLE volunteer.volunteer;
