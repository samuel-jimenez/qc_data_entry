package main

import (
	"fmt"
	"strconv"
	"strings"

	"codeberg.org/go-pdf/fpdf"
	"github.com/samuel-jimenez/winc"
)

type FrictionReducerProduct struct {
	base_product Product
	sg           float64
	string_test  float64
	viscosity    float64
}

func newFrictionReducerProduct(product_field *winc.Edit, lot_field *winc.Edit, visual_field *winc.CheckBox, viscosity_field *winc.Edit, mass_field *winc.Edit, string_field *winc.Edit) FrictionReducerProduct {
	base_product := Product{product_field.Text(), lot_field.Text(), visual_field.Checked()}
	viscosity, _ := strconv.ParseFloat(viscosity_field.Text(), 64)
	mass, _ := strconv.ParseFloat(mass_field.Text(), 64)
	// if !err.Error(){fmt.Println("error",err)}
	string_test, _ := strconv.ParseFloat(string_field.Text(), 64)
	sg := mass / SAMPLE_VOLUME

	return FrictionReducerProduct{base_product, sg, string_test, viscosity}

}

func (product FrictionReducerProduct) get_pdf_name() string {
	return product.base_product.get_pdf_name()
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
	pdf.Cell(field_width, field_height, strings.ToUpper(product.base_product.product_type))

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
	pdf.Cell(field_width, field_height, strings.ToUpper(product.base_product.lot_number))

	err := pdf.OutputFileAndClose(product.get_pdf_name())
	return err
}

func show_fr(mainWindow *winc.Form) {
	label_col := 10
	field_col := 120

	product_row := 20
	lot_row := 45

	visual_row := 75
	viscosity_row := 100
	mass_row := 125
	string_row := 150
	submit_row := 175

	// group_row := 120

	product_text := "Product"
	lot_text := "Lot Number"

	visual_text := "Visual Inspection"
	viscosity_text := "Viscosity"
	mass_text := "Mass"
	string_text := "String"

	submit_button := winc.NewPushButton(mainWindow)
	string_field := show_edit(mainWindow, label_col, field_col, string_row, string_text)
	mass_field := show_edit(mainWindow, label_col, field_col, mass_row, mass_text)
	viscosity_field := show_edit(mainWindow, label_col, field_col, viscosity_row, viscosity_text)
	visual_field := show_checkbox(mainWindow, label_col, field_col, visual_row, visual_text)
	lot_field := show_edit(mainWindow, label_col, field_col, lot_row, lot_text)
	product_field := show_edit(mainWindow, label_col, field_col, product_row, product_text)

	// product_text := "Product"
	// product_field := show_edit(mainWindow, label_col, field_col, product_row, product_text)

	submit_button.SetText("Submit")
	submit_button.SetPos(40, submit_row) // (x, y)
	submit_button.SetSize(100, 40)       // (width, height)
	submit_button.OnClick().Bind(func(e *winc.Event) {

		product := newFrictionReducerProduct(product_field, lot_field, visual_field, viscosity_field, mass_field, string_field)

		if product.check_data() {
			fmt.Println("data", product)
			product.print()
		}
	})

	// visual_field.SetFocus()
}


