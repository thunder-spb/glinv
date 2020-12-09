CREATE TABLE inventory_gvars
(
    id      SERIAL NOT NULL PRIMARY KEY,
    name    VARCHAR (100) NOT NULL,
    value   VARCHAR (100) NOT NULL,
    meta    BOOL DEFAULT false
);