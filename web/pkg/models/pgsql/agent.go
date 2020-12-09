package pgsql

import (
	"database/sql"

	"glinv/pkg/models"
)

// ServerAgentModel type which wraps a sql.DB connection pool.
type ServerAgentModel struct {
	DB *sql.DB
}

// GetServer will return a specific server based on its id.
func (m *ServerAgentModel) GetServer(id int) (*models.ServerAgent, error) {
	stmt := `SELECT id, hostname, ip, os_name, os_architecture, kernel_release, cpu_model, cpu_num, ram, hdd, uptime, resolv, mtu, delete, status_hard, status_resolv, created, updated 
	FROM server_hard_agent WHERE id = $1`

	server := &models.ServerAgent{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&server.ID,
		&server.Hostname,
		&server.IP,
		&server.OSName,
		&server.OSArchitecture,
		&server.KernelRelease,
		&server.ModelCPU,
		&server.NumCPU,
		&server.RAM,
		&server.HDD,
		&server.Uptime,
		&server.ResolvConf,
		&server.MTU,
		&server.Delete,
		&server.StatusHard,
		&server.StatusResolv,
		&server.Created,
		&server.Updated,
	)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	// If everything went OK then return the server object.
	return server, nil
}

// GetPackages will ...
func (m *ServerAgentModel) GetPackages(id int) ([]*models.PackageServer, error) {
	stmt := `SELECT package, version, status_package
	FROM server_soft_agent WHERE server_id = $1 ORDER BY status_package ASC`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pkgs := []*models.PackageServer{}

	for rows.Next() {
		pkg := &models.PackageServer{}
		err := rows.Scan(
			&pkg.Package,
			&pkg.Version,
			&pkg.StatusPackage,
		)

		if err != nil {
			return nil, err
		}

		pkgs = append(pkgs, pkg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the servers slice.
	return pkgs, nil
}

// GetSysCtl will ...
func (m *ServerAgentModel) GetSysCtl(idServer int) ([]*models.SysCtlServer, error) {
	stmt := `SELECT name, value, status_sysctl
	FROM server_sysctl_agent WHERE server_id = $1`

	rows, err := m.DB.Query(stmt, idServer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sysctl := []*models.SysCtlServer{}

	for rows.Next() {
		s := &models.SysCtlServer{}
		err := rows.Scan(
			&s.Name,
			&s.Value,
			&s.StatusSysCtl,
		)

		if err != nil {
			return nil, err
		}

		sysctl = append(sysctl, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the servers slice.
	return sysctl, nil
}

// GetAll will return created servers.
func (m *ServerAgentModel) GetAll() ([]*models.ServerAgent, error) {
	stmt := `SELECT id, hostname, ip, os_name, os_architecture, kernel_release, cpu_model, cpu_num, ram, hdd, uptime, resolv, mtu, delete, status_hard, status_resolv, created, updated 
	FROM server_hard_agent ORDER BY created DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	servers := []*models.ServerAgent{}

	for rows.Next() {
		server := &models.ServerAgent{}
		err := rows.Scan(
			&server.ID,
			&server.Hostname,
			&server.IP,
			&server.OSName,
			&server.OSArchitecture,
			&server.KernelRelease,
			&server.ModelCPU,
			&server.NumCPU,
			&server.RAM,
			&server.HDD,
			&server.Uptime,
			&server.ResolvConf,
			&server.MTU,
			&server.Delete,
			&server.StatusHard,
			&server.StatusResolv,
			&server.Created,
			&server.Updated,
		)

		if err != nil {
			return nil, err
		}

		servers = append(servers, server)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the servers slice.
	return servers, nil
}

// GetStatuses ...
func (m *ServerAgentModel) GetStatuses() (map[int]int, error) {
	stmt := `SELECT server_id, host_id FROM compare_server_host`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statuses := make(map[int]int)

	for rows.Next() {
		var s, h int
		if err := rows.Scan(&s, &h); err != nil {
			return nil, err
		}

		statuses[s] = h
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return statuses, nil
}

// GetCountServers ...
func (m *ServerAgentModel) GetCountServers() int {
	stmt := `SELECT COUNT(*) FROM server_hard_agent;`
	var count int
	m.DB.QueryRow(stmt).Scan(&count)
	return count
}

// GetCountSofts ...
func (m *ServerAgentModel) GetCountSofts() int {
	stmt := `SELECT COUNT(*) FROM server_soft_agent;`
	var count int
	m.DB.QueryRow(stmt).Scan(&count)
	return count
}

// AlertDisable ...
func (m *ServerAgentModel) AlertDisable(id int) (int, error) {
	stmt := `INSERT INTO server_alerts (server_id, alert) VALUES($1, $2) RETURNING server_id`

	var lastInsertID int
	if err := m.DB.QueryRow(stmt, id, "false").Scan(&lastInsertID); err != nil {
		return lastInsertID, err
	}

	return lastInsertID, nil
}

// AlertEnable ...
func (m *ServerAgentModel) AlertEnable(id int) error {
	stmt := `DELETE FROM server_alerts WHERE server_id = $1;`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// AlertCheck ...
func (m *ServerAgentModel) AlertCheck(id int) (bool, error) {
	stmt := `SELECT alert FROM server_alerts WHERE server_id = $1 limit 1;`
	alert, err := m.DB.Query(stmt, id)

	if !alert.Next() {
		return false, err
	}

	return true, nil
}
