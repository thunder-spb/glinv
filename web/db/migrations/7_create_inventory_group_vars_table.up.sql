CREATE TABLE inventory_group_vars
(
    id SERIAL NOT NULL PRIMARY KEY,
    group_id INT NOT NULL,
    var_id INT NOT NULL
);