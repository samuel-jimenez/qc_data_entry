package main

import (
	"log"
	"strconv"

	"github.com/samuel-jimenez/windigo"
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

func newWaterBasedProduct(base_product BaseProduct, visual_field *windigo.CheckBox, sg_field windigo.LabeledEdit, ph_field windigo.LabeledEdit) Product {

	base_product.Visual = visual_field.Checked()
	sg, _ := strconv.ParseFloat(sg_field.Text(), 64)
	// if !err.Error(){log.Println("error",err)}
	ph, _ := strconv.ParseFloat(ph_field.Text(), 64)

	return WaterBasedProduct{base_product, sg, ph}.toProduct()

}

func (product WaterBasedProduct) check_data() bool {
	return true
}

func show_water_based(parent windigo.AutoPanel, create_new_product_cb func() BaseProduct) {

	label_width := 110
	field_width := 200
	field_height := 22

	top_spacer_height := 25
	top_spacer_width := 10
	inter_spacer_height := 5
	bottom_spacer_height := BUTTON_SPACER_HEIGHT

	group_width := 300
	group_height := 170
	group_margin := 5

	button_width := 100
	button_height := 40
	// 	200
	// 50

	visual_text := "Visual Inspection"
	sg_text := "SG"
	ph_text := "pH"

	group_panel := windigo.NewAutoPanel(parent)
	group_panel.SetSize(group_width, group_height)

	group_panel.SetPaddings(top_spacer_height, inter_spacer_height, inter_spacer_height, top_spacer_width)
	group_panel.SetMargins(group_margin, 0, 0, group_margin)

	// visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)

	visual_field := windigo.NewCheckBox(group_panel)
	visual_field.SetText(visual_text)
	sg_field := show_edit(group_panel, label_width, field_width, field_height, sg_text)
	ph_field := show_edit(group_panel, label_width, field_width, field_height, ph_text)

	sg_field.SetMarginTop(inter_spacer_height)
	ph_field.SetMarginTop(inter_spacer_height)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(sg_field, windigo.Top)
	group_panel.Dock(ph_field, windigo.Top)

	submit_cb := func() {
		product := newWaterBasedProduct(create_new_product_cb(), visual_field, sg_field, ph_field)
		if product.check_data() {
			log.Println("data", product)
			product.save()
			product.output()
		}
	}

	clear_cb := func() {
		visual_field.SetChecked(false)
		sg_field.SetText("")
		ph_field.SetText("")
	}

	log_cb := func() {
		product := newWaterBasedProduct(create_new_product_cb(), visual_field, sg_field, ph_field)
		if product.check_data() {
			log.Println("data", product)
			product.save()
			product.export_json()
		}
	}

	button_dock := build_marginal_button_dock(parent, button_width, button_height, []string{"Submit", "Clear", "Log"}, []int{40, 0, 10}, []func(){submit_cb, clear_cb, log_cb})
	button_dock.SetMarginBtm(bottom_spacer_height)

	parent.Dock(button_dock, windigo.Bottom)
	parent.Dock(group_panel, windigo.Left)

}
