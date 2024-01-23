package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func (m *Migrator) createPostgresSchema(table string, model reflect.Type) error {
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
		if strings.Contains(constraint, "auto_increment") {
			// replace 'auto_increment' with 'nextval('table_name_id_seq')' for postgres compatibility
			tableMigration = strings.Replace(tableMigration, "INT", "SERIAL", -1)
		} else {
			tableMigration += constraint + " "
		}
	}
	tableMigration += "\n);"
	_, err := m.DB.Exec(tableMigration)

	return err
}

func (m *Migrator) generatePostgresColumnMigration(table string, params map[string]string) error {
	infos, err := m.getPostgresSchemaInformation(table, params["column"])
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if infos != nil && convertSqlDataType(infos.DataType) != strings.ToLower(params["type"]) {
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
	infos, err = m.getPostgresSchemaInformation(table, params["column"])
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	constraints, hasConstraint := params["constraints"]
	if hasConstraint {
		for _, constraint := range strings.Split(constraints, ",") {
			if !checkConstraint(constraint) {
				fmt.Printf("[WARN] constraint %s is not valid and was ignored\n", constraint)
				continue
			}
			if infos != nil {
				if constraint == "not null" && infos.IsNullable {
					continue
				}
			}
			query = fmt.Sprintf(
				"ALTER TABLE %s ",
				table,
			)
			switch constraint {
			case "unique":
				query += fmt.Sprintf(
					"ADD CONSTRAINT %s UNIQUE(%s);\n",
					fmt.Sprintf("unique_%s_%s", table, params["column"]),
					params["column"],
				)
				break
			case "not null":
				query += fmt.Sprintf("ALTER COLUMN %s SET NOT NULL;\n", params["column"])
				break
			default:
				fmt.Printf("unknown constraint : %s\n", constraints)
				continue
			}
			_, err = m.DB.Exec(query)
			if err != nil {
				return err
			}
		}
	}
	infos, err = m.getPostgresSchemaInformation(table, params["column"])
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	defaultValue, hasDefaultValue := params["default"]
	if hasDefaultValue && defaultValue != infos.Default {
		if strings.Contains(params["type"], "VARCHAR") || strings.Contains(params["type"], "TEXT") {
			defaultValue = "'" + defaultValue + "'"
		}
		query = fmt.Sprintf(
			"ALTER TABLE %s ALTER COLUMN %s SET DEFAULT %s;\n",
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
		var indexInfo *PostgresIndexInfo
		indexInfo, err = m.verifyPostgresIndexExists(table, params["column"])
		if err != nil {
			return err
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if strings.Compare(indexInfo.IndexName, fmt.Sprintf("index_%s", params["column"])) != 0 || errors.Is(err, sql.ErrNoRows) {
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

type PostgresTableInfo struct {
	ColumnName string
	DataType   string
	IsNullable bool
	Default    interface{}
}

func (m *Migrator) getPostgresSchemaInformation(table, column string) (*PostgresTableInfo, error) {
	query := fmt.Sprintf(
		`select column_name, data_type, column_default, is_nullable 
				from INFORMATION_SCHEMA.COLUMNS where table_name = '%s' and column_name = '%s' ;`,
		table,
		column,
	)
	var nullable string
	var result PostgresTableInfo
	err := m.DB.QueryRow(query).Scan(&result.ColumnName, &result.DataType, &result.Default, &nullable)
	if err != nil {
		return nil, err
	}
	result.IsNullable = strings.Contains(nullable, "YES")

	return &result, err
}

type PostgresIndexInfo struct {
	TableName  string
	IndexName  string
	ColumnName string
}

func (m *Migrator) verifyPostgresIndexExists(table, column string) (*PostgresIndexInfo, error) {
	query := fmt.Sprintf(`select
				t.relname as table_name,
				i.relname as index_name,
				a.attname as column_name
			from
				pg_class t,
				pg_class i,
				pg_index ix,
				pg_attribute a
			where
				t.oid = ix.indrelid
			  and i.oid = ix.indexrelid
			  and a.attrelid = t.oid
			  and a.attnum = ANY(ix.indkey)
			  and t.relkind = 'r'
			  and t.relname = '%s'
			  and a.attname = '%s'
			order by
				t.relname,
				i.relname;`,
		table,
		column,
	)
	var result PostgresIndexInfo
	err := m.DB.QueryRow(query).Scan(&result.TableName, &result.IndexName, &result.ColumnName)
	if err != nil {
		return nil, err
	}

	return &result, err
}
