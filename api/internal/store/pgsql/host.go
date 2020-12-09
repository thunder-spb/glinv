package pgsql

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"restapi/internal/model"
	"strconv"
)

// HostRepository ...
type HostRepository struct {
	store *Store
}

// GetAll ...
func (r *HostRepository) GetAll() ([]*model.Host, error) {
	stmt := `SELECT id, hostname, ip, environment
	FROM inventory_hosts WHERE approved=true ORDER BY hostname ASC`

	hosts := []*model.Host{}

	rows, err := r.store.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		host := &model.Host{}
		err := rows.Scan(
			&host.ID,
			&host.Hostname,
			&host.IP,
			&host.Environment,
		)

		if err != nil {
			return nil, err
		}

		hosts = append(hosts, host)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

// GetHostVars ... ?
func (r *HostRepository) GetHostVars(hostID int) ([]*model.HVar, error) {
	// Get slice ID of vars
	stmtHostVar := `SELECT var_id
		FROM inventory_host_vars WHERE host_id = $1`

	rows, err := r.store.db.Query(stmtHostVar, hostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var idHostVars []int

	for rows.Next() {
		var v int
		err := rows.Scan(&v)

		if err != nil {
			log.Println(err)
			continue
		}

		idHostVars = append(idHostVars, v)
	}

	if idHostVars == nil {
		idHostVars = []int{0}
	}

	// Get values of vars
	stmtVarValue := `SELECT name, value
	FROM inventory_hvars WHERE id IN (`

	buf := bytes.NewBufferString(stmtVarValue)

	for i, v := range idHostVars {
		if i > 0 {
			buf.WriteString(",")
		}

		buf.WriteString(strconv.Itoa(v))
	}

	buf.WriteString(")")

	rows, err = r.store.db.Query(buf.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	varValues := []*model.HVar{}

	for rows.Next() {
		vn := &model.HVar{}
		err := rows.Scan(
			&vn.Name,
			&vn.Value,
		)

		if err != nil {
			return nil, err
		}

		varValues = append(varValues, vn)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return varValues, nil
}

// GetHostByEnv ...
func (r *HostRepository) GetHostByEnv(env string) ([]*model.Host, error) {
	stmt := `SELECT id, hostname, ip, environment
	FROM inventory_hosts WHERE approved=true AND environment=$1 ORDER BY hostname ASC`

	hosts := []*model.Host{}

	rows, err := r.store.db.Query(stmt, env)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		host := &model.Host{}
		err := rows.Scan(
			&host.ID,
			&host.Hostname,
			&host.IP,
			&host.Environment,
		)

		if err != nil {
			return nil, err
		}

		hosts = append(hosts, host)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

// GetIDTagsByValueTags ...
func (r *HostRepository) GetIDTagsByValueTags(valueTags []string) ([]int, error) {
	stmtIDTags := `SELECT id FROM inventory_htags WHERE value IN (`
	buf := bytes.NewBufferString(stmtIDTags)

	for i, v := range valueTags {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf("'%s'", v))
	}
	buf.WriteString(")")

	rows, err := r.store.db.Query(buf.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var idTags []int

	for rows.Next() {
		var v int
		err := rows.Scan(&v)

		if err != nil {
			log.Println(err)
			continue
		}

		idTags = append(idTags, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return idTags, nil
}

// GetAllIDHostsByIDTags ...
func (r *HostRepository) GetAllIDHostsByIDTags(idTags []int) ([]int, error) {
	stmtIDHosts := `SELECT host_id FROM inventory_host_tags WHERE tag_id IN (`
	buf := bytes.NewBufferString(stmtIDHosts)

	for i, v := range idTags {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(v))
	}
	buf.WriteString(")")

	rows, err := r.store.db.Query(buf.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var idHosts []int

	for rows.Next() {
		var v int
		err := rows.Scan(&v)

		if err != nil {
			log.Println(err)
			continue
		}

		idHosts = append(idHosts, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return idHosts, nil
}

// GetIDVarsByValueVars ...
func (r *HostRepository) GetIDVarsByValueVars(typeVars string, valueVars []string) ([]int, error) {
	stmt := `SELECT id FROM inventory_hvars WHERE name=$1 AND value IN (`
	buf := bytes.NewBufferString(stmt)

	for i, v := range valueVars {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf("'%s'", v))
	}
	buf.WriteString(")")

	rows, err := r.store.db.Query(buf.String(), typeVars)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var idVars []int

	for rows.Next() {
		var v int
		err := rows.Scan(&v)

		if err != nil {
			log.Println(err)
			continue
		}

		idVars = append(idVars, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return idVars, nil
}

// GetAllIDHostsByIDVars ...
func (r *HostRepository) GetAllIDHostsByIDVars(idVars []int) ([]int, error) {
	stmtIDHosts := `SELECT host_id FROM inventory_host_vars WHERE var_id IN (`
	buf := bytes.NewBufferString(stmtIDHosts)

	for i, v := range idVars {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(v))
	}
	buf.WriteString(")")

	rows, err := r.store.db.Query(buf.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var idHosts []int

	for rows.Next() {
		var v int
		err := rows.Scan(&v)

		if err != nil {
			log.Println(err)
			continue
		}

		idHosts = append(idHosts, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return idHosts, nil
}

// GetIDHostsByEnv ...
func (r *HostRepository) GetIDHostsByEnv(env string, idHosts []int) ([]int, error) {
	var stmt string
	var ids []int

	if env != "all" {
		stmt = `SELECT id FROM inventory_hosts WHERE approved=true AND environment=$1 AND id IN (`
		buf := bytes.NewBufferString(stmt)

		for i, v := range idHosts {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(strconv.Itoa(v))
		}
		buf.WriteString(")")

		rows, err := r.store.db.Query(buf.String(), env)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var v int
			err := rows.Scan(&v)

			if err != nil {
				log.Println(err)
				continue
			}

			ids = append(ids, v)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	} else {
		ids = idHosts
	}

	return ids, nil
}

// GetHostsByIDHosts ...
func (r *HostRepository) GetHostsByIDHosts(idHosts []int) ([]*model.Host, error) {
	hosts := []*model.Host{}

	stmt := `SELECT id, hostname, ip, environment
	FROM inventory_hosts WHERE approved=true AND id IN (`

	buf := bytes.NewBufferString(stmt)

	for i, v := range idHosts {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(v))
	}
	buf.WriteString(")")

	rows, err := r.store.db.Query(buf.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		host := &model.Host{}
		err := rows.Scan(
			&host.ID,
			&host.Hostname,
			&host.IP,
			&host.Environment,
		)

		if err != nil {
			return nil, err
		}

		hosts = append(hosts, host)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

// GetHostByFilters ...
func (r *HostRepository) GetHostByFilters(env string, values url.Values, sortKeys []string) ([]*model.Host, error) {
	tag := "tag"
	hosts := []*model.Host{}
	var idSlice []int

	for _, k := range sortKeys {
		values := StrWithCommaToArrStr(values[k])

		if k == tag {
			idTags, err := r.GetIDTagsByValueTags(values)
			if err != nil {
				return nil, err
			}

			allTagHostIDs, err := r.GetAllIDHostsByIDTags(idTags)
			if err != nil {
				return nil, err
			}

			allTagHostIDsByEnv, err := r.GetIDHostsByEnv(env, allTagHostIDs)
			if err != nil {
				return nil, err
			}

			// Write it in a slice
			idSlice = allTagHostIDsByEnv
		} else {
			idVars, err := r.GetIDVarsByValueVars(k, values)
			if err != nil {
				return nil, err
			}

			allVarHostIDs, err := r.GetAllIDHostsByIDVars(idVars)
			if err != nil {
				return nil, err
			}

			allVarHostIDsByEnv, err := r.GetIDHostsByEnv(env, allVarHostIDs)
			if err != nil {
				return nil, err
			}

			// Checking for occurrences in the slice
			var newSlice []int
			if len(idSlice) != 0 {
				for _, id := range allVarHostIDsByEnv {
					if Find(idSlice, id) {
						newSlice = append(newSlice, id)
					}
				}

			} else {
				newSlice = allVarHostIDsByEnv
			}

			// Updating data in the slice
			idSlice = newSlice
		}
	}

	hosts, err := r.GetHostsByIDHosts(idSlice)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}
