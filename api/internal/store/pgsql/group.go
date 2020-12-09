package pgsql

import "restapi/internal/model"

// GroupRepository ...
type GroupRepository struct {
	store *Store
}

// GetAll ...
func (r *GroupRepository) GetAll() ([]*model.Group, error) {
	stmt := `SELECT id, value
	FROM inventory_groups ORDER BY value ASC`

	groups := []*model.Group{}

	rows, err := r.store.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		group := &model.Group{}
		err := rows.Scan(
			&group.ID,
			&group.Value,
		)

		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// func (r *GroupRepository) GetGVars(hostID int) ([]*model.GVar, error) {
// 	// Get slice ID of vars
// 	stmtHostVar := `SELECT var_id
// 		FROM inventory_host_vars WHERE host_id = $1`

// 	rows, err := r.store.db.Query(stmtHostVar, hostID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var idHostVars []int

// 	for rows.Next() {
// 		var v int
// 		err := rows.Scan(&v)

// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}

// 		idHostVars = append(idHostVars, v)
// 	}

// 	if idHostVars == nil {
// 		idHostVars = []int{0}
// 	}

// 	// Get values of vars
// 	stmtVarValue := `SELECT name, value
// 	FROM inventory_hvars WHERE id IN (`

// 	buf := bytes.NewBufferString(stmtVarValue)

// 	for i, v := range idHostVars {
// 		if i > 0 {
// 			buf.WriteString(",")
// 		}

// 		buf.WriteString(strconv.Itoa(v))
// 	}

// 	buf.WriteString(")")

// 	rows, err = r.store.db.Query(buf.String())
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	varValues := []*model.HVar{}

// 	for rows.Next() {
// 		vn := &model.HVar{}
// 		err := rows.Scan(
// 			&vn.Name,
// 			&vn.Value,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		varValues = append(varValues, vn)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return varValues, nil
// }

// func (r *GroupRepository) GetGroupByEnv(env string) ([]*model.Group, error) {
// 	stmt := `SELECT id, hostname, ip, environment
// 	FROM inventory_hosts WHERE approved=true AND environment=$1 ORDER BY hostname ASC`

// 	hosts := []*model.Host{}

// 	rows, err := r.store.db.Query(stmt, env)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		host := &model.Host{}
// 		err := rows.Scan(
// 			&host.ID,
// 			&host.Hostname,
// 			&host.IP,
// 			&host.Environment,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		hosts = append(hosts, host)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return hosts, nil
// }
