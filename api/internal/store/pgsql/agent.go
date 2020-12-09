package pgsql

import (
	"fmt"
	"reflect"
	"restapi/internal/model"
	"sort"
	"strings"
	"time"
)

// AgentRepository ...
type AgentRepository struct {
	store *Store
}

// Insert ...
func (r *AgentRepository) Insert(data model.Data) error {
	// Insert data of server
	stmtHard := `INSERT INTO server_hard_agent (hostname, ip, os_name, os_architecture, kernel_release, cpu_model, cpu_num, ram, hdd, uptime, resolv, mtu)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`

	var lastServerID int
	if err := r.store.db.QueryRow(stmtHard,
		data["hard"]["hostname"],
		data["hard"]["ip"],
		data["hard"]["osName"],
		data["hard"]["osArch"],
		data["hard"]["kernelRelease"],
		data["hard"]["modelCPU"],

		data["hard"]["numCPU"],
		data["hard"]["ram"],
		data["hard"]["hdd"],
		data["hard"]["uptime"],
		data["hard"]["resolv.conf"],
		0,
	).Scan(&lastServerID); err != nil {
		return err
	}

	// Insert sysctl of server
	stmtSysCtl := `INSERT INTO server_sysctl_agent (server_id, name, value) VALUES `

	sysctls := []interface{}{}
	for name, value := range data["sysctl"] {
		stmtSysCtl += "(?, ?, ?),"
		sysctls = append(sysctls, lastServerID, name, value)
	}

	stmtSysCtl = strings.TrimSuffix(stmtSysCtl, ",")
	stmtSysCtl = ReplaceSQL(stmtSysCtl, "?") //Replacing ? with $n for PostgreSQL

	prepSysCtl, err := r.store.db.Prepare(stmtSysCtl)
	if err != nil {
		return err
	}

	if _, err := prepSysCtl.Exec(sysctls...); err != nil {
		return err
	}

	// Insert packages of server
	stmtSoft := `INSERT INTO server_soft_agent (server_id, package, version) VALUES `

	soft := []interface{}{}
	for pkg, ver := range data["pkg"] {
		stmtSoft += "(?, ?, ?),"
		soft = append(soft, lastServerID, pkg, ver)
	}

	stmtSoft = strings.TrimSuffix(stmtSoft, ",")
	stmtSoft = ReplaceSQL(stmtSoft, "?") //Replacing ? with $n for PostgreSQL

	prepSoft, err := r.store.db.Prepare(stmtSoft)
	if err != nil {
		return err
	}

	if _, err := prepSoft.Exec(soft...); err != nil {
		return err
	}

	return nil
}

// Update ...
func (r *AgentRepository) Update(id int, data model.Data) error {
	//
	// Update host
	//

	stmtHard := `UPDATE server_hard_agent
	SET hostname=$2, ip=$3, os_name=$4, os_architecture=$5, kernel_release=$6, cpu_model=$7, cpu_num=$8, ram=$9, hdd=$10, uptime=$11, resolv=$12, mtu=$13
	WHERE id = $1;`

	_, err := r.store.db.Exec(stmtHard, id,
		data["hard"]["hostname"],
		data["hard"]["ip"],
		data["hard"]["osName"],
		data["hard"]["osArch"],
		data["hard"]["kernelRelease"],
		data["hard"]["modelCPU"],

		data["hard"]["numCPU"],
		data["hard"]["ram"],
		data["hard"]["hdd"],
		data["hard"]["uptime"],
		data["hard"]["resolv.conf"],
		0,
	)
	if err != nil {
		return err
	}

	//
	// Update sysctl
	//

	// get data from the DB and fill in the map
	stmtSysctl := `SELECT name, value, updated FROM server_sysctl_agent WHERE server_id = $1`

	rows, err := r.store.db.Query(stmtSysctl, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	sysctlDB := map[string]string{}

	for rows.Next() {
		var name, value, updated string
		err := rows.Scan(&name, &value, &updated)
		if err != nil {
			return err
		}

		sysctlDB[name] = value + "|" + updated
	}

	if err = rows.Err(); err != nil {
		return err
	}

	// get data from the agent and fill in the map
	sysctlAgent := map[string]string{}
	for name, value := range data["sysctl"] {
		sysctlAgent[name] = value
	}

	// sort sysctlDB
	var namesDB []string
	for name := range sysctlDB {
		namesDB = append(namesDB, name)
	}
	sort.Strings(namesDB)

	// sort sysctlAgent
	var namesAgent []string
	for name := range sysctlAgent {
		namesAgent = append(namesAgent, name)
	}
	sort.Strings(namesAgent)

	// compare the slices of the key maps
	compareNames := reflect.DeepEqual(namesDB, namesAgent)
	//fmt.Println("sysctl : ", compareNames, "in db:", len(namesDB), "from agent:", len(namesAgent))

	// if the keys match, check the key values, and make
	// an update of those keys whose values do not match
	if compareNames {
		for name, valueDB := range sysctlDB {
			vDB := strings.Split(valueDB, "|")
			check := false
			for _, valueAgent := range sysctlAgent {
				if valueAgent == vDB[0] {
					check = true
				}
			}

			if !check {
				// string to time
				timeUpdated, err := time.Parse(time.RFC3339, vDB[1])
				if err != nil {
					return err
				}

				// difference between the current time and the record update time
				dur := time.Since(timeUpdated)

				// value of the time variable 1h
				var i time.Duration
				diff := (60 - i) * time.Minute

				// if the difference is more than 1h then perform an update
				if dur > diff {
					//fmt.Println(check, name, vDB[0], vDB[1])
					stmtSysctlUpdate := `UPDATE server_sysctl_agent SET value = $1 WHERE name = $2 AND server_id = $3;`

					_, err := r.store.db.Exec(stmtSysctlUpdate, vDB[0], name, id)
					if err != nil {
						return err
					}
				}
			}
		}
	} else {
		// Delete sysctl of server
		stmtSysctlDelete := `DELETE FROM server_sysctl_agent WHERE server_id = $1;`

		if _, err := r.store.db.Exec(stmtSysctlDelete, id); err != nil {
			return err
		}

		// Insert sysctl of server
		stmtSysctlInsert := `INSERT INTO server_sysctl_agent (server_id, name, value) VALUES `

		sysctls := []interface{}{}
		for name, value := range data["sysctl"] {
			stmtSysctlInsert += "(?, ?, ?),"
			sysctls = append(sysctls, id, name, value)
		}

		stmtSysctlInsert = strings.TrimSuffix(stmtSysctlInsert, ",")
		stmtSysctlInsert = ReplaceSQL(stmtSysctlInsert, "?") //Replacing ? with $n for PostgreSQL

		prepSysCtl, err := r.store.db.Prepare(stmtSysctlInsert)
		if err != nil {
			return err
		}

		if _, err := prepSysCtl.Exec(sysctls...); err != nil {
			return err
		}

	}

	//
	// Update packages
	//

	// get data from the DB and fill in the map
	stmtSoft := `SELECT package, version FROM server_soft_agent WHERE server_id = $1`

	rows, err = r.store.db.Query(stmtSoft, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	softDB := map[string]string{}

	for rows.Next() {
		var pkg, ver string
		err := rows.Scan(&pkg, &ver)
		if err != nil {
			return err
		}

		softDB[pkg] = ver
	}

	if err = rows.Err(); err != nil {
		return err
	}

	// get data from the agent and fill in the map
	softAgent := map[string]string{}
	for pkg, ver := range data["pkg"] {
		softAgent[pkg] = ver
	}

	// sort softDB
	var pkgsDB []string
	for pkg := range softDB {
		pkgsDB = append(pkgsDB, pkg)
	}
	sort.Strings(pkgsDB)

	// sort softAgent
	var pkgsAgent []string
	for pkg := range softAgent {
		pkgsAgent = append(pkgsAgent, pkg)
	}
	sort.Strings(pkgsAgent)

	// compare the slices of the key maps
	comparePkgs := reflect.DeepEqual(pkgsDB, pkgsAgent)
	//fmt.Println("pkgs: ", comparePkgs, "in db:", len(pkgsDB), "from agent:", len(pkgsAgent))

	// if the keys match, check the key values, and make
	// an update of those keys whose values do not match
	if comparePkgs {
		for pkg, verAgent := range softAgent {
			check := false
			for _, verDB := range softDB {
				if verAgent == verDB {
					check = true
				}
			}

			if !check {
				fmt.Println(pkg, "=", verAgent)
				stmtSoftUpdate := `UPDATE server_soft_agent SET version = $1 WHERE package = $2 AND server_id = $3;`

				_, err := r.store.db.Exec(stmtSoftUpdate, verAgent, pkg, id)
				if err != nil {
					return err
				}
			}

		}
	} else {
		// Delete soft of server
		stmtSoftDelete := `DELETE FROM server_soft_agent WHERE server_id = $1;`

		if _, err := r.store.db.Exec(stmtSoftDelete, id); err != nil {
			return err
		}

		// Insert soft of server
		stmtSoftInsert := `INSERT INTO server_soft_agent (server_id, package, version) VALUES `

		soft := []interface{}{}
		for pkg, ver := range data["pkg"] {
			stmtSoftInsert += "(?, ?, ?),"
			soft = append(soft, id, pkg, ver)
		}

		stmtSoftInsert = strings.TrimSuffix(stmtSoftInsert, ",")
		stmtSoftInsert = ReplaceSQL(stmtSoftInsert, "?") //Replacing ? with $n for PostgreSQL

		prepSoft, err := r.store.db.Prepare(stmtSoftInsert)
		if err != nil {
			return err
		}

		if _, err := prepSoft.Exec(soft...); err != nil {
			return err
		}

	}

	return nil
}

// CheckServerExistsByHostname ...
func (r *AgentRepository) CheckServerExistsByHostname(hostname string) (bool, error) {
	stmtExists := `SELECT EXISTS(SELECT ip FROM server_hard_agent WHERE hostname = $1)`

	var exists bool

	row := r.store.db.QueryRow(stmtExists, hostname)

	if err := row.Scan(&exists); err != nil {
		return false, err
	} else if !exists {
		return exists, nil
	}

	return exists, nil
}

// GetIDByHostname ...
func (r *AgentRepository) GetIDByHostname(hostname string) (int, error) {
	stmt := `SELECT id FROM server_hard_agent WHERE hostname = $1`

	var id int

	if err := r.store.db.QueryRow(stmt, hostname).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// DeletePackages ...
func (r *AgentRepository) DeletePackages(id int) error {
	stmt := `DELETE FROM server_soft_agent WHERE hard_agent_id = $1;`

	if _, err := r.store.db.Exec(stmt, id); err != nil {
		return err
	}

	return nil
}

// GetByHostname Server
func (r *AgentRepository) GetByHostname(hostname string) (*model.Server, error) {
	stmt := `SELECT id, hostname, ip, os_name, os_architecture, kernel_release, cpu_model
	FROM server_hard_agent WHERE hostname = $1`

	s := &model.Server{}

	if err := r.store.db.QueryRow(stmt, hostname).Scan(
		&s.ID, &s.Hostname, &s.IP, &s.OSName, &s.OSArchitecture, &s.KernelRelease, &s.CPUModel,
	); err != nil {
		return nil, err
	}

	return s, nil
}

// GetAll ...
func (r *AgentRepository) GetAll() ([]*model.Server, error) {
	return nil, nil
}
