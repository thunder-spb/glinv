package pgsql

import (
	"bytes"
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/lib/pq"

	"glinv/pkg/models"
)

// InventoryHostModel type which wraps a sql.DB connection pool.
type InventoryHostModel struct {
	DB *sql.DB
}

// Get will return a specific host based on its id.
func (m *InventoryHostModel) Get(hostID int) (*models.InventoryHost, []*models.InventoryGroup, []*models.InventoryHVar, []*models.InventoryHTag, error) {
	stmtHosts := `SELECT id, hostname, ip, environment, approved, status, description, properties, created, updated
	FROM inventory_hosts WHERE id = $1`

	inventoryHost := &models.InventoryHost{}

	row := m.DB.QueryRow(stmtHosts, hostID)
	err := row.Scan(
		&inventoryHost.ID,
		&inventoryHost.Hostname,
		&inventoryHost.IP,
		&inventoryHost.Environment,
		&inventoryHost.Approved,
		&inventoryHost.Status,
		&inventoryHost.Description,
		&inventoryHost.Properties,
		&inventoryHost.Created,
		&inventoryHost.Updated,
	)

	if err == sql.ErrNoRows {
		return nil, nil, nil, nil, models.ErrNoRecord
	} else if err != nil {
		return nil, nil, nil, nil, err
	}

	// Group
	// Get slice ID of groups
	stmtHostsGroups := `SELECT group_id
	FROM inventory_hosts_groups WHERE host_id = $1`

	rows, err := m.DB.Query(stmtHostsGroups, hostID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	var idHostGroups []int
	for rows.Next() {
		var v int
		err := rows.Scan(&v)
		if err != nil {
			log.Println(err)
			continue
		}
		idHostGroups = append(idHostGroups, v)
	}

	if idHostGroups == nil {
		idHostGroups = []int{0}
	}

	// Get values of groups
	stmtGroupValue := `SELECT id, value
	FROM inventory_groups WHERE id IN (`

	buf := bytes.NewBufferString(stmtGroupValue)
	for i, v := range idHostGroups {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(v))
	}
	buf.WriteString(")")

	rows, err = m.DB.Query(buf.String())
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	groupValues := []*models.InventoryGroup{}

	for rows.Next() {
		gn := &models.InventoryGroup{}
		err := rows.Scan(
			&gn.ID,
			&gn.Value,
		)

		if err != nil {
			return nil, nil, nil, nil, err
		}

		groupValues = append(groupValues, gn)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, nil, nil, err
	}

	//Var
	// Get slice ID of vars
	stmtHostVars := `SELECT var_id
		FROM inventory_host_vars WHERE host_id = $1`

	rows, err = m.DB.Query(stmtHostVars, hostID)
	if err != nil {
		return nil, nil, nil, nil, err
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
	stmtVarValue := `SELECT id, name, value
	FROM inventory_hvars WHERE id IN (`

	buf = bytes.NewBufferString(stmtVarValue)
	for i, v := range idHostVars {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(v))
	}
	buf.WriteString(")")

	rows, err = m.DB.Query(buf.String())
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	hvarValues := []*models.InventoryHVar{}

	for rows.Next() {
		vn := &models.InventoryHVar{}
		err := rows.Scan(
			&vn.ID,
			&vn.Name,
			&vn.Value,
		)

		if err != nil {
			return nil, nil, nil, nil, err
		}

		hvarValues = append(hvarValues, vn)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, nil, nil, err
	}

	//Tag
	// Get slice ID of tags
	stmtHostTags := `SELECT tag_id
		FROM inventory_host_tags WHERE host_id = $1`

	rows, err = m.DB.Query(stmtHostTags, hostID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	var idHostTags []int
	for rows.Next() {
		var t int
		err := rows.Scan(&t)
		if err != nil {
			log.Println(err)
			continue
		}
		idHostTags = append(idHostTags, t)
	}

	if idHostTags == nil {
		idHostTags = []int{0}
	}

	// Get values of tags
	stmtTagValue := `SELECT id, value
	FROM inventory_htags WHERE id IN (`

	buf = bytes.NewBufferString(stmtTagValue)
	for i, t := range idHostTags {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(t))
	}
	buf.WriteString(")")

	rows, err = m.DB.Query(buf.String())
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	htagValues := []*models.InventoryHTag{}

	for rows.Next() {
		tag := &models.InventoryHTag{}
		err := rows.Scan(
			&tag.ID,
			&tag.Value,
		)

		if err != nil {
			return nil, nil, nil, nil, err
		}

		htagValues = append(htagValues, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, nil, nil, err
	}

	return inventoryHost, groupValues, hvarValues, htagValues, nil
}

// GetAll will return the created hosts.
func (m *InventoryHostModel) GetAll() ([]*models.InventoryHost, error) {
	stmt := `SELECT id, hostname, ip, environment, approved, status, description, properties, created, updated
	FROM inventory_hosts WHERE approved=true ORDER BY hostname ASC`

	inventoryHosts := []*models.InventoryHost{}

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		inventoryHost := &models.InventoryHost{}
		err := rows.Scan(
			&inventoryHost.ID,
			&inventoryHost.Hostname,
			&inventoryHost.IP,
			&inventoryHost.Environment,
			&inventoryHost.Approved,
			&inventoryHost.Status,
			&inventoryHost.Description,
			&inventoryHost.Properties,
			&inventoryHost.Created,
			&inventoryHost.Updated,
		)

		if err != nil {
			return nil, err
		}

		inventoryHosts = append(inventoryHosts, inventoryHost)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inventoryHosts, nil
}

// GetStatuses ...
func (m *InventoryHostModel) GetStatuses() (map[int]int, error) {
	stmt := `SELECT host_id, server_id FROM compare_server_host`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statuses := make(map[int]int)

	for rows.Next() {
		var h, s int
		if err := rows.Scan(&h, &s); err != nil {
			return nil, err
		}

		statuses[h] = s
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return statuses, nil
}

// GetHostsByEnv will return the created hosts with filtering by environment.
func (m *InventoryHostModel) GetHostsByEnv(env string) ([]*models.InventoryHost, error) {
	stmtAll := `SELECT id, hostname, ip, environment, approved, delete, status, description, properties, created, updated
	FROM inventory_hosts ORDER BY hostname ASC`

	stmtEnv := `SELECT id, hostname, ip, environment, approved, delete, status, description, properties, created, updated
	FROM inventory_hosts WHERE environment=$1 ORDER BY hostname ASC`

	inventoryHosts := []*models.InventoryHost{}

	if env == "" || env == "all" {
		rows, err := m.DB.Query(stmtAll)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			inventoryHost := &models.InventoryHost{}
			err := rows.Scan(
				&inventoryHost.ID,
				&inventoryHost.Hostname,
				&inventoryHost.IP,
				&inventoryHost.Environment,
				&inventoryHost.Approved,
				&inventoryHost.Delete,
				&inventoryHost.Status,
				&inventoryHost.Description,
				&inventoryHost.Properties,
				&inventoryHost.Created,
				&inventoryHost.Updated,
			)

			if err != nil {
				return nil, err
			}

			inventoryHost.CountServices = m.GetCountServicesByHost(inventoryHost.Hostname)

			inventoryHosts = append(inventoryHosts, inventoryHost)
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
			inventoryHost := &models.InventoryHost{}
			err := rows.Scan(
				&inventoryHost.ID,
				&inventoryHost.Hostname,
				&inventoryHost.IP,
				&inventoryHost.Environment,
				&inventoryHost.Approved,
				&inventoryHost.Delete,
				&inventoryHost.Status,
				&inventoryHost.Description,
				&inventoryHost.Properties,
				&inventoryHost.Created,
				&inventoryHost.Updated,
			)

			if err != nil {
				return nil, err
			}

			inventoryHost.CountServices = m.GetCountServicesByHost(inventoryHost.Hostname)

			inventoryHosts = append(inventoryHosts, inventoryHost)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return inventoryHosts, nil
}

// GetHostsByIDGroup ...
func (m *InventoryHostModel) GetHostsByIDGroup(idGroup int) ([]*models.InventoryHost, error) {
	// Get ID Hosts
	idHosts, _ := m.GetIDHostsByIDGroup(idGroup)

	// Get Hosts by Group
	inventoryHosts := []*models.InventoryHost{}

	if len(idHosts) != 0 {
		stmt := `SELECT id, hostname, ip, environment, approved, created, updated FROM inventory_hosts WHERE id IN (`

		buf := bytes.NewBufferString(stmt)

		for i, v := range idHosts {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(strconv.Itoa(v))
		}
		buf.WriteString(")")

		rows, err := m.DB.Query(buf.String())
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			host := &models.InventoryHost{}
			err := rows.Scan(
				&host.ID,
				&host.Hostname,
				&host.IP,
				&host.Environment,
				&host.Approved,
				&host.Created,
				&host.Updated,
			)

			if err != nil {
				return nil, err
			}

			inventoryHosts = append(inventoryHosts, host)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return inventoryHosts, nil
}

// GetIDHostsByIDGroup ...
func (m *InventoryHostModel) GetIDHostsByIDGroup(idGroup int) ([]int, error) {
	stmt := `SELECT host_id FROM inventory_hosts_groups WHERE group_id  = $1`

	var idHosts []int

	rows, err := m.DB.Query(stmt, idGroup)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var value int
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err)
			continue
		}
		idHosts = append(idHosts, value)
	}

	return idHosts, nil
}

// GetHostsByIDVar ...
func (m *InventoryHostModel) GetHostsByIDVar(idVar int) ([]*models.InventoryHost, error) {
	// Get ID Hosts
	idHosts, _ := m.GetIDHostsByIDVar(idVar)

	// Get Hosts by Var
	inventoryHosts := []*models.InventoryHost{}

	if len(idHosts) != 0 {
		stmt := `SELECT id, hostname, ip, environment, approved, created, updated FROM inventory_hosts WHERE id IN (`

		buf := bytes.NewBufferString(stmt)

		for i, v := range idHosts {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(strconv.Itoa(v))
		}
		buf.WriteString(")")

		rows, err := m.DB.Query(buf.String())
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			host := &models.InventoryHost{}
			err := rows.Scan(
				&host.ID,
				&host.Hostname,
				&host.IP,
				&host.Environment,
				&host.Approved,
				&host.Created,
				&host.Updated,
			)

			if err != nil {
				return nil, err
			}

			inventoryHosts = append(inventoryHosts, host)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return inventoryHosts, nil
}

// GetIDHostsByIDVar ...
func (m *InventoryHostModel) GetIDHostsByIDVar(idVar int) ([]int, error) {
	stmt := `SELECT host_id FROM inventory_host_vars WHERE var_id  = $1`

	var idHosts []int

	rows, err := m.DB.Query(stmt, idVar)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var value int
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err)
			continue
		}
		idHosts = append(idHosts, value)
	}

	return idHosts, nil
}

// GetHostsByIDTag ...
func (m *InventoryHostModel) GetHostsByIDTag(idTag int) ([]*models.InventoryHost, error) {
	// Get ID Hosts
	idHosts, _ := m.GetIDHostsByIDTag(idTag)

	// Get Hosts by Var
	inventoryHosts := []*models.InventoryHost{}

	if len(idHosts) != 0 {
		stmt := `SELECT id, hostname, ip, environment, approved, created, updated FROM inventory_hosts WHERE id IN (`

		buf := bytes.NewBufferString(stmt)

		for i, v := range idHosts {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(strconv.Itoa(v))
		}
		buf.WriteString(")")

		rows, err := m.DB.Query(buf.String())
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			host := &models.InventoryHost{}
			err := rows.Scan(
				&host.ID,
				&host.Hostname,
				&host.IP,
				&host.Environment,
				&host.Approved,
				&host.Created,
				&host.Updated,
			)

			if err != nil {
				return nil, err
			}

			inventoryHosts = append(inventoryHosts, host)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return inventoryHosts, nil
}

// GetIDHostsByIDTag ...
func (m *InventoryHostModel) GetIDHostsByIDTag(idVar int) ([]int, error) {
	stmt := `SELECT host_id FROM inventory_host_tags WHERE tag_id  = $1`

	var idHosts []int

	rows, err := m.DB.Query(stmt, idVar)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var value int
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err)
			continue
		}
		idHosts = append(idHosts, value)
	}

	return idHosts, nil
}

// GetIDHostGroups ...
func (m *InventoryHostModel) GetIDHostGroups(hostID int) ([]int, error) {
	stmt := `SELECT group_id FROM inventory_hosts_groups WHERE host_id = $1`

	var idGroups []int

	rows, err := m.DB.Query(stmt, hostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var value int
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err)
			continue
		}
		idGroups = append(idGroups, value)
	}

	return idGroups, nil
}

// GetIDHostVars ...
func (m *InventoryHostModel) GetIDHostVars(hostID int) ([]int, error) {
	stmt := `SELECT var_id FROM inventory_host_vars WHERE host_id = $1`

	var idHVars []int

	rows, err := m.DB.Query(stmt, hostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var value int
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err)
			continue
		}
		idHVars = append(idHVars, value)
	}

	return idHVars, nil
}

// GetIDHostTags ...
func (m *InventoryHostModel) GetIDHostTags(hostID int) ([]int, error) {
	stmt := `SELECT tag_id FROM inventory_host_tags WHERE host_id = $1`

	var idHTags []int

	rows, err := m.DB.Query(stmt, hostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var value int
		err := rows.Scan(&value)
		if err != nil {
			log.Println(err)
			continue
		}
		idHTags = append(idHTags, value)
	}

	return idHTags, nil
}

// Update will Update a current host into the database.
func (m *InventoryHostModel) Update(id, hostname, ip, description string, newCheckGroups, newCheckVars, newCheckTags map[string]string) (int, error) {
	//
	// Update Host
	//
	stmtHost := `UPDATE inventory_hosts SET hostname = $2, ip = $3, description = $4, approved = $5
	WHERE id = $1`

	idHost, _ := strconv.Atoi(id)

	_, err := m.DB.Exec(stmtHost, idHost, hostname, ip, description, false)
	if err != nil {
		return 0, err
	}

	//
	// Update Groups
	//
	stmtDeleteGroups := `DELETE FROM inventory_hosts_groups WHERE host_id = $1;`
	_, err = m.DB.Exec(stmtDeleteGroups, idHost)
	if err != nil {
		return 0, err
	}

	// inserting all new entries for this host
	if len(newCheckGroups) != 0 {
		stmtInsertGroups := `INSERT INTO inventory_hosts_groups (host_id, group_id) VALUES `

		IDs := []interface{}{}

		for ID := range newCheckGroups {
			stmtInsertGroups += "(?, ?),"
			IDs = append(IDs, idHost, ID)
		}

		stmtInsertGroups = strings.TrimSuffix(stmtInsertGroups, ",")
		stmtInsertGroups = ReplaceSQL(stmtInsertGroups, "?") //Replacing ? with $n for PostgreSQL

		prepGroups, _ := m.DB.Prepare(stmtInsertGroups)
		if err != nil {
			return 0, err
		}

		if _, err = prepGroups.Exec(IDs...); err != nil {
			return 0, err
		}
	}

	//
	// Update HVars
	//
	currentCheckVars, _ := m.GetIDHostVars(idHost)

	// deleting all entries for this host
	if len(currentCheckVars) != 0 {
		stmtDeleteHVars := `DELETE FROM inventory_host_vars WHERE host_id = $1;`
		_, err = m.DB.Exec(stmtDeleteHVars, idHost)
		if err != nil {
			return 0, err
		}
	}

	// inserting all new entries for this host
	if len(newCheckVars) != 0 {
		stmtInsertHVars := `INSERT INTO inventory_host_vars (host_id, var_id) VALUES `

		IDs := []interface{}{}

		for ID := range newCheckVars {
			stmtInsertHVars += "(?, ?),"
			IDs = append(IDs, idHost, ID)
		}

		stmtInsertHVars = strings.TrimSuffix(stmtInsertHVars, ",")
		stmtInsertHVars = ReplaceSQL(stmtInsertHVars, "?") //Replacing ? with $n for PostgreSQL

		prepHVars, _ := m.DB.Prepare(stmtInsertHVars)
		if err != nil {
			return 0, err
		}

		if _, err = prepHVars.Exec(IDs...); err != nil {
			return 0, err
		}
	}

	//
	// Update HTags
	//
	currentCheckTags, _ := m.GetIDHostTags(idHost)

	// deleting all entries for this host
	if len(currentCheckTags) != 0 {
		stmtDeleteHTags := `DELETE FROM inventory_host_tags WHERE host_id = $1;`
		_, err = m.DB.Exec(stmtDeleteHTags, idHost)
		if err != nil {
			return 0, err
		}
	}

	// inserting all new entries for this host
	if len(newCheckTags) != 0 {
		stmtInsertHTags := `INSERT INTO inventory_host_tags (host_id, tag_id) VALUES `

		IDs := []interface{}{}

		for ID := range newCheckTags {
			stmtInsertHTags += "(?, ?),"
			IDs = append(IDs, idHost, ID)
		}

		stmtInsertHTags = strings.TrimSuffix(stmtInsertHTags, ",")
		stmtInsertHTags = ReplaceSQL(stmtInsertHTags, "?") //Replacing ? with $n for PostgreSQL

		prepHTags, _ := m.DB.Prepare(stmtInsertHTags)
		if err != nil {
			return 0, err
		}

		if _, err = prepHTags.Exec(IDs...); err != nil {
			return 0, err
		}
	}

	return idHost, nil
}

// Insert will insert a new host into the database.
func (m *InventoryHostModel) Insert(hostname, ip, environment string, checkGroups, checkVars, checkTags map[string]string) (int, error) {
	stmtHost := `INSERT INTO inventory_hosts (hostname, ip, environment)
	VALUES($1, $2, $3) RETURNING id`

	var lastInsertHostID int
	err := m.DB.QueryRow(stmtHost, hostname, ip, environment).Scan(&lastInsertHostID)
	if err != nil {
		// TODO
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && strings.Contains(pqErr.Code.Name(), "unique_violation") {
				return lastInsertHostID, models.ErrDuplicateEmail
			}
		}
	}

	// 1 host = N group (checkGroups) M:M
	stmtGroup := `INSERT INTO inventory_hosts_groups (host_id, group_id) VALUES `

	if len(checkGroups) != 0 {
		groupIDs := []interface{}{}
		for groupID := range checkGroups {
			stmtGroup += "(?, ?),"
			groupIDs = append(groupIDs, lastInsertHostID, groupID)
		}

		stmtGroup = strings.TrimSuffix(stmtGroup, ",")
		stmtGroup = ReplaceSQL(stmtGroup, "?") //Replacing ? with $n for PostgreSQL

		prepGroup, _ := m.DB.Prepare(stmtGroup)
		if err != nil {
			panic(err)
		}

		if _, err = prepGroup.Exec(groupIDs...); err != nil {
			panic(err)
		}
	}

	// 1 host = N var (checkVars) M:M
	stmtVar := `INSERT INTO inventory_host_vars (host_id, var_id)
	VALUES `

	varIDs := []interface{}{}
	if len(checkVars) != 0 {
		for varID := range checkVars {
			stmtVar += "(?, ?),"
			varIDs = append(varIDs, lastInsertHostID, varID)
		}

		stmtVar = strings.TrimSuffix(stmtVar, ",")
		stmtVar = ReplaceSQL(stmtVar, "?") //Replacing ? with $n for PostgreSQL

		prepVar, _ := m.DB.Prepare(stmtVar)
		if err != nil {
			panic(err)
		}

		if _, err = prepVar.Exec(varIDs...); err != nil {
			panic(err)
		}
	}

	// 1 host = N tag (checkTags) M:M
	stmtTag := `INSERT INTO inventory_host_tags (host_id, tag_id)
		VALUES `

	tagIDs := []interface{}{}
	if len(checkTags) != 0 {
		for tagID := range checkTags {
			stmtTag += "(?, ?),"
			tagIDs = append(tagIDs, lastInsertHostID, tagID)
		}

		stmtTag = strings.TrimSuffix(stmtTag, ",")
		stmtTag = ReplaceSQL(stmtTag, "?") //Replacing ? with $n for PostgreSQL

		prepTag, _ := m.DB.Prepare(stmtTag)
		if err != nil {
			panic(err)
		}

		if _, err = prepTag.Exec(tagIDs...); err != nil {
			panic(err)
		}
	}

	return lastInsertHostID, err
}

// GetAllForApprove will ...
func (m *InventoryHostModel) GetAllForApprove() ([]*models.InventoryHost, error) {
	stmt := `SELECT id, hostname, ip, environment, approved, description, properties, created, updated
	FROM inventory_hosts WHERE approved=false`

	inventoryHosts := []*models.InventoryHost{}

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		inventoryHost := &models.InventoryHost{}
		err := rows.Scan(
			&inventoryHost.ID,
			&inventoryHost.Hostname,
			&inventoryHost.IP,
			&inventoryHost.Environment,
			&inventoryHost.Approved,
			&inventoryHost.Description,
			&inventoryHost.Properties,
			&inventoryHost.Created,
			&inventoryHost.Updated,
		)

		if err != nil {
			return nil, err
		}

		inventoryHosts = append(inventoryHosts, inventoryHost)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inventoryHosts, nil
}

// GetAllForDelete will ...
func (m *InventoryHostModel) GetAllForDelete() ([]*models.InventoryHost, error) {
	stmt := `SELECT id, hostname, ip, environment, delete, description, properties, created, updated
	FROM inventory_hosts WHERE delete=true`

	inventoryHosts := []*models.InventoryHost{}

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		inventoryHost := &models.InventoryHost{}
		err := rows.Scan(
			&inventoryHost.ID,
			&inventoryHost.Hostname,
			&inventoryHost.IP,
			&inventoryHost.Environment,
			&inventoryHost.Delete,
			&inventoryHost.Description,
			&inventoryHost.Properties,
			&inventoryHost.Created,
			&inventoryHost.Updated,
		)

		if err != nil {
			return nil, err
		}

		inventoryHosts = append(inventoryHosts, inventoryHost)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inventoryHosts, nil
}

// ApproveHost ...
func (m *InventoryHostModel) ApproveHost(id string) error {
	stmt := `UPDATE inventory_hosts SET approved = $1 WHERE id = $2;`

	_, err := m.DB.Exec(stmt, "true", id)
	if err != nil {
		return err
	}
	return nil
}

// Delete ...
func (m *InventoryHostModel) Delete(id string) error {
	// Group
	stmtHostsGroups := `DELETE FROM inventory_hosts_groups WHERE host_id = $1;`
	_, err := m.DB.Exec(stmtHostsGroups, id)
	if err != nil {
		return err
	}

	// HVar
	stmtHostVars := `DELETE FROM inventory_host_vars WHERE host_id = $1;`
	_, err = m.DB.Exec(stmtHostVars, id)
	if err != nil {
		return err
	}

	// Tag
	stmtHostTags := `DELETE FROM inventory_host_tags WHERE host_id = $1;`
	_, err = m.DB.Exec(stmtHostTags, id)
	if err != nil {
		return err
	}

	// Host
	stmtHosts := `DELETE FROM inventory_hosts WHERE id = $1;`
	_, err = m.DB.Exec(stmtHosts, id)
	if err != nil {
		return err
	}

	return nil
}

// MarkHostForDeletion ...
func (m *InventoryHostModel) MarkHostForDeletion(id int) error {
	stmt := `UPDATE inventory_hosts SET delete = $1 WHERE id = $2;`

	_, err := m.DB.Exec(stmt, "true", id)
	if err != nil {
		return err
	}
	return nil
}

// CancelDeletion ...
func (m *InventoryHostModel) CancelDeletion(id int) error {
	stmt := `UPDATE inventory_hosts SET delete = $1 WHERE id = $2;`

	_, err := m.DB.Exec(stmt, "false", id)
	if err != nil {
		return err
	}
	return nil
}

//
// Counts
//

// GetCountHosts ...
func (m *InventoryHostModel) GetCountHosts() int {
	stmt := `SELECT COUNT(*) FROM inventory_hosts;`
	var count int
	m.DB.QueryRow(stmt).Scan(&count)
	return count
}

// GetCountHostsVars ...
func (m *InventoryHostModel) GetCountHostsVars() int {
	stmt := `SELECT COUNT(*) FROM inventory_hvars;`
	var count int
	m.DB.QueryRow(stmt).Scan(&count)
	return count
}

// GetCountHostsTags ...
func (m *InventoryHostModel) GetCountHostsTags() int {
	stmt := `SELECT COUNT(*) FROM inventory_htags;`
	var count int
	m.DB.QueryRow(stmt).Scan(&count)
	return count
}

// GetCountServicesByHost ...
func (m *InventoryHostModel) GetCountServicesByHost(hostname string) int {
	stmt := `SELECT COUNT(*) FROM inventory_host_services WHERE hostname = $1;`
	var count int
	m.DB.QueryRow(stmt, hostname).Scan(&count)
	return count
}

// GetCountUnapprovedHosts ...
func (m *InventoryHostModel) GetCountUnapprovedHosts() int {
	stmtApproved := `SELECT COUNT(*) FROM inventory_hosts WHERE approved = $1;`
	var countApproved int
	m.DB.QueryRow(stmtApproved, "false").Scan(&countApproved)

	stmtDelete := `SELECT COUNT(*) FROM inventory_hosts WHERE delete = $1;`
	var countDelete int
	m.DB.QueryRow(stmtDelete, "true").Scan(&countDelete)

	return countApproved + countDelete
}
