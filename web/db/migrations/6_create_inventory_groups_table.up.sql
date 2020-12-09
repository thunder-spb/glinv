CREATE TABLE inventory_groups (
	id 			SERIAL NOT NULL PRIMARY KEY,
    environment VARCHAR (10) NOT NULL,
    parent_id 	INTEGER DEFAULT 0 NOT NULL,
	value 		VARCHAR (100) NOT NULL,
    approved    BOOL DEFAULT false,
    delete      BOOL DEFAULT false,
    description TEXT,
    tags        JSONB,
	meta		JSONB,
    properties  JSONB,
    created             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON inventory_groups
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

INSERT INTO inventory_groups (environment, parent_id, value)
VALUES
   ('all', 0, 'all');

