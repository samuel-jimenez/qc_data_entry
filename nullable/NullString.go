package nullable

import (
	"database/sql"
	"encoding/json"
	"strings"
)

/* NullString
 *
 */
type NullString struct {
	sql.NullString
}

func NewNullString(value string) NullString {
	valid := value != ""
	return NullString{sql.NullString{String: value, Valid: valid}}
}

// Comparable
func (object NullString) Compare(b_n NullString) int {
	var a, b string
	if object.Valid {
		a = object.String
	}
	if b_n.Valid {
		b = b_n.String
	}
	return strings.Compare(a, b)
}

// JSON
func (object NullString) MarshalJSON() ([]byte, error) {
	if object.Valid {
		return json.Marshal(object.String)
	}
	return json.Marshal(nil)
}

func (object *NullString) UnmarshalJSON(data []byte) error {
	var f *string
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		object.Valid = true
		object.String = *f
	} else {
		object.Valid = false
	}
	return nil
}
