package viewer

import (
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/windigo"
)

/*
 * NullStringViewer
 *
 */
type NullStringViewer interface {
	windigo.Editable
	Get() nullable.NullString
	Clear()
}

/*
 * NullStringView
 *
 */
type NullStringView struct {
	*windigo.Edit
}

func NewNullStringView(parent windigo.Controller) NullStringView {
	edit_field := windigo.NewEdit(parent)

	edit_field.OnKillFocus().Bind(func(e *windigo.Event) {
		edit_field.SetText(strings.TrimSpace(edit_field.Text()))
	})

	return NullStringView{edit_field}
}

func (control NullStringView) Get() nullable.NullString {
	var nullString nullable.NullString
	nullString.String = control.Text()
	nullString.Valid = nullString.String != ""
	return nullString
}

func (control NullStringView) Clear() {
	control.SetText("")
}
