package GUI

import (
	"database/sql"
	"strings"

	"github.com/samuel-jimenez/windigo"
)

/*
 * SearchBox
 *
 * TODO reconcile with DiscreteSearchView
 *
 * ComboBoxable ?
 */
type SearchBox struct {
	ComboBox
	entries  []string
	onChange windigo.EventManager
}

func NewSearchBox(parent windigo.Controller) *SearchBox {
	return NewLabeledSearchBox(parent, "")
}

func NewLabeledSearchBox(parent windigo.Controller, Label string) *SearchBox {
	data_view := new(SearchBox)

	data_view.ComboBox = *ComboBox_from_new(parent, Label)
	data_view.ComboBox.OnChange().Bind(func(e *windigo.Event) {
		start, _ := data_view.Selected()

		text := strings.ToUpper(data_view.Text())
		terms := strings.Split(text, " ")
		data_view.Search(terms...)

		data_view.SetText(text)

		data_view.SelectText(start, -1)
		data_view.onChange.Fire(e)
	})

	return data_view
}

func NewSearchBoxWithLabels(parent windigo.Controller, labels []string) *SearchBox {
	data_view := NewSearchBox(parent)
	data_view.Update(labels)
	return data_view
}

func NewSearchBoxFromQuery(parent windigo.Controller, select_statement *sql.Stmt, args ...any) *SearchBox {
	return NewLabeledSearchBoxFromQuery(parent, "", select_statement, args...)
}

func NewLabeledSearchBoxFromQuery(parent windigo.Controller, Label string, select_statement *sql.Stmt, args ...any) *SearchBox {
	data_view := NewLabeledSearchBox(parent, Label)
	data_view.FromQuery(select_statement, args...)
	return data_view
}

func NewLabeledListSearchBox(parent windigo.Controller, Label string) *SearchBox {
	// data_view := NewLabeledSearchBox(parent, Label)
	data_view := new(SearchBox)

	data_view.ComboBox = *ComboBox_from_new(parent, Label)
	data_view.ComboBox.OnChange().Bind(func(e *windigo.Event) {
		start, _ := data_view.Selected()

		text := strings.ToUpper(data_view.Text())
		terms := strings.Split(text, " ")
		data_view.Search(terms...)

		// TODO split
		data_view.ShowDropdown(false)
		data_view.SetText(text)
		data_view.ShowDropdown(true)

		data_view.SelectText(start, -1)
		data_view.onChange.Fire(e)
	})
	return data_view
}

func NewListSearchBox(parent windigo.Controller) *SearchBox {
	return NewLabeledListSearchBox(parent, "")
}

func NewListSearchBoxWithLabels(parent windigo.Controller, labels []string) *SearchBox {
	data_view := NewListSearchBox(parent)
	data_view.Update(labels)
	return data_view
}

func NewListSearchBoxFromQuery(parent windigo.Controller, select_statement *sql.Stmt, args ...any) *SearchBox {
	data_view := NewListSearchBox(parent)
	data_view.FromQuery(select_statement, args...)
	return data_view
}

func NewLabeledListSearchBoxFromQuery(parent windigo.Controller, Label string, select_statement *sql.Stmt, args ...any) *SearchBox {
	data_view := NewLabeledListSearchBox(parent, Label)
	data_view.FromQuery(select_statement, args...)
	return data_view
}

func (data_view *SearchBox) Update(set []string) {
	data_view.entries = set
	data_view.DeleteAllItems()
	for _, name := range set {
		data_view.ComboBox.AddItem(name)
	}
}

// func (control *SearchBox) Height() int {
// 	rect := w32.GetWindowRect(control.hwnd)
// 	// return int(rect.Bottom - rect.Top)
// 	log.Println("ComboBox-Height", int(rect.Bottom-rect.Top), int(w32.SendMessage(control.hwnd, w32.CB_GETITEMHEIGHT, 0, 0)))
// 	// return int(w32.SendMessage(control.hwnd, w32.CB_GETITEMHEIGHT, uintptr(-1),0))
// 	return int(w32.SendMessage(control.hwnd, w32.CB_GETITEMHEIGHT, 0, 0))
// }

// func (control *SearchBox) Height() int {
// 	rect := w32.GetWindowRect(control.Handle())
// 	log.Println("SearchBox-Height", int(rect.Bottom-rect.Top), int(w32.SendMessage(control.Handle(), w32.CB_GETITEMHEIGHT, 0, 0)))
// 	MaxUint := uintptr(^uint32(0)) // uintptr(-1)
//
// 	// log.Println("LotView-Height", int(rect.Bottom-rect.Top), int(w32.SendMessage(control.Handle(), w32.CB_GETITEMHEIGHT, 0, 0)))
// 	log.Println("SearchBox-Height", MaxUint, int(rect.Bottom-rect.Top), int(w32.SendMessage(control.Handle(), w32.CB_GETITEMHEIGHT, MaxUint, 0)))
// 	// return int(w32.SendMessage(control.hwnd, w32.CB_GETITEMHEIGHT, uintptr(-1),0))
// 	// return int(w32.SendMessage(control.Handle(), w32.CB_GETITEMHEIGHT, 0, 0))
// 	return int(rect.Bottom - rect.Top)
// }

func (data_view *SearchBox) Entries() []string {
	return data_view.entries
}

func (data_view *SearchBox) AddItem(entry string) bool {
	val := data_view.ComboBox.AddItem(entry)
	if val {
		data_view.entries = append(data_view.entries, entry)
	}
	return val
}

func (data_view *SearchBox) FromQuery(select_statement *sql.Stmt, args ...any) {
	data_view.Fill_FromFnQuery(
		data_view.query_fn,
		select_statement, args...)
}

func (data_view *SearchBox) query_fn(id int, name string) {
	data_view.AddItem(name)
}

func (data_view *SearchBox) Fill_FromFnQuery(fn func(int, string), select_statement *sql.Stmt, args ...any) {
	data_view.entries = nil
	Fill_combobox_from_query_fn(
		data_view,
		fn,
		select_statement, args...)
}

// TODO this one doesn't autoselect
// go though zand make suere it cant be flattened
// TODO c.f Fill_combobox_from_query_fn
// data_view.entries = nil
// Fill_combobox_from_query_fn(data_view, fn, select_statement, args...)
// func (data_view *SearchBox) Fill_FromFnQuery(fn func(int, string), select_statement *sql.Stmt, args ...any) {
// 	DB.Forall_err("SearchBox.FromFnQuery",
// 		func() {
// 			data_view.entries = nil
// 			data_view.DeleteAllItems()
// 		},
// 		func(row *sql.Rows) error {
// 			var (
// 				id   int
// 				name string
// 			)
//
// 			if err := row.Scan(
// 				&id, &name,
// 			); err != nil {
// 				return err
// 			}
// 			fn(id, name)
// 			return nil
// 		},
// 		select_statement, args...)
// }

// func Fill_combobox_from_query_fn(control windigo.ComboBoxable, fn func(int, string), select_statement *sql.Stmt, args ...any) {
// 	i := 0
// 	DB.Forall("fill_combobox_from_query",
// 		func() {
// 			control.DeleteAllItems()
// 		},
// 		func(row *sql.Rows) {
// 			var (
// 				id   int
// 				name string
// 			)
//
// 			if err := row.Scan(
// 				&id, &name,
// 			); err == nil {
// 				fn(id, name)
// 			} else {
// 				log.Printf("error: [%s]: %q\n", "fill_combobox_from_query", err)
// 				// return -1
// 			}
// 			i++
// 		},
// 		select_statement, args...)
// 	if i == 1 {
// 		control.SetSelectedItem(0)
// 	}
// }

func (data_view SearchBox) Search(terms ...string) {
	data_view.DeleteAllItems()

	for _, entry := range data_view.entries {
		matched := true
		upcased := strings.ToUpper(entry)

		for _, term := range terms {
			if !strings.Contains(upcased, term) {
				matched = false
				break
			}
		}
		if matched {
			data_view.ComboBox.AddItem(entry)
		}
	}
}

func (control *SearchBox) OnChange() *windigo.EventManager {
	return &control.onChange
}
