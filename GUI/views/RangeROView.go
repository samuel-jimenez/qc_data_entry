package views

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/datatypes"
	"github.com/samuel-jimenez/windigo"
)

/*
 * RangeROView
 *
 */

type RangeROView struct {
	*windigo.AutoPanel
	field_data datatypes.Range
	min_field,
	min_field_spacer,
	target_field,
	max_field_spacer,
	max_field *NullFloat64ROView
	data_map func(float64) float64
}

func RangeROViewMap_from_new(parent windigo.Controller, field_text string, field_data datatypes.Range, format func(float64) string, data_map func(float64) float64) *RangeROView {

	panel := windigo.NewAutoPanel(parent)
	//TODO toolti[p]
	// label := windigo.NewLabel(panel)
	// label.SetText(field_text)

	spacer_format := "≤"
	min_field := NullFloat64ROView_from_new(panel, field_data.Min, format)
	min_field_spacer := NullFloat64View_Spacer_from_new(panel, field_data.Min, spacer_format)
	target_field := NullFloat64ROView_from_new(panel, field_data.Target, format)
	max_field := NullFloat64ROView_from_new(panel, field_data.Max, format)
	max_field_spacer := NullFloat64View_Spacer_from_new(panel, field_data.Max, spacer_format)

	// panel.Dock(label, windigo.Left)
	panel.Dock(min_field, windigo.Left)
	panel.Dock(min_field_spacer, windigo.Left)
	panel.Dock(target_field, windigo.Left)
	panel.Dock(max_field_spacer, windigo.Left)
	panel.Dock(max_field, windigo.Left)

	return &RangeROView{panel,
		field_data,
		min_field,
		min_field_spacer,
		target_field,
		max_field_spacer,
		max_field, data_map}
}

func RangeROView_from_new(parent windigo.Controller, field_text string, field_data datatypes.Range, format func(float64) string) *RangeROView {
	return RangeROViewMap_from_new(parent, field_text, field_data, format, nil)
}

func (data_view *RangeROView) SetFont(font *windigo.Font) {
	data_view.min_field.SetFont(font)
	data_view.min_field_spacer.SetFont(font)
	data_view.target_field.SetFont(font)
	data_view.max_field_spacer.SetFont(font)
	data_view.max_field.SetFont(font)
}

func (data_view *RangeROView) RefreshSize() {
	data_view.SetSize(GUI.OFF_AXIS, GUI.RANGES_RO_FIELD_HEIGHT)
	data_view.min_field.Refresh()
	data_view.min_field_spacer.Refresh()
	data_view.target_field.Refresh()
	data_view.max_field_spacer.Refresh()
	data_view.max_field.Refresh()
}

// TODO As-is here?
func (data_view *RangeROView) Update(update_data datatypes.Range) {
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

func (data_view *RangeROView) Clear() {
	data_view.min_field.Ok()
	data_view.max_field.Ok()
}

func (data_view *RangeROView) Check(data float64) bool {
	data_view.Clear()
	if data_view.field_data.Min.Valid &&
		data_view.field_data.Min.Float64 > data {
		data_view.min_field.Error()
	} else if data_view.field_data.Max.Valid &&
		data_view.field_data.Max.Float64 < data {
		data_view.max_field.Error()
	}

	return data_view.field_data.Check(data)
}

func (data_view *RangeROView) CheckAll(data ...float64) (result []bool) {
	data_view.Clear()
	for _, val := range data {
		if data_view.field_data.Min.Valid &&
			data_view.field_data.Min.Float64 > val {
			data_view.min_field.Error()
		} else if data_view.field_data.Max.Valid &&
			data_view.field_data.Max.Float64 < val {
			data_view.max_field.Error()
		}
		result = append(result, data_view.field_data.Check(val))
	}

	return result
}
