package main

import "github.com/samuel-jimenez/windigo"

/*
 * Range
 *
 */

type Range struct {
	Min    NullFloat64
	Target NullFloat64
	Max    NullFloat64
}

func NewRange(
	min NullFloat64,
	target NullFloat64,
	max NullFloat64,
) Range {
	return Range{min, target, max}
}

// func (field_data Range) Check(data NullFloat64) bool {
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

/*
 * RangeView
 *
 */

type RangeView struct {
	*windigo.AutoPanel
	Get func() Range
}

func BuildNewRangeView(parent windigo.Controller, field_text string, field_data Range, format func(float64) string) RangeView {

	panel := windigo.NewAutoPanel(parent)
	panel.SetSize(OFF_AXIS, RANGES_FIELD_HEIGHT)
	// panel.SetMarginsAll(10)
	panel.SetPaddingsAll(10)
	label := windigo.NewLabel(panel)
	label.SetText(field_text)

	min_field := BuildNewNullFloat64View(panel, field_data.Min, format)
	target_field := BuildNewNullFloat64View(panel, field_data.Target, format)
	max_field := BuildNewNullFloat64View(panel, field_data.Max, format)

	panel.Dock(label, windigo.Left)
	panel.Dock(min_field, windigo.Left)
	panel.Dock(target_field, windigo.Left)

	panel.Dock(max_field, windigo.Left)

	get := func() Range {
		return NewRange(
			min_field.Get(),
			target_field.Get(),
			max_field.Get(),
		)
	}

	return RangeView{panel, get}
}
