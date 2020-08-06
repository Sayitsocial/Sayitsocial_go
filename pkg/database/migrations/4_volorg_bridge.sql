-- +migrate Up
CREATE TABLE IF NOT EXISTS bridge.volorg (
                                                   volunteer_id varchar NOT NULL ,
                                                   organisation_id varchar NOT NULL ,
                                                   PRIMARY KEY(volunteer_id, organisation_id)


);


-- +migrate Down
DROP TABLE  bridge.volorg;
