package GUI

import (
	"database/sql"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/windigo"
	"github.com/samuel-jimenez/windigo/w32"
)

var (
	BASE_FONT_SIZE = 10
	OFF_AXIS       = 0
	RANGES_PADDING = 10
	LABEL_WIDTH    = 100

	ERROR_MARGIN = 3

	ErroredPen = windigo.NewPen(w32.PS_GEOMETRIC, 2, windigo.NewSolidColorBrush(windigo.RGB(255, 0, 64)))
	OKPen      = windigo.NewPen(w32.PS_GEOMETRIC, 2, windigo.NewSystemColorBrush(w32.COLOR_BTNFACE))
)

var (
	WINDOW_WIDTH,
	WINDOW_HEIGHT,
	WINDOW_FUDGE_MARGIN,

	GROUPBOX_CUSHION,
	TOP_SPACER_WIDTH,
	TOP_SPACER_HEIGHT,
	INTER_SPACER_HEIGHT,
	BTM_SPACER_WIDTH,
	BTM_SPACER_HEIGHT,

	TOP_PANEL_WIDTH,
	REPRINT_BUTTON_WIDTH,
	TOP_PANEL_INTER_SPACER_WIDTH,

	RANGE_WIDTH,
	PRODUCT_TYPE_WIDTH,

	PRODUCT_FIELD_WIDTH,
	PRODUCT_FIELD_HEIGHT,
	DATA_FIELD_WIDTH,
	DATA_SUBFIELD_WIDTH,
	DATA_UNIT_WIDTH,
	DATA_MARGIN,
	SOURCES_LABEL_WIDTH,
	SOURCES_FIELD_WIDTH,
	EDIT_FIELD_HEIGHT,
	NUM_FIELDS,

	CLOCK_WIDTH,

	ACCEPT_BUTTON_WIDTH,
	CANCEL_BUTTON_WIDTH,

	DISCRETE_FIELD_WIDTH,
	DISCRETE_FIELD_HEIGHT,

	RANGES_WINDOW_WIDTH,
	RANGES_WINDOW_HEIGHT,
	RANGES_WINDOW_PADDING,
	RANGES_FIELD_HEIGHT,
	RANGES_FIELD_SMALL_HEIGHT,
	RANGES_FIELD_WIDTH,
	RANGES_BUTTON_WIDTH,
	RANGES_BUTTON_HEIGHT,
	RANGES_BUTTON_MARGIN,

	RANGES_RO_PADDING,
	RANGES_RO_FIELD_WIDTH,
	RANGES_RO_SPACER_WIDTH,
	RANGES_RO_FIELD_HEIGHT,

	SMOL_BUTTON_WIDTH,
	BUTTON_WIDTH,
	BUTTON_HEIGHT,
	BUTTON_MARGIN,

	GROUP_WIDTH,
	GROUP_HEIGHT,
	GROUP_MARGIN int
)

func Fill_combobox_from_query_rows(control windigo.ComboBoxable, fn func(row *sql.Rows) error, select_statement *sql.Stmt, args ...any) {
	i := 0
	DB.Forall_err("fill_combobox_from_query",
		func() {
			control.DeleteAllItems()
		},
		func(row *sql.Rows) error {
			i++
			return fn(row)
		},
		select_statement, args...)
	if i == 1 {
		control.SetSelectedItem(0)
	}
}

func Fill_combobox_from_query_fn(control windigo.ComboBoxable, fn func(int, string), select_statement *sql.Stmt, args ...any) {
	i := 0
	DB.Forall_err("fill_combobox_from_query",
		func() {
			control.DeleteAllItems()
		},
		func(row *sql.Rows) error {
			var (
				id   int
				name string
			)

			if err := row.Scan(
				&id, &name,
			); err != nil {
				return err
			}
			fn(id, name)
			i++
			return nil
		},
		select_statement, args...)
	if i == 1 {
		control.SetSelectedItem(0)
	}
}

func Fill_combobox_from_query(control windigo.ComboBoxable, select_statement *sql.Stmt, args ...any) {
	Fill_combobox_from_query_fn(control, func(id int, name string) { control.AddItem(name) }, select_statement, args...)
}

/* ComboBoxable
 *
 */
type ComboBoxable interface {
	windigo.ComboBoxable
}

type ComboBox struct {
	*windigo.LabeledComboBox
}

func NewComboBox(parent windigo.Controller, field_text string) *ComboBox {
	combobox_field := &ComboBox{windigo.NewLabeledComboBox(parent, field_text)}
	combobox_field.OnKillFocus().Bind(func(e *windigo.Event) {
		combobox_field.SetText(strings.ToUpper(strings.TrimSpace(combobox_field.Text())))
	})

	return combobox_field
}

func NewListComboBox(parent windigo.Controller, field_text string) *ComboBox {
	combobox_field := &ComboBox{windigo.NewLabeledListComboBox(parent, field_text)}
	combobox_field.OnKillFocus().Bind(func(e *windigo.Event) {
		combobox_field.SetText(strings.ToUpper(strings.TrimSpace(combobox_field.Text())))
	})

	return combobox_field
}

func NewSizedListComboBox(parent windigo.Controller, label_width, control_width, height int, field_text string) *ComboBox {
	combobox_field := &ComboBox{windigo.NewSizedLabeledListComboBox(parent, label_width, control_width, height, field_text)}
	combobox_field.OnKillFocus().Bind(func(e *windigo.Event) {
		combobox_field.SetText(strings.ToUpper(strings.TrimSpace(combobox_field.Text())))
	})

	return combobox_field
}

/*
* SearchBox
*
* TODO reconcile with DiscreteSearchView

* ComboBoxable ?
 */
type SearchBox struct {
	*ComboBox
	entries  []string
	onChange windigo.EventManager
}

func (data_view *SearchBox) Update(set []string) {
	data_view.entries = set
	data_view.DeleteAllItems()
	for _, name := range set {
		data_view.AddItem(name)
	}
}

func (data_view *SearchBox) AddEntry(entry string) {
	data_view.AddItem(entry)
	data_view.entries = append(data_view.entries, entry)
}

func (data_view *SearchBox) FromQuery(select_statement *sql.Stmt, args ...any) {
	data_view.entries = nil
	Fill_combobox_from_query_fn(data_view, func(id int, name string) {
		data_view.AddItem(name)
		data_view.entries = append(data_view.entries, name)
	}, select_statement, args...)
}

// TODO c.f Fill_combobox_from_query_fn, Fill_combobox_from_query_rows_0
// TODO? maybe rmove
func (data_view *SearchBox) _FromFn_Query_xx(fn func(int, string), select_statement *sql.Stmt, args ...any) {
	DB.Forall_err("SearchBox.FromFnQuery",
		func() {
			data_view.entries = nil
			data_view.DeleteAllItems()
		},
		func(row *sql.Rows) error {
			var (
				id   int
				name string
			)

			if err := row.Scan(
				&id, &name,
			); err != nil {
				return err
			}
			//??TODO do we needthis??
			// data_view.AddItem(name)
			data_view.entries = append(data_view.entries, name)
			fn(id, name)
			return nil
		},
		select_statement, args...)
}

// TODO c.f Fill_combobox_from_query_fn, Fill_combobox_from_query_rows_0
// data_view.entries = nil
// Fill_combobox_from_query_fn(data_view, fn, select_statement, args...)
func (data_view *SearchBox) Fill_FromFnQuery(fn func(int, string), select_statement *sql.Stmt, args ...any) {
	DB.Forall_err("SearchBox.FromFnQuery",
		func() {
			data_view.entries = nil
			data_view.DeleteAllItems()
		},
		func(row *sql.Rows) error {
			var (
				id   int
				name string
			)

			if err := row.Scan(
				&id, &name,
			); err != nil {
				return err
			}
			fn(id, name)
			return nil
		},
		select_statement, args...)
}

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

func (data_view SearchBox) Search(terms []string) {
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
			data_view.AddItem(entry)
		}
	}
}

func (control *SearchBox) OnChange() *windigo.EventManager {
	return &control.onChange
}

func NewSearchBox(parent windigo.Controller) *SearchBox {
	return NewLabeledSearchBox(parent, "")
}

func NewLabeledSearchBox(parent windigo.Controller, Label string) *SearchBox {
	data_view := new(SearchBox)

	data_view.ComboBox = NewComboBox(parent, Label)
	data_view.ComboBox.OnChange().Bind(func(e *windigo.Event) {

		start, _ := data_view.Selected()

		text := strings.ToUpper(data_view.Text())
		terms := strings.Split(text, " ")
		data_view.Search(terms)

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

	data_view.ComboBox = NewComboBox(parent, Label)
	data_view.ComboBox.OnChange().Bind(func(e *windigo.Event) {

		start, _ := data_view.Selected()

		text := strings.ToUpper(data_view.Text())
		terms := strings.Split(text, " ")
		data_view.Search(terms)

		//TODO split
		data_view.ShowDropdown(false)
		data_view.SetText(text)
		data_view.ShowDropdown(true)

		data_view.SelectText(start, -1)
		data_view.onChange.Fire(e)
	})
	return data_view

}

func NewSizedLabeledListSearchBox(parent windigo.Controller, label_width, control_width, height int, field_text string) *SearchBox {
	data_view := NewLabeledListSearchBox(parent, field_text)
	data_view.SetLabeledSize(label_width, control_width, height)
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
