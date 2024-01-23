# go-db-migration

Go package to generate and execute migration schemas using structure tags

## Release notes

* **Release v2.1.1**
  * Remove SQLite from config and roadmap *(because it's a mess with no column update method)*.
  * Fix Postgres columns update due to type comparison error.
  * Don't update default value if it did not change for Postgres and MySQL.
  * Close database clients in tests *(sorry I forgot, I know it's bad but I made this at 2am)*.
* **Release v2.1.0**
  * Update MySQL migration to separate schema creation.
  * Add Postgres support.
  * Update tests.
  * Add Postgres container to docker-compose.
* **Release v2.0.0**:
  * GitHub action for test.
  * Migrator config.
  * Add database drivers *(MySQL, Postgres, SQlite)* to config.
  * Separate MySQL migration support from Migrator instance.
  * Add docker-compose CI for tests and GitHub workflow.

## Documentation

### Installation

To add package to your go mod run :
````bash
go get github.com/euphoria-laxis/go-db-migration@v2.1.1
````
To generate the schema add the `migration` tag to your model structure.

````go
type model1 struct {
    ID        int          `migration:"constraints:primary key,not null,unique,auto_increment;index"`
    Username  string       `migration:"constraints:not null,unique;index"`
    CreatedAt time.Time    `migration:"default:now()"`
    UpdatedAt time.Time    `migration:"default:now()"`
    DeletedAt sql.NullTime `migration:"default:now()"`
    Name      string       `migration:"constraint:not null"`
    Content   string       `migration:"type:text;constraints:not null"`
    Role      string       `migration:"constraints:not null;default:user"`
}
type Model2 struct {
    ID        int          `migration:"constraints:primary key,not null,unique,auto_increment;index"`
    Username  string       `migration:"constraints:not null,unique;index"`
    CreatedAt time.Time    `migration:"default:now()"`
    UpdatedAt time.Time    `migration:"default:now()"`
    DeletedAt sql.NullTime `migration:"default:now()"`
    Name      string       `migration:"constraints:not null"`
    Content   string       `migration:"type:text;constraints:not null"`
    Role      string       `migration:"constraints:not null;default:user"`
    Valid     bool         `migration:"default:false"`
}
````

Then play the migrations 

````go
// create model
m1 := model1{}
m2 := Model2{}
// connect to database
cfg := mysql.Config{
    User:                 "migration_test",
    Passwd:               "password@123",
    Net:                  "tcp",
    Addr:                 "127.0.0.1:3306",
    DBName:               "migration",
    AllowNativePasswords: true,
}
db, _ := sql.Open("mysql", cfg.FormatDSN())
// set migration options
migrator := NewMigrator(
    SetDB(db),
    SetTablePrefix("app_"),
    WithForeignKeys(true),
    WithSnakeCase(true),
    SetDefaultTextSize(128),
    SetDriver("mysql"),
)
// generate and execute schemas
err := migrator.MigrateModels(m1, m2)
if err != nil {
    panic(err)
}
````

This will generate the migration files and execute them to update the database.

### Usage

The column type will be determined by the type used in the structure, for **TEXT** datatype you
must set in the structure tag the text type. 

#### Tags

|       Tag       |         Usage          |                   Values                   |
|:---------------:|:----------------------:|:------------------------------------------:|
| **constraints** | Add column constraints | primary key,not null,unique,auto_increment |
|    **index**    |      Create index      |                                            |
|   **default**   |   Add default value    |         float, int, bool or string         |
|    **type**     |    Set column type     |                    text                    |

#### Drivers

|    Driver    |     Available      |             Availability status             |
|:------------:|:------------------:|:-------------------------------------------:|
|  **MySQL**   | :white_check_mark: |                  Available                  |
| **Postgres** | :white_check_mark: |                  Available                  |
| **MariaDB**  |     :warning:      | Use MySQL driver for MariaDB *(not tested)* |

## Roadmap

### Planned features

* Foreign keys creation and updates.
* Handling more datatypes:
  * Postgres:
    * uuid
    * date without timezone, timestamps.
    * float4
    * int(n)
    * path
    * json
    * array
    * bytea
    * macaddr, macaddr8
  * MySQL:
    * uuid
    * date, time, timestamp, year
    * json
    * binary, varbinary
    * bit
    * blob
    * enum
    * spatial data types
* Database drivers:
  * MariaDB support.
* Soft delete (managed by a SQL function).
* Postgres check.
* Mysql column value range.

## Contribute

### Run containers

You can run containers to run packages test :
````bash
docker-compose up # containers available: mysql, postgres
````

### Add tests

Update existing tests in [`migration_test.go`](./v2/migration/migration_test.go) or create a new
file. It must be validated by [`go.yml`](./.github/workflows/go.yml) workflow to validate the pull
request.

### Submitting your contribution

Create a Pull Request on a branch following the naming convention as following :
`staging/v + $RELEASE_VERSION` increment the minor version and rebase it to the superior version
starting with `staging/`.

Your commits will be reviewed and your pull request confirmed if everything is ok. If your PR
requires modifications you will be contacted to apply them and resubmit your work.

### Contributors

* Euphoria Laxis
  * Role : Maintainer, owner
  * Contact : [euphoria.laxis@euphoria-laxis.com](mailto:euphoria.laxis@euphoria-laxis.com)