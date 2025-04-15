package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/samuel-jimenez/winc"
)

var SAMPLE_VOLUME = 83.2
var LB_PER_GAL = 8.345 // g/mL
var LABEL_PATH = "C:/Users/QC Lab/Documents/golang/qc_data_entry/labels"

type BaseProduct struct {
	product_type string
	lot_number   string
	sample_point string
	visual       bool
	product_id   int64
	lot_id       int64
}

// TODO fix this
func newProduct_0(product_field *winc.Edit, lot_field *winc.Edit) BaseProduct {
	return BaseProduct{strings.ToUpper(product_field.Text()), strings.ToUpper(lot_field.Text()), "", false, -1, -1}
}

func newProduct_1(product_field *winc.Edit, lot_field *winc.Edit,
	visual_field *winc.CheckBox) BaseProduct {
	return BaseProduct{strings.ToUpper(product_field.Text()), strings.ToUpper(lot_field.Text()), "", visual_field.Checked(), -1, -1}
}

// func newProduct_2(product_field *winc.Edit, lot_field *winc.Edit) BaseProduct {
// 	return BaseProduct{strings.ToUpper(product_field.Text()), strings.ToUpper(lot_field.Text()), false, -1, -1}.insel_all()
// }

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
	fmt.Println("insel_product_id product_name", product_name)

	product.product_id = insel_product_id(product_name)
	fmt.Println("insel_product_id", product)
	// return product

}

func (product *BaseProduct) insel_lot_id(lot_name string) {
	product.lot_id = insel_lot_id(lot_name, product.product_id)
	fmt.Println("lot_id", product.lot_id)

}

// func (product *BaseProduct) insel_product_self() {
func (product *BaseProduct) insel_product_self() *BaseProduct {
	product.insel_product_id(product.product_type)
	return product

}

func (product *BaseProduct) insel_lot_self() *BaseProduct {
	product.insel_lot_id(product.lot_number)
	fmt.Println("insel_lot_self", product.lot_id)
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
	fmt.Println("copy_ids lot_id", product.lot_id)

}

func (product BaseProduct) toProduct() Product {
	return Product{product, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}}
	//TODO Option?
	// NullFloat64

}

func main() {
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
	// 	fmt.Println(id, name)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	show_window()
}

// TDFO add clear

func show_window() {

	fmt.Println("Process started")
	// DEBUG
	// fmt.Println(time.Now().UTC().UnixNano())

	mainWindow := winc.NewForm(nil)
	mainWindow.SetSize(800, 600) // (width, height)
	mainWindow.SetText("QC Data Entry")

	dock := winc.NewSimpleDock(mainWindow)

	parent := winc.NewPanel(mainWindow)

	label_col := 10
	field_col := 120

	product_row := 20
	lot_row := 45
	// 		lot_row := 45
	// 		sample_row := 70
	//
	// 		visual_row := 125
	// 		viscosity_row := 150
	// 		mass_row := 175
	// 		string_row := 200
	// group_row := 120

	product_text := "Product"
	lot_text := "Lot Number"
	// sample_text := "Sample Point"

	// var product_id, lot_id int64
	var product_lot BaseProduct

	//TODO InsertItem NewComboBox ComboBox
	product_field := show_edit(parent, label_col, field_col, product_row, product_text)
	product_field.OnKillFocus().Bind(func(e *winc.Event) {
		product_field.SetText(strings.ToUpper(strings.TrimSpace(product_field.Text())))
		if product_field.Text() != "" {
			product_lot.insel_product_id(product_field.Text())
			fmt.Println("product_field started", product_lot)

		}
	})

	lot_field := show_edit_with_lose_focus(parent, label_col, field_col, lot_row, lot_text, strings.ToUpper)
	lot_field.OnKillFocus().Bind(func(e *winc.Event) {
		lot_field.SetText(strings.ToUpper(strings.TrimSpace(lot_field.Text())))
		if lot_field.Text() != "" && product_field.Text() != "" {
			product_lot.insel_lot_id(lot_field.Text())
		}
	})

	new_product_cb := func() BaseProduct {
		// return newProduct_0(product_field, lot_field).copy_ids(product_lot)
		fmt.Println("product_field new_product_cb", product_lot)

		base_product := newProduct_0(product_field, lot_field)
		base_product.copy_ids(product_lot)
		fmt.Println("base_product new_product_cb", base_product)

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

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}
