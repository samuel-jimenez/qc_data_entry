package main

import (
	"strings"

	"github.com/samuel-jimenez/windigo"
)


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

func show_text(parent windigo.AutoPanel, label_width, field_width, field_height int, field_text string, field_units string) *windigo.Label {
	panel := windigo.NewAutoPanel(parent)
	panel.SetSize(label_width+2*field_width, field_height)

	text_label := windigo.NewLabel(panel)
	text_label.SetSize(label_width, field_height)
	text_label.SetText(field_text)

	text_field := windigo.NewLabel(panel)
	text_field.SetSize(field_width, field_height)
	text_field.SetText("0.000")

	text_units := windigo.NewLabel(panel)
	text_units.SetSize(field_width, field_height)
	text_units.SetText(field_units)

	panel.Dock(text_label, windigo.Left)
	panel.Dock(text_field, windigo.Left)
	panel.Dock(text_units, windigo.Left)
	parent.Dock(panel, windigo.Bottom)

	return text_field
}

func show_mass_sg(parent windigo.AutoPanel, label_width, control_width, height int, field_text string) windigo.LabeledEdit {

	field_width := 50

	sg_text := "Specific Gravity"
	density_text := "Density"

	sg_units := "g/mL"
	density_units := "lb/gal"

	mass_field := show_edit(parent, label_width, control_width, height, field_text)

	//PUSH TO BOTTOM
	density_field := show_text(parent, label_width, field_width, height, density_text, density_units)
	sg_field := show_text(parent, label_width, field_width, height, sg_text, sg_units)

	mass_field.OnChange().Bind(func(e *windigo.Event) {
		start, end := mass_field.Selected()
		// IndexAny(s, chars string) int
		mass_field.SetText(strings.TrimSpace(mass_field.Text()))
		// sg := sg_from_mass(mass_field)
		// mass_field.SelectText(-1, -1)
		mass_field.SelectText(start, end)
		// mass_field.SelectText(-1, 0)

		sg := sg_from_mass(mass_field)

		density := density_from_sg(sg)
		sg_field.SetText(format_sg(sg, false))
		density_field.SetText(format_density(density))
	})

	return mass_field
}

func show_status_bar(message string) {
	status_bar.SetText(message)
	// TODO TIMER
}
