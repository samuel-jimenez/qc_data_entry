package viewer

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/* SQLFilterViewDiscreteSearch
 *
 */
type SQLFilterViewDiscreteSearch struct {
	SQLFilterView
	selection_options *DiscreteSearchView
}

func NewSQLFilterViewDiscreteSearch(parent windigo.Controller, key, label string, set []string) *SQLFilterViewDiscreteSearch {
	view := new(SQLFilterViewDiscreteSearch)

	panel := windigo.NewAutoPanel(parent)
	view.AutoPanel = panel

	panel.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)

	panel_label := NewSQLFilterViewHeader(view, label)
	panel_label.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)

	selection_options := BuildNewDiscreteSearchView(view, set)
	panel_label.SetHidePanel(selection_options.AutoPanel)
	selection_options.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)

	panel.Dock(panel_label, windigo.Top)
	panel.Dock(selection_options, windigo.Top)

	view.label = panel_label
	view.selection_options = selection_options
	view.key = key

	return view
}

func (view *SQLFilterViewDiscreteSearch) Get() SQLFilter {
	return SQLFilterDiscrete{view.key,
		view.selection_options.Get()}
}

func (view *SQLFilterViewDiscreteSearch) Update(set []string) {
	view.selection_options.Update(set)
}

func (view *SQLFilterViewDiscreteSearch) DelItem(entry string) {
	view.selection_options.DelItem(entry)
}

func (view *SQLFilterViewDiscreteSearch) RefreshSize() {
	view.SQLFilterView.RefreshSize()
	view.selection_options.RefreshSize()
}

func (view *SQLFilterViewDiscreteSearch) SetFont(font *windigo.Font) {
	view.SQLFilterView.SetFont(font)
	view.selection_options.SetFont(font)
}

func (view *SQLFilterViewDiscreteSearch) Clear() {
	view.SQLFilterView.Clear()
	view.selection_options.Clear()
}
