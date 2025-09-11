package viewer

import (
	"strings"
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
