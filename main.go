package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"codeberg.org/go-pdf/fpdf"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/samuel-jimenez/winc"
)

var SAMPLE_VOLUME = 83.2
var LB_PER_GAL = 8.345 // g/mL
var LABEL_PATH = "C:/Users/QC Lab/Documents/golang/qc_data_entry/labels"
var DB_PATH = "C:/Users/QC Lab/Documents/golang/qc_data_entry/qc.sqlite3"

var qc_db *sql.DB
var db_get_product, db_insert_product,
	db_get_lot, db_insert_lot *sql.Stmt
var err error

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
	return fmt.Sprintf("%s/%s-%s.pdf", LABEL_PATH, strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.product_type)), " ", "_"), strings.ToUpper(product.lot_number))
}

type AllProduct struct {
	Product
	sg           sql.NullFloat64
	ph           sql.NullFloat64
	density      sql.NullFloat64
	string_test  sql.NullFloat64
	viscosity    sql.NullFloat64
	sample_point sql.NullString
}

func (product Product) toAllProduct() AllProduct {
	return AllProduct{product, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullString{"", false}}
	//TODO Option?
	// NullFloat64

}

func (product AllProduct) save() error {

	stmt, err := qc_db.Prepare(`insert into
qc_samples (qc_id, batch_id, sample_point text, time_stamp integer, specific_gravity real,  ph real,   string_test real,   viscosity real, );
foo(id, name)
values(?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	// for i := 0; i < 100; i++ {
	// 	_, err = stmt.Exec(i, fmt.Sprintf("こんにちは世界%03d", i))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		// 	}
	// 	}
	// }
	return err
}

func (product AllProduct) print() error {
	var label_width, label_height,
		field_width, field_height,
		label_col,
		// field_col,
		product_row,
		curr_row,
		curr_row_delta,
		lot_row float64

	label_width = 40
	label_height = 10

	field_width = 40
	field_height = 10

	label_col = 10
	// field_col = 120

	product_row = 0
	lot_row = 45

	pdf := fpdf.New("L", "mm", "A7", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(label_col, product_row)
	pdf.Cell(field_width, field_height, strings.ToUpper(product.product_type))

	if product.density.Valid {
		curr_row = 5
		curr_row_delta = 5

	} else {
		curr_row = 10
		curr_row_delta = 10

	}
	//TODO unit

	if product.sg.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "SG")
		pdf.Cell(field_width, field_height, strconv.FormatFloat(product.sg.Float64, 'f', 4, 64))
	}

	if product.ph.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "pH")
		pdf.Cell(field_width, field_height, strconv.FormatFloat(product.ph.Float64, 'f', 2, 64))
	}

	if product.density.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "DENSITY")
		pdf.Cell(field_width, field_height, strconv.FormatFloat(product.density.Float64, 'f', 3, 64))
	}

	if product.string_test.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "STRING")
		pdf.Cell(field_width, field_height, strconv.FormatFloat(product.string_test.Float64, 'f', 3, 64))
	}
	if product.viscosity.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "VISCOSITY")
		pdf.Cell(field_width, field_height, strconv.FormatFloat(product.viscosity.Float64, 'f', 3, 64))
	}

	fmt.Println(curr_row)

	pdf.SetXY(label_col, lot_row)
	pdf.Cell(field_width, field_height, strings.ToUpper(product.lot_number))
	pdf.CellFormat(field_width, field_height, strings.ToUpper(product.sample_point.String), "", 0, "R", false, 0, "")

	fmt.Println("saving to: ", product.get_pdf_name())
	err := pdf.OutputFileAndClose(product.get_pdf_name())
	return err
}

func main() {
	//open_db
	qc_db, err := sql.Open("sqlite3", DB_PATH)
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
func dbinit(db *sql.DB) {

	sqlStmt := `
PRAGMA foreign_keys = ON;
create table bs.product_line (product_id integer not null primary key, product_name text unique);
create table bs.product_lot (lot_id integer not null primary key, lot_name text, product_id references product, unique (lot_name,product_id));
create table bs.qc_samples (qc_id integer not null primary key, lot_id references product_lot, sample_point text, time_stamp integer, specific_gravity real,  ph real,   string_test real,   viscosity real);
`

	// 	sqlStmt := `
	// PRAGMA foreign_keys = ON;
	// create table product_lot (lot_id integer not null primary key, lot_name text, product_id references product);
	// create table qc_samples (qc_id integer not null primary key, lot_id references product_lot, sample_point text, time_stamp integer, specific_gravity real,  ph real,   string_test real,   viscosity real);
	// `

	/*
	   	sqlStmt := `
	   drop table qc_samples
	   create table qc_samples (qc_id integer not null primary key, batch_id references product_batch, sample_point text, time_stamp integer, specific_gravity real,  ph real,   string_test real,   viscosity real, );
	   `*/
	// foreign key(trackartist) references product(product_id)
	// references product
	// recipe
	//qc_values

	db.Exec(sqlStmt)

	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Printf("%q: %s\n", err, sqlStmt)
	// 	return
	// }

	get_product_statement := `select product_id from bs.product_line where product_name = ?`
	db_get_product, err = db.Prepare(get_product_statement)
	if err != nil {
		log.Printf("%q: %s\n", err, get_product_statement)
		return
	}

	insert_product_statement := `insert into bs.product_line (product_name) values (?) returning product_id`
	db_insert_product, err = db.Prepare(insert_product_statement)
	if err != nil {
		log.Printf("%q: %s\n", err, insert_product_statement)
		return
	}

	get_lot_statement := `select lot_id from bs.product_lot join product_line using (product_id) where lot_name = ? and product_name = ?`
	db_get_lot, err = db.Prepare(get_lot_statement)
	if err != nil {
		log.Printf("%q: %s\n", err, get_lot_statement)
		return
	}

	insert_lot_statement := `insert into bs.product_lot (lot_name) values (?) returning lot_id`
	db_insert_lot, err = db.Prepare(insert_lot_statement)
	if err != nil {
		log.Printf("%q: %s\n", err, insert_lot_statement)
		return
	}
}

func show_window() {

	fmt.Println("Process started")
	// DEBUG
	// fmt.Println(time.Now().UTC().UnixNano())

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
