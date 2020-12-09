package pgsql

import (
	"bytes"
	"database/sql"
	"log"
	"strconv"

	"glinv/pkg/models"
)

// InventoryGVarModel type which wraps a sql.DB connection pool.
type InventoryGVarModel struct {
	DB *sql.DB
}

// Get will return a specific ...
func (m *InventoryGVarModel) Get(id int) (*models.InventoryGVar, error) {
	stmt := `SELECT id, name, value FROM inventory_gvars WHERE id = $1`

	inventoryGVar := &models.InventoryGVar{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&inventoryGVar.ID,
		&inventoryGVar.Name,
		&inventoryGVar.Value,
	)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return inventoryGVar, nil
}

// GetAll will return the all inventory vars for list.
func (m *InventoryGVarModel) GetAll() ([]*models.InventoryGVar, error) {
	stmt := `SELECT id, name, value
	FROM inventory_gvars ORDER BY value ASC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gvars := []*models.InventoryGVar{}

	for rows.Next() {
		gvar := &models.InventoryGVar{}
		err := rows.Scan(
			&gvar.ID,
			&gvar.Name,
			&gvar.Value,
		)

		if err != nil {
			return nil, err
		}

		gvars = append(gvars, gvar)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return gvars, nil
}

// GetGVarsByGroupID ...
func (m *InventoryGVarModel) GetGVarsByGroupID(id int) ([]*models.InventoryGVar, error) {
	stmtGroupVars := `SELECT var_id FROM inventory_group_vars WHERE group_id = $1`

	rows, err := m.DB.Query(stmtGroupVars, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var idGVars []int
	for rows.Next() {
		var v int
		err := rows.Scan(&v)
		if err != nil {
			log.Println(err)
			continue
		}
		idGVars = append(idGVars, v)
	}

	if idGVars == nil {
		idGVars = []int{0}
	}

	stmtGvars := `SELECT id, name, value FROM inventory_gvars WHERE id IN (`

	buf := bytes.NewBufferString(stmtGvars)
	for i, v := range idGVars {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(v))
	}
	buf.WriteString(")")

	rows, err = m.DB.Query(buf.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gvars := []*models.InventoryGVar{}

	for rows.Next() {
		g := &models.InventoryGVar{}
		err := rows.Scan(
			&g.ID,
			&g.Name,
			&g.Value,
		)

		if err != nil {
			return nil, err
		}

		gvars = append(gvars, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return gvars, nil
}

// Insert will insert a new var of group into the database.
func (m *InventoryGVarModel) Insert(name, value string) (int, error) {
	stmt := `INSERT INTO inventory_gvars (name, value)
	VALUES($1, $2) RETURNING id`

	var lastInsertID int
	if err := m.DB.QueryRow(stmt, name, value).Scan(&lastInsertID); err != nil {
		panic(err)
	}

	return lastInsertID, nil
}

// Update ...
func (m *InventoryGVarModel) Update(id, value string) error {
	stmt := `UPDATE inventory_gvars SET value = $2
	WHERE id = $1`

	_, err := m.DB.Exec(stmt, id, value)
	if err != nil {
		return err
	}
	return nil
}

// Delete ...
func (m *InventoryGVarModel) Delete(id int) error {
	stmtGroupVars := `DELETE FROM inventory_group_vars WHERE var_id = $1;`

	_, err := m.DB.Exec(stmtGroupVars, id)
	if err != nil {
		return err
	}

	stmtGVars := `DELETE FROM inventory_gvars WHERE id = $1;`

	_, err = m.DB.Exec(stmtGVars, id)
	if err != nil {
		return err
	}

	return nil
}

// TreeGVars will return a specific vars based on its id.
func (m *InventoryGVarModel) TreeGVars() (map[int][]*models.InventoryGVar, error) {
	stmtGroupsID := `SELECT group_id
	FROM inventory_group_vars`

	// Getting the ID of all groups
	rows, err := m.DB.Query(stmtGroupsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupsID []int
	for rows.Next() {
		var v int
		err := rows.Scan(&v)
		if err != nil {
			log.Println(err)
			continue
		}
		groupsID = append(groupsID, v)
	}

	// Excluding duplicates from the selection
	var uniqueGroupsID = Unique(groupsID)

	mapGroupGVars := map[int][]*models.InventoryGVar{}

	stmtGroupVars := `SELECT var_id
	FROM inventory_group_vars WHERE group_id = $1`

	stmtGVarValue := `SELECT id, name, value
	FROM inventory_gvars WHERE id IN (`

	for _, id := range uniqueGroupsID {
		// Get slice id vars of group
		rows, err = m.DB.Query(stmtGroupVars, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var idGroupVars []int
		for rows.Next() {
			var v int
			err := rows.Scan(&v)
			if err != nil {
				log.Println(err)
				continue
			}
			idGroupVars = append(idGroupVars, v)
		}

		// Get values of vars
		buf := bytes.NewBufferString(stmtGVarValue)
		for i, v := range idGroupVars {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(strconv.Itoa(v))
		}
		buf.WriteString(")")

		rows, err = m.DB.Query(buf.String())
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		gvarValues := []*models.InventoryGVar{}

		for rows.Next() {
			gv := &models.InventoryGVar{}
			err := rows.Scan(
				&gv.ID,
				&gv.Name,
				&gv.Value,
			)

			if err != nil {
				return nil, err
			}

			gvarValues = append(gvarValues, gv)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		mapGroupGVars[id] = gvarValues
	}

	return mapGroupGVars, err
}
