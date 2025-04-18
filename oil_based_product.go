package main

import (
	"log"
	"strconv"

	"github.com/samuel-jimenez/winc"
)

type OilBasedProduct struct {
	BaseProduct
	sg float64
}

func (product OilBasedProduct) toProduct() Product {
	return Product{product.toBaseProduct(), NewNullFloat64(product.sg, true), NewNullFloat64(0, false), NewNullFloat64(0, false), NewNullFloat64(0, false), NewNullFloat64(0, false)}

	//TODO Option?
}

func newOilBasedProduct(base_product BaseProduct,
	visual_field *winc.CheckBox, mass_field *winc.Edit) Product {
	base_product.Visual = visual_field.Checked()
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

	submit_col := SUBMIT_COL
	submit_button_width := 100
	submit_button_height := 40

	visual_text := "Visual Inspection"
	mass_text := "Mass"

	visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)
	mass_field := show_mass_sg(parent, label_col, field_col, mass_row, mass_text)

	submit_button := winc.NewPushButton(parent)
	submit_button.SetText("Submit")
	submit_button.SetPos(submit_col, SUBMIT_ROW) // (x, y)
	// submit_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	submit_button.SetSize(submit_button_width, submit_button_height) // (width, height)
	submit_button.OnClick().Bind(func(e *winc.Event) {
		product := newOilBasedProduct(create_new_product_cb(), visual_field, mass_field)
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
		mass_field.SetText("")
		mass_field.OnKillFocus().Fire(nil)
	}

	clear_button.SetText("Clear")
	clear_button.SetPos(clear_button_col, clear_button_row) // (x, y)
	// clear_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	clear_button.SetSize(clear_button_width, clear_button_height) // (width, height)
	clear_button.OnClick().Bind(func(e *winc.Event) { clear_cb() })

	// visual_field.SetFocus()
}
