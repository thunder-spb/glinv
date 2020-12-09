CREATE TABLE inventory_hosts (
    id                      SERIAL NOT NULL PRIMARY KEY,
    hostname                VARCHAR (100) UNIQUE NOT NULL,
    ip                      VARCHAR (15) NOT NULL,
    environment             VARCHAR (10) NOT NULL,
    approved                BOOL DEFAULT false,
    delete                  BOOL DEFAULT false,
    reg_to_consul           BOOL DEFAULT false,
    status_in_consul        BOOL DEFAULT false,
    method_check_consul     VARCHAR (20),
    status                  INTEGER DEFAULT 0 NOT NULL,
    description             TEXT,
    tags                    JSONB,
    meta		            JSONB,
    properties              JSONB,
    created                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated                 TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON inventory_hosts
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();