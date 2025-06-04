package main

import (
	"strconv"
	"strings"

	"github.com/samuel-jimenez/windigo"
)

/*
 * NumberEditViewable
 *
 */
type NumberEditViewable interface {
	windigo.ComponentFrame
	windigo.Controller
	Get() float64
	GetFixed() float64
	Clear()
	Check(bool)
}

/*
 * NumberEditView
 *
 */
type NumberEditView struct {
	ErrableView
	*windigo.Edit
}

func (control *NumberEditView) Get() float64 {
	val, _ := strconv.ParseFloat(strings.TrimSpace(control.Text()), 64)
	return val
}

func (control *NumberEditView) GetFixed() float64 {
	start, end := control.Selected()
	// IndexAny(s, chars string) int
	control.SetText(strings.TrimSpace(control.Text()))
	// mass_field.SelectText(-1, -1)
	control.SelectText(start, end)
	// mass_field.SelectText(-1, 0)

	return control.Get()
}

func (control *NumberEditView) Clear() {
	control.SetText("")
	control.Ok()
}

func (control *NumberEditView) Check(test bool) {
	if test {
		control.Ok()
	} else {
		control.Error()
	}
}

func NewNumberEditViewFromLabeledEdit(label windigo.LabeledEdit) *NumberEditView {
	return &NumberEditView{&View{label.ComponentFrame}, label.Edit}
}

func NewNumberEditView(parent windigo.Controller, label_width, control_width, height int, field_text string) *NumberEditView {
	edit_field := NewNumberEditViewFromLabeledEdit(windigo.NewLabeledEdit(parent, label_width, control_width, height, field_text))
	edit_field.SetPaddingsAll(ERROR_MARGIN)
	return edit_field
}

func NewNumberEditViewWithChange(parent windigo.Controller, label_width, control_width, height int, field_text string, range_field *RangeROView) *NumberEditView {

	edit_field := NewNumberEditView(parent, label_width, control_width, height, field_text)
	edit_field.OnChange().Bind(func(e *windigo.Event) {
		edit_field.Check(range_field.Check(edit_field.GetFixed()))
	})
	return edit_field
}

func NewNumberEditViewWithUnits(parent windigo.AutoPanel, label_width, field_width, unit_width, height int, field_text, field_units string) *NumberEditView {

	margin := 10
	field_width = field_width - margin
	unit_width = unit_width - margin

	panel := windigo.NewAutoPanel(parent)
	panel.SetSize(label_width+field_width+unit_width+margin, height)
	panel.SetPaddingsAll(ERROR_MARGIN)

	text_label := windigo.NewLabel(panel)
	text_label.SetSize(label_width, height)
	text_label.SetText(field_text)
	text_label.SetMarginTop(ERROR_MARGIN)

	text_field := windigo.NewEdit(panel)
	text_field.SetSize(field_width, height)
	text_field.SetText("0.000")
	text_field.SetMarginRight(margin)

	text_units := windigo.NewLabel(panel)
	text_units.SetSize(unit_width, height)
	text_units.SetText(field_units)
	text_units.SetMarginTop(ERROR_MARGIN)

	panel.Dock(text_label, windigo.Left)
	panel.Dock(text_field, windigo.Left)
	panel.Dock(text_units, windigo.Left)
	parent.Dock(panel, windigo.Bottom)

	return &NumberEditView{ErrableView: &View{panel}, Edit: text_field}
}
