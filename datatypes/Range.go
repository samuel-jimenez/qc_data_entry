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
	Min       nullable.NullFloat64
	Target    nullable.NullFloat64
	Max       nullable.NullFloat64
	Method    nullable.NullString
	Valid     bool
	Publish_p bool
}

func NewRange(
	min nullable.NullFloat64,
	target nullable.NullFloat64,
	max nullable.NullFloat64,
	Method nullable.NullString,
	Valid, Publish_p bool,
) Range {
	Valid = Valid || Publish_p || min.Valid || target.Valid || max.Valid
	return Range{min, target, max, Method, Valid, Publish_p}
}

// func (field_data Range) Check(data nullable.NullFloat64) bool {
func (field_data Range) Check(data float64) bool {
	return (!field_data.Min.Valid ||
		field_data.Min.Float64 <= data) && (!field_data.Max.Valid ||
		data <= field_data.Max.Float64)
}

func (field_data Range) Map(data_map func(float64) float64) Range {
	return NewRange(
		field_data.Min.Map(data_map),
		field_data.Target.Map(data_map),
		field_data.Max.Map(data_map),
		field_data.Method,
		field_data.Valid,
		field_data.Publish_p,
	)
}

func (field_data Range) Print(format func(float64) string, coa_p bool) string {
	ASIS := "As-is" //TODO datatypes.ASIS
	if coa_p && !field_data.Publish_p {
		return ASIS
	}
	if field_data.Min.Valid {
		if field_data.Max.Valid {
			return fmt.Sprintf("%s - %s", format(field_data.Min.Float64), format(field_data.Max.Float64))
		}
		return fmt.Sprintf("≥ %s", format(field_data.Min.Float64))
	} else if field_data.Max.Valid {
		return fmt.Sprintf("≤ %s", format(field_data.Max.Float64))
	}

	return ASIS
}

func (field_data Range) CoA(format func(float64) string) string {
	return field_data.Print(format, true)
}

func (field_data Range) QC(format func(float64) string) string {
	return field_data.Print(format, false)
}
