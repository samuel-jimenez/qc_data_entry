package main

import (
	"fmt"
	"strings"

	"github.com/samuel-jimenez/winc"


)


var SAMPLE_VOLUME = 83.2
var LB_PER_GAL = 8.345 // g/mL

type Product struct {
	product_type string
	lot_number   string
	visual       bool
}

func (product Product) get_pdf_name() string {
	return strings.TrimSpace(product.product_type) + "-" + product.lot_number + ".pdf"
}


func main() {
	show_window()
}
func show_window() {

	fmt.Println("Process started")
	// DEBUG

	mainWindow := winc.NewForm(nil)
	mainWindow.SetSize(400, 300) // (width, height)
	mainWindow.SetText("QC Data Entry")
	// show_water_based(mainWindow)
	show_fr(mainWindow)
	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)
	mainWindow.RunMainLoop() // Must call to start event loop.
}


func show_checkbox(parent *winc.Form, x_label_pos, x_field_pos, y_pos int, field_text string) *winc.CheckBox {
	checkbox_label := winc.NewLabel(parent)
	checkbox_label.SetPos(x_label_pos, y_pos)

	checkbox_label.SetText(field_text)

	checkbox_field := winc.NewCheckBox(parent)
	checkbox_field.SetText("")

	checkbox_field.SetPos(x_field_pos, y_pos)
	// visual_label.OnClick().Bind(func(e *winc.Event) {
	// 		visual_field.SetFocus()
	// })
	return checkbox_field
}

func show_edit(parent *winc.Form, x_label_pos, x_field_pos, y_pos int, field_text string) *winc.Edit {
	edit_label := winc.NewLabel(parent)
	edit_label.SetPos(x_label_pos, y_pos)
	edit_label.SetText(field_text)

	// edit_field := edit_label.NewEdit(mainWindow)
	edit_field := winc.NewEdit(parent)
	edit_field.SetPos(x_field_pos, y_pos)
	// Most Controls have default size unless SetSize is called.
	edit_field.SetText("")
	// edit_field.SetParent(edit_label)
	// edit_label.SetParent(edit_field)

	// edit_label.OnClick().Bind(func(e *winc.Event) {
	// 		edit_field.SetFocus()
	// })
	return edit_field
}

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}
