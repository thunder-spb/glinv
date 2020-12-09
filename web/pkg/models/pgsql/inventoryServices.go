package pgsql

import (
	"database/sql"

	"glinv/pkg/models"
)

// InventoryServiceModel type which wraps a sql.DB connection pool.
type InventoryServiceModel struct {
	DB *sql.DB
}

// Delete ...
func (m *InventoryServiceModel) Delete(id string) error {
	stmt := `DELETE FROM inventory_services WHERE id = $1;`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// MarkServiceForDeletion ...
func (m *InventoryServiceModel) MarkServiceForDeletion(id int) error {
	stmt := `UPDATE inventory_services SET delete = $1 WHERE id = $2;`

	_, err := m.DB.Exec(stmt, "true", id)
	if err != nil {
		return err
	}
	return nil
}

// CancelDeletion ...
func (m *InventoryServiceModel) CancelDeletion(id int) error {
	stmt := `UPDATE inventory_services SET delete = $1 WHERE id = $2;`

	_, err := m.DB.Exec(stmt, "false", id)
	if err != nil {
		return err
	}
	return nil
}

// Update ...
func (m *InventoryServiceModel) Update(kind, techname, domain, placement, team, resp, title, value, location, reg, method, description, port string, idHost, idService int) error {
	stmt := `UPDATE inventory_services SET host_id = $2, type =$3, location = $4, title = $5, techname = $6, domain = $7, placement = $8, team = $9, resp = $10, value = $11, reg_to_consul = $12, method_check_consul = $13, port = $14, approved = $15, description = $16
	WHERE id = $1`

	_, err := m.DB.Exec(stmt, idService, idHost, kind, location, title, techname, domain, placement, team, resp, value, reg, method, port, "false", description)
	if err != nil {
		return err
	}

	return nil
}

// Insert will insert a new service into the database.
func (m *InventoryServiceModel) Insert(ids []string, location, kind, title, techname, domain, value, placement, team, resp, reg, method, port, description string) error {
	stmt := `INSERT INTO inventory_services (host_id, location, type, title, techname, domain, value, placement, team, resp, reg_to_consul, method_check_consul, port, description)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	for _, v := range ids {

		_, err := m.DB.Exec(stmt, v, location, kind, title, techname, domain, value, placement, team, resp, reg, method, port, description)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get will return a specific server based on its id.
func (m *InventoryServiceModel) Get(id int) (*models.InventoryService, error) {
	stmt := `SELECT id, type, location, title, techname, domain, value, placement, team, resp, port, reg_to_consul, status_in_consul, method_check_consul, approved, delete, description, created, updated
	FROM inventory_services WHERE id = $1`

	service := &models.InventoryService{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&service.ID,
		&service.Type,
		&service.Location,
		&service.Title,
		&service.TechName,
		&service.Domain,
		&service.Value,
		&service.Placement,
		&service.Team,
		&service.Resp,
		&service.Port,
		&service.RegToConsul,
		&service.StatusInConsul,
		&service.MethodCheckConsul,
		&service.Approved,
		&service.Delete,
		&service.Description,
		&service.Created,
		&service.Updated,
	)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	// If everything went OK then return the server object.
	return service, nil
}

// GetAll will return the created services with filtering by environment.
func (m *InventoryServiceModel) GetAll() ([]*models.InventoryService, error) {
	stmt := `SELECT id, type, location, title, techname, domain, value, placement, team, resp, port, approved, delete, description, properties, created, updated
	FROM inventory_services ORDER BY title ASC `

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := []*models.InventoryService{}

	for rows.Next() {
		service := &models.InventoryService{}
		err := rows.Scan(
			&service.ID,
			&service.Type,
			&service.Location,
			&service.Title,
			&service.TechName,
			&service.Domain,
			&service.Value,
			&service.Placement,
			&service.Team,
			&service.Resp,
			&service.Port,
			&service.Approved,
			&service.Delete,
			&service.Description,
			&service.Properties,
			&service.Created,
			&service.Updated,
		)

		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// GetAllForApprove will ...
func (m *InventoryServiceModel) GetAllForApprove() ([]*models.InventoryService, error) {
	stmt := `SELECT id, title, techname, value, port, approved, description, properties, created, updated
	FROM inventory_services WHERE approved = false`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := []*models.InventoryService{}

	for rows.Next() {
		service := &models.InventoryService{}
		err := rows.Scan(
			&service.ID,
			&service.Title,
			&service.TechName,
			&service.Value,
			&service.Port,
			&service.Approved,
			&service.Description,
			&service.Properties,
			&service.Created,
			&service.Updated,
		)

		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// GetAllForDelete will ...
func (m *InventoryServiceModel) GetAllForDelete() ([]*models.InventoryService, error) {
	stmt := `SELECT id, title, techname, value, port, delete, description, properties, created, updated
	FROM inventory_services WHERE delete=true`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := []*models.InventoryService{}

	for rows.Next() {
		service := &models.InventoryService{}
		err := rows.Scan(
			&service.ID,
			&service.Title,
			&service.TechName,
			&service.Value,
			&service.Port,
			&service.Delete,
			&service.Description,
			&service.Properties,
			&service.Created,
			&service.Updated,
		)

		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// GetServicesByEnv will return from view host_services.
func (m *InventoryServiceModel) GetServicesByEnv(env string) ([]*models.InventoryService, error) {
	stmtAll := `SELECT id, environment, hostname, ip, type, location, title, techname, domain, value, placement, team, resp, port, approved, status_in_consul, delete
	FROM inventory_host_services ORDER BY title ASC`

	stmtEnv := `SELECT id, environment, hostname, ip, type, location, title, techname, domain, value, placement, team, resp, port, approved, status_in_consul, delete
	FROM inventory_host_services WHERE environment=$1 ORDER BY title ASC`

	hostServices := []*models.InventoryService{}

	if env == "" || env == "all" {
		rows, err := m.DB.Query(stmtAll)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			hostService := &models.InventoryService{}
			err := rows.Scan(
				&hostService.ID,
				&hostService.Host.Environment,
				&hostService.Host.Hostname,
				&hostService.Host.IP,
				&hostService.Type,
				&hostService.Location,
				&hostService.Title,
				&hostService.TechName,
				&hostService.Domain,
				&hostService.Value,
				&hostService.Placement,
				&hostService.Team,
				&hostService.Resp,
				&hostService.Port,
				&hostService.Approved,
				&hostService.StatusInConsul,
				&hostService.Delete,
			)

			if err != nil {
				return nil, err
			}

			hostServices = append(hostServices, hostService)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	} else {
		rows, err := m.DB.Query(stmtEnv, env)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			hostService := &models.InventoryService{}
			err := rows.Scan(
				&hostService.ID,
				&hostService.Host.Environment,
				&hostService.Host.Hostname,
				&hostService.Host.IP,
				&hostService.Type,
				&hostService.Location,
				&hostService.Title,
				&hostService.TechName,
				&hostService.Domain,
				&hostService.Value,
				&hostService.Placement,
				&hostService.Team,
				&hostService.Resp,
				&hostService.Port,
				&hostService.Approved,
				&hostService.StatusInConsul,
				&hostService.Delete,
			)

			if err != nil {
				return nil, err
			}

			hostServices = append(hostServices, hostService)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return hostServices, nil
}

// GetServicesByHostname will return from view host_services.
func (m *InventoryServiceModel) GetServicesByHostname(hostname string) ([]*models.InventoryService, error) {
	stmt := `SELECT id, environment, hostname, ip, title, techname, value, port, approved, delete, status_in_consul
	FROM inventory_host_services WHERE hostname = $1 ORDER BY title ASC`

	hostServices := []*models.InventoryService{}
	rows, err := m.DB.Query(stmt, hostname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		hostService := &models.InventoryService{}
		err := rows.Scan(
			&hostService.ID,
			&hostService.Host.Environment,
			&hostService.Host.Hostname,
			&hostService.Host.IP,
			&hostService.Title,
			&hostService.TechName,
			&hostService.Value,
			&hostService.Port,
			&hostService.Approved,
			&hostService.Delete,
			&hostService.StatusInConsul,
		)

		if err != nil {
			return nil, err
		}

		hostServices = append(hostServices, hostService)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return hostServices, nil
}

// GetHostsService ...
func (m *InventoryServiceModel) GetHostsService(id int) ([]*models.InventoryHost, error) {
	stmt := `SELECT host_id FROM inventory_services WHERE id = $1`

	var hosts []*models.InventoryHost

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		host := &models.InventoryHost{}
		err := rows.Scan(&host.ID)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

// ApproveService ...
func (m *InventoryServiceModel) ApproveService(id string) error {
	stmt := `UPDATE inventory_services SET approved = $1 WHERE id = $2;`
	_, err := m.DB.Exec(stmt, "true", id)
	if err != nil {
		return err
	}
	return nil
}

// GetCountServices ...
func (m *InventoryServiceModel) GetCountServices() int {
	stmt := `SELECT COUNT(*) FROM inventory_services;`
	var count int
	m.DB.QueryRow(stmt).Scan(&count)
	return count
}

// GetCountServicesByHost ...
func (m *InventoryServiceModel) GetCountServicesByHost(hostname string) int {
	stmt := `SELECT COUNT(*) FROM inventory_host_services WHERE hostname = $1;`
	var count int
	m.DB.QueryRow(stmt, hostname).Scan(&count)
	return count
}

// GetCountUnapprovedServices ...
func (m *InventoryServiceModel) GetCountUnapprovedServices() int {
	stmtApproved := `SELECT COUNT(*) FROM inventory_services WHERE approved = $1;`
	var countApproved int
	m.DB.QueryRow(stmtApproved, "false").Scan(&countApproved)

	stmtDelete := `SELECT COUNT(*) FROM inventory_services WHERE delete = $1;`
	var countDelete int
	m.DB.QueryRow(stmtDelete, "true").Scan(&countDelete)

	return countApproved + countDelete
}
