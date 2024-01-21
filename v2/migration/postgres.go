package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

func (m *Migrator) createPostgresSchema(table string, model reflect.Type) error {
	return nil
}

func (m *Migrator) generatePostgresColumnMigration(table string, params map[string]string) error {
	return nil
}

func (m *Migrator) getPostgresSchemaInformation(table, column string) (*interface{}, error) {
	return nil, nil
}
