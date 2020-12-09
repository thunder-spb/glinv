CREATE TABLE server_sysctl_agent (
	id SERIAL NOT NULL PRIMARY KEY,
	server_id INTEGER references server_hard_agent (id) ON DELETE CASCADE,
	name VARCHAR,
	value VARCHAR,
	status_sysctl INTEGER DEFAULT 0 NOT NULL,
	created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON server_sysctl_agent
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();