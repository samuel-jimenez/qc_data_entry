package main

import (
	"database/sql"
	"encoding/json"
)

// Nullable Float64 that overrides sql.NullFloat64
type NullFloat64 struct {
	sql.NullFloat64
}

func NewNullFloat64(Float64 float64, Valid bool) NullFloat64 {
	return NullFloat64{sql.NullFloat64{Float64: Float64, Valid: Valid}}
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.Float64)
	}
	return json.Marshal(nil)
}

func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		nf.Valid = true
		nf.Float64 = *f
	} else {
		nf.Valid = false
	}
	return nil
}

func (nullfloat NullFloat64) Map(data_map func(float64) float64) NullFloat64 {
	var return_value NullFloat64
	if nullfloat.Valid {
		return NewNullFloat64(data_map(nullfloat.Float64), nullfloat.Valid)
	}
	return return_value
}
