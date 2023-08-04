package migration

import (
	"fmt"
	"regexp"
	"strings"
)

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
		return "DATETIME"
	}
	switch kind {
	case "int":
		return "INT"
	case "float":
		return "FLOAT"
	case "string":
		return fmt.Sprintf("VARCHAR(%d)", m.DefaultTextSize)
	case "bool":
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
