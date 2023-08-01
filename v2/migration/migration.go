package migration

import (
	"errors"
	"fmt"
	"os"
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

func createDirectory() error {
	path := "sql"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(path, os.ModePerm)
	}

	return nil
}

func parseTag(tag string) map[string]string {
	parsed := make(map[string]string)
	items := strings.Split(tag, ";")
	for _, item := range items {
		s := strings.Split(item, ":")
		if len(s) != 2 {
			continue
		}
		key := s[0]
		value := s[1]
		parsed[key] = value
	}

	return parsed
}

func convertType(kind string) string {
	r, err := regexp.MatchString("Time$", kind)
	if err != nil {
		return ""
	}
	if r {
		return "DATETIME"
	}
	switch kind {
	case "int":
		return "INT"
	case "float":
		return "FLOAT"
	case "string":
		return "VARCHAR(255)"
	case "time.Time":
		return "DATETIME"
	case "sql.NullTime":
		return "DATETIME"
	case "null.Time":
		return "DATETIME"
	case "bool":
		return "BOOL"
	default:
		return ""
	}
}

func (m *Migrator) generateColumnMigration(table string, params map[string]string) error {
	columnMigration := fmt.Sprintf(
		"ALTER TABLE %s ADD COLUMN %s %s;\n",
		table,
		params["column"],
		params["type"],
	)
	_, err := m.DB.Exec(columnMigration)
	if err != nil {
		return err
	}
	constraints, hasConstraint := params["constraints"]
	if hasConstraint {
		for _, constraint := range strings.Split(constraints, ",") {
			columnMigration = fmt.Sprintf(
				"ALTER TABLE %s ",
				table,
			)
			columnMigration += fmt.Sprintf(
				"MODIFY %s %s %s;\n",
				params["column"],
				params["type"],
				constraint,
			)
			_, err = m.DB.Exec(columnMigration)
			if err != nil {
				return err
			}
		}
	}
	defaultValue, hasDefaultValue := params["default"]
	if hasDefaultValue {
		if strings.Contains(params["type"], "VARCHAR") {
			defaultValue = "'" + defaultValue + "'"
		}
		columnMigration = fmt.Sprintf(
			"ALTER TABLE %s ALTER %s SET DEFAULT %s;\n",
			table,
			params["column"],
			defaultValue,
		)

		_, err = m.DB.Exec(columnMigration)
		if err != nil {
			return err
		}
	}
	_, isIndex := params["index"]
	if isIndex {
		columnMigration = fmt.Sprintf(
			"CREATE INDEX index_%s ON %s (%s);\n",
			params["column"],
			table,
			params["column"],
		)
		_, err = m.DB.Exec(columnMigration)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) generateTableMigration(model reflect.Type) error {
	var table string
	if model.Kind() == reflect.Ptr {
		table = model.Elem().Name()
	} else {
		table = model.Name()
	}
	table = toSnakeCase(table)

	// ID must be first property of model structure
	IDfield := model.Field(0)

	tableMigration := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s\n(\n",
		table,
	)
	tableMigration += "		"
	tableMigration += toSnakeCase(IDfield.Name) + " "
	tableMigration += convertType(IDfield.Type.String()) + " "
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
			values["type"] = convertType(kind)
		}
		err = m.generateColumnMigration(table, values)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) MigrateModels(models ...interface{}) error {
	for _, model := range models {
		reflection := reflect.TypeOf(model)
		err := m.generateTableMigration(reflection)
		if err != nil {
			return err
		}
	}

	return nil
}
