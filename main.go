package main

import (
	"fmt"
	"strings"

	"github.com/samuel-jimenez/winc"
)

var SAMPLE_VOLUME = 83.2
var LB_PER_GAL = 8.345 // g/mL
var LABEL_PATH="C:/Users/QC Lab/Documents/golang/qc_data_entry/labels"

type Product struct {
	product_type string
	lot_number   string
	visual       bool
}

func newProduct_0(product_field *winc.Edit, lot_field *winc.Edit) Product {
	return Product{strings.ToUpper(product_field.Text()), strings.ToUpper(lot_field.Text()), false}
}

func newProduct_1(product_field *winc.Edit, lot_field *winc.Edit,
		visual_field *winc.CheckBox) Product {
			return Product{strings.ToUpper(product_field.Text()), strings.ToUpper(lot_field.Text()), visual_field.Checked()}
		}



func (product Product) get_pdf_name() string {
	return fmt.Sprintf("%s/%s-%s.pdf",LABEL_PATH, strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.product_type)), " ", "_"), strings.ToUpper(product.lot_number))
}



func main() {
	show_window()
}
func show_window() {

	fmt.Println("Process started")
	// DEBUG

	mainWindow := winc.NewForm(nil)
	mainWindow.SetSize(800, 600) // (width, height)
	mainWindow.SetText("QC Data Entry")


	dock := winc.NewSimpleDock(mainWindow)

	tabs := winc.NewTabView(mainWindow)
	// tabs.SetPos(20, 20)
	// tabs.SetSize(750, 500)
	tab_wb := tabs.AddPanel("Water Based")
	tab_oil := tabs.AddPanel("Oil Based")
	tab_fr := tabs.AddPanel("Friction Reducer")

	show_water_based(tab_wb)
	show_oil_based(tab_oil)
	show_fr(tab_fr)

	// dock.Dock(quux, winc.Top)        // toolbars always dock to the top
	dock.Dock(tabs, winc.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), winc.Fill) // tab panels dock just below tabs and fill area

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)
	mainWindow.RunMainLoop() // Must call to start event loop.
}

func show_checkbox(parent winc.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *winc.CheckBox {
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

func show_edit(parent winc.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *winc.Edit {
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
