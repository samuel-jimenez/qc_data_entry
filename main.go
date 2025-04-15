package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/samuel-jimenez/winc"
	"github.com/samuel-jimenez/winc/w32"
)

var LABEL_PATH = "C:/Users/QC Lab/Documents/golang/qc_data_entry/labels"

var SUBMIT_ROW = 200

type BaseProduct struct {
	product_type string
	lot_number   string
	sample_point string
	visual       bool
	product_id   int64
	lot_id       int64
}

// TODO fix this
func newProduct_0(product_field winc.Controller, lot_field *winc.Edit) BaseProduct {
	return BaseProduct{strings.ToUpper(product_field.Text()), strings.ToUpper(lot_field.Text()), "", false, -1, -1}
}

func newProduct_1(product_field *winc.Edit, lot_field *winc.Edit,
	visual_field *winc.CheckBox) BaseProduct {
	return BaseProduct{strings.ToUpper(product_field.Text()), strings.ToUpper(lot_field.Text()), "", visual_field.Checked(), -1, -1}
}

// func newProduct_2(product_field *winc.Edit, lot_field *winc.Edit) BaseProduct {
// 	return BaseProduct{strings.ToUpper(product_field.Text()), strings.ToUpper(lot_field.Text()), false, -1, -1}.insel_all()
// }

func newProduct_3(product_field winc.Controller, lot_field winc.Controller, sample_field winc.Controller) BaseProduct {
	return BaseProduct{strings.ToUpper(product_field.Text()), strings.ToUpper(lot_field.Text()), strings.ToUpper(sample_field.Text()), false, -1, -1}
}

func (product BaseProduct) toBaseProduct() BaseProduct {
	return product
}

func (product BaseProduct) get_pdf_name() string {
	// if (product.sample_point.Valid) {
	if product.sample_point != "" {
		return fmt.Sprintf("%s/%s-%s-%s.pdf", LABEL_PATH, strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.product_type)), " ", "_"), strings.ToUpper(product.lot_number), product.sample_point)
	}

	return fmt.Sprintf("%s/%s-%s.pdf", LABEL_PATH, strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.product_type)), " ", "_"), strings.ToUpper(product.lot_number))
}

func (product *BaseProduct) insel_product_id(product_name string) {
	log.Println("insel_product_id product_name", product_name)

	product.product_id = insel_product_id(product_name)
	log.Println("insel_product_id", product)
	// return product

}

func (product *BaseProduct) insel_lot_id(lot_name string) {
	product.lot_id = insel_lot_id(lot_name, product.product_id)
	log.Println("lot_id", product.lot_id)

}

// func (product *BaseProduct) insel_product_self() {
func (product *BaseProduct) insel_product_self() *BaseProduct {
	product.insel_product_id(product.product_type)
	return product

}

func (product *BaseProduct) insel_lot_self() *BaseProduct {
	product.insel_lot_id(product.lot_number)
	log.Println("insel_lot_self", product.lot_id)
	return product

}

// TODO use this one
// func (product BaseProduct) insel_all() BaseProduct {
func (product BaseProduct) insel_all() *BaseProduct {
	// func (product *BaseProduct) insel_all() *BaseProduct {
	// func (product *BaseProduct) insel_all() {
	return product.insel_product_self().insel_lot_self()

}

func (product *BaseProduct) copy_ids(product_lot BaseProduct) {
	if product_lot.product_id > 0 {
		product.product_id = product_lot.product_id
		if product_lot.lot_id > 0 {
			product.lot_id = product_lot.lot_id
		} else {
			product.insel_lot_self()
		}
	} else {
		product.insel_all()
	}
	log.Println("copy_ids lot_id", product.lot_id)

}

func (product BaseProduct) toProduct() Product {
	return Product{product, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}}
	//TODO Option?
	// NullFloat64

}

func main() {

	//log to file
	log_file, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer log_file.Close()

	log.SetOutput(log_file)
	log.Println("This is a test log entry")

	//open_db
	// qc_db, err := sql.Open("sqlite3", DB_PATH)
	qc_db, err := sql.Open("sqlite3", ":memory:")
	qc_db.Exec("attach ? as 'bs'", DB_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer qc_db.Close()

	dbinit(qc_db)

	// tx, err := qc_db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// for i := 0; i < 100; i++ {
	// 	_, err = stmt.Exec(i, fmt.Sprintf("こんにちは世界%03d", i))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// err = tx.Commit()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// rows, err := qc_db.Query("select id, name from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	err = rows.Scan(&id, &name)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Println(id, name)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	show_window()
}

// TDFO add clear

func show_window() {

	log.Println("Process started")
	// DEBUG
	// log.Println(time.Now().UTC().UnixNano())

	mainWindow := winc.NewForm(nil)
	mainWindow.SetSize(800, 600) // (width, height)
	mainWindow.SetText("QC Data Entry")

	dock := winc.NewSimpleDock(mainWindow)

	parent := winc.NewPanel(mainWindow)
	parent.SetSize(750, 100)

	label_col := 10
	field_col := 120

	product_row := 20
	lot_row := 45
	// 		lot_row := 45
	sample_row := 70
	//
	// 		visual_row := 125
	// 		viscosity_row := 150
	// 		mass_row := 175
	// 		string_row := 200
	// group_row := 120

	product_text := "Product"
	lot_text := "Lot Number"
	sample_text := "Sample Point"

	// var product_id, lot_id int64
	var product_lot BaseProduct

	//TODO InsertItem NewComboBox ComboBox
	// db_select_all_product
	// product_field := show_edit(parent, label_col, field_col, product_row, product_text)
	product_field := show_combobox(parent, label_col, field_col, product_row, product_text)
	product_field.OnKillFocus().Bind(func(e *winc.Event) {
		product_field.SetText(strings.ToUpper(strings.TrimSpace(product_field.Text())))
		if product_field.Text() != "" {
			product_lot.insel_product_id(product_field.Text())
			log.Println("product_field started", product_lot)

		}
	})

	lot_field := show_edit_with_lose_focus(parent, label_col, field_col, lot_row, lot_text, strings.ToUpper)
	lot_field.OnKillFocus().Bind(func(e *winc.Event) {
		lot_field.SetText(strings.ToUpper(strings.TrimSpace(lot_field.Text())))
		if lot_field.Text() != "" && product_field.Text() != "" {
			product_lot.insel_lot_id(lot_field.Text())
		}
	})

	sample_field := show_edit(mainWindow, label_col, field_col, sample_row, sample_text)

	new_product_cb := func() BaseProduct {
		// return newProduct_0(product_field, lot_field).copy_ids(product_lot)
		log.Println("product_field new_product_cb", product_lot)

		// base_product := newProduct_0(product_field, lot_field)
		base_product := newProduct_3(product_field, lot_field, sample_field)

		base_product.copy_ids(product_lot)
		log.Println("base_product new_product_cb", base_product)

		return base_product
	}

	tabs := winc.NewTabView(mainWindow)
	// tabs.SetPos(20, 20)
	// tabs.SetSize(750, 500)
	tab_wb := tabs.AddPanel("Water Based")
	tab_oil := tabs.AddPanel("Oil Based")
	tab_fr := tabs.AddPanel("Friction Reducer")

	show_water_based(tab_wb, new_product_cb)
	show_oil_based(tab_oil, new_product_cb)
	show_fr(tab_fr, new_product_cb)

	// dock.Dock(quux, winc.Top)        // toolbars always dock to the top
	dock.Dock(parent, winc.Top)         // tabs should prefer docking at the top
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

func NewComboBox(parent winc.Controller) *winc.ComboBox {
	cb := new(winc.ComboBox)

	cb.InitControl("COMBOBOX", parent, 0, w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.WS_VSCROLL|w32.CBS_DROPDOWN|w32.CBS_UPPERCASE)
	// cb.InitControl("COMBOBOX", parent, 0, w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.WS_VSCROLL|w32.CBS_DROPDOWN)
	winc.RegMsgHandler(cb)

	cb.SetFont(winc.DefaultFont)
	cb.SetSize(200, 400)
	return cb
}

// TODO InsertItem NewComboBox ComboBox
func show_combobox(parent winc.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *winc.ComboBox {
	// func show_combobox(parent winc.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *winc.Edit {
	combobox_label := winc.NewLabel(parent)
	combobox_label.SetPos(x_label_pos, y_pos)
	combobox_label.SetText(field_text)

	// combobox_field := combobox_label.NewEdit(mainWindow)
	// combobox_field := winc.NewEdit(parent)

	combobox_field := winc.NewComboBox(parent)

	rows, err := db_select_all_product.Query()
	if err != nil {
		log.Printf("%q: %s\n", err, "insel_lot_id")
		// return -1
	}
	for rows.Next() {
		var (
			id    uint8
			value string
		)

		if error := rows.Scan(&id, &value); error == nil {
			// data[id] = value
			combobox_field.AddItem(value)
		}
	}

	combobox_field.SetPos(x_field_pos, y_pos)
	// Most Controls have default size unless SetSize is called.
	combobox_field.SetText("")
	// combobox_field.SetParent(combobox_label)
	// combobox_label.SetParent(combobox_field)

	// combobox_label.OnClick().Bind(func(e *winc.Event) {
	// 		combobox_field.SetFocus()
	// })
	combobox_field.OnKillFocus().Bind(func(e *winc.Event) {
		combobox_field.SetText(strings.TrimSpace(combobox_field.Text()))
	})

	return combobox_field
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
	edit_field.OnKillFocus().Bind(func(e *winc.Event) {
		edit_field.SetText(strings.TrimSpace(edit_field.Text()))
	})

	return edit_field
}

func show_edit_with_lose_focus(parent winc.Controller, x_label_pos, x_field_pos, y_pos int, field_text string, focus_cb func(string) string) *winc.Edit {
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
	edit_field.OnKillFocus().Bind(func(e *winc.Event) {
		edit_field.SetText(focus_cb(strings.TrimSpace(edit_field.Text())))
	})

	return edit_field
}

func show_text(parent winc.Controller, x_label_pos, x_field_pos, y_pos, field_width int, field_text string, field_units string) *winc.Label {
	field_height := 20
	text_label := winc.NewLabel(parent)
	text_label.SetPos(x_label_pos, y_pos)
	text_label.SetText(field_text)

	// text_field := text_label.NewEdit(mainWindow)
	text_field := winc.NewLabel(parent)
	text_field.SetPos(x_field_pos, y_pos)
	text_field.SetSize(field_width, field_height)
	// Most Controls have default size unless  is called.
	text_field.SetText("0.000")

	text_units := winc.NewLabel(parent)
	text_units.SetPos(x_field_pos+field_width, y_pos)
	text_units.SetText(field_units)

	return text_field
}

func show_mass_sg(parent winc.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *winc.Edit {

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

	mass_field.OnKillFocus().Bind(func(e *winc.Event) {
		mass_field.SetText(strings.TrimSpace(mass_field.Text()))
		sg := sg_from_mass(mass_field)
		density := density_from_sg(sg)
		sg_field.SetText(format_sg(sg))
		density_field.SetText(format_density(density))
	})

	return mass_field
}

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}
