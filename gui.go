package main

import (
	"database/sql"
	"log"
	"strconv"
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
	TOP_SPACER_HEIGHT    = 22
	INTER_SPACER_HEIGHT  = 2
	BTM_SPACER_WIDTH     = 2
	BTM_SPACER_HEIGHT    = 2
	BUTTON_SPACER_HEIGHT = 195

	erroredPen = windigo.NewPen(w32.PS_GEOMETRIC, 2, windigo.NewSolidColorBrush(windigo.RGB(255, 0, 64)))
	okPen      = windigo.NewPen(w32.PS_GEOMETRIC, 2, windigo.NewSystemColorBrush(w32.COLOR_BTNFACE))
)

func parse_field(field windigo.Controller) float64 {
	val, _ := strconv.ParseFloat(strings.TrimSpace(field.Text()), 64)
	return val
}

func edit_and_parse_field(field windigo.LabeledEdit) float64 {
	start, end := field.Selected()
	// IndexAny(s, chars string) int
	field.SetText(strings.TrimSpace(field.Text()))
	// mass_field.SelectText(-1, -1)
	field.SelectText(start, end)
	// mass_field.SelectText(-1, 0)

	return parse_field(field)
}

func check_or_error(field windigo.Bordered, test bool) {
	if test {
		field.SetBorder(okPen)
	} else {
		field.SetBorder(erroredPen)
	}
}

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

func show_checkbox(parent windigo.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *windigo.CheckBox {
	checkbox_label := windigo.NewLabel(parent)
	checkbox_label.SetPos(x_label_pos, y_pos)

	checkbox_label.SetText(field_text)

	checkbox_field := windigo.NewCheckBox(parent)
	checkbox_field.SetText("")

	checkbox_field.SetPos(x_field_pos, y_pos)
	// visual_label.OnClick().Bind(func(e *windigo.Event) {
	// 		visual_field.SetFocus()
	// })
	return checkbox_field
}

func show_combobox(parent windigo.Controller, label_width, control_width, height int, field_text string) windigo.LabeledComboBox {
	combobox_field := windigo.NewLabeledComboBox(parent, label_width, control_width, height, field_text)

	combobox_field.OnKillFocus().Bind(func(e *windigo.Event) {
		combobox_field.SetText(strings.ToUpper(strings.TrimSpace(combobox_field.Text())))
	})

	return combobox_field
}

func show_edit(parent windigo.Controller, label_width, control_width, height int, field_text string) windigo.LabeledEdit {
	return windigo.NewLabeledEdit(parent, label_width, control_width, height, field_text)
}

func show_number_edit(parent windigo.Controller, label_width, control_width, height int, field_text string, range_field *RangeROView) windigo.LabeledEdit {

	edit_field := show_edit(parent, label_width, control_width, height, field_text)
	edit_field.SetPaddingsAll(ERROR_MARGIN)
	edit_field.OnChange().Bind(func(e *windigo.Event) {
		check_or_error(edit_field, range_field.Check(edit_and_parse_field(edit_field)))
	})
	return edit_field

}

func show_text(parent windigo.AutoPanel, label_width, field_width, field_height int, field_text string, field_units string) windigo.LabeledEdit {

	margin := 10
	panel := windigo.NewAutoPanel(parent)
	panel.SetSize(label_width+2*field_width, field_height)

	text_label := windigo.NewLabel(panel)
	text_label.SetSize(label_width, field_height)
	text_label.SetText(field_text)

	text_field := windigo.NewEdit(panel)
	text_field.SetSize(field_width-margin, field_height)
	text_field.SetText("0.000")
	text_field.SetMarginRight(margin)

	text_units := windigo.NewLabel(panel)
	text_units.SetSize(field_width, field_height)
	text_units.SetText(field_units)

	panel.Dock(text_label, windigo.Left)
	panel.Dock(text_field, windigo.Left)
	panel.Dock(text_units, windigo.Left)
	parent.Dock(panel, windigo.Bottom)

	return windigo.LabeledEdit{ComponentFrame: panel, Edit: text_field}
}

// TODO View
func clear_field(field windigo.LabeledEdit) {
	field.SetText("")
	field.SetBorder(okPen)
}

// func show_mass_sg(parent windigo.AutoPanel, label_width, control_width, height int, field_text string, ranges_panel MassRangesView) windigo.LabeledEdit {
func show_mass_sg(parent windigo.AutoPanel, label_width, control_width, height int, field_text string, ranges_panel MassRangesView) MassDataView {

	field_width := DATA_FIELD_WIDTH

	sg_text := "Specific Gravity"
	density_text := "Density"

	sg_units := "g/mL"
	density_units := "lb/gal"

	mass_field := show_edit(parent, label_width, control_width, height, field_text)
	mass_field.SetPaddingsAll(ERROR_MARGIN)

	//PUSH TO BOTTOM
	density_field := show_text(parent, label_width, field_width, height, density_text, density_units)
	density_field.SetPaddingsAll(ERROR_MARGIN)

	sg_field := show_text(parent, label_width, field_width, height, sg_text, sg_units)
	sg_field.SetPaddingsAll(ERROR_MARGIN)

	check_or_error_mass := func(mass, sg, density float64) {
		check_or_error(mass_field, ranges_panel.CheckMass(mass))
		check_or_error(sg_field, ranges_panel.CheckSG(sg))
		check_or_error(density_field, ranges_panel.CheckDensity(density))
	}

	mass_field.OnChange().Bind(func(e *windigo.Event) {
		mass := edit_and_parse_field(mass_field)
		sg := sg_from_mass(mass)
		density := density_from_sg(sg)

		sg_field.SetText(format_sg(sg, false))
		density_field.SetText(format_density(density))

		check_or_error_mass(mass, sg, density)

	})

	sg_field.OnChange().Bind(func(e *windigo.Event) {
		sg := edit_and_parse_field(sg_field)
		mass := mass_from_sg(sg)
		density := density_from_sg(sg)

		mass_field.SetText(format_mass(mass))
		density_field.SetText(format_density(density))

		check_or_error_mass(mass, sg, density)
	})

	density_field.OnChange().Bind(func(e *windigo.Event) {

		density := edit_and_parse_field(density_field)
		sg := sg_from_density(density)
		mass := mass_from_sg(sg)

		sg_field.SetText(format_sg(sg, false))
		mass_field.SetText(format_mass(mass))

		check_or_error_mass(mass, sg, density)

	})

	Clear := func() {
		// mass_field.Clear()
		// sg_field.Clear()
		// density_field.Clear()
		clear_field(mass_field)
		clear_field(sg_field)
		clear_field(density_field)
	}

	return MassDataView{mass_field, Clear}
}
