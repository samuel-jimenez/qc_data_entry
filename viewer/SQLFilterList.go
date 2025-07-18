package viewer

import (
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/* SQLFilterList
 *
 */
type SQLFilterList struct {
	Filters []SQLFilter
}

func (data_view SQLFilterList) Get() string {

	var selected []string
	do_where := true
	for _, filter := range data_view.Filters {
		selection := filter.Get()
		if selection != "" {
			if do_where {
				selected = append(selected, "where ")
				do_where = false
			} else {
				selected = append(selected, "\nand ")
			}
			selected = append(selected, selection)
		}
	}
	if do_where {
		selected = append(selected, "where true")
	}

	return strings.Join(selected, "")
}

func NewSQLFilterList() *SQLFilterList {
	view := new(SQLFilterList)
	return view
}

/* SQLFilterListView
 *
 */
type SQLFilterListView struct {
	*windigo.AutoPanel
	// FilterList *SQLFilterList
	Filters map[string]SQLFilterView

	// Key() string
	// Get() func SQLFilterList
}

// func (filter SQLFilterListView) Add()  { return filter.key }

func (data_view *SQLFilterListView) addFilter(key string) {
	data_view.Dock(data_view.Filters[key], windigo.Top)
	data_view.SetSize(data_view.ClientWidth(), data_view.Height()+HEADER_HEIGHT)
}

func (data_view *SQLFilterListView) AddContinuous(key, label string) {

	data_view.Filters[key] = *NewSQLFilterViewContinuous(data_view, key, label)
	data_view.addFilter(key)
}

func (data_view *SQLFilterListView) AddDiscreteMulti(key, label string,
	set []string) {
	data_view.Filters[key] = *NewSQLFilterViewDiscreteMulti(data_view, key, label, set)
	data_view.addFilter(key)
}

func (data_view *SQLFilterListView) AddDiscreteSearch(key, label string,
	set []string) {
	data_view.Filters[key] = *NewSQLFilterViewDiscreteSearch(data_view, key, label, set)
	data_view.addFilter(key)
}

func (data_view *SQLFilterListView) Update(key string,
	set []string) {
	old_height := data_view.Filters[key].Height()
	data_view.Filters[key].Update(set)
	data_view.SetSize(data_view.ClientWidth(), data_view.Height()+data_view.Filters[key].Height()-old_height)

}

func (data_view SQLFilterListView) Get() *SQLFilterList {
	selected := NewSQLFilterList()
	for _, filter := range data_view.Filters {
		selected.Filters = append(selected.Filters, filter.Get())
	}
	return selected
}

func NewSQLFilterListView(parent windigo.Controller) *SQLFilterListView {
	view := new(SQLFilterListView)
	view.AutoPanel = windigo.NewAutoPanel(parent)
	view.Filters = make(map[string]SQLFilterView)
	// view.FilterList = NewSQLFilterList()
	view.SetSize(GUI.OFF_AXIS, FIELD_HEIGHT)

	return view
}
