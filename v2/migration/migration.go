package migration

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func toSnakeCase(pattern string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(pattern, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}

func (m *Migrator) migrateModel(model reflect.Type) error {
	var table string
	if model.Kind() == reflect.Ptr {
		table = model.Elem().Name()
	} else {
		table = model.Name()
	}
	if m.TablePrefix != "" {
		table = m.TablePrefix + table
	}
	if m.SnakeCase {
		table = toSnakeCase(table)
	}

	switch m.Driver {
	case DBDriverMySQL:
		err := m.createMySqlSchemas(table, model)
		if err != nil {
			return err
		}
	case DBDriverPostgres:
		err := m.createPostgresSchema(table, model)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown driver: %v, allowed drivers: [mysql,postgres]", m.Driver)
	}

	for i := 1; i < model.NumField(); i++ {
		column := model.Field(i)
		if strings.Compare(column.Name, "-") == 0 {
			continue
		}
		values := parseTag(column.Tag.Get("migration"))
		values["column"] = toSnakeCase(column.Name)
		_, hasType := values["type"]
		if !hasType {
			kind := column.Type.String()
			values["type"] = m.convertType(kind)
		}
		switch m.Driver {
		case DBDriverMySQL:
			err := m.generateMySqlColumnMigration(table, values)
			if err != nil {
				return err
			}
		case DBDriverPostgres:
			err := m.generatePostgresColumnMigration(table, values)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown driver: %v, allowed drivers: [mysql,postgres]", m.Driver)
		}
	}

	return nil
}

func (m *Migrator) MigrateModels(models ...interface{}) error {
	for _, model := range models {
		reflection := reflect.TypeOf(model)
		err := m.migrateModel(reflection)
		if err != nil {
			return err
		}
	}

	return nil
}
