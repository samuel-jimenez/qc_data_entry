package main

import (
	"fmt"
	"strconv"
	"strings"

	"codeberg.org/go-pdf/fpdf"
	"github.com/samuel-jimenez/winc"
)


type WaterBasedProduct struct {
	base_product Product
	sg           float64
	ph           float64
}

func newWaterBasedProduct(product_field *winc.Edit, lot_field *winc.Edit, visual_field *winc.CheckBox, sg_field *winc.Edit, ph_field *winc.Edit) WaterBasedProduct {
	base_product := Product{product_field.Text(), lot_field.Text(), visual_field.Checked()}
	sg, _ := strconv.ParseFloat(sg_field.Text(), 64)
	// if !err.Error(){fmt.Println("error",err)}
	ph, _ := strconv.ParseFloat(ph_field.Text(), 64)

	return WaterBasedProduct{base_product, sg, ph}

}

func (product WaterBasedProduct) get_pdf_name() string {
	return product.base_product.get_pdf_name()
}


func (product WaterBasedProduct) check_data() bool {
	return true
}

func (product WaterBasedProduct) print() error {
	var label_width, label_height,
	field_width, field_height,
	label_col,
	// field_col,
	product_row,
	sg_row,
	ph_row,
	lot_row float64

	label_width = 50
	label_height = 10

	field_width = 50
	field_height = 10

	label_col = 10
	// field_col = 120

	product_row = 0
	sg_row = 20
	ph_row = 30
	lot_row = 45

	pdf := fpdf.New("L", "mm", "A7", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(label_col, product_row)
	pdf.Cell(field_width, field_height, strings.ToUpper(product.base_product.product_type))

	pdf.SetXY(label_col, sg_row)
	pdf.Cell(label_width, label_height, "SG")
	pdf.Cell(field_width, field_height, strconv.FormatFloat(product.sg, 'f', 5, 64))

	pdf.SetXY(label_col, ph_row)
	pdf.Cell(label_width, label_height, "pH")
	pdf.Cell(field_width, field_height, strconv.FormatFloat(product.ph, 'f', 2, 64))

	pdf.SetXY(label_col, lot_row)
	pdf.Cell(field_width, field_height, strings.ToUpper(product.base_product.lot_number))

	err := pdf.OutputFileAndClose(product.get_pdf_name())
	return err
}


func show_water_based(mainWindow *winc.Form) {
	label_col := 10
	field_col := 120

	product_row := 20
	lot_row := 45

	visual_row := 75
	sg_row := 100
	ph_row := 125
	submit_row := 175

	// group_row := 120

	product_text := "Product"
	lot_text := "Lot Number"

	visual_text := "Visual Inspection"
	sg_text := "SG"
	ph_text := "pH"

	submit_button := winc.NewPushButton(mainWindow)
	ph_field := show_edit(mainWindow, label_col, field_col, ph_row, ph_text)
	sg_field := show_edit(mainWindow, label_col, field_col, sg_row, sg_text)
	visual_field := show_checkbox(mainWindow, label_col, field_col, visual_row, visual_text)
	lot_field := show_edit(mainWindow, label_col, field_col, lot_row, lot_text)
	product_field := show_edit(mainWindow, label_col, field_col, product_row, product_text)

	// product_text := "Product"
	// product_field := show_edit(mainWindow, label_col, field_col, product_row, product_text)

	submit_button.SetText("Submit")
	submit_button.SetPos(40, submit_row) // (x, y)
	submit_button.SetSize(100, 40)       // (width, height)
	submit_button.OnClick().Bind(func(e *winc.Event) {

		product := newWaterBasedProduct(product_field, lot_field, visual_field, sg_field, ph_field)

		if product.check_data() {
			fmt.Println("data", product)
			product.print()
		}
	})

	// visual_field.SetFocus()
}
