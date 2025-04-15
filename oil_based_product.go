package main

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/samuel-jimenez/winc"
)

type OilBasedProduct struct {
	BaseProduct
	sg float64
}

func (product OilBasedProduct) toProduct() Product {
	return Product{product.toBaseProduct(), sql.NullFloat64{product.sg, true}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}}

	//TODO Option?
}

func newOilBasedProduct(base_product BaseProduct,
	visual_field *winc.CheckBox, mass_field *winc.Edit) Product {
	base_product.visual = visual_field.Checked()
	mass, _ := strconv.ParseFloat(mass_field.Text(), 64)
	// if !err.Error(){log.Println("error",err)}
	sg := mass / SAMPLE_VOLUME

	return OilBasedProduct{base_product, sg}.toProduct()

}

func (product OilBasedProduct) check_data() bool {
	return true
}

func show_oil_based(parent winc.Controller, create_new_product_cb func() BaseProduct) {

	label_col := 10
	field_col := 120

	visual_row := 25
	mass_row := 50

	submit_col := 40
	submit_button_width := 100
	submit_button_height := 40

	visual_text := "Visual Inspection"
	mass_text := "Mass"

	// sample_field := show_edit(mainWindow, label_col, field_col, sample_row, sample_text)

	visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)
	mass_field := show_mass_sg(parent, label_col, field_col, mass_row, mass_text)

	// 	product_row := 20
	// product_text := "Product"
	// product_field := show_edit(mainWindow, label_col, field_col, product_row, product_text)

	submit_button := winc.NewPushButton(parent)

	submit_button.SetText("Submit")
	submit_button.SetPos(submit_col, SUBMIT_ROW) // (x, y)
	// submit_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	submit_button.SetSize(submit_button_width, submit_button_height) // (width, height)
	submit_button.OnClick().Bind(func(e *winc.Event) {

		// product := newOilBasedProduct(product_field, lot_field, sample_field, visual_field, mass_field)
		product := newOilBasedProduct(create_new_product_cb(), visual_field, mass_field)

		if product.check_data() {
			log.Println("data", product)
			product.save()
			product.print()
		}
	})

	visual_field.SetFocus()
}
