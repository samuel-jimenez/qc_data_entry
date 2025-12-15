package qc_ui

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/windigo"
)

/*
 * NullFloat64ROView
 *
 */

type NullFloat64ROView struct {
	*GUI.View
	Update  func(field_data nullable.NullFloat64)
	SetFont func(*windigo.Font)
	Refresh func()
}

func NullFloat64ROView_from_new(parent windigo.Controller, field_data nullable.NullFloat64, format func(float64) string) *NullFloat64ROView {
	data_field := windigo.LabeledLabel_from_new(parent, "")

	update := func(field_data nullable.NullFloat64) {
		data_field.SetText(field_data.Format(format))
	}
	update(field_data)

	refresh := func() {
		data_field.SetSize(GUI.RANGES_RO_FIELD_WIDTH, GUI.OFF_AXIS)
		data_field.SetPaddingsAll(GUI.ERROR_MARGIN)
		data_field.SetPaddingTop(2 * GUI.ERROR_MARGIN)
	}

	return &NullFloat64ROView{&GUI.View{ComponentFrame: data_field}, update, data_field.SetFont, refresh}
}

func NullFloat64View_Spacer_from_new(parent windigo.Controller, field_data nullable.NullFloat64, format string) *NullFloat64ROView {
	data_field := windigo.LabeledLabel_from_new(parent, "")

	update := func(field_data nullable.NullFloat64) {
		if field_data.Valid {
			data_field.SetText(format)
		} else {
			data_field.SetText("")
		}
	}
	update(field_data)

	refresh := func() {
		data_field.SetSize(GUI.RANGES_RO_SPACER_WIDTH, GUI.OFF_AXIS)
		data_field.SetPaddingsAll(GUI.ERROR_MARGIN)
		data_field.SetPaddingTop(2 * GUI.ERROR_MARGIN)
	}

	return &NullFloat64ROView{&GUI.View{ComponentFrame: data_field}, update, data_field.SetFont, refresh}
}
