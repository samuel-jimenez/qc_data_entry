package main

import (
	"strings"

	"github.com/samuel-jimenez/windigo"
)

func build_text_dock(parent windigo.Controller, labels []string) *windigo.Panel {
	panel := windigo.NewPanel(parent)
	panel.SetSize(50, 50)

	dock := windigo.NewSimpleDock(panel)
	dock.SetMargins(10)

	for _, label_text := range labels {
		label := windigo.NewLabel(panel)
		label.SetSize(200, 25)
		label.SetText(label_text)
		dock.Dock(label, windigo.Left)

	}
	return panel
}

func build_button_dock(parent windigo.Controller, labels []string, onclicks []func()) *windigo.Panel {
	assertEqual(len(labels), len(onclicks))
	panel := windigo.NewPanel(parent)
	panel.SetSize(50, 50)

	dock := windigo.NewSimpleDock(panel)
	dock.SetMargins(10)

	for i, label_text := range labels {
		button := windigo.NewPushButton(panel)
		button.SetSize(200, 25)
		button.SetMargins(10)

		button.OnClick().Bind(func(e *windigo.Event) {
			onclicks[i]()

		})
		button.SetText(label_text)
		dock.Dock(button, windigo.Left)

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

func show_combobox(parent windigo.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *windigo.ComboBox {
	// func show_combobox(parent windigo.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *windigo.Edit {
	combobox_label := windigo.NewLabel(parent)
	combobox_label.SetPos(x_label_pos, y_pos)
	combobox_label.SetText(field_text)

	// combobox_field := combobox_label.NewEdit(mainWindow)
	// combobox_field := windigo.NewEdit(parent)

	combobox_field := windigo.NewComboBox(parent)

	combobox_field.SetPos(x_field_pos, y_pos)
	// Most Controls have default size unless SetSize is called.
	combobox_field.SetText("")
	// combobox_field.SetParent(combobox_label)
	// combobox_label.SetParent(combobox_field)

	// combobox_label.OnClick().Bind(func(e *windigo.Event) {
	// 		combobox_field.SetFocus()
	// })
	combobox_field.OnKillFocus().Bind(func(e *windigo.Event) {
		combobox_field.SetText(strings.TrimSpace(combobox_field.Text()))
	})

	return combobox_field
}

func show_edit(parent windigo.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *windigo.Edit {
	edit_label := windigo.NewLabel(parent)
	edit_label.SetPos(x_label_pos, y_pos)
	edit_label.SetText(field_text)

	// edit_field := edit_label.NewEdit(mainWindow)
	edit_field := windigo.NewEdit(parent)
	edit_field.SetPos(x_field_pos, y_pos)
	// Most Controls have default size unless SetSize is called.
	edit_field.SetText("")
	// edit_field.SetParent(edit_label)
	// edit_label.SetParent(edit_field)

	// edit_label.OnClick().Bind(func(e *windigo.Event) {
	// 		edit_field.SetFocus()
	// })
	edit_field.OnKillFocus().Bind(func(e *windigo.Event) {
		edit_field.SetText(strings.TrimSpace(edit_field.Text()))
	})

	return edit_field
}

func show_edit_with_lose_focus(parent windigo.Controller, x_label_pos, x_field_pos, y_pos int, field_text string, focus_cb func(string) string) *windigo.Edit {
	edit_label := windigo.NewLabel(parent)
	edit_label.SetPos(x_label_pos, y_pos)
	edit_label.SetText(field_text)

	// edit_field := edit_label.NewEdit(mainWindow)
	edit_field := windigo.NewEdit(parent)
	edit_field.SetPos(x_field_pos, y_pos)
	// Most Controls have default size unless SetSize is called.
	edit_field.SetText("")
	// edit_field.SetParent(edit_label)
	// edit_label.SetParent(edit_field)

	// edit_label.OnClick().Bind(func(e *windigo.Event) {
	// 		edit_field.SetFocus()
	// })
	edit_field.OnKillFocus().Bind(func(e *windigo.Event) {
		edit_field.SetText(focus_cb(strings.TrimSpace(edit_field.Text())))
	})

	return edit_field
}

func show_text(parent windigo.Controller, x_label_pos, x_field_pos, y_pos, field_width int, field_text string, field_units string) *windigo.Label {
	field_height := 20
	text_label := windigo.NewLabel(parent)
	text_label.SetPos(x_label_pos, y_pos)
	text_label.SetText(field_text)

	// text_field := text_label.NewEdit(mainWindow)
	text_field := windigo.NewLabel(parent)
	text_field.SetPos(x_field_pos, y_pos)
	text_field.SetSize(field_width, field_height)
	// Most Controls have default size unless  is called.
	text_field.SetText("0.000")

	text_units := windigo.NewLabel(parent)
	text_units.SetPos(x_field_pos+field_width, y_pos)
	text_units.SetText(field_units)

	return text_field
}

func show_mass_sg(parent windigo.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *windigo.Edit {

	field_width := 50

	density_row := 125
	sg_row := 150

	sg_text := "Specific Gravity"
	density_text := "Density"

	sg_units := "g/mL"
	density_units := "lb/gal"

	mass_field := show_edit(parent, x_label_pos, x_field_pos, y_pos, field_text)

	sg_field := show_text(parent, x_label_pos, x_field_pos, sg_row, field_width, sg_text, sg_units)
	density_field := show_text(parent, x_label_pos, x_field_pos, density_row, field_width, density_text, density_units)

	mass_field.OnKillFocus().Bind(func(e *windigo.Event) {
		mass_field.SetText(strings.TrimSpace(mass_field.Text()))
		sg := sg_from_mass(mass_field)
		density := density_from_sg(sg)
		sg_field.SetText(format_sg(sg))
		density_field.SetText(format_density(density))
	})

	return mass_field
}

func show_status_bar(message string) {
	status_bar.SetText(message)
	// TODO TIMER
}
