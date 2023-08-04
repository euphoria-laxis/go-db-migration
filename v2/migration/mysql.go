package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Statistic struct {
	NonUnique bool
	IndexName string
	Nullable  bool
}

func (m *Migrator) generateMySqlColumnMigration(table string, params map[string]string) error {
	infos, err := m.getMySqlSchemaInformation(table, params["column"])
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

func (m *Migrator) getMySqlSchemaInformation(table, column string) (*Result, error) {
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
