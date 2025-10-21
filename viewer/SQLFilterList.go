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

// TODO extract to external, then call it
func FilterSQL(Filters ...string) string {
	var selected []string
	do_where := true
	for _, filter := range Filters {
		selection := filter
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

func (data_view SQLFilterList) Get() string {
	var selected []string
	for _, filter := range data_view.Filters {
		selected = append(selected, filter.Get())
	}
	return FilterSQL(selected...)
}

func NewSQLFilterList() *SQLFilterList {
	view := new(SQLFilterList)
	return view
}
