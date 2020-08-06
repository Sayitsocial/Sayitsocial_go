-- +migrate Up
CREATE TABLE IF NOT EXISTS organisation.organisation (
                                        organisation_id varchar primary key NOT NULL,
                                        display_name varchar,
                                        locality varchar,
                                        registration_no  varchar,
                                        contact_email varchar,
                                        contact_phone varchar,
                                        description varchar,
                                        owner varchar,
                                        achievements varchar


);


-- +migrate Down
DROP TABLE organisation.organisation;
