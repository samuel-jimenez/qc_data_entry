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
func (a_n NullString) Compare(b_n NullString) int {
	var a, b string
	if a_n.Valid {
		a = a_n.String
	}
	if b_n.Valid {
		b = b_n.String
	}
	return strings.Compare(a, b)
}

// JSON
func (nf NullString) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.String)
	}
	return json.Marshal(nil)
}

func (nf *NullString) UnmarshalJSON(data []byte) error {
	var f *string
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		nf.Valid = true
		nf.String = *f
	} else {
		nf.Valid = false
	}
	return nil
}
