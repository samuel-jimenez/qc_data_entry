package GUI

import "github.com/samuel-jimenez/windigo"

/*
 * NumbestEditView
 *
 */
type NumbestEditView struct {
	NumberEditView
}

func NumbestEditView_from_new(parent windigo.Controller, field_text string) *NumbestEditView {
	edit_field := NumberEditView_from_new(parent, field_text)
	return &NumbestEditView{*edit_field}
}

func (control *NumbestEditView) SetLabeledSize(label_width, control_width, height int) {
	control.SetSize(label_width+control_width, height)
	control.Label().SetSize(label_width, height)
}
