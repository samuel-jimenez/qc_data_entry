package main

import (
	"log"
	"strconv"

	"github.com/samuel-jimenez/winc"
)

type WaterBasedProduct struct {
	BaseProduct
	sg float64
	ph float64
}

func (product WaterBasedProduct) toProduct() Product {
	return Product{product.toBaseProduct(), NewNullFloat64(product.sg, true), NewNullFloat64(product.ph, true), NewNullFloat64(0, false), NewNullFloat64(0, false), NewNullFloat64(0, false)}

	//TODO Option?
}

func newWaterBasedProduct(base_product BaseProduct, visual_field *winc.CheckBox, sg_field *winc.Edit, ph_field *winc.Edit) Product {

	base_product.Visual = visual_field.Checked()
	sg, _ := strconv.ParseFloat(sg_field.Text(), 64)
	// if !err.Error(){log.Println("error",err)}
	ph, _ := strconv.ParseFloat(ph_field.Text(), 64)

	return WaterBasedProduct{base_product, sg, ph}.toProduct()

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

	submit_col := SUBMIT_COL
	submit_button_width := 100
	submit_button_height := 40

	visual_text := "Visual Inspection"
	sg_text := "SG"
	ph_text := "pH"

	visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)
	sg_field := show_edit(parent, label_col, field_col, sg_row, sg_text)
	ph_field := show_edit(parent, label_col, field_col, ph_row, ph_text)

	submit_button := winc.NewPushButton(parent)
	submit_button.SetText("Submit")
	submit_button.SetPos(submit_col, SUBMIT_ROW)                     // (x, y)
	submit_button.SetSize(submit_button_width, submit_button_height) // (width, height)
	submit_button.OnClick().Bind(func(e *winc.Event) {
		product := newWaterBasedProduct(create_new_product_cb(), visual_field, sg_field, ph_field)
		if product.check_data() {
			log.Println("data", product)
			product.save()
			product.output()
		}
	})

	clear_button_col := CLEAR_COL
	clear_button_row := SUBMIT_ROW
	clear_button_width := 100
	clear_button_height := 40
	clear_button := winc.NewPushButton(parent)
	clear_cb := func() {
		visual_field.SetChecked(false)
		sg_field.SetText("")
		ph_field.SetText("")
	}

	clear_button.SetText("Clear")
	clear_button.SetPos(clear_button_col, clear_button_row) // (x, y)
	// clear_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	clear_button.SetSize(clear_button_width, clear_button_height) // (width, height)
	clear_button.OnClick().Bind(func(e *winc.Event) { clear_cb() })

	log_button_col := 250
	log_button_row := SUBMIT_ROW
	log_button_width := 100
	log_button_height := 40
	log_button := winc.NewPushButton(parent)

	log_button.SetText("Log")
	log_button.SetPos(log_button_col, log_button_row) // (x, y)
	// log_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	log_button.SetSize(log_button_width, log_button_height) // (width, height)
	log_button.OnClick().Bind(func(e *winc.Event) {
		product := newWaterBasedProduct(create_new_product_cb(), visual_field, sg_field, ph_field)
		if product.check_data() {
			log.Println("data", product)
			product.save()
			product.export_json()
		}
	})

}
