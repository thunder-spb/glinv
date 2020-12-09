CREATE TABLE inventory_hvars
(
    id      SERIAL NOT NULL PRIMARY KEY,
    name    VARCHAR (100) NOT NULL,
    value   VARCHAR (100) NOT NULL,
    meta    BOOL DEFAULT false
);