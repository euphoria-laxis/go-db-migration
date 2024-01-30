package types

import (
	"database/sql/driver"
	"strconv"
)

type Year int64

// String format Year to string
func (y *Year) String() string {
	r := int64(*y)

	return strconv.Itoa(int(r))
}

// MarshalJSON parse Year to JSON
func (y *Year) MarshalJSON() ([]byte, error) {
	r := int64(*y)

	return []byte(strconv.FormatInt(r, 64)), nil
}

// UnmarshalJSON parse JSON to Year
func (y *Year) UnmarshalJSON(b []byte) error {
	i, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*y = Year(int64(i))

	return nil
}

// Scan data received from DB
func (y *Year) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	i, err := strconv.Atoi(string(bytes))
	if err != nil {
		return err
	}
	*y = Year(int64(i))

	return nil
}

// Value data are saved into DB
func (y *Year) Value() (driver.Value, error) {
	return int64(*y), nil
}
