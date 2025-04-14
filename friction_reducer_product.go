package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"codeberg.org/go-pdf/fpdf"
	"github.com/samuel-jimenez/winc"
	"github.com/samuel-jimenez/winc/w32"
)

type FrictionReducerProduct struct {
	BaseProduct
	sg           float64
	string_test  float64
	viscosity    float64
	sample_point string
}

func (product FrictionReducerProduct) toAllProduct() Product {
	return Product{BaseProduct{product.product_type, product.lot_number, product.visual}, sql.NullFloat64{product.sg, true}, sql.NullFloat64{0, false}, sql.NullFloat64{product.sg * LB_PER_GAL, true}, sql.NullFloat64{product.string_test, true}, sql.NullFloat64{product.viscosity, true}, sql.NullString{product.sample_point, true}}

}

func newFrictionReducerProduct(base_product BaseProduct, sample_point string, viscosity_field *winc.Edit, mass_field *winc.Edit, string_field *winc.Edit) FrictionReducerProduct {
	viscosity, _ := strconv.ParseFloat(strings.TrimSpace(viscosity_field.Text()), 64)
	mass, _ := strconv.ParseFloat(strings.TrimSpace(mass_field.Text()), 64)
	// if !err.Error(){fmt.Println("error",err)}
	string_test, _ := strconv.ParseFloat(strings.TrimSpace(string_field.Text()), 64)
	sg := mass / SAMPLE_VOLUME

	return FrictionReducerProduct{base_product, sg, string_test, viscosity, sample_point}

}

func (product FrictionReducerProduct) get_pdf_name() string {
	// if (product.sample_point.Valid) {
	if product.sample_point != "" {

		return fmt.Sprintf("%s/%s-%s-%s.pdf", LABEL_PATH, strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.product_type)), " ", "_"), strings.ToUpper(product.lot_number), product.sample_point)
	}

	return fmt.Sprintf("%s/%s-%s.pdf", LABEL_PATH, strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.product_type)), " ", "_"), strings.ToUpper(product.lot_number))
}

func (product FrictionReducerProduct) check_data() bool {
	return true
}

func (product FrictionReducerProduct) print() error {
	var label_width, label_height,
		field_width, field_height,
		label_col,
		// field_col,
		product_row,
		density_row,
		sg_row,
		string_row,
		viscosity_row,
		lot_row float64

	label_width = 40
	label_height = 10

	field_width = 40
	field_height = 10

	label_col = 10
	// field_col = 120

	product_row = 0
	sg_row = 10
	density_row = 15
	string_row = 20
	viscosity_row = 25
	lot_row = 45

	pdf := fpdf.New("L", "mm", "A7", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(label_col, product_row)
	pdf.Cell(field_width, field_height, strings.ToUpper(product.product_type))

	pdf.SetXY(label_col, sg_row)
	pdf.Cell(label_width, label_height, "SG")
	pdf.Cell(field_width, field_height, strconv.FormatFloat(product.sg, 'f', 3, 64))

	pdf.SetXY(label_col, density_row)
	pdf.Cell(label_width, label_height, "DENSITY")
	pdf.Cell(field_width, field_height, strconv.FormatFloat(product.sg*LB_PER_GAL, 'f', 3, 64))

	pdf.SetXY(label_col, string_row)
	pdf.Cell(label_width, label_height, "STRING")
	pdf.Cell(field_width, field_height, strconv.FormatFloat(product.string_test, 'f', 0, 64))

	pdf.SetXY(label_col, viscosity_row)
	pdf.Cell(label_width, label_height, "VISCOSITY")
	pdf.Cell(field_width, field_height, strconv.FormatFloat(product.viscosity, 'f', 0, 64))

	pdf.SetXY(label_col, lot_row)
	pdf.Cell(field_width, field_height, strings.ToUpper(product.lot_number))
	pdf.CellFormat(field_width, field_height, strings.ToUpper(product.sample_point), "", 0, "R", false, 0, "")

	err := pdf.OutputFileAndClose(product.get_pdf_name())
	return err
}



// create table product_line (product_id integer not null primary key, product_name text);
func show_fr(parent winc.Controller) {

	top_col := 10
	bottom_col := 320
	// field_col := 120
	label_col := 10
	field_col := 120

	product_row := 20
	lot_row := 45
	group_row := 70
	// 		lot_row := 45
	// 		sample_row := 70
	//
	// 		visual_row := 125
	// 		viscosity_row := 150
	// 		mass_row := 175
	// 		string_row := 200
	// group_row := 120

	submit_col := 40
	submit_row := 225
	submit_button_width := 100
	submit_button_height := 40

	product_text := "Product"
	lot_text := "Lot Number"
	// visual_text := "Visual Inspection"
	// viscosity_text := "Viscosity"
	// mass_text := "Mass"
	// string_text := "String"
	// sample_text := "Sample Point"
	group_width := 300
	group_height := 120

	top_text := "Top"
	// bottom_text := "Bottom"
	bottom_text := "Btm"

	//TODO EXTRACT
	var product_id, lot_id int64

	product_field := show_edit(parent, label_col, field_col, product_row, product_text)
	product_field.OnKillFocus().Bind(func(e *winc.Event) {
		product_field.SetText(strings.ToUpper(strings.TrimSpace(product_field.Text())))
		if product_field.Text() != "" {
			product_id = insel_product_id(product_field.Text())
			fmt.Println("product_id", product_id)
		}
	})

	lot_field := show_edit_with_lose_focus(parent, label_col, field_col, lot_row, lot_text, strings.ToUpper)
	lot_field.OnKillFocus().Bind(func(e *winc.Event) {
		lot_field.SetText(strings.ToUpper(strings.TrimSpace(lot_field.Text())))
		if lot_field.Text() != "" && product_field.Text() != "" {
			lot_id = insel_lot_id(lot_field.Text(), product_id)
			fmt.Println("lot_id", lot_id)
		}
	})

	top_group_cb := show_fr_sample_group(parent, top_text, top_col, group_row, group_width, group_height)
	// show_fr_sample_group(parent, top_text, top_col, group_row, group_width, group_height)

	// bottom_group := show_fr_sample_group(parent, bottom_text, bottom_col, group_row, group_width, group_height, top_group)
	// show_fr_sample_group(parent, bottom_text, bottom_col, group_row, group_width, group_height, top_group)
	bottom_group_cb := show_fr_sample_group(parent, bottom_text, bottom_col, group_row, group_width, group_height)

	// string_field := show_edit(mainWindow, label_col, field_col, string_row, string_text)
	// mass_field := show_edit(mainWindow, label_col, field_col, mass_row, mass_text)
	// viscosity_field := show_edit(mainWindow, label_col, field_col, viscosity_row, viscosity_text)
	// visual_field := show_checkbox(mainWindow, label_col, field_col, visual_row, visual_text)
	// sample_field := show_edit(mainWindow, label_col, field_col, sample_row, sample_text)

	// lot_field :=
	// show_edit(mainWindow, label_col, field_col, lot_row, lot_text)
	// // product_field :=
	// show_edit(mainWindow, label_col, field_col, product_row, product_text)

	//

	// 	product_row := 20
	// product_text := "Product"
	// product_field := show_edit(mainWindow, label_col, field_col, product_row, product_text)

	submit_button := winc.NewPushButton(parent)

	submit_button.SetText("Submit")
	submit_button.SetPos(submit_col, submit_row) // (x, y)
	// submit_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	submit_button.SetSize(submit_button_width, submit_button_height) // (width, height)
	submit_button.OnClick().Bind(func(e *winc.Event) {

		base_product := newProduct_0(product_field, lot_field)
		top_product := top_group_cb(base_product)
		bottom_product := bottom_group_cb(base_product)
		fmt.Println("top", top_product)
		fmt.Println("btm", bottom_product)
		if top_product.check_data() {
			fmt.Println("data", top_product)
			top_product.print()
			top_product.toAllProduct().save()
		}
		if bottom_product.check_data() {
			fmt.Println("data", bottom_product)
			bottom_product.print()
		}
	})

	top_button_col := 150
	top_button_row := 225
	top_button_width := 100
	top_button_height := 40
	top_button := winc.NewPushButton(parent)

	top_button.SetText("Accept Top")
	top_button.SetPos(top_button_col, top_button_row) // (x, y)
	// top_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	top_button.SetSize(top_button_width, top_button_height) // (width, height)
	top_button.OnClick().Bind(func(e *winc.Event) {

		base_product := newProduct_0(product_field, lot_field)
		top_product := top_group_cb(base_product)
		top_product.sample_point = ""
		if top_product.check_data() {
			fmt.Println("data", top_product)
			top_product.print()
		}

	})

	btm_button_col := 250
	btm_button_row := 225
	btm_button_width := 100
	btm_button_height := 40
	btm_button := winc.NewPushButton(parent)

	btm_button.SetText("Accept Btm")
	btm_button.SetPos(btm_button_col, btm_button_row) // (x, y)
	// btm_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	btm_button.SetSize(btm_button_width, btm_button_height) // (width, height)
	btm_button.OnClick().Bind(func(e *winc.Event) {

		base_product := newProduct_0(product_field, lot_field)
		bottom_product := bottom_group_cb(base_product)
		bottom_product.sample_point = ""
		if bottom_product.check_data() {
			fmt.Println("data", bottom_product)
			bottom_product.print()
		}
	})

}

// func show_fr_sample_group(parent winc.Controller, sample_point string, x_pos, y_pos, group_width, group_height int) winc.Controller {

func show_fr_sample_group(parent winc.Controller, sample_point string, x_pos, y_pos, group_width, group_height int) func(base_product BaseProduct) FrictionReducerProduct {

	// func show_fr_sample_group(parent winc.Controller, sample_point string, x_pos, y_pos, group_width, group_height int, after winc.Controller) winc.Controller {

	sample_group := winc.NewPanel(parent)
	sample_group.SetAndClearStyleBits(w32.WS_TABSTOP, 0)
	sample_group.SetPos(x_pos, y_pos)
	sample_group.SetSize(group_width, group_height)
	sample_group.SetText(sample_point)

	bottom_group := winc.NewGroupBox(parent)
	bottom_group.SetPos(x_pos-5, y_pos-5)
	bottom_group.SetSize(group_width+10, group_height+10)
	bottom_group.SetText(sample_point)

	return show_fr_sample(sample_group, sample_point)
	// return sample_group

}

func show_fr_sample(parent winc.Controller, sample_point string) func(base_product BaseProduct) FrictionReducerProduct {
	label_col := 10
	field_col := 120

	visual_row := 25
	viscosity_row := 50
	mass_row := 75
	string_row := 100

	// group_row := 120

	visual_text := "Visual Inspection"
	viscosity_text := "Viscosity"
	mass_text := "Mass"
	string_text := "String"

	visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)

	viscosity_field := show_edit(parent, label_col, field_col, viscosity_row, viscosity_text)
	mass_field := show_edit(parent, label_col, field_col, mass_row, mass_text)
	string_field := show_edit(parent, label_col, field_col, string_row, string_text)

	// parent.Bind(w32.WM_COPYDATA, func(arg *EventArg) {
	// 	sender := arg.Sender()
	// 	if data, ok := arg.Data().(*gform.RawMsg); ok {
	// 		println(data.Hwnd, data.Msg, data.WParam, data.LParam)
	// 	}
	// }

	visual_field.SetFocus()

	return func(base_product BaseProduct) FrictionReducerProduct {
		base_product.visual = visual_field.Checked()
		return newFrictionReducerProduct(base_product, sample_point, viscosity_field, mass_field, string_field)
	}

}
