package views

import (
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/windigo"
)

/*
 * NullStringComboBox
 *
 */
type NullStringComboBox struct {
	windigo.ComponentFrame
	*windigo.ComboBox
}

func NullStringComboBox_from_new(parent windigo.Controller) *NullStringComboBox {

	view := new(NullStringComboBox)
	panel := windigo.NewAutoPanel(parent)
	view.ComponentFrame = panel
	view.ComboBox = windigo.NewComboBox(panel)
	view.ComboBox.OnKillFocus().Bind(view.ComboBox_OnKillFocus)

	panel.Dock(view.ComboBox, windigo.Top)

	return view
}

func (view *NullStringComboBox) ComboBox_OnKillFocus(e *windigo.Event) {
	view.ComboBox.SetText(strings.TrimSpace(view.ComboBox.Text()))
}

func (view *NullStringComboBox) Get() nullable.NullString {
	return nullable.NewNullString(view.ComboBox.Text())
}

func (view *NullStringComboBox) Set(field_data nullable.NullString) {
	var field_text string
	if field_data.Valid {
		field_text = field_data.String
	}
	view.ComboBox.SetText(field_text)
}
