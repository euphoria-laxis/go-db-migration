package migration

import (
	"database/sql"
	"errors"
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

type Statistic struct {
	NonUnique bool
	IndexName string
	Nullable  bool
}

func (m *Migrator) generateColumnMigration(table string, params map[string]string) error {
	infos, err := m.getSchemaInformation(table, params["column"])
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if infos != nil && convertSqlDataType(infos.Type) != strings.ToLower(params["type"]) {
		query := fmt.Sprintf(
			"ALTER TABLE %s DROP COLUMN %s;",
			table,
			params["column"],
		)
		_, err = m.DB.Exec(query)
		if err != nil {
			return err
		}
	}
	query := fmt.Sprintf(
		"ALTER TABLE %s ADD COLUMN %s %s;\n",
		table,
		params["column"],
		params["type"],
	)
	_, err = m.DB.Exec(query)
	if err != nil {
		rowExist := strings.Contains(err.Error(), "Duplicate column name")
		if !rowExist {
			return err
		}
	}
	constraints, hasConstraint := params["constraints"]
	if hasConstraint {
		for _, constraint := range strings.Split(constraints, ",") {
			if !checkConstraint(constraint) {
				fmt.Printf("[WARN] constraint %s is not valid and was ignored\n", constraint)
				continue
			}
			if infos != nil {
				if constraint == "not null" && infos.Null == "YES" {
					continue
				} else if constraint == "unique" && infos.Key == "UNI" {
					continue
				}
			}
			query = fmt.Sprintf(
				"ALTER TABLE %s ",
				table,
			)
			query += fmt.Sprintf(
				"MODIFY %s %s %s;\n",
				params["column"],
				params["type"],
				constraint,
			)
			_, err = m.DB.Exec(query)
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
		query = fmt.Sprintf(
			"ALTER TABLE %s ALTER %s SET DEFAULT %s;\n",
			table,
			params["column"],
			defaultValue,
		)

		_, err = m.DB.Exec(query)
		if err != nil {
			return err
		}
	}
	_, isIndex := params["index"]
	if isIndex {
		query = fmt.Sprintf(
			"SELECT NON_UNIQUE, INDEX_NAME, NULLABLE FROM information_schema.statistics WHERE table_name = '%s' AND column_name = '%s';",
			table,
			params["column"],
		)
		var statistic Statistic
		err = m.DB.QueryRow(query).Scan(&statistic.NonUnique, &statistic.IndexName, &statistic.Nullable)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if errors.Is(err, sql.ErrNoRows) {
			query = fmt.Sprintf(
				"CREATE INDEX index_%s ON %s (%s);\n",
				params["column"],
				table,
				params["column"],
			)
			_, err = m.DB.Exec(query)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type Result struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Extra   string
	Default interface{}
}

func (m *Migrator) getSchemaInformation(table, column string) (*Result, error) {
	query := fmt.Sprintf(
		"select COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, EXTRA, COLUMN_DEFAULT from information_schema.COLUMNS where table_name = '%s' and column_name = '%s' ;",
		table,
		column,
	)
	var result Result
	err := m.DB.QueryRow(query).Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Extra, &result.Default)
	if err != nil {
		return nil, err
	}

	return &result, err
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
		err := m.migrateModel(reflection)
		if err != nil {
			return err
		}
	}

	return nil
}
