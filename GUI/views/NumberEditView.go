package views

// TODO todo NumberEditView

import (
	"strconv"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/*
 * NumberEditViewable
 *
 */
type NumberEditViewable interface {
	GUI.ErrableView
	windigo.Editable
	windigo.DiffLabelable
	Get() float64
	GetFixed() float64
	Set(val float64)
	Clear()
	Check(bool)
}

// TODO combine NumbEditView nmaybe?

/*
 * NumberEditView
 *
 */
type NumberEditView struct {
	GUI.ErrableView
	*windigo.Edit
	*windigo.Labeled
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
func (control *NumberEditView) Set(val float64) {
	start, end := control.Selected()
	control.SetText(strconv.FormatFloat(val, 'f', 2, 64))
	control.SelectText(start, end)
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

func (control *NumberEditView) SetFont(font *windigo.Font) {
	control.Edit.SetFont(font)
	control.Label().SetFont(font)
}

func (control *NumberEditView) SetLabeledSize(label_width, control_width, height int) {
	control.SetSize(label_width+control_width, height)
	// control.Edit.SetSize(control_width, height)
	control.Label().SetSize(label_width, height)
	control.SetPaddingsAll(GUI.ERROR_MARGIN)

}

func NewNumberEditViewFromLabeledEdit(label *windigo.LabeledEdit) *NumberEditView {
	return &NumberEditView{&GUI.View{ComponentFrame: label.ComponentFrame}, label.Edit, &windigo.Labeled{FieldLabel: label.Label()}}
}

func NewNumberEditView(parent windigo.Controller, field_text string) *NumberEditView {
	edit_field := NewNumberEditViewFromLabeledEdit(windigo.NewLabeledEdit(parent, field_text))
	return edit_field
}

func NewNumberEditViewWithChange(parent windigo.Controller, field_text string, range_field *RangeROView) *NumberEditView {

	edit_field := NewNumberEditView(parent, field_text)
	edit_field.OnChange().Bind(func(e *windigo.Event) {
		edit_field.Check(range_field.Check(edit_field.GetFixed()))
	})
	return edit_field
}

/*
 * NumberUnitsEditView
 *
 */
type NumberUnitsEditView struct {
	*NumberEditView
	SetFont        func(font *windigo.Font)
	SetLabeledSize func(label_width, control_width, unit_width, height int)
}

//	func (control *NumberUnitsEditView) SetLabeledSize(label_width, control_width, height int) {
//		control.NumberEditView.SetLabeledSize(label_width,control_width, height)
//		control.SetSize(label_width+control_width, height)
//		control.Label().SetSize(label_width, height)
//	}
func NewNumberEditViewWithUnits(parent *windigo.AutoPanel, field_text, field_units string) *NumberUnitsEditView {

	panel := windigo.NewAutoPanel(parent)

	text_label := windigo.NewLabel(panel)
	text_label.SetText(field_text)

	text_field := windigo.NewEdit(panel)
	text_field.SetText("0.000")

	text_units := windigo.NewLabel(panel)
	text_units.SetText(field_units)

	panel.Dock(text_label, windigo.Left)
	panel.Dock(text_field, windigo.Left)
	panel.Dock(text_units, windigo.Left)

	setFont := func(font *windigo.Font) {
		text_label.SetFont(font)
		text_field.SetFont(font)
		text_units.SetFont(font)
	}
	setLabeledSize := func(label_width, control_width, unit_width, height int) {

		panel.SetSize(label_width+control_width+unit_width, height)
		panel.SetPaddingsAll(GUI.ERROR_MARGIN)

		text_label.SetSize(label_width, height)
		text_label.SetMarginTop(GUI.ERROR_MARGIN)

		text_field.SetSize(control_width, height)
		text_field.SetMarginRight(GUI.DATA_MARGIN)

		text_units.SetSize(unit_width, height)
		text_units.SetMarginTop(GUI.ERROR_MARGIN)

	}

	return &NumberUnitsEditView{&NumberEditView{&GUI.View{ComponentFrame: panel}, text_field, &windigo.Labeled{FieldLabel: text_label}}, setFont, setLabeledSize}
}
