CREATE TABLE inventory_host_vars
(
    id SERIAL NOT NULL PRIMARY KEY,
    host_id INT NOT NULL,
    var_id INT NOT NULL
);