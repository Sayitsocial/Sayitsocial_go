-- +migrate Up
CREATE TABLE IF NOT EXISTS volunteer.volunteer (
                                                         volunteer_id varchar primary key NOT NULL,
                                                         display_name varchar ,
                                                         contact_email varchar,
                                                         contact_phone varchar,
                                                         bio varchar


);


-- +migrate Down
DROP TABLE volunteer.volunteer;
