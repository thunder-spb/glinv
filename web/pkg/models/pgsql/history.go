package pgsql

import (
	"database/sql"
	"glinv/pkg/models"
)

// HistoryModel type which wraps a sql.DB connection pool.
type HistoryModel struct {
	DB *sql.DB
}

// Event ...
func (m *HistoryModel) Event(userID, userEmail, entityID, entity, action, event string) {
	stmt := `INSERT INTO history (user_id, user_email, entity_id, entity, action, event) VALUES($1, $2, $3, $4, $5, $6) RETURNING id`
	m.DB.QueryRow(stmt, userID, userEmail, entityID, entity, action, event)
}

// GetAll will return the all ...
func (m *HistoryModel) GetAll() ([]*models.History, error) {
	stmt := `SELECT id, user_id, user_email, entity_id, entity, action, event, created FROM history ORDER BY created ASC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	all := []*models.History{}

	for rows.Next() {
		h := &models.History{}
		err := rows.Scan(
			&h.ID,
			&h.UserID,
			&h.UserEmail,
			&h.EntityID,
			&h.Entity,
			&h.Action,
			&h.Event,
			&h.Created,
		)

		if err != nil {
			return nil, err
		}

		all = append(all, h)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return all, nil
}
