# go-db-migration

Go package to generate and execute migration schemas using structure tags

## Documentation

### Installation

To add package to your go mod run :
````bash
go get github.com/euphoria-laxis/go-db-migration@v2.1.0
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

|    Driver    |     Available      |     Availability status      |
|:------------:|:------------------:|:----------------------------:|
|  **MySQL**   | :white_check_mark: |          Available           |
| **Postgres** | :white_check_mark: |          Available           |
|  **SQLite**  |   :construction:   |       Work In Progress       |
| **MariaDB**  |     :warning:      | Use MySQL driver for MariaDB |

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

## License

This project is under [MIT license](./LICENSE).
