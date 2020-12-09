CREATE VIEW inventory_host_services
  AS SELECT s.id, h.environment, h.hostname, h.ip, s.type, s.location, s.title, s.link1, s.techname, s.link2, s.domain, s.placement, s.publictime, s.team, s.resp, s.value, s.port, s.approved, s.status_in_consul, s.delete
    FROM inventory_hosts h, inventory_services s
    WHERE h.id = s.host_id;      