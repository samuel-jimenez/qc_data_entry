package views

import (
	"strings"

	"github.com/samuel-jimenez/windigo"
)

/*
 * NullView
 *
 */
type NullView struct {
	*windigo.Edit
}

func NullView_from_new(parent windigo.Controller) *NullView {

	view := new(NullView)

	view.Edit = windigo.NewEdit(parent)
	view.Edit.OnKillFocus().Bind(view.Edit_OnKillFocus)

	return view
}

func (view *NullView) Edit_OnKillFocus(e *windigo.Event) {
	view.Edit.SetText(strings.TrimSpace(view.Edit.Text()))
}
