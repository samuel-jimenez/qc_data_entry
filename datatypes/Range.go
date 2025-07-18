package datatypes

import (
	"fmt"

	"github.com/samuel-jimenez/qc_data_entry/nullable"
)

/*
 * Range
 *
 */

type Range struct {
	Min    nullable.NullFloat64
	Target nullable.NullFloat64
	Max    nullable.NullFloat64
}

func NewRange(
	min nullable.NullFloat64,
	target nullable.NullFloat64,
	max nullable.NullFloat64,
) Range {
	return Range{min, target, max}
}

// func (field_data Range) Check(data nullable.NullFloat64) bool {
func (field_data Range) Check(data float64) bool {
	return (!field_data.Min.Valid ||
		field_data.Min.Float64 <= data) && (!field_data.Max.Valid ||
		data <= field_data.Max.Float64)
}

func (field_data Range) Map(data_map func(float64) float64) Range {
	return Range{field_data.Min.Map(data_map),
		field_data.Target.Map(data_map),
		field_data.Max.Map(data_map)}
}

func (field_data Range) CoA(format func(float64) string) string {
	if field_data.Min.Valid {
		if field_data.Max.Valid {
			return fmt.Sprintf("%s - %s", format(field_data.Min.Float64), format(field_data.Max.Float64))
		}
		return fmt.Sprintf("≥ %s", format(field_data.Min.Float64))
	} else if field_data.Max.Valid {
		return fmt.Sprintf("≤ %s", format(field_data.Max.Float64))
	}

	return "As-is"
}
