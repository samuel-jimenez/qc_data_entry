package main

import (
	"strconv"
	"strings"

	"github.com/samuel-jimenez/windigo"
	"github.com/samuel-jimenez/windigo/w32"
)

var (
	GROUP_WIDTH  = 300
	GROUP_HEIGHT = 170
	GROUP_MARGIN = 5

	LABEL_WIDTH  = 100
	FIELD_WIDTH  = 200
	FIELD_HEIGHT = 28

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
		combobox_field.SetText(strings.TrimSpace(combobox_field.Text()))
	})

	return combobox_field
}

func show_edit(parent windigo.Controller, label_width, control_width, height int, field_text string) windigo.LabeledEdit {
	noop := func(str string) string { return str }
	return show_edit_with_lose_focus(parent, label_width, control_width, height, field_text, noop)
}

func show_edit_with_lose_focus(parent windigo.Controller, label_width, control_width, height int, field_text string, focus_cb func(string) string) windigo.LabeledEdit {

	edit_field := windigo.NewLabeledEdit(parent, label_width, control_width, height, field_text)

	edit_field.OnKillFocus().Bind(func(e *windigo.Event) {
		edit_field.SetText(focus_cb(strings.TrimSpace(edit_field.Text())))
	})

	return edit_field
}

func check_or_error(field windigo.Bordered, test bool) {
	if test {
		field.SetBorder(okPen)
	} else {
		field.SetBorder(erroredPen)
	}
}

func show_number_edit(parent windigo.Controller, label_width, control_width, height int, field_text string, range_field *RangeROView) windigo.LabeledEdit {

	edit_field := windigo.NewLabeledEdit(parent, label_width, control_width, height, field_text)
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

	return windigo.LabeledEdit{panel, text_field}
}

type DerivedMassRangesView struct {
	mass_field,
	sg_field,
	density_field *RangeROView
}

func (view DerivedMassRangesView) CheckMass(data float64) bool {
	return view.mass_field.Check(data)
}

func (view DerivedMassRangesView) CheckSG(data float64) bool {
	return view.sg_field.Check(data)
}

func (view DerivedMassRangesView) CheckDensity(data float64) bool {
	return view.density_field.Check(data)
}

type MassRangesView interface {
	CheckMass(data float64) bool
	CheckSG(data float64) bool
	CheckDensity(data float64) bool
}

func (data_view DerivedMassRangesView) Clear() {
	data_view.mass_field.Clear()
	data_view.sg_field.Clear()
	data_view.density_field.Clear()
}

type MassDataView struct {
	windigo.LabeledEdit
	Clear func()
}

func clear_field(field windigo.LabeledEdit) {
	field.SetText("")
	field.SetBorder(okPen)
}

// type MassDataView interface {
// 	windigo.LabeledEdit
// 	Clear()
// }

// func (data_view MassDataView) Clear() {
// 	data_view.mass_field.Clear()
// 	data_view.sg_field.Clear()
// 	data_view.density_field.Clear()
// }

// func (data_view windigo.LabeledEdit) Clear() {
// 	data_view.mass_field.Clear()
// 	data_view.sg_field.Clear()
// 	data_view.density_field.Clear()
// }

// func show_mass_sg(parent windigo.AutoPanel, label_width, control_width, height int, field_text string, ranges_panel MassRangesView) windigo.LabeledEdit {
func show_mass_sg(parent windigo.AutoPanel, label_width, control_width, height int, field_text string, ranges_panel MassRangesView) MassDataView {

	field_width := 60

	sg_text := "Specific Gravity"
	density_text := "Density"

	sg_units := "g/mL"
	density_units := "lb/gal"

	mass_field := show_edit(parent, label_width, control_width, height, field_text)
	// mass_field := show_number_edit(parent, label_width, control_width, height, field_text, ranges_panel.mass_field)
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

func show_status_bar(message string) {
	status_bar.SetText(message)
	// TODO TIMER
}
