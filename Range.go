package main

import (
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/windigo"
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

/*
 * RangeView
 *
 */

type RangeView struct {
	*windigo.AutoPanel
	Get            func() Range
	SetLabeledSize func(label_width, control_width, height int)
}

// func (control RangeView) SetDockSize(width, height int) {
// 	control.SetSize(width, height)
// 	for _, label := range control.labels {
// 		label.SetSize(width, height)
// 	}
// }

func BuildNewRangeView(parent windigo.Controller, field_text string, field_data Range, format func(float64) string) RangeView {

	panel := windigo.NewAutoPanel(parent)
	// panel.SetMarginsAll(10)
	panel.SetPaddingsAll(RANGES_PADDING)
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

	setLabeledSize := func(label_width, control_width, height int) {
		panel.SetSize(control_width, height)
		label.SetSize(label_width, height)
		min_field.SetSize(control_width, height)
		target_field.SetSize(control_width, height)
		max_field.SetSize(control_width, height)
	}

	return RangeView{panel, get, setLabeledSize}
}
