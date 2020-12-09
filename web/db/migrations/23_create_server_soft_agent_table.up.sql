CREATE TABLE server_soft_agent (
	server_id INTEGER references server_hard_agent (id) ON DELETE CASCADE,
	package VARCHAR (255),
	version VARCHAR (255),
	status_package INTEGER DEFAULT 0 NOT NULL,
	created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TRIGGER set_timestamp
BEFORE UPDATE ON server_soft_agent
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();