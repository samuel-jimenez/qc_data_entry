package GUI

import (
	"database/sql"
	"log"
	"strings"

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
	FIELD_HEIGHT,
	NUM_FIELDS,

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

	BUTTON_WIDTH,
	BUTTON_HEIGHT,
	BUTTON_MARGIN,

	GROUP_WIDTH,
	GROUP_HEIGHT,
	GROUP_MARGIN int
)

func Fill_combobox_from_query_rows(control windigo.ComboBoxable, selected_rows *sql.Rows, err error, fn func(*sql.Rows)) {

	if err != nil {
		log.Printf("error: %q: %s\n", err, "fill_combobox_from_query")
		// return -1
	}
	control.DeleteAllItems()
	i := 0
	for selected_rows.Next() {
		fn(selected_rows)
		i++
	}
	if i == 1 {
		control.SetSelectedItem(0)
	}
}
func Fill_combobox_from_query_0_fn(control windigo.ComboBoxable, select_statement *sql.Stmt, fn func(*sql.Rows)) {
	rows, err := select_statement.Query()
	Fill_combobox_from_query_rows(control, rows, err, fn)
}

func Fill_combobox_from_query_1_fn(control windigo.ComboBoxable, select_statement *sql.Stmt, select_id int64, fn func(*sql.Rows)) {
	rows, err := select_statement.Query(select_id)
	Fill_combobox_from_query_rows(control, rows, err, fn)
}

func Fill_combobox_from_query_0_2(control windigo.ComboBoxable, select_statement *sql.Stmt, fn func(int, string)) {
	Fill_combobox_from_query_0_fn(control, select_statement, func(rows *sql.Rows) {
		var (
			id   int
			name string
		)

		if err := rows.Scan(&id, &name); err == nil {
			fn(id, name)
		} else {
			log.Printf("error: %q: %s\n", err, "fill_combobox_from_query")
			// return -1
		}
	})
}

func Fill_combobox_from_query_1_2(control windigo.ComboBoxable, select_statement *sql.Stmt, select_id int64, fn func(int, string)) {
	Fill_combobox_from_query_1_fn(control, select_statement, select_id, func(rows *sql.Rows) {
		var (
			id   int
			name string
		)

		if err := rows.Scan(&id, &name); err == nil {
			fn(id, name)
		} else {
			log.Printf("error: %q: %s\n", err, "fill_combobox_from_query")
			// return -1
		}
	})
}

func Fill_combobox_from_query(control windigo.ComboBoxable, select_statement *sql.Stmt, select_id int64) {
	Fill_combobox_from_query_1_2(control, select_statement, select_id, func(id int, name string) {
		control.AddItem(name)
	})
}

type ComboBox struct {
	*windigo.LabeledComboBox
}

func (control *ComboBox) SetFont(font *windigo.Font) {
	control.ComboBox.SetFont(font)
	control.Label().SetFont(font)
}

func (control *ComboBox) SetLabeledSize(label_width, control_width, height int) {
	control.SetSize(label_width+control_width, height)
	control.Label().SetSize(label_width, height)
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
	data_view := new(SearchBox)

	data_view.ComboBox = NewComboBox(parent, "")
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
