package main

import "github.com/samuel-jimenez/windigo"

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
	windigo.LabeledEdit
}

func (control NumberEditView) Get() float64 {
	return parse_field(control.LabeledEdit)
}
func (control NumberEditView) GetFixed() float64 {
	return edit_and_parse_field(control.LabeledEdit)
}

func (control NumberEditView) Clear() {
	control.SetText("")
	control.SetBorder(okPen)
}

func BuildNewNumberEditView(parent windigo.Controller, label_width, control_width, height int, field_text string, range_field *RangeROView) NumberEditView {

	edit_field := NumberEditView{show_edit(parent, label_width, control_width, height, field_text)}
	edit_field.SetPaddingsAll(ERROR_MARGIN)
	edit_field.OnChange().Bind(func(e *windigo.Event) {
		check_or_error(edit_field, range_field.Check(edit_field.GetFixed()))
	})
	return edit_field

}
