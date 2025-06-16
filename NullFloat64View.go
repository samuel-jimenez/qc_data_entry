package main

import (
	"strconv"
	"strings"

	"github.com/samuel-jimenez/windigo"
)

/*
 * NullFloat64View
 *
 */

type NullFloat64View struct {
	*windigo.Edit
	Get func() NullFloat64
}

func NewNullFloat64View(control *windigo.Edit, get func() NullFloat64) NullFloat64View {
	return NullFloat64View{control, get}
}

func BuildNewNullFloat64View(parent windigo.Controller, field_data NullFloat64, format func(float64) string) NullFloat64View {
	edit_field := windigo.NewEdit(parent)

	if field_data.Valid {
		edit_field.SetText(format(field_data.Float64))
	}
	edit_field.OnKillFocus().Bind(func(e *windigo.Event) {
		edit_field.SetText(strings.TrimSpace(edit_field.Text()))
	})

	get := func() NullFloat64 {
		field_text := edit_field.Text()
		var value float64
		valid := field_text != ""
		if valid {
			value, err = strconv.ParseFloat(field_text, 64)
			valid = err == nil
		}
		//only valid if we have some text which parses correctly
		return NewNullFloat64(value, valid)
	}

	return NewNullFloat64View(edit_field, get)
}

/*
 * NullFloat64ROView
 *
 */

type NullFloat64ROView struct {
	*View
	Update  func(field_data NullFloat64)
	SetFont func(*windigo.Font)
	Refresh func()
}

func BuildNewNullFloat64ROView(parent windigo.Controller, field_data NullFloat64, format func(float64) string) *NullFloat64ROView {
	data_field := windigo.NewLabeledLabel(parent, "")

	update := func(field_data NullFloat64) {
		if field_data.Valid {
			data_field.SetText(format(field_data.Float64))
		} else {
			data_field.SetText("")
		}
	}
	update(field_data)

	refresh := func() {
		data_field.SetSize(RANGES_RO_FIELD_WIDTH, OFF_AXIS)
		data_field.SetPaddingsAll(ERROR_MARGIN)
		data_field.SetPaddingTop(2 * ERROR_MARGIN)
	}

	return &NullFloat64ROView{&View{data_field}, update, data_field.SetFont, refresh}
}

func BuildNullFloat64SpacerView(parent windigo.Controller, field_data NullFloat64, format string) *NullFloat64ROView {
	data_field := windigo.NewLabeledLabel(parent, "")

	update := func(field_data NullFloat64) {
		if field_data.Valid {
			data_field.SetText(format)
		} else {
			data_field.SetText("")
		}
	}
	update(field_data)

	refresh := func() {
		data_field.SetSize(RANGES_RO_SPACER_WIDTH, OFF_AXIS)
		data_field.SetPaddingsAll(ERROR_MARGIN)
		data_field.SetPaddingTop(2 * ERROR_MARGIN)
	}

	return &NullFloat64ROView{&View{data_field}, update, data_field.SetFont, refresh}
}
