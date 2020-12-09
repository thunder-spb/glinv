CREATE TABLE server_soft_base (
	id SERIAL NOT NULL PRIMARY KEY,
	template_id INT references server_templates (id),
	title VARCHAR,
	value VARCHAR,
	created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON server_soft_base
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();