package viewer

import (
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/windigo"
)

/*
 * NullStringView
 *
 */
type NullStringView struct {
	*windigo.Edit
	Get func() nullable.NullString
}

/*
func NewNullStringView(parent windigo.Controller, field_data nullable.NullString) NullStringView {
	edit_field := windigo.NewEdit(parent)

	if field_data.Valid {
		edit_field.SetText(field_data.String)
	}
	edit_field.OnKillFocus().Bind(func(e *windigo.Event) {
		edit_field.SetText(strings.TrimSpace(edit_field.Text()))
	})

	get := func() nullable.NullString {
		var nullString nullable.NullString
		nullString.String = edit_field.Text()
		nullString.Valid = nullString.String != ""

		// value := edit_field.Text()
		// valid := value != ""

		// return NewNullString(value, valid)
		// return NullString{value, valid}
		return nullString
	}

	return NullStringView{edit_field, get}
}*/

func NewNullStringView(parent windigo.Controller) NullStringView {
	edit_field := windigo.NewEdit(parent)

	edit_field.OnKillFocus().Bind(func(e *windigo.Event) {
		edit_field.SetText(strings.TrimSpace(edit_field.Text()))
	})

	get := func() nullable.NullString {
		var nullString nullable.NullString
		nullString.String = edit_field.Text()
		nullString.Valid = nullString.String != ""

		// value := edit_field.Text()
		// valid := value != ""

		// return NewNullString(value, valid)
		// return NullString{value, valid}
		return nullString
	}

	return NullStringView{edit_field, get}
}
