CREATE TABLE inventory_global_vars (
	id 		SERIAL NOT NULL PRIMARY KEY,
	title 	VARCHAR (20) UNIQUE NOT NULL,
	value 	VARCHAR (20) UNIQUE NOT NULL,
	meta	BOOL DEFAULT false
);

INSERT INTO inventory_global_vars (title, value)
VALUES
   ('connection', 'ssh');

INSERT INTO inventory_global_vars (title, value)
VALUES
   ('port', 22);