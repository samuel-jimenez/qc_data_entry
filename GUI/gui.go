package GUI

import (
	"database/sql"

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
	HPANEL_WIDTH,
	REPRINT_BUTTON_WIDTH,
	REPRINT_BUTTON_MARGIN_L,

	TOP_PANEL_INTER_SPACER_WIDTH,
	HPANEL_MARGIN,

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
	CLOCK_TIMER_WIDTH,
	CLOCK_TIMER_OFFSET_H,

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
