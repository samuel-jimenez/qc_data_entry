package GUI

import (
	"strconv"
	"strings"
)

func Format_float(val float64) string {
	return strconv.FormatFloat(val, 'f', 2, 64)
}

func Format_int(val float64) string {
	return strconv.FormatFloat(val, 'f', 0, 64)
}

/*
 * NumbSearchView
 *
 */
type NumbSearchView struct {
	*SearchBox
}

func NumbSearchView_From_SearchBox(SearchBox *SearchBox) *NumbSearchView {
	control := new(NumbSearchView)
	control.SearchBox = SearchBox
	return control
}

func (control *NumbSearchView) Get() float64 {
	val, _ := strconv.ParseFloat(strings.TrimSpace(control.GetSelectedItem()), 64)
	return val
}

func (control *NumbSearchView) Add(val float64) {
	control.AddItem(Format_float(val))
}

func (control *NumbSearchView) AddInt(val float64) {
	control.AddItem(Format_int(val))
}

func (control *NumbSearchView) Set(val float64) {
	control.Search()
	control.SetText(Format_float(val))
}

func (control *NumbSearchView) SetInt(val float64) {
	control.Search()
	control.SetText(Format_int(val))
}
