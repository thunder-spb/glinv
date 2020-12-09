package pgsql

import (
	"restapi/internal/model"
)

// ServiceRepository ...
type ServiceRepository struct {
	store *Store
}

// GetServiceByEnv ...
func (r *ServiceRepository) GetServiceByEnv(env string) ([]*model.Service, error) {
	stmt := `SELECT id, environment, hostname, ip, location, techname, value, port, approved
	FROM inventory_host_services WHERE approved=true AND environment=$1`

	services := []*model.Service{}

	rows, err := r.store.db.Query(stmt, env)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		service := &model.Service{}
		err := rows.Scan(
			&service.ID,
			&service.Environment,
			&service.Hostname,
			&service.IP,
			&service.Location,
			&service.Title,
			&service.Value,
			&service.Port,
			&service.Approved,
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

// GetServiceByLocation ...
func (r *ServiceRepository) GetServiceByLocation(env, location string) ([]*model.Service, error) {
	stmt := `SELECT id, environment, hostname, ip, location, techname, value, port, approved
	FROM inventory_host_services WHERE approved=true AND environment=$1 AND location=$2`

	services := []*model.Service{}

	rows, err := r.store.db.Query(stmt, env, location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		service := &model.Service{}
		err := rows.Scan(
			&service.ID,
			&service.Environment,
			&service.Hostname,
			&service.IP,
			&service.Location,
			&service.Title,
			&service.Value,
			&service.Port,
			&service.Approved,
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
