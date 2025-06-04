package main

import (
	"strings"

	"github.com/samuel-jimenez/windigo"
)

/*
 * NumberEditableView
 *
 */
type NumberEditableView interface {
	windigo.LabeledEdit
	Get() float64
	Clear()
}

/*
 * NumberEditView
 *
 */
type NumberEditView struct {
	ErrableView
	*windigo.Edit
}

func (control NumberEditView) Get() float64 {
	return parse_field(control)
}

func (control NumberEditView) GetFixed() float64 {
	start, end := control.Selected()
	// IndexAny(s, chars string) int
	control.SetText(strings.TrimSpace(control.Text()))
	// mass_field.SelectText(-1, -1)
	control.SelectText(start, end)
	// mass_field.SelectText(-1, 0)

	return parse_field(control)
}

func (control NumberEditView) Clear() {
	control.SetText("")
	control.Ok()
}

func NewNumberEditViewFromLabeledEdit(label windigo.LabeledEdit) NumberEditView {
	return NumberEditView{&View{label.ComponentFrame}, label.Edit}
}

func NewNumberEditView(parent windigo.Controller, label_width, control_width, height int, field_text string) NumberEditView {
	edit_field := NewNumberEditViewFromLabeledEdit(windigo.NewLabeledEdit(parent, label_width, control_width, height, field_text))
	edit_field.SetPaddingsAll(ERROR_MARGIN)
	return edit_field
}

func BuildNewNumberEditView(parent windigo.Controller, label_width, control_width, height int, field_text string, range_field *RangeROView) NumberEditView {

	edit_field := NewNumberEditView(parent, label_width, control_width, height, field_text)
	edit_field.OnChange().Bind(func(e *windigo.Event) {
		check_or_error(edit_field, range_field.Check(edit_field.GetFixed()))
	})
	return edit_field

}

func NewNumberEditViewWithUnits(parent windigo.AutoPanel, label_width, field_width, field_height int, field_text, field_units string) NumberEditView {

	margin := 10
	panel := windigo.NewAutoPanel(parent)
	panel.SetSize(label_width+2*field_width, field_height)
	panel.SetPaddingsAll(ERROR_MARGIN)

	text_label := windigo.NewLabel(panel)
	text_label.SetSize(label_width, field_height)
	text_label.SetText(field_text)

	text_field := windigo.NewEdit(panel)
	text_field.SetSize(field_width-margin, field_height)
	text_field.SetText("0.000")
	text_field.SetMarginRight(margin)

	text_units := windigo.NewLabel(panel)
	text_units.SetSize(field_width, field_height)
	text_units.SetText(field_units)

	panel.Dock(text_label, windigo.Left)
	panel.Dock(text_field, windigo.Left)
	panel.Dock(text_units, windigo.Left)
	parent.Dock(panel, windigo.Bottom)

	return NumberEditView{ErrableView: &View{panel}, Edit: text_field}
}
