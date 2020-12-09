CREATE TABLE history (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id INTEGER DEFAULT 0 NOT NULL,
    user_email VARCHAR (100),
    entity_id INTEGER DEFAULT 0 NOT NULL,
    entity VARCHAR (25),
    action VARCHAR (25),
    event VARCHAR,
    description VARCHAR,
    created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON history
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();