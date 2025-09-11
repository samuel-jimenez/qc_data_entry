package viewer

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

//TODO add collapse/expand all

/* SQLFilterListView
 *
 */
type SQLFilterListView struct {
	*windigo.AutoPanel
	// FilterList *SQLFilterList
	Filters map[string]SQLFilterViewable

	// Key() string
	// Get() func SQLFilterList
}

func NewSQLFilterListView(parent windigo.Controller) *SQLFilterListView {
	view := new(SQLFilterListView)
	view.AutoPanel = windigo.NewAutoPanel(parent)
	view.Filters = make(map[string]SQLFilterViewable)
	// view.FilterList = NewSQLFilterList()
	view.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)

	return view
}

func (view *SQLFilterListView) RefreshSize() {
	client_height := 0
	for _, control := range view.Filters {
		client_height -= control.Height()
		control.RefreshSize()
		client_height += control.Height()
	}
	// view.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)
	view.SetSize(view.Width(), view.Height()+client_height)
}

func (view *SQLFilterListView) SetFont(font *windigo.Font) {
	view.AutoPanel.SetFont(font)
	for _, control := range view.Filters {
		control.SetFont(font)
	}
}

// func (filter SQLFilterListView) Add()  { return filter.key }

func (view *SQLFilterListView) addFilter(key string) {
	view.Dock(view.Filters[key], windigo.Top)
	view.SetSize(view.ClientWidth(), view.Height()+HEADER_HEIGHT)
}

func (view *SQLFilterListView) AddContinuous(key, label string) {

	view.Filters[key] = NewSQLFilterViewContinuous(view, key, label)
	view.addFilter(key)
}

func (view *SQLFilterListView) AddDiscreteMulti(key, label string,
	set []string) {
	view.Filters[key] = NewSQLFilterViewDiscreteMulti(view, key, label, set)
	view.addFilter(key)
}

func (view *SQLFilterListView) AddDiscreteSearch(key, label string,
	set []string) {
	view.Filters[key] = NewSQLFilterViewDiscreteSearch(view, key, label, set)
	view.addFilter(key)
}

func (view *SQLFilterListView) Update(key string,
	set []string) {
	old_height := view.Filters[key].Height()
	view.Filters[key].Update(set)
	view.SetSize(view.ClientWidth(), view.Height()+view.Filters[key].Height()-old_height)

}

func (view SQLFilterListView) Get() *SQLFilterList {
	selected := NewSQLFilterList()
	for _, filter := range view.Filters {
		selected.Filters = append(selected.Filters, filter.Get())
	}
	return selected
}

func (view *SQLFilterListView) Clear() {
	for _, filter := range view.Filters {
		filter.Clear()
	}
}
