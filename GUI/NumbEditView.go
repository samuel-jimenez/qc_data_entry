package GUI

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
	GetInt() int
	Set(float64)
	SetInt(val float64)
}

/*
 * NumbEditView
 *
 */
type NumbEditView struct {
	// ErrableView
	*windigo.Edit
}

func NumbEditView_from_new(parent windigo.Controller) *NumbEditView {
	edit_field := new(NumbEditView)
	edit_field.Edit = windigo.NewEdit(parent)
	return edit_field
}

func (control *NumbEditView) Get() float64 {
	val, _ := strconv.ParseFloat(strings.TrimSpace(control.Text()), 64)
	return val
}

func (control *NumbEditView) GetInt() int {
	val, _ := strconv.Atoi(strings.TrimSpace(control.Text()))
	return val
}

func (control *NumbEditView) Set(val float64) {
	start, end := control.Selected()
	control.SetText(Format_float(val))
	control.SelectText(start, end)
}

func (control *NumbEditView) Setf(val float64, Format func(float64) string) {
	start, end := control.Selected()
	control.SetText(Format(val))
	control.SelectText(start, end)
}

func (control *NumbEditView) SetInt(val float64) {
	start, end := control.Selected()
	control.SetText(Format_int(val))
	control.SelectText(start, end)
}
