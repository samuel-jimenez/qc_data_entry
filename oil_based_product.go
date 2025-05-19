package main

import (
	"log"

	"github.com/samuel-jimenez/windigo"
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
	visual_field *windigo.CheckBox, mass_field windigo.LabeledEdit) Product {
	base_product.Visual = visual_field.Checked()
	sg := sg_from_mass(mass_field)

	return OilBasedProduct{base_product, sg}.toProduct()

}

func (product OilBasedProduct) check_data() bool {
	return true
}

func show_oil_based(parent windigo.AutoPanel, create_new_product_cb func() BaseProduct) {

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
	mass_text := "Mass"

	group_panel := windigo.NewAutoPanel(parent)
	group_panel.SetSize(group_width, group_height)

	group_panel.SetPaddings(top_spacer_height, inter_spacer_height, inter_spacer_height, top_spacer_width)
	group_panel.SetMargins(group_margin, 0, 0, group_margin)

	// visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)

	visual_field := windigo.NewCheckBox(group_panel)
	visual_field.SetText(visual_text)

	mass_field := show_mass_sg(group_panel, label_width, field_width, field_height, mass_text)

	mass_field.SetMarginTop(inter_spacer_height)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(mass_field, windigo.Top)

	submit_cb := func() {
		product := newOilBasedProduct(create_new_product_cb(), visual_field, mass_field)
		if product.check_data() {
			log.Println("data", product)
			product.save()
			product.output()
		}
	}

	clear_cb := func() {
		visual_field.SetChecked(false)

		mass_field.SetText("")
		mass_field.OnChange().Fire(nil)
	}

	log_cb := func() {
		product := newOilBasedProduct(create_new_product_cb(), visual_field, mass_field)
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
