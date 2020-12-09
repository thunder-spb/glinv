CREATE TABLE inventory_htags
(
    id      SERIAL NOT NULL PRIMARY KEY,
    value   VARCHAR (100) UNIQUE NOT NULL,
    meta    BOOL DEFAULT false
);