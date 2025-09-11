package viewer

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/* SQLFilterViewDiscreteMulti
 *
 */
type SQLFilterViewDiscreteMulti struct {
	SQLFilterView
	selection_options *DiscreteMultiView
}

func NewSQLFilterViewDiscreteMulti(parent windigo.Controller, key, label string, set []string) *SQLFilterViewDiscreteMulti {
	view := new(SQLFilterViewDiscreteMulti)
	panel := windigo.NewAutoPanel(parent)
	view.AutoPanel = panel
	panel.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)

	panel_label := NewSQLFilterViewHeader(view, label)
	panel_label.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)

	selection_options := BuildNewDiscreteMultiView(panel, set)
	panel_label.SetHidePanel(selection_options.AutoPanel)

	panel.Dock(panel_label, windigo.Top)
	panel.Dock(selection_options, windigo.Top)

	view.label = panel_label
	view.selection_options = selection_options
	view.key = key

	return view

}

func (view *SQLFilterViewDiscreteMulti) Get() SQLFilter {
	return SQLFilterDiscrete{view.key,
		view.selection_options.Get()}
}

func (view *SQLFilterViewDiscreteMulti) Update(set []string) {
	view.selection_options.Close()
	view.selection_options = BuildNewDiscreteMultiView(view.AutoPanel, set)
	if view.label.Hidden {
		view.label.SetHidePanel(view.selection_options.AutoPanel)
	} else {
		view.label.SetShowPanel(view.selection_options.AutoPanel)
		view.AutoPanel.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT+view.selection_options.Height())
	}
	view.AutoPanel.Dock(view.selection_options, windigo.Top)
}

func (view *SQLFilterViewDiscreteMulti) RefreshSize() {
	view.SQLFilterView.RefreshSize()
	view.selection_options.RefreshSize()
}

func (view *SQLFilterViewDiscreteMulti) SetFont(font *windigo.Font) {
	view.SQLFilterView.SetFont(font)
	view.selection_options.SetFont(font)
}
func (view *SQLFilterViewDiscreteMulti) Clear() {
	view.SQLFilterView.Clear()
	view.selection_options.Clear()
}
