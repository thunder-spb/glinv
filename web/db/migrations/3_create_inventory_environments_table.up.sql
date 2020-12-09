CREATE TABLE inventory_environments (
	id 		SERIAL NOT NULL PRIMARY KEY,
	value 	VARCHAR (20) UNIQUE NOT NULL,
	meta	BOOL DEFAULT false
);