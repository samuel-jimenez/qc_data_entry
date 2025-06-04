package main

import (
	"database/sql"
	"log"
	"strings"

	"github.com/samuel-jimenez/windigo"
	"github.com/samuel-jimenez/windigo/w32"
)

var (
	RANGE_WIDTH        = 200
	GROUP_WIDTH        = 210
	GROUP_HEIGHT       = 170
	GROUP_MARGIN       = 5
	PRODUCT_TYPE_WIDTH = 150

	LABEL_WIDTH         = 100
	PRODUCT_FIELD_WIDTH = 150
	DATA_FIELD_WIDTH    = 60
	FIELD_HEIGHT        = 28

	RANGES_PADDING         = 5
	RANGES_FIELD_HEIGHT    = 50
	OFF_AXIS               = 0
	RANGES_RO_FIELD_WIDTH  = 40
	RANGES_RO_SPACER_WIDTH = 20
	RANGES_RO_FIELD_HEIGHT = FIELD_HEIGHT

	BUTTON_WIDTH  = 100
	BUTTON_HEIGHT = 40
	// 	200
	// 50

	ERROR_MARGIN = 3

	TOP_SPACER_WIDTH     = 7
	TOP_SPACER_HEIGHT    = 17
	INTER_SPACER_HEIGHT  = 2
	BTM_SPACER_WIDTH     = 2
	BTM_SPACER_HEIGHT    = 2
	BUTTON_SPACER_HEIGHT = 195

	erroredPen = windigo.NewPen(w32.PS_GEOMETRIC, 2, windigo.NewSolidColorBrush(windigo.RGB(255, 0, 64)))
	okPen      = windigo.NewPen(w32.PS_GEOMETRIC, 2, windigo.NewSystemColorBrush(w32.COLOR_BTNFACE))
)

func fill_combobox_from_query_rows(control windigo.LabeledComboBox, selected_rows *sql.Rows, err error, fn func(*sql.Rows)) {

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

func fill_combobox_from_query_fn(control windigo.LabeledComboBox, select_statement *sql.Stmt, select_id int64, fn func(*sql.Rows)) {
	rows, err := select_statement.Query(select_id)
	fill_combobox_from_query_rows(control, rows, err, fn)
}

func fill_combobox_from_query(control windigo.LabeledComboBox, select_statement *sql.Stmt, select_id int64) {
	fill_combobox_from_query_fn(control, select_statement, select_id, func(rows *sql.Rows) {
		var (
			id   uint8
			name string
		)

		if err := rows.Scan(&id, &name); err == nil {
			// data[id] = value
			control.AddItem(name)
		} else {
			log.Printf("error: %q: %s\n", err, "fill_combobox_from_query")
			// return -1
		}
	})
}

func build_text_dock(parent windigo.Controller, labels []string) windigo.Pane {
	panel := windigo.NewAutoPanel(parent)
	panel.SetSize(50, 50)
	panel.SetPaddingsAll(10)

	for _, label_text := range labels {
		label := windigo.NewLabel(panel)
		label.SetSize(200, 25)
		label.SetText(label_text)
		panel.Dock(label, windigo.Left)

	}
	return panel
}

func build_button_dock(parent windigo.Controller, labels []string, onclicks []func()) windigo.Pane {
	assertEqual(len(labels), len(onclicks))
	panel := windigo.NewAutoPanel(parent)
	panel.SetSize(50, 50)
	panel.SetPaddingLeft(10)

	for i, label_text := range labels {
		button := windigo.NewPushButton(panel)
		button.SetSize(200, 25)
		button.SetMarginsAll(10)

		button.OnClick().Bind(func(e *windigo.Event) {
			onclicks[i]()

		})
		button.SetText(label_text)
		panel.Dock(button, windigo.Left)

	}
	return panel
}

func build_marginal_button_dock(parent windigo.Controller, width, height int, labels []string, margins []int, onclicks []func()) windigo.Pane {
	assertEqual(len(labels), len(margins))
	assertEqual(len(labels), len(onclicks))
	panel := windigo.NewAutoPanel(parent)
	panel.SetSize(width, height)

	for i, label_text := range labels {
		button := windigo.NewPushButton(panel)
		button.SetSize(width, height)
		button.SetMarginLeft(margins[i])

		button.OnClick().Bind(func(e *windigo.Event) {
			onclicks[i]()

		})
		button.SetText(label_text)
		panel.Dock(button, windigo.Left)

	}
	return panel
}

func show_checkbox(parent windigo.Controller, width, height int, field_text string) windigo.LabeledCheckBox {
	checkbox_field := windigo.NewLabeledCheckBox(parent, width, height, field_text)
	checkbox_field.SetPaddingsAll(ERROR_MARGIN)
	return checkbox_field
}

func show_combobox(parent windigo.Controller, label_width, control_width, height int, field_text string) windigo.LabeledComboBox {
	combobox_field := windigo.NewLabeledComboBox(parent, label_width, control_width, height, field_text)

	combobox_field.OnKillFocus().Bind(func(e *windigo.Event) {
		combobox_field.SetText(strings.ToUpper(strings.TrimSpace(combobox_field.Text())))
	})

	return combobox_field
}
