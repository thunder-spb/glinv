package pgsql

import (
	"database/sql"

	"glinv/pkg/models"
)

// InventoryHVarModel type which wraps a sql.DB connection pool.
type InventoryHVarModel struct {
	DB *sql.DB
}

// Get will return a specific ...
func (m *InventoryHVarModel) Get(id int) (*models.InventoryHVar, error) {
	stmt := `SELECT id, name, value FROM inventory_hvars WHERE id = $1`

	inventoryHVar := &models.InventoryHVar{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&inventoryHVar.ID,
		&inventoryHVar.Name,
		&inventoryHVar.Value,
	)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return inventoryHVar, nil
}

// GetAll will return the all inventory vars for list.
func (m *InventoryHVarModel) GetAll() ([]*models.InventoryHVar, error) {
	stmt := `SELECT id, name, value
	FROM inventory_hvars`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hvars := []*models.InventoryHVar{}

	for rows.Next() {
		hvar := &models.InventoryHVar{}
		err := rows.Scan(
			&hvar.ID,
			&hvar.Name,
			&hvar.Value,
		)

		if err != nil {
			return nil, err
		}

		hvars = append(hvars, hvar)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return hvars, nil
}

// Insert will insert a new var of host into the database.
func (m *InventoryHVarModel) Insert(name, value string) (int, error) {
	stmt := `INSERT INTO inventory_hvars (name, value)
	VALUES($1, $2) RETURNING id`

	var lastInsertID int
	if err := m.DB.QueryRow(stmt, name, value).Scan(&lastInsertID); err != nil {
		panic(err)
	}

	return lastInsertID, nil
}

// Update ...
func (m *InventoryHVarModel) Update(id, value string) error {
	stmt := `UPDATE inventory_hvars SET value = $2
	WHERE id = $1`

	_, err := m.DB.Exec(stmt, id, value)
	if err != nil {
		return err
	}
	return nil
}

// Delete ...
func (m *InventoryHVarModel) Delete(id int) error {
	stmtHostVars := `DELETE FROM inventory_host_vars WHERE var_id = $1;`

	_, err := m.DB.Exec(stmtHostVars, id)
	if err != nil {
		return err
	}

	stmtHVars := `DELETE FROM inventory_hvars WHERE id = $1;`

	_, err = m.DB.Exec(stmtHVars, id)
	if err != nil {
		return err
	}

	return nil
}
