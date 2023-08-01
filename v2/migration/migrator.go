package migration

import "database/sql"

const (
	DBDriverPostgres DBDriver = iota
	DBDriverMySQL
	DBDriverMariaDB
	DBDriverSQLite
	DBDriverSQLite3
)

type DBDriver int

func (d DBDriver) String() string {
	switch d {
	case DBDriverPostgres:
		return "postgres"
	case DBDriverMySQL:
		return "mysql"
	case DBDriverMariaDB:
		return "mysql"
	case DBDriverSQLite:
		return "sqlite"
	case DBDriverSQLite3:
		return "sqlite3"
	default:
		return "sqlite3"
	}
}

func NewDBDriver(driver string) DBDriver {
	switch driver {
	case "postgres":
		return DBDriverPostgres
	case "mysql":
		return DBDriverMySQL
	case "mariadb":
		return DBDriverMariaDB
	case "sqlite":
		return DBDriverSQLite
	case "sqlite3":
		return DBDriverSQLite3
	default:
		return DBDriverSQLite3
	}
}

type Options struct {
	SnakeCase bool
	DB        *sql.DB
	DBName    string
}

type OptFunc func(*Options)

var defaultOptions = Options{
	SnakeCase: true,
	DB:        nil,
	DBName:    "",
}

func WithSnakeCase(active bool) OptFunc {
	return func(opts *Options) {
		opts.SnakeCase = active
	}
}

func SetDBName(name string) OptFunc {
	return func(opts *Options) {
		opts.DBName = name
	}
}

func SetDB(db *sql.DB) OptFunc {
	return func(opts *Options) {
		opts.DB = db
	}
}

type Migrator struct {
	SnakeCase bool
	DB        *sql.DB
	Name      string
}

func NewMigrator(opts ...OptFunc) *Migrator {
	o := defaultOptions
	for _, fn := range opts {
		fn(&o)
	}
	migrator := Migrator{
		DB:        o.DB,
		Name:      o.DBName,
		SnakeCase: o.SnakeCase,
	}

	return &migrator
}
