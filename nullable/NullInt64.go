package nullable

import (
	"cmp"
	"database/sql"
	"encoding/json"
	"strconv"
)

/* NullInt64
 * Nullable Int64 that overrides sql.NullInt64
 *
 */
type NullInt64 struct {
	sql.NullInt64
}

func NullInt64Default() NullInt64 {
	return NullInt64{}
}

func NewNullInt64(val int64) NullInt64 {
	return NullInt64{sql.NullInt64{Int64: val, Valid: true}}
}

func (object NullInt64) Diff(b_n NullInt64) int64 {
	var a, b int64
	if object.Valid {
		a = object.Int64
	}
	if b_n.Valid {
		b = b_n.Int64
	}
	return a - b
}

// Comparable
func (object NullInt64) Compare(b_n NullInt64) int {
	var a, b int64
	if object.Valid {
		a = object.Int64
	}
	if b_n.Valid {
		b = b_n.Int64
	}
	return cmp.Compare(a, b)
}

func (object NullInt64) String() (output string) {
	if object.Valid {
		output = strconv.FormatInt(object.Int64, 10)
	}
	return
}

// JSON
func (object NullInt64) MarshalJSON() ([]byte, error) {
	if object.Valid {
		return json.Marshal(object.Int64)
	}
	return json.Marshal(nil)
}

func (object *NullInt64) UnmarshalJSON(data []byte) error {
	var f *int64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		object.Valid = true
		object.Int64 = *f
	} else {
		object.Valid = false
	}
	return nil
}
