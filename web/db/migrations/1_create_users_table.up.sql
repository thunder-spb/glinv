CREATE TABLE users (
    id              SERIAL NOT NULL PRIMARY KEY,
    user_name       VARCHAR (50) NOT NULL,
    user_role       INTEGER DEFAULT 0 NOT NULL,
    email           VARCHAR (100) UNIQUE NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);