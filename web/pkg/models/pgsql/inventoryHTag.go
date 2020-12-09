package pgsql

import (
	"database/sql"
	"strings"

	"glinv/pkg/models"

	"github.com/lib/pq"
)

// InventoryHTagModel type which wraps a sql.DB connection pool.
type InventoryHTagModel struct {
	DB *sql.DB
}

// Get will return a specific ...
func (m *InventoryHTagModel) Get(id int) (*models.InventoryHTag, error) {
	stmt := `SELECT id, value FROM inventory_htags WHERE id = $1`

	inventoryHTag := &models.InventoryHTag{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&inventoryHTag.ID,
		&inventoryHTag.Value,
	)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return inventoryHTag, nil
}

// GetAll will return the all inventory tags for list.
func (m *InventoryHTagModel) GetAll() ([]*models.InventoryHTag, error) {
	stmt := `SELECT id, value
	FROM inventory_htags`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	htags := []*models.InventoryHTag{}

	for rows.Next() {
		htag := &models.InventoryHTag{}
		err := rows.Scan(
			&htag.ID,
			&htag.Value,
		)

		if err != nil {
			return nil, err
		}

		htags = append(htags, htag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return htags, nil
}

// Insert will insert a new var of host into the database.
func (m *InventoryHTagModel) Insert(value string) (int, error) {
	stmt := `INSERT INTO inventory_htags (value)
	VALUES($1) RETURNING id`

	var lastInsertID int
	err := m.DB.QueryRow(stmt, value).Scan(&lastInsertID)
	if err != nil {
		// TODO
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && strings.Contains(pqErr.Code.Name(), "unique_violation") {
				return lastInsertID, models.ErrDuplicateEmail
			}
		}
	}

	return lastInsertID, nil
}

// Update ...
func (m *InventoryHTagModel) Update(id, value string) error {
	stmt := `UPDATE inventory_htags SET value = $2 WHERE id = $1`

	_, err := m.DB.Exec(stmt, id, value)
	if err != nil {
		return err
	}
	return nil
}

// Delete ...
func (m *InventoryHTagModel) Delete(id int) error {
	stmtHostTags := `DELETE FROM inventory_host_tags WHERE tag_id = $1;`

	_, err := m.DB.Exec(stmtHostTags, id)
	if err != nil {
		return err
	}

	stmtHTags := `DELETE FROM inventory_htags WHERE id = $1;`

	_, err = m.DB.Exec(stmtHTags, id)
	if err != nil {
		return err
	}

	return nil
}
