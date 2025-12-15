package nullable

import (
	"cmp"
	"database/sql"
	"encoding/json"
)

/* NullFloat64
 * Nullable Float64 that overrides sql.NullFloat64
 *
 */
type NullFloat64 struct {
	sql.NullFloat64
}

func NewNullFloat64(Float64 float64, Valid bool) NullFloat64 {
	return NullFloat64{sql.NullFloat64{Float64: Float64, Valid: Valid}}
}

// Comparable
func (object NullFloat64) Compare(b_n NullFloat64) int {
	var a, b float64
	if object.Valid {
		a = object.Float64
	}
	if b_n.Valid {
		b = b_n.Float64
	}
	return cmp.Compare(a, b)
}

// JSON
func (object NullFloat64) MarshalJSON() ([]byte, error) {
	if object.Valid {
		return json.Marshal(object.Float64)
	}
	return json.Marshal(nil)
}

func (object *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		object.Valid = true
		object.Float64 = *f
	} else {
		object.Valid = false
	}
	return nil
}

func (object NullFloat64) Map(data_map func(float64) float64) (output NullFloat64) {
	if object.Valid {
		output = NewNullFloat64(data_map(object.Float64), object.Valid)
	}
	return
}

func (object NullFloat64) Format(format func(float64) string) (output string) {
	if object.Valid {
		output = format(object.Float64)
	}
	return
}
