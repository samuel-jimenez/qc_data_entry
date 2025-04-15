package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/samuel-jimenez/winc"
)

type WaterBasedProduct struct {
	BaseProduct
	sg float64
	ph float64
}

func (product WaterBasedProduct) toProduct() Product {
	return Product{product.toBaseProduct(), sql.NullFloat64{product.sg, true}, sql.NullFloat64{product.ph, true}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}, sql.NullFloat64{0, false}}

	//TODO Option?
}

func newWaterBasedProduct(base_product BaseProduct, visual_field *winc.CheckBox, sg_field *winc.Edit, ph_field *winc.Edit) WaterBasedProduct {

	base_product.visual = visual_field.Checked()
	sg, _ := strconv.ParseFloat(sg_field.Text(), 64)
	// if !err.Error(){fmt.Println("error",err)}
	ph, _ := strconv.ParseFloat(ph_field.Text(), 64)

	return WaterBasedProduct{base_product, sg, ph}

}

func (product WaterBasedProduct) check_data() bool {
	return true
}

func show_water_based(parent winc.Controller, create_new_product_cb func() BaseProduct) {
	label_col := 10
	field_col := 120

	visual_row := 25
	sg_row := 50
	ph_row := 75

	submit_col := 40
	submit_row := 180
	submit_button_width := 100
	submit_button_height := 40

	// group_row := 120

	visual_text := "Visual Inspection"
	sg_text := "SG"
	ph_text := "pH"

	// sample_row := 70
	// sample_text := "Sample Point"
	// sample_field := show_edit(mainWindow, label_col, field_col, sample_row, sample_text)

	visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)
	sg_field := show_edit(parent, label_col, field_col, sg_row, sg_text)
	ph_field := show_edit(parent, label_col, field_col, ph_row, ph_text)
	submit_button := winc.NewPushButton(parent)

	// product_text := "Product"
	// product_field := show_edit(mainWindow, label_col, field_col, product_row, product_text)

	submit_button.SetText("Submit")
	submit_button.SetPos(submit_col, submit_row)                     // (x, y)
	submit_button.SetSize(submit_button_width, submit_button_height) // (width, height)
	submit_button.OnClick().Bind(func(e *winc.Event) {

		// base_product := create_new_product_cb()

		product := newWaterBasedProduct(create_new_product_cb(), visual_field, sg_field, ph_field).toProduct()

		if product.check_data() {
			fmt.Println("data", product)
			product.save()
			product.print()
		}
	})

	// visual_field.SetFocus()
}
