package migration

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func createDirectory() error {
	path := "sql"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(path, os.ModePerm)
	}

	return nil
}

func NewMigrator() *Migrator {
	return new(Migrator)
}

func (db *Migrator) GenerateDatabaseMigration(dbname string) error {
	err := createDirectory()
	if err != nil {
		return err
	}
	sqlFile, err := os.Create("sql/database_" + strconv.Itoa(int(time.Now().Unix())) + ".sql")
	if err != nil {
		return ErrDatabaseSQLFile
	}
	defer sqlFile.Close()
	dbClean := fmt.Sprintf("DROP DATABASE IF EXISTS %s;\n", dbname)
	sqlFile.WriteString(dbClean)
	dbDesc := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s "+
		"CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin;", dbname)
	sqlFile.WriteString(dbDesc)
	sqlFile.Sync()

	return err
}

func (db *Migrator) GenerateMigration(models ...interface{}) error {
	err := createDirectory()
	if err != nil {
		return err
	}
	for _, model := range models {
		structReflect := reflect.TypeOf(model)
		tableNameStructField, _ := structReflect.FieldByName("TableDesc")
		tableName := tableNameStructField.Tag.Get("table")
		dbName := tableNameStructField.Tag.Get("database")
		tableDesc := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s\n(\n",
			dbName, tableName)
		sqlFile, err := os.Create("sql/" + tableName + "_schema_" +
			strconv.Itoa(int(time.Now().Unix())) + ".sql")
		if err != nil {
			log.Println(err)

			return err
		}
		defer sqlFile.Close()
		sqlFile.WriteString(tableDesc)
		var indexColumns []string
		for i := 1; i < structReflect.NumField(); i++ {
			fieldTag := structReflect.Field(i).Tag
			columnName := fieldTag.Get("column")
			if strings.Compare(columnName, "-") == 0 {
				continue
			}
			columnType := fieldTag.Get("datatype")
			columnConstraint := fieldTag.Get("constraint")
			columnDefault := fieldTag.Get("default")
			var columnDesc string
			columnDesc = "\t" + columnName + " " + columnType
			if strings.Compare(columnConstraint, "") != 0 && strings.Compare(columnConstraint, " ") != 0 {
				columnDesc += " " + columnConstraint
			}
			if strings.Compare(columnDefault, "") != 0 && strings.Compare(columnDefault, " ") != 0 {
				columnDesc += " DEFAULT " + columnDefault
			}
			if strings.Compare(fieldTag.Get("index"), "true") == 0 {
				indexColumns = append(indexColumns, columnName)
			}
			if i != structReflect.NumField()-1 || len(indexColumns) > 0 {
				columnDesc += ","
			}
			columnDesc += "\n"
			sqlFile.WriteString(columnDesc)
		}
		for j := 0; j < len(indexColumns); j++ {
			indexDesc := fmt.Sprintf("\tINDEX %s_per_%s (%s)", tableName,
				indexColumns[j], indexColumns[j])
			if j != len(indexColumns)-1 {
				indexDesc += ","
			}
			indexDesc += "\n"
			sqlFile.WriteString(indexDesc)
		}
		sqlFile.WriteString(");")
		sqlFile.Sync()
	}

	return nil
}
