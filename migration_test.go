package go_db_migration

import (
	"database/sql"
	"github.com/euphoria-laxis/go-db-migration/v1/migration"
	"testing"
)

func TestGenerateDatabaseMigration(t *testing.T) {
	migrator := migration.NewMigrator()
	err := migrator.GenerateDatabaseMigration("migration")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateMigration(t *testing.T) {
	type model1 struct {
		TableDesc string       `json:"-" table:"model1" database:"migration"`
		ID        int          `json:"id" column:"user_id" datatype:"INT" constraint:"PRIMARY KEY NOT NULL UNIQUE AUTO_INCREMENT" index:"true"`
		Username  string       `json:"username" column:"username" datatype:"VARCHAR(100)" constraint:"UNIQUE NOT NULL" index:"true"`
		CreatedAt sql.NullTime `json:"created_at" column:"created_at" datatype:"DATETIME" default:"NOW()"`
		UpdatedAt sql.NullTime `json:"updated_at" column:"updated_at" datatype:"DATETIME" default:"NOW()"`
		DeletedAt sql.NullTime `json:"deleted_at" column:"deleted_at" datatype:"DATETIME"`
		Name      string       `json:"name" column:"name" datatype:"VARCHAR(128)" constraint:"NOT NULL"`
		Content   string       `json:"content" column:"content" datatype:"TEXT" constraint:"NOT NULL"`
		Role      string       `json:"role" column:"role" datatype:"VARCHAR(255)" constraint:"NOT NULL" DEFAULT:"user"`
	}
	type model2 struct {
		TableDesc string       `json:"-" table:"model2" database:"migration"`
		ID        int          `json:"id" column:"user_id" datatype:"INT" constraint:"PRIMARY KEY NOT NULL UNIQUE AUTO_INCREMENT" index:"true"`
		Username  string       `json:"username" column:"username" datatype:"VARCHAR(100)" constraint:"UNIQUE NOT NULL" index:"true"`
		CreatedAt sql.NullTime `json:"created_at" column:"created_at" datatype:"DATETIME" default:"NOW()"`
		UpdatedAt sql.NullTime `json:"updated_at" column:"updated_at" datatype:"DATETIME" default:"NOW()"`
		DeletedAt sql.NullTime `json:"deleted_at" column:"deleted_at" datatype:"DATETIME"`
		Name      string       `json:"name" column:"name" datatype:"VARCHAR(128)" constraint:"NOT NULL"`
		Content   string       `json:"content" column:"content" datatype:"TEXT" constraint:"NOT NULL"`
		Role      string       `json:"role" column:"role" datatype:"VARCHAR(255)" constraint:"NOT NULL" DEFAULT:"user"`
	}
	m1 := model1{}
	m2 := model2{}
	migrator := migration.NewMigrator()
	err := migrator.GenerateMigration(m1, m2)
	if err != nil {
		t.Fatal(err)
	}
}
