package migration

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGenerateMySQLMigrations(t *testing.T) {
	type model1 struct {
		ID        int          `json:"id" migration:"constraints:primary key,not null,unique,auto_increment;index"`
		Username  string       `json:"username" migration:"constraints:not null,unique;index"`
		CreatedAt time.Time    `json:"created_at" migration:"default:now()"`
		UpdatedAt time.Time    `json:"updated_at" migration:"default:now()"`
		DeletedAt sql.NullTime `json:"deleted_at"`
		Name      string       `json:"name" migration:"constraint:not null"`
		Content   string       `json:"content" migration:"type:text;constraints:not null"`
		Role      string       `json:"role" migration:"constraints:not null;default:user"`
	}
	type model2 struct {
		ID        int          `json:"id" migration:"constraints:primary key,not null,unique,auto_increment;index"`
		Username  string       `json:"username" migration:"constraints:not null,unique;index"`
		CreatedAt time.Time    `json:"created_at" migration:"default:now()"`
		UpdatedAt time.Time    `json:"updated_at" migration:"default:now()"`
		DeletedAt sql.NullTime `json:"deleted_at"`
		Name      string       `json:"name" migration:"constraints:not null"`
		Content   string       `json:"content" migration:"type:text;constraints:not null"`
		Role      string       `json:"role" migration:"constraints:not null;default:user"`
		Valid     bool         `json:"valid" migration:"default:false"`
	}
	user := "migration_test"
	passwd := "password@123"
	dbname := "migration"
	if strings.Contains(os.Getenv("ENVIRONMENT"), "test") {
		user = os.Getenv("MYSQL_USER")
		passwd = os.Getenv("MYSQL_PASSWORD")
		dbname = os.Getenv("MYSQL_DATABASE")
	}
	// Generate MySQL config
	cfg := mysql.Config{
		User:                 user,
		Passwd:               passwd,
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               dbname,
		AllowNativePasswords: true,
	}
	// Connect to MySQL
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
	// Check if MySQL server is accessible
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
	migrator := NewMigrator(
		SetDB(db),
		SetTablePrefix("app_"),
		WithForeignKeys(true),
		WithSnakeCase(true),
		SetDefaultTextSize(128),
		SetDriver("mysql"),
	)
	err = migrator.MigrateModels(model1{}, model2{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGeneratePostgresMigrations(t *testing.T) {
	type model1 struct {
		ID        int          `json:"id" migration:"constraints:primary key,not null,unique,auto_increment;index"`
		Username  string       `json:"username" migration:"constraints:not null,unique;index"`
		CreatedAt time.Time    `json:"created_at" migration:"default:now()"`
		UpdatedAt time.Time    `json:"updated_at" migration:"default:now()"`
		DeletedAt sql.NullTime `json:"deleted_at"`
		Name      string       `json:"name" migration:"constraint:not null"`
		Content   string       `json:"content" migration:"type:text;constraints:not null"`
		Role      string       `json:"role" migration:"constraints:not null;default:user"`
	}
	type model2 struct {
		ID        int          `json:"id" migration:"constraints:primary key,not null,unique,auto_increment;index"`
		Username  string       `json:"username" migration:"constraints:not null,unique;index"`
		CreatedAt time.Time    `json:"created_at" migration:"default:now()"`
		UpdatedAt time.Time    `json:"updated_at" migration:"default:now()"`
		DeletedAt sql.NullTime `json:"deleted_at"`
		Name      string       `json:"name" migration:"constraints:not null"`
		Content   string       `json:"content" migration:"type:text;constraints:not null"`
		Role      string       `json:"role" migration:"constraints:not null;default:user"`
		Valid     bool         `json:"valid" migration:"default:false"`
	}
	host := "localhost"
	port := 5432
	user := "migration"
	password := "password@123"
	dbname := "migration_test"
	if strings.Contains(os.Getenv("ENVIRONMENT"), "test") {
		user = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname = os.Getenv("POSTGRES_DB")
	}
	// Create Postgres DSN
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)
	// Connect to Postgres
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	// Check if Postgres server is accessible
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
	migrator := NewMigrator(
		SetDB(db),
		SetTablePrefix("app_"),
		WithForeignKeys(true),
		WithSnakeCase(true),
		SetDefaultTextSize(128),
		SetDriver("postgres"),
	)
	err = migrator.MigrateModels(model1{}, model2{})
	if err != nil {
		t.Fatal(err)
	}
}
