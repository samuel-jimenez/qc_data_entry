package view

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/datatypes"
	"github.com/samuel-jimenez/windigo"
)

/*
 * RangeView
 *
 */

type RangeView struct {
	*windigo.AutoPanel
	Get            func() datatypes.Range
	Set            func(datatypes.Range)
	SetLabeledSize func(label_width, control_width, height int)
}

// func (control RangeView) SetDockSize(width, height int) {
// 	control.SetSize(width, height)
// 	for _, label := range control.labels {
// 		label.SetSize(width, height)
// 	}
// }

func BuildNewRangeView(parent windigo.Controller, field_text string, field_data datatypes.Range, format func(float64) string) RangeView {

	panel := windigo.NewAutoPanel(parent)
	// panel.SetMarginsAll(10)
	panel.SetPaddingsAll(GUI.RANGES_PADDING)
	label := windigo.NewLabel(panel)
	label.SetText(field_text)

	min_field := BuildNewNullFloat64View(panel, field_data.Min, format)
	target_field := BuildNewNullFloat64View(panel, field_data.Target, format)
	max_field := BuildNewNullFloat64View(panel, field_data.Max, format)

	panel.Dock(label, windigo.Left)
	panel.Dock(min_field, windigo.Left)
	panel.Dock(target_field, windigo.Left)

	panel.Dock(max_field, windigo.Left)

	get := func() datatypes.Range {
		return datatypes.NewRange(
			min_field.Get(),
			target_field.Get(),
			max_field.Get(),
		)
	}

	set := func(field_data datatypes.Range) {
		min_field.Set(field_data.Min)
		target_field.Set(field_data.Target)
		max_field.Set(field_data.Max)
	}

	setLabeledSize := func(label_width, control_width, height int) {
		panel.SetSize(control_width, height)
		label.SetSize(label_width, height)
		min_field.SetSize(control_width, height)
		target_field.SetSize(control_width, height)
		max_field.SetSize(control_width, height)
	}

	return RangeView{panel, get, set, setLabeledSize}
}
