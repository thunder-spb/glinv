CREATE TABLE server_hard_agent (
	id SERIAL NOT NULL PRIMARY KEY,
	hostname VARCHAR (100) NOT NULL,
	ip VARCHAR (15)  NOT NULL,
	os_name VARCHAR (100),
	os_vendor VARCHAR (100),
	os_version VARCHAR (100),
	os_release VARCHAR (100),
	os_architecture VARCHAR (100),
	kernel_release VARCHAR (100),
	cpu_model VARCHAR (100),
	cpu_num VARCHAR (5),
	ram VARCHAR (255),
	hdd VARCHAR (255),
	uptime VARCHAR (50),
	resolv VARCHAR (255),
	timezone VARCHAR (100),
	mtu INTEGER DEFAULT 0 NOT NULL,
	delete INTEGER DEFAULT 0 NOT NULL,
	status_hard INTEGER DEFAULT 0 NOT NULL,
	status_resolv INTEGER DEFAULT 0 NOT NULL,
	created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON server_hard_agent
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();