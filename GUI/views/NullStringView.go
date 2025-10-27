package views

import (
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/windigo"
)

/*
 * NullStringView
 *
 */
type NullStringView struct {
	*NullView
	// *windigo.Edit
}

func NullStringView_from_new(parent windigo.Controller, field_data nullable.NullString) *NullStringView {

	view := new(NullStringView)
	view.NullView = NullView_from_new(parent)

	view.Set(field_data)

	return view
}

func (view *NullStringView) Get() nullable.NullString {
	return nullable.NewNullString(view.Edit.Text())
}

func (view *NullStringView) Set(field_data nullable.NullString) {
	var field_text string
	if field_data.Valid {
		field_text = field_data.String
	}
	view.Edit.SetText(field_text)
}
