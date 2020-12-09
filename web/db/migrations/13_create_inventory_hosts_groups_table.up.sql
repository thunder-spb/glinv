CREATE TABLE inventory_hosts_groups
(
    id          SERIAL NOT NULL PRIMARY KEY,
    host_id     INT NOT NULL,
    group_id    INT NOT NULL
);