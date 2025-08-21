package views

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/*
 * BoolCheckboxViewable
 *
 */
type BoolCheckboxViewable interface {
	GUI.ErrableView
	windigo.Buttonable
	Get() bool
	Clear()
}

/*
 * BoolCheckboxView
 *
 */
type BoolCheckboxView struct {
	GUI.ErrableView
	*windigo.CheckBox
}

func NewBoolCheckboxViewFromLabeledCheckBox(label *windigo.LabeledCheckBox) *BoolCheckboxView {
	return &BoolCheckboxView{&GUI.View{ComponentFrame: label.ComponentFrame}, label.CheckBox}
}

func NewBoolCheckboxView(parent windigo.Controller, field_text string) *BoolCheckboxView {
	edit_field := NewBoolCheckboxViewFromLabeledCheckBox(windigo.NewLabeledCheckBox(parent, field_text))
	edit_field.SetPaddingsAll(GUI.ERROR_MARGIN)
	return edit_field
}

func (control *BoolCheckboxView) Get() bool {
	return control.Checked()
}

func (control *BoolCheckboxView) Clear() {
	control.SetChecked(false)
}
