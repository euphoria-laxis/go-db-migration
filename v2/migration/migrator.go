package migration

import "database/sql"

const (
	DBDriverPostgres DBDriver = iota
	DBDriverMySQL
)

type DBDriver int

func (d DBDriver) String() string {
	switch d {
	case DBDriverPostgres:
		return "postgres"
	case DBDriverMySQL:
		return "mysql"
	default:
		return "mysql"
	}
}

func NewDBDriver(driver string) DBDriver {
	switch driver {
	case "postgres":
		return DBDriverPostgres
	case "mysql":
		return DBDriverMySQL
	default:
		return DBDriverMySQL
	}
}

type Options struct {
	Driver            DBDriver
	SnakeCase         bool
	DB                *sql.DB
	DefaultTextSize   uint8
	IgnoreForeignKeys bool
	TablePrefix       string
}

type OptFunc func(*Options)

var defaultOptions = Options{
	SnakeCase:         true,
	DB:                nil,
	DefaultTextSize:   255,
	IgnoreForeignKeys: true,
	TablePrefix:       "",
}

func SetDriver(driver string) OptFunc {
	return func(opts *Options) {
		d := NewDBDriver(driver)
		opts.Driver = d
	}
}

func WithSnakeCase(active bool) OptFunc {
	return func(opts *Options) {
		opts.SnakeCase = active
	}
}

func SetDB(db *sql.DB) OptFunc {
	return func(opts *Options) {
		opts.DB = db
	}
}

func SetDefaultTextSize(size uint8) OptFunc {
	return func(opts *Options) {
		if size > 0 {
			opts.DefaultTextSize = size
		} else {
			opts.DefaultTextSize = 255
		}
	}
}

func WithForeignKeys(foreignKeys bool) OptFunc {
	return func(opts *Options) {
		opts.IgnoreForeignKeys = !foreignKeys
	}
}

func SetTablePrefix(prefix string) OptFunc {
	return func(opts *Options) {
		opts.TablePrefix = prefix
	}
}

type Migrator struct {
	Driver            DBDriver
	SnakeCase         bool
	DB                *sql.DB
	DefaultTextSize   uint8
	IgnoreForeignKeys bool
	TablePrefix       string
}

func NewMigrator(opts ...OptFunc) *Migrator {
	o := defaultOptions
	for _, fn := range opts {
		fn(&o)
	}
	migrator := Migrator{
		Driver:            o.Driver,
		DB:                o.DB,
		SnakeCase:         o.SnakeCase,
		DefaultTextSize:   o.DefaultTextSize,
		IgnoreForeignKeys: o.IgnoreForeignKeys,
		TablePrefix:       o.TablePrefix,
	}

	return &migrator
}
