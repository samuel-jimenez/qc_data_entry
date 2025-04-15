package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/samuel-jimenez/winc"
)

type OilBasedProduct struct {
	BaseProduct
	sg float64
}

func (product OilBasedProduct) toProduct() Product {
	return Product{BaseProduct{product.product_type, product.lot_number, product.visual, product.product_id, product.lot_id}, sql.NullFloat64{product.sg, true}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullString{"", false}}

	//TODO Option?
}

func newOilBasedProduct(base_product BaseProduct,
	visual_field *winc.CheckBox, mass_field *winc.Edit) OilBasedProduct {
	base_product.visual = visual_field.Checked()
	mass, _ := strconv.ParseFloat(mass_field.Text(), 64)
	// if !err.Error(){fmt.Println("error",err)}
	sg := mass / SAMPLE_VOLUME

	return OilBasedProduct{base_product, sg}

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
	submit_row := 180
	submit_button_width := 100
	submit_button_height := 40

	visual_text := "Visual Inspection"
	mass_text := "Mass"

	// sample_field := show_edit(mainWindow, label_col, field_col, sample_row, sample_text)

	visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)
	mass_field := show_edit(parent, label_col, field_col, mass_row, mass_text)

	// 	product_row := 20
	// product_text := "Product"
	// product_field := show_edit(mainWindow, label_col, field_col, product_row, product_text)

	submit_button := winc.NewPushButton(parent)

	submit_button.SetText("Submit")
	submit_button.SetPos(submit_col, submit_row) // (x, y)
	// submit_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	submit_button.SetSize(submit_button_width, submit_button_height) // (width, height)
	submit_button.OnClick().Bind(func(e *winc.Event) {

		// product := newOilBasedProduct(product_field, lot_field, sample_field, visual_field, mass_field)
		product := newOilBasedProduct(create_new_product_cb(), visual_field, mass_field).toProduct()

		if product.check_data() {
			fmt.Println("data", product)
			product.print()
		}
	})

	visual_field.SetFocus()
}
