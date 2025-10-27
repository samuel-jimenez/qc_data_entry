package views

import (
	"strconv"

	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/windigo"
)

var (
	err error
)

/*
 * NullFloat64Viewer
 *
 */
type NullFloat64Viewer interface {
	Get() nullable.NullFloat64
	Set(nullable.NullFloat64)
}

/*
 * NullFloat64View
 *
 */
type NullFloat64View struct {
	*NullView
	// *windigo.Edit
	format func(float64) string
}

func NullFloat64View_from_new(parent windigo.Controller, field_data nullable.NullFloat64, format func(float64) string) NullFloat64View {

	view := new(NullFloat64View)
	view.format = format
	view.NullView = NullView_from_new(parent)

	view.Set(field_data)

	return *view
}

func (view *NullFloat64View) Get() nullable.NullFloat64 {
	field_text := view.Edit.Text()
	var value float64
	valid := field_text != ""
	if valid {
		value, err = strconv.ParseFloat(field_text, 64)
		valid = err == nil
	}
	//only valid if we have some text which parses correctly
	return nullable.NewNullFloat64(value, valid)
}

func (view *NullFloat64View) Set(field_data nullable.NullFloat64) {
	var field_text string
	if field_data.Valid {
		field_text = view.format(field_data.Float64)
	}
	view.Edit.SetText(field_text)
}
