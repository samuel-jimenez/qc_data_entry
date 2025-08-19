package views

import (
	"strconv"
	"strings"

	"github.com/samuel-jimenez/windigo"
)

/*
 * NumbEditViewer
 *
 */
type NumbEditViewer interface {
	Get() float64
	Set(float64)
}

// TODO combine NumberEditView nmaybe?

/*
* NumbEditView
* 	// cf NumberEditView

*
 */
type NumbEditView struct {
	// GUI.ErrableView
	*windigo.Edit
}

func NewNumbEditView(parent windigo.Controller) *NumbEditView {
	edit_field := new(NumbEditView)
	edit_field.Edit = windigo.NewEdit(parent)
	return edit_field
}

func (control *NumbEditView) Get() float64 {
	val, _ := strconv.ParseFloat(strings.TrimSpace(control.Text()), 64)
	return val
}

func (control *NumbEditView) Set(val float64) {
	start, end := control.Selected()
	control.SetText(strconv.FormatFloat(val, 'f', 2, 64))
	control.SelectText(start, end)
}

func (control *NumbEditView) SetInt(val float64) {
	start, end := control.Selected()
	control.SetText(strconv.FormatFloat(val, 'f', 0, 64))
	control.SelectText(start, end)
}
