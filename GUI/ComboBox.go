package GUI

import (
	"strings"

	"github.com/samuel-jimenez/windigo"
)

/* ComboBoxable
 *
 */
type ComboBoxable interface {
	windigo.ComboBoxable
	ErrableView
	Alert()
}

type ComboBox struct {
	*windigo.LabeledComboBox
}

func ComboBox_from_new(parent windigo.Controller, field_text string) *ComboBox {
	combobox_field := &ComboBox{windigo.NewLabeledComboBox(parent, field_text)}
	combobox_field.OnKillFocus().Bind(func(e *windigo.Event) {
		combobox_field.SetText(strings.ToUpper(strings.TrimSpace(combobox_field.Text())))
	})

	return combobox_field
}

func List_ComboBox_from_new(parent windigo.Controller, field_text string) *ComboBox {
	combobox_field := &ComboBox{windigo.NewLabeledListComboBox(parent, field_text)}
	combobox_field.OnKillFocus().Bind(func(e *windigo.Event) {
		combobox_field.SetText(strings.ToUpper(strings.TrimSpace(combobox_field.Text())))
	})

	return combobox_field
}

func (control *ComboBox) Ok() {
	// control.SetBorder(nil)
	control.ClearFGColor()
	control.ClearBGColor()
}

func (control *ComboBox) Error() {
	// control.SetBorder(ErroredPen)
	control.SetFGColor(windigo.RGB(255, 255, 255))
	control.SetBGColor(windigo.RGB(255, 0, 0))
}

func (control *ComboBox) Alert() {
	control.Error()
	control.ShowDropdown(true)
	control.SetFocus()
}
