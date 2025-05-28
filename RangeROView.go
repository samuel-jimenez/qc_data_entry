package main

import "github.com/samuel-jimenez/windigo"

/*
 * RangeROView
 *
 */

type RangeROView struct {
	windigo.AutoPanel
	field_data Range
	min_field,
	min_field_spacer,
	target_field,
	max_field_spacer,
	max_field NullFloat64ROView
	data_map func(float64) float64
}

func (data_view *RangeROView) Update(update_data Range) {
	if data_view.data_map != nil {
		update_data = update_data.Map(data_view.data_map)
	}

	data_view.field_data = update_data
	data_view.min_field.Update(update_data.Min)
	data_view.min_field_spacer.Update(update_data.Min)
	data_view.target_field.Update(update_data.Target)
	data_view.max_field.Update(update_data.Max)
	data_view.max_field_spacer.Update(update_data.Max)
}

func (data_view RangeROView) Clear() {
	data_view.min_field.Ok()
	data_view.max_field.Ok()
}

func (data_view RangeROView) Check(data float64) bool {
	data_view.Clear()
	if data_view.field_data.Min.Valid &&
		data_view.field_data.Min.Float64 >= data {
		data_view.min_field.Error()
	} else if data_view.field_data.Max.Valid &&
		data_view.field_data.Max.Float64 <= data {
		data_view.max_field.Error()
	}

	return data_view.field_data.Check(data)
}

func BuildNewRangeROViewMap(parent windigo.Controller, field_text string, field_data Range, format func(float64) string, data_map func(float64) float64) RangeROView {

	panel := windigo.NewAutoPanel(parent)
	panel.SetSize(22, 22)
	panel.SetMarginTop(5)
	//TODO toolti[p]
	// label := windigo.NewLabel(panel)
	// label.SetText(field_text)

	spacer_format := "≤"
	min_field := BuildNewNullFloat64ROView(panel, field_data.Min, format)
	min_field_spacer := BuildNullFloat64SpacerView(panel, field_data.Min, spacer_format)
	target_field := BuildNewNullFloat64ROView(panel, field_data.Target, format)
	max_field := BuildNewNullFloat64ROView(panel, field_data.Max, format)
	max_field_spacer := BuildNullFloat64SpacerView(panel, field_data.Max, spacer_format)

	// panel.Dock(label, windigo.Left)
	panel.Dock(min_field, windigo.Left)
	panel.Dock(min_field_spacer, windigo.Left)
	panel.Dock(target_field, windigo.Left)
	panel.Dock(max_field_spacer, windigo.Left)
	panel.Dock(max_field, windigo.Left)

	return RangeROView{panel,
		field_data,
		min_field,
		min_field_spacer,
		target_field,
		max_field_spacer,
		max_field, data_map}
}

func BuildNewRangeROView(parent windigo.Controller, field_text string, field_data Range, format func(float64) string) RangeROView {
	return BuildNewRangeROViewMap(parent, field_text, field_data, format, nil)
}
