package migration

import (
	"fmt"
	"regexp"
	"strings"
)

/**** POSTGRES TYPE LIST ****/
/**
bigint
bigserial 	serial8
bit [ (n) ]
bit varying [ (n) ] 	varbit [ (n) ]
boolean 	bool
box
bytea
character [ (n) ] 	char [ (n) ]
character varying [ (n) ] 	varchar [ (n) ]
cidr
circle
date
double precision 	float8
inet
integer
interval [ fields ] [ (p) ]
json
jsonb
line
lseg
macaddr
macaddr8
money
numeric [ (p, s) ] 	decimal [ (p, s) ]
path
pg_lsn
pg_snapshot
point
polygon
real 	float4
smallint 	int2
smallserial 	serial2
serial 	serial4
text
time [ (p) ] [ without time zone ]
time [ (p) ] with time zone 	timetz
timestamp [ (p) ] [ without time zone ]
timestamp [ (p) ] with time zone 	timestamptz
tsquery
tsvector
txid_snapshot
uuid
xml
**/

// convertType convert go type to SQL datatype
func (m *Migrator) convertType(kind string) string {
	isTime, err := regexp.MatchString("Time$", kind)
	if err != nil {
		return ""
	}
	isNullTime, err := regexp.MatchString("NullTime$", kind)
	if err != nil {
		return ""
	}
	if isTime || isNullTime {
		switch m.Driver {
		case DBDriverPostgres:
			return "TIMETZ"
		case DBDriverMySQL:
			return "DATETIME"
		}
	}
	isUuid, err := regexp.MatchString("UUID", kind)
	if err != nil {
		return ""
	}
	if isUuid {
		switch m.Driver {
		case DBDriverPostgres:
			return "UUID"
		case DBDriverMySQL:
			return "binary(16)"
		}
	}
	isBool, err := regexp.MatchString("Bool", kind)
	if err != nil {
		return ""
	}
	if isBool {
		switch m.Driver {
		case DBDriverPostgres:
			return "BOOL"
		case DBDriverMySQL:
			return "TINYINT(1)"
		}
	}
	isYear, err := regexp.MatchString("Year", kind)
	if err != nil {
		return ""
	}
	if isYear {
		switch m.Driver {
		case DBDriverPostgres:
			return "INTERVAL"
		case DBDriverMySQL:
			return "YEAR"
		}
	}
	switch kind {
	case "int":
		return "INT"
	case "int8":
		switch m.Driver {
		case DBDriverPostgres:
			return "INT8"
		case DBDriverMySQL:
			return "TINYINT"
		}
	case "uint8":
		switch m.Driver {
		case DBDriverPostgres:
			return "INT8"
		case DBDriverMySQL:
			return "TINYINT"
		}
	case "int16":
		return "SMALLINT"
	case "uint16":
		return "SMALLINT"
	case "int32":
		switch m.Driver {
		case DBDriverPostgres:
			return "INT32"
		case DBDriverMySQL:
			return "INT"
		}
	case "uint32":
		switch m.Driver {
		case DBDriverPostgres:
			return "INT32"
		case DBDriverMySQL:
			return "INT"
		}
	case "int64":
		return "BIGINT"
	case "uint64":
		return "BIGINT"
	case "float32":
		switch m.Driver {
		case DBDriverPostgres:
			return "DOUBLE PRECISION"
		case DBDriverMySQL:
			return "FLOAT"
		}
	case "float64":
		switch m.Driver {
		case DBDriverPostgres:
			return "NUMERIC"
		case DBDriverMySQL:
			return "DOUBLE"
		}
	case "string":
		return fmt.Sprintf("VARCHAR(%d)", m.DefaultTextSize)
	case "bool":
		return "BOOL"
	case "[]byte":
		switch m.Driver {
		case DBDriverPostgres:
			return "BYTEA"
		case DBDriverMySQL:
			return "BLOB"
		}
	default:
		return ""
	}

	return ""
}

func convertPostgresSqlType(datatype string) string {
	if strings.Contains(datatype, "int") {
		return "INT"
	} else if strings.Contains(datatype, "float") {
		return "FLOAT"
	}
	switch datatype {
	case "character varying":
		return "VARCHAR"
	case "time with time zone":
		return "TIMETZ"
	case "time without time zone":
		return "TIME"
	case "text":
		return "TEXT"
	case "boolean":
		return "BOOL"
	default:
		return ""
	}
}

func convertSqlDataType(datatype string) string {
	d := strings.ToLower(datatype)
	if strings.Compare(d, "tinyint(1)") == 0 {
		return "bool"
	} else {
		return d
	}
}

func checkConstraint(constraint string) bool {
	switch constraint {
	case "not null":
		return true
	case "unique":
		return true
	default:
		return false
	}
}
