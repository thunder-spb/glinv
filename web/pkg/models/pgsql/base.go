package pgsql

import (
	"database/sql"

	"glinv/pkg/models"
)

// BaseTemplateModel type which wraps a sql.DB connection pool.
type BaseTemplateModel struct {
	DB *sql.DB
}

// InsertBaseTpl ...
func (m *BaseTemplateModel) InsertBaseTpl(typeTpl, valueTpl string) (int, error) {
	stmt := `INSERT INTO server_templates (type, value) VALUES($1, $2) RETURNING id`

	var lastInsertID int
	if err := m.DB.QueryRow(stmt, typeTpl, valueTpl).Scan(&lastInsertID); err != nil {
		return lastInsertID, err
	}

	return lastInsertID, nil
}

// GetBaseTpl ...
func (m *BaseTemplateModel) GetBaseTpl(id string) (*models.BaseTemplate, error) {
	stmt := `SELECT id, type, value FROM server_templates  WHERE id = $1`

	baseTemplate := &models.BaseTemplate{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&baseTemplate.ID,
		&baseTemplate.Type,
		&baseTemplate.Value,
	)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return baseTemplate, nil
}

// GetAllBaseTpl ...
func (m *BaseTemplateModel) GetAllBaseTpl() ([]*models.BaseTemplate, error) {
	stmt := `SELECT id, type, value, created, updated 
	FROM server_templates ORDER BY created DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	baseTemplates := []*models.BaseTemplate{}

	for rows.Next() {
		bt := &models.BaseTemplate{}
		err := rows.Scan(
			&bt.ID,
			&bt.Type,
			&bt.Value,
			&bt.Created,
			&bt.Updated,
		)

		if err != nil {
			return nil, err
		}

		bt.Parameters = m.GetCountBaseTplItems(bt.ID, bt.Type)

		baseTemplates = append(baseTemplates, bt)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return baseTemplates, nil
}

// EditBaseTpl ...
func (m *BaseTemplateModel) EditBaseTpl(idTemplate, value string) error {
	stmt := `UPDATE server_templates SET value = $2 WHERE id = $1`
	_, err := m.DB.Exec(stmt, idTemplate, value)
	if err != nil {
		return err
	}

	return nil
}

//
// Hard
//

// DeleteBaseTplHard ...
func (m *BaseTemplateModel) DeleteBaseTplHard(idTemplate string) error {
	stmt := `DELETE FROM server_hard_base WHERE template_id = $1;`
	_, err := m.DB.Exec(stmt, idTemplate)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM server_templates WHERE id = $1;`
	_, err = m.DB.Exec(stmt, idTemplate)
	if err != nil {
		return err
	}

	return nil
}

// InsertBaseTplHardItem ...
func (m *BaseTemplateModel) InsertBaseTplHardItem(idTemplate, title, value string) (int, error) {
	stmt := `INSERT INTO server_hard_base (template_id, title, value) VALUES($1, $2, $3) RETURNING id`

	var lastInsertID int
	if err := m.DB.QueryRow(stmt, idTemplate, title, value).Scan(&lastInsertID); err != nil {
		return lastInsertID, err
	}

	return lastInsertID, nil
}

// DeleteBaseTplHardItem ...
func (m *BaseTemplateModel) DeleteBaseTplHardItem(idTemplate, id string) error {
	stmt := `DELETE FROM server_hard_base WHERE template_id = $1 AND id = $2 ;`

	_, err := m.DB.Exec(stmt, idTemplate, id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllBaseTplHardItem ...
func (m *BaseTemplateModel) GetAllBaseTplHardItem(idTemplate string) ([]*models.BaseTemplateItem, error) {
	stmt := `SELECT id, template_id, title, value, created, updated 
	FROM server_hard_base WHERE template_id = $1 ORDER BY created DESC`

	rows, err := m.DB.Query(stmt, idTemplate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	BaseTemplatesItem := []*models.BaseTemplateItem{}

	for rows.Next() {
		bth := &models.BaseTemplateItem{}
		err := rows.Scan(
			&bth.ID,
			&bth.IDTpl,
			&bth.Title,
			&bth.Value,
			&bth.Created,
			&bth.Updated,
		)

		if err != nil {
			return nil, err
		}

		BaseTemplatesItem = append(BaseTemplatesItem, bth)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return BaseTemplatesItem, nil
}

//
// Package
//

// DeleteBaseTplPackage ...
func (m *BaseTemplateModel) DeleteBaseTplPackage(idTemplate string) error {
	stmt := `DELETE FROM server_soft_base WHERE template_id = $1;`
	_, err := m.DB.Exec(stmt, idTemplate)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM server_templates WHERE id = $1;`
	_, err = m.DB.Exec(stmt, idTemplate)
	if err != nil {
		return err
	}

	return nil
}

// InsertBaseTplPackageItem ...
func (m *BaseTemplateModel) InsertBaseTplPackageItem(idTemplate, title, value string) (int, error) {
	stmt := `INSERT INTO server_soft_base (template_id, title, value) VALUES($1, $2, $3) RETURNING id`

	var lastInsertID int
	if err := m.DB.QueryRow(stmt, idTemplate, title, value).Scan(&lastInsertID); err != nil {
		return lastInsertID, err
	}

	return lastInsertID, nil
}

// DeleteBaseTplPackageItem ...
func (m *BaseTemplateModel) DeleteBaseTplPackageItem(idTemplate, id string) error {
	stmt := `DELETE FROM server_soft_base WHERE template_id = $1 AND id = $2 ;`

	_, err := m.DB.Exec(stmt, idTemplate, id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllBaseTplPackageItem ...
func (m *BaseTemplateModel) GetAllBaseTplPackageItem(idTemplate string) ([]*models.BaseTemplateItem, error) {
	stmt := `SELECT id, template_id, title, value, created, updated 
	FROM server_soft_base WHERE template_id = $1 ORDER BY created DESC`

	rows, err := m.DB.Query(stmt, idTemplate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	BaseTemplatesItem := []*models.BaseTemplateItem{}

	for rows.Next() {
		bth := &models.BaseTemplateItem{}
		err := rows.Scan(
			&bth.ID,
			&bth.IDTpl,
			&bth.Title,
			&bth.Value,
			&bth.Created,
			&bth.Updated,
		)

		if err != nil {
			return nil, err
		}

		BaseTemplatesItem = append(BaseTemplatesItem, bth)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return BaseTemplatesItem, nil
}

//
// ResolvConf
//

// DeleteBaseTplResolvConf ...
func (m *BaseTemplateModel) DeleteBaseTplResolvConf(idTemplate string) error {
	stmt := `DELETE FROM server_resolvconf_base WHERE template_id = $1;`
	_, err := m.DB.Exec(stmt, idTemplate)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM server_templates WHERE id = $1;`
	_, err = m.DB.Exec(stmt, idTemplate)
	if err != nil {
		return err
	}

	return nil
}

// InsertBaseTplResolvConfItem ...
func (m *BaseTemplateModel) InsertBaseTplResolvConfItem(idTemplate, title, value string) (int, error) {
	stmt := `INSERT INTO server_resolvconf_base (template_id, title, value) VALUES($1, $2, $3) RETURNING id`

	var lastInsertID int
	if err := m.DB.QueryRow(stmt, idTemplate, title, value).Scan(&lastInsertID); err != nil {
		return lastInsertID, err
	}

	return lastInsertID, nil
}

// DeleteBaseTplResolvConfItem ...
func (m *BaseTemplateModel) DeleteBaseTplResolvConfItem(idTemplate, id string) error {
	stmt := `DELETE FROM server_resolvconf_base WHERE template_id = $1 AND id = $2 ;`

	_, err := m.DB.Exec(stmt, idTemplate, id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllBaseTplResolvConfItem ...
func (m *BaseTemplateModel) GetAllBaseTplResolvConfItem(idTemplate string) ([]*models.BaseTemplateItem, error) {
	stmt := `SELECT id, template_id, title, value, created, updated 
	FROM server_resolvconf_base WHERE template_id = $1 ORDER BY created DESC`

	rows, err := m.DB.Query(stmt, idTemplate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	BaseTemplatesItem := []*models.BaseTemplateItem{}

	for rows.Next() {
		bth := &models.BaseTemplateItem{}
		err := rows.Scan(
			&bth.ID,
			&bth.IDTpl,
			&bth.Title,
			&bth.Value,
			&bth.Created,
			&bth.Updated,
		)

		if err != nil {
			return nil, err
		}

		BaseTemplatesItem = append(BaseTemplatesItem, bth)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return BaseTemplatesItem, nil
}

//
// SysCtl
//

// DeleteBaseTplSysCtl ...
func (m *BaseTemplateModel) DeleteBaseTplSysCtl(idTemplate string) error {
	stmt := `DELETE FROM server_sysctl_base WHERE template_id = $1;`
	_, err := m.DB.Exec(stmt, idTemplate)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM server_templates WHERE id = $1;`
	_, err = m.DB.Exec(stmt, idTemplate)
	if err != nil {
		return err
	}

	return nil
}

// InsertBaseTplSysCtlItem ...
func (m *BaseTemplateModel) InsertBaseTplSysCtlItem(idTemplate, title, value string) (int, error) {
	stmt := `INSERT INTO server_sysctl_base (template_id, title, value) VALUES($1, $2, $3) RETURNING id`

	var lastInsertID int
	if err := m.DB.QueryRow(stmt, idTemplate, title, value).Scan(&lastInsertID); err != nil {
		return lastInsertID, err
	}

	return lastInsertID, nil
}

// DeleteBaseTplSysCtlItem ...
func (m *BaseTemplateModel) DeleteBaseTplSysCtlItem(idTemplate, id string) error {
	stmt := `DELETE FROM server_sysctl_base WHERE template_id = $1 AND id = $2 ;`

	_, err := m.DB.Exec(stmt, idTemplate, id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllBaseTplSysCtlItem ...
func (m *BaseTemplateModel) GetAllBaseTplSysCtlItem(idTemplate string) ([]*models.BaseTemplateItem, error) {
	stmt := `SELECT id, template_id, title, value, created, updated 
	FROM server_sysctl_base WHERE template_id = $1 ORDER BY created DESC`

	rows, err := m.DB.Query(stmt, idTemplate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	BaseTemplatesItem := []*models.BaseTemplateItem{}

	for rows.Next() {
		bth := &models.BaseTemplateItem{}
		err := rows.Scan(
			&bth.ID,
			&bth.IDTpl,
			&bth.Title,
			&bth.Value,
			&bth.Created,
			&bth.Updated,
		)

		if err != nil {
			return nil, err
		}

		BaseTemplatesItem = append(BaseTemplatesItem, bth)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return BaseTemplatesItem, nil
}

//
// Count
//

// GetCountBaseTemplates ...
func (m *ServerAgentModel) GetCountBaseTemplates() int {
	stmt := `SELECT COUNT(*) FROM server_templates;`
	var count int
	m.DB.QueryRow(stmt).Scan(&count)
	return count
}

// GetCountBaseTplItems ...
func (m *BaseTemplateModel) GetCountBaseTplItems(id int, typeTpl string) int {
	var count int

	if typeTpl == "hard" {
		stmt := `SELECT COUNT(*) FROM server_hard_base WHERE template_id = $1;`
		m.DB.QueryRow(stmt, id).Scan(&count)
	}

	if typeTpl == "packaege" {
		stmt := `SELECT COUNT(*) FROM server_soft_base WHERE template_id = $1;`
		m.DB.QueryRow(stmt, id).Scan(&count)
	}

	if typeTpl == "resolvconf" {
		stmt := `SELECT COUNT(*) FROM server_resolvconf_base WHERE template_id = $1;`
		m.DB.QueryRow(stmt, id).Scan(&count)
	}

	if typeTpl == "sysctl" {
		stmt := `SELECT COUNT(*) FROM server_sysctl_base WHERE template_id = $1;`
		m.DB.QueryRow(stmt, id).Scan(&count)
	}

	return count
}
