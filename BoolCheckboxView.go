package main

import (
	"github.com/samuel-jimenez/windigo"
)

/*
 * BoolCheckboxViewable
 *
 */
type BoolCheckboxViewable interface {
	ErrableView
	windigo.Controller
	Get() bool
	Clear()
}

/*
 * BoolCheckboxView
 *
 */
type BoolCheckboxView struct {
	ErrableView
	*windigo.CheckBox
}

func (control *BoolCheckboxView) Get() bool {
	return control.Checked()
}

func (control *BoolCheckboxView) Clear() {
	control.SetChecked(false)
}

func NewBoolCheckboxViewFromLabeledCheckBox(label windigo.LabeledCheckBox) *BoolCheckboxView {
	return &BoolCheckboxView{&View{label.ComponentFrame}, label.CheckBox}
}

func NewBoolCheckboxView(parent windigo.Controller, width, height int, field_text string) *BoolCheckboxView {
	edit_field := NewBoolCheckboxViewFromLabeledCheckBox(windigo.NewLabeledCheckBox(parent, width, height, field_text))
	edit_field.SetPaddingsAll(ERROR_MARGIN)
	return edit_field
}
