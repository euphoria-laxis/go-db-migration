package migration

import "reflect"

func (m *Migrator) createSqliteSchema(table string, model reflect.Type) error {
	return nil
}

func (m *Migrator) generateSqliteColumnMigration(table string, params map[string]string) error {
	return nil
}

func (m *Migrator) getSqliteSchemaInformation(table, column string) (*interface{}, error) {
	return nil, nil
}
