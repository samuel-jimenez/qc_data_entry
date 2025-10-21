package views

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
	label *windigo.Label
	min_field,
	target_field,
	max_field NullFloat64View
	valid_field,
	export_field *windigo.CheckBox
}

func RangeView_from_new(parent windigo.Controller, field_text string, field_data datatypes.Range, format func(float64) string) *RangeView {

	view := new(RangeView)
	view.AutoPanel = windigo.NewAutoPanel(parent)
	// view.AutoPanel.SetMarginsAll(10)
	view.AutoPanel.SetPaddingsAll(GUI.RANGES_PADDING)
	view.label = windigo.NewLabel(view.AutoPanel)
	view.label.SetText(field_text)

	view.min_field = NullFloat64View_from_new(view.AutoPanel, field_data.Min, format)
	view.target_field = NullFloat64View_from_new(view.AutoPanel, field_data.Target, format)
	view.max_field = NullFloat64View_from_new(view.AutoPanel, field_data.Max, format)

	view.valid_field = windigo.NewCheckBox(view)
	view.valid_field.SetMarginsAll(GUI.ERROR_MARGIN + GUI.RANGES_PADDING)
	view.valid_field.SetText("")
	view.valid_field.SetChecked(field_data.Valid)

	view.export_field = windigo.NewCheckBox(view)
	view.export_field.SetMarginsAll(GUI.ERROR_MARGIN)
	view.export_field.SetText("")
	view.export_field.SetChecked(field_data.Publish_p)

	view.AutoPanel.Dock(view.label, windigo.Left)
	view.AutoPanel.Dock(view.min_field, windigo.Left)
	view.AutoPanel.Dock(view.target_field, windigo.Left)
	view.AutoPanel.Dock(view.max_field, windigo.Left)
	view.AutoPanel.Dock(view.valid_field, windigo.Left)
	view.AutoPanel.Dock(view.export_field, windigo.Left)

	return view
}

func (view RangeView) Get() datatypes.Range {
	return datatypes.NewRange(
		view.min_field.Get(),
		view.target_field.Get(),
		view.max_field.Get(),
		view.valid_field.Checked(),
		view.export_field.Checked(),
	)
}

func (view RangeView) Set(field_data datatypes.Range) {
	view.min_field.Set(field_data.Min)
	view.target_field.Set(field_data.Target)
	view.max_field.Set(field_data.Max)
	view.valid_field.SetChecked(field_data.Valid)
	view.export_field.SetChecked(field_data.Publish_p)
}

func (view RangeView) SetLabeledSize(label_width, control_width, height int) {
	view.AutoPanel.SetSize(control_width, height)
	view.label.SetSize(label_width, height)
	view.min_field.SetSize(control_width, height)
	view.target_field.SetSize(control_width, height)
	view.max_field.SetSize(control_width, height)
	view.valid_field.SetSize(control_width, height)
	view.valid_field.SetMarginLeft(control_width / 2)
	view.export_field.SetSize(control_width, height)
}

// func (control RangeView) SetDockSize(width, height int) {
// 	control.SetSize(width, height)
// 	for _, label := range control.labels {
// 		label.SetSize(width, height)
// 	}
// }
