package views

import (
	"strconv"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/windigo"
)

var (
	err error
)

/*
 * NullFloat64View
 *
 */

type NullFloat64View struct {
	*windigo.Edit
	Get func() nullable.NullFloat64
	Set func(nullable.NullFloat64)
}

func BuildNewNullFloat64View(parent windigo.Controller, field_data nullable.NullFloat64, format func(float64) string) NullFloat64View {
	edit_field := windigo.NewEdit(parent)

	if field_data.Valid {
		edit_field.SetText(format(field_data.Float64))
	}
	edit_field.OnKillFocus().Bind(func(e *windigo.Event) {
		edit_field.SetText(strings.TrimSpace(edit_field.Text()))
	})

	get := func() nullable.NullFloat64 {
		field_text := edit_field.Text()
		var value float64
		valid := field_text != ""
		if valid {
			value, err = strconv.ParseFloat(field_text, 64)
			valid = err == nil
		}
		//only valid if we have some text which parses correctly
		return nullable.NewNullFloat64(value, valid)
	}
	set := func(field_data nullable.NullFloat64) {
		field_text := ""
		if field_data.Valid {
			field_text = format(field_data.Float64)
		}
		edit_field.SetText(field_text)
	}

	return NullFloat64View{edit_field, get, set}
}

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

func BuildNewNullFloat64ROView(parent windigo.Controller, field_data nullable.NullFloat64, format func(float64) string) *NullFloat64ROView {
	data_field := windigo.NewLabeledLabel(parent, "")

	update := func(field_data nullable.NullFloat64) {
		if field_data.Valid {
			data_field.SetText(format(field_data.Float64))
		} else {
			data_field.SetText("")
		}
	}
	update(field_data)

	refresh := func() {
		data_field.SetSize(GUI.RANGES_RO_FIELD_WIDTH, GUI.OFF_AXIS)
		data_field.SetPaddingsAll(GUI.ERROR_MARGIN)
		data_field.SetPaddingTop(2 * GUI.ERROR_MARGIN)
	}

	return &NullFloat64ROView{&GUI.View{ComponentFrame: data_field}, update, data_field.SetFont, refresh}
}

func BuildNullFloat64SpacerView(parent windigo.Controller, field_data nullable.NullFloat64, format string) *NullFloat64ROView {
	data_field := windigo.NewLabeledLabel(parent, "")

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
