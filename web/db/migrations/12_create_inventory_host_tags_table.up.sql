CREATE TABLE inventory_host_tags
(
    id SERIAL NOT NULL PRIMARY KEY,
    host_id INT NOT NULL,
    tag_id INT NOT NULL
);