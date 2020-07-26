-- +migrate Up
CREATE TABLE IF NOT EXISTS organisation.organisation (
                                        display_name varchar primary key,
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
