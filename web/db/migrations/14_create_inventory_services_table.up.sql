CREATE TABLE inventory_services
(
    id                  SERIAL NOT NULL PRIMARY KEY,
    host_id             INT references inventory_hosts (id),
    type                VARCHAR (50) NOT NULL,
    location            VARCHAR (50) NOT NULL,
    title               VARCHAR (100),
    link1               VARCHAR (255),
    techname            VARCHAR (100) NOT NULL,
    link2               VARCHAR (255),
    domain              VARCHAR (100),
    placement           VARCHAR (100),
    publictime          VARCHAR (255),
    team                VARCHAR (100),
    resp                VARCHAR (100),
    value               VARCHAR (100) NOT NULL,
    port                INT NOT NULL,
    approved            BOOL DEFAULT false,
    delete              BOOL DEFAULT false,
    reg_to_consul       BOOL DEFAULT false,
    status_in_consul    BOOL DEFAULT false,
    method_check_consul VARCHAR (20),
    description         TEXT,
    tags                JSONB,
    meta		        JSONB,
    variables           JSONB,
    properties          JSONB,
    created             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON inventory_services
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();