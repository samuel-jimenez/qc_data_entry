package viewer

import (
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
	view.key = key

	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.label = SQLFilterViewHeader_from_new(view, label)

	view.selection_options = BuildNewDiscreteSearchView(view, set)
	view.label.SetHidePanel(view.selection_options.AutoPanel)

	view.AutoPanel.Dock(view.label, windigo.Top)
	view.AutoPanel.Dock(view.selection_options, windigo.Top)

	return view
}

func (view *SQLFilterViewDiscreteSearch) Get() string {
	return SQLFilterDiscrete{view.key,
		view.selection_options.Get()}.Get()
}

func (view *SQLFilterViewDiscreteSearch) Update(set []string) {
	view.selection_options.Update(set)
}

func (view *SQLFilterViewDiscreteSearch) AddItem(entry string) {
	if view.selection_options.Contains(entry) {
		return
	}

	view.SQLFilterView.AddItem(entry)
	view.selection_options.AddItem(entry)
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
