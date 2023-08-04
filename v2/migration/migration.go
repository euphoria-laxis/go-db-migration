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

	// ID must be first property of model structure
	IDfield := model.Field(0)

	tableMigration := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s\n(\n",
		table,
	)
	tableMigration += "		"
	tableMigration += toSnakeCase(IDfield.Name) + " "
	tableMigration += m.convertType(IDfield.Type.String()) + " "
	idValues := parseTag(IDfield.Tag.Get("migration"))
	for _, constraint := range strings.Split(idValues["constraints"], ",") {
		tableMigration += constraint + " "
	}
	tableMigration += "\n);"

	_, err := m.DB.Exec(tableMigration)
	if err != nil {
		return err
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
			err = m.generateMySqlColumnMigration(table, values)
			if err != nil {
				return err
			}
			break
		case DBDriverPostgres:
			err = m.generatePostgresColumnMigration(table, values)
			if err != nil {
				return err
			}
			break
		case DBDriverSQLite:
			err = m.generateSqliteColumnMigration(table, values)
			if err != nil {
				return err
			}
			break
		default:
			return fmt.Errorf("unknown driver: %v, allowed drivers: [mysql,postgres,sqlite]", m.Driver)
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
