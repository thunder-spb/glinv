package pgsql

import (
	"bytes"
	"database/sql"
	"log"
	"strconv"
	"strings"

	"glinv/pkg/models"

	"github.com/lib/pq"
)

// InventoryGroupModel type which wraps a sql.DB connection pool.
type InventoryGroupModel struct {
	DB *sql.DB
}

// Delete ...
func (m *InventoryGroupModel) Delete(id int) error {
	stmtGroupVars := `DELETE FROM inventory_group_vars WHERE group_id = $1;`
	_, err := m.DB.Exec(stmtGroupVars, id)
	if err != nil {
		return err
	}

	stmtGroup := `DELETE FROM inventory_groups WHERE id = $1;`
	_, err = m.DB.Exec(stmtGroup, id)
	if err != nil {
		return err
	}

	return nil
}

// Update ...
func (m *InventoryGroupModel) Update(id, value string, gvars map[string]string) error {
	stmt := `UPDATE inventory_groups SET value = $2 WHERE id = $1`
	_, err := m.DB.Exec(stmt, id, value)
	if err != nil {
		return err
	}

	stmtDeleteGVars := `DELETE FROM inventory_group_vars WHERE group_id = $1;`
	_, err = m.DB.Exec(stmtDeleteGVars, id)
	if err != nil {
		return err
	}

	stmtGVar := `INSERT INTO inventory_group_vars (group_id, var_id) VALUES `
	if len(gvars) != 0 {
		gvarIDs := []interface{}{}
		for varID := range gvars {
			stmtGVar += "(?, ?),"
			gvarIDs = append(gvarIDs, id, varID)
		}

		stmtGVar = strings.TrimSuffix(stmtGVar, ",")
		stmtGVar = ReplaceSQL(stmtGVar, "?") //Replacing ? with $n for PostgreSQL

		prepGVar, err := m.DB.Prepare(stmtGVar)
		if err != nil {
			return err
		}

		if _, err := prepGVar.Exec(gvarIDs...); err != nil {
			return err
		}
	}

	return nil
}

// Insert will insert a new host into the database.
func (m *InventoryGroupModel) Insert(environment, parent, newgroup, description string, gvars map[string]string) (int, error) {
	// stmtSearchParentID := `SELECT id
	// FROM inventory_groups WHERE value = $1`

	// var parentID int
	// row := m.DB.QueryRow(stmtSearchParentID, parent)
	// err := row.Scan(&parentID)
	// if err != nil {
	// 	return 0, err
	// }
	// fmt.Println(environment, parentID, parent, newgroup, description)

	stmtGroup := `INSERT INTO inventory_groups (environment, parent_id, value, description)
	VALUES($1, $2, $3, $4) RETURNING id`

	var lastInsertGroupID int
	err := m.DB.QueryRow(stmtGroup, environment, parent, newgroup, description).Scan(&lastInsertGroupID)
	if err != nil {
		// TODO
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && strings.Contains(pqErr.Code.Name(), "unique_violation") {
				return lastInsertGroupID, models.ErrDuplicateEmail
			}
		}
	}

	// 1 group = N var (gvars)
	stmtGVar := `INSERT INTO inventory_group_vars (group_id, var_id)
	VALUES `

	if len(gvars) != 0 {
		gvarIDs := []interface{}{}
		for varID := range gvars {
			stmtGVar += "(?, ?),"
			gvarIDs = append(gvarIDs, lastInsertGroupID, varID)
		}

		stmtGVar = strings.TrimSuffix(stmtGVar, ",")
		stmtGVar = ReplaceSQL(stmtGVar, "?") //Replacing ? with $n for PostgreSQL

		prepGVar, err := m.DB.Prepare(stmtGVar)
		if err != nil {
			panic(err)
		}

		if _, err := prepGVar.Exec(gvarIDs...); err != nil {
			panic(err)
		}
	}

	return lastInsertGroupID, err
}

// Get will return a specific host based on its id.
func (m *InventoryGroupModel) Get(id int) (*models.InventoryGroup, error) {
	stmt := `SELECT id, environment, parent_id, value
	FROM inventory_groups WHERE id = $1`

	inventoryGroup := &models.InventoryGroup{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&inventoryGroup.ID,
		&inventoryGroup.Environment,
		&inventoryGroup.ParentID,
		&inventoryGroup.Value,
	)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return inventoryGroup, nil
}

// GetAll will return the all inventory groups for tree.
func (m *InventoryGroupModel) GetAll() ([]*models.InventoryGroup, error) {
	stmt := `SELECT id, environment, parent_id, value
	FROM inventory_groups`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventoryGroups := []*models.InventoryGroup{}

	for rows.Next() {
		inventoryGroup := &models.InventoryGroup{}
		err := rows.Scan(
			&inventoryGroup.ID,
			&inventoryGroup.Environment,
			&inventoryGroup.ParentID,
			&inventoryGroup.Value,
		)

		if err != nil {
			return nil, err
		}

		inventoryGroups = append(inventoryGroups, inventoryGroup)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inventoryGroups, nil
}

// GetGroupsByEnv will return the all inventory groups for tree.
func (m *InventoryGroupModel) GetGroupsByEnv(environment string) ([]*models.InventoryGroup, error) {
	stmt := `SELECT id, environment, parent_id, value
	FROM inventory_groups WHERE environment IN ($1, $2) ORDER BY created`

	rows, err := m.DB.Query(stmt, environment, "all")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventoryGroups := []*models.InventoryGroup{}

	for rows.Next() {
		inventoryGroup := &models.InventoryGroup{}
		err := rows.Scan(
			&inventoryGroup.ID,
			&inventoryGroup.Environment,
			&inventoryGroup.ParentID,
			&inventoryGroup.Value,
		)

		if err != nil {
			return nil, err
		}

		inventoryGroups = append(inventoryGroups, inventoryGroup)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inventoryGroups, nil
}

// GetIDGroupsByIDVar ...
func (m *InventoryGroupModel) GetIDGroupsByIDVar(id int) ([]int, error) {
	stmt := `SELECT group_id FROM inventory_group_vars WHERE var_id = $1`

	var idGroups []int

	rows, err := m.DB.Query(stmt, id)
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

// GetGroupsByID ...
func (m *InventoryGroupModel) GetGroupsByID(id int) ([]*models.InventoryGroup, error) {
	idGroups, err := m.GetIDGroupsByIDVar(id)
	if err != nil {
		return nil, err
	}

	inventoryGroups := []*models.InventoryGroup{}

	if len(idGroups) != 0 {
		stmt := `SELECT id, value, created, updated FROM inventory_groups WHERE id IN (`

		buf := bytes.NewBufferString(stmt)

		for i, v := range idGroups {
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
			group := &models.InventoryGroup{}
			err := rows.Scan(
				&group.ID,
				&group.Value,
				&group.Created,
				&group.Updated,
			)

			if err != nil {
				return nil, err
			}

			inventoryGroups = append(inventoryGroups, group)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

	}

	return inventoryGroups, nil
}

// CheckSubGroups ...
func (m *InventoryGroupModel) CheckSubGroups(id int) bool {
	stmt := `SELECT parent_id FROM inventory_groups WHERE parent_id = $1`

	err := m.DB.QueryRow(stmt, id).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Print(err)
		}

		return false
	}

	return true
}

// GetCountGroups ...
func (m *InventoryGroupModel) GetCountGroups() int {
	stmt := `SELECT COUNT(*) FROM inventory_groups;`
	var count int
	m.DB.QueryRow(stmt).Scan(&count)
	return count
}

// GetCountGroupsVars ...
func (m *InventoryGroupModel) GetCountGroupsVars() int {
	stmt := `SELECT COUNT(*) FROM inventory_gvars;`
	var count int
	m.DB.QueryRow(stmt).Scan(&count)
	return count
}
