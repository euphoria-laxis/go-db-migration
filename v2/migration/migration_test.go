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
	type Model2 struct {
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
	m1 := model1{}
	m2 := Model2{}
	user := "migration_test"
	password := "password@123"
	dbname := "migration"
	if strings.Compare(os.Getenv("ENVIRONMENT"), "test") == 0 {
		user = os.Getenv("MYSQL_USER")
		password = os.Getenv("MYSQL_PASSWORD")
		dbname = os.Getenv("MYSQL_DATABASE")
	}
	cfg := mysql.Config{
		User:                 user,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               dbname,
		AllowNativePasswords: true,
	}
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		t.Fatal(pingErr)
	}
	migrator := NewMigrator(
		SetDB(db),
		SetTablePrefix("app_"),
		WithForeignKeys(true),
		WithSnakeCase(true),
		SetDefaultTextSize(128),
		SetDriver("mysql"),
	)
	err = migrator.MigrateModels(m1, m2)
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
	type Model2 struct {
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
	m1 := model1{}
	m2 := Model2{}
	host := "localhost"
	port := 5432
	user := "migration_test"
	password := "password@123"
	dbname := "migration"
	if strings.Compare(os.Getenv("ENVIRONMENT"), "test") == 0 {
		user = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname = os.Getenv("POSTGRES_DB")
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)
	// Connect to database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
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
	err = migrator.MigrateModels(m1, m2)
	if err != nil {
		t.Fatal(err)
	}
}
