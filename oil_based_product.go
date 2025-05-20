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
	sg := sg_from_mass_field(mass_field)

	return OilBasedProduct{base_product, sg}.toProduct()

}

func (product OilBasedProduct) check_data() bool {
	return true
}

func show_oil_based(parent windigo.AutoPanel, qc_product QCProduct, create_new_product_cb func() BaseProduct) func(qc_product QCProduct) {

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

	ranges_panel := BuildNewOilBasedProductRangesView(parent, qc_product, group_width, group_height)
	ranges_panel.SetMarginTop(group_margin)

	// visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)

	visual_field := windigo.NewCheckBox(group_panel)
	visual_field.SetText(visual_text)

	mass_field := show_mass_sg(group_panel, label_width, field_width, field_height, mass_text, ranges_panel)

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
	parent.Dock(ranges_panel, windigo.Right)

	return ranges_panel.Update

}

type OilBasedProductRangesView struct {
	windigo.AutoPanel
	DerivedMassRangesView

	Update func(qc_product QCProduct)
}

func BuildNewOilBasedProductRangesView(parent windigo.AutoPanel, qc_product QCProduct, group_width, group_height int) OilBasedProductRangesView {

	top_spacer_height := 25
	top_spacer_width := 10
	inter_spacer_height := 5

	visual_text := "Visual Inspection"
	mass_text := "Mass"

	sg_text := "Specific Gravity"
	density_text := "Density"

	group_panel := windigo.NewAutoPanel(parent)
	group_panel.SetSize(group_width, group_height)
	group_panel.SetPaddings(top_spacer_height, inter_spacer_height, inter_spacer_height, top_spacer_width)

	visual_field := windigo.NewLabel(group_panel)
	visual_field.SetText(visual_text)

	mass_field := BuildNewRangeROViewMap(group_panel, mass_text, qc_product.SG, format_mass, mass_from_sg)

	sg_field := BuildNewRangeROView(group_panel, sg_text, qc_product.SG, format_ranges_sg)
	density_field := BuildNewRangeROView(group_panel, density_text, qc_product.Density, format_ranges_density)

	sg_field.SetMarginTop(inter_spacer_height)
	density_field.SetMarginTop(inter_spacer_height)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(mass_field, windigo.Top)
	group_panel.Dock(density_field, windigo.Bottom)
	group_panel.Dock(sg_field, windigo.Bottom)

	update := func(qc_product QCProduct) {
		log.Println("update BuildNewOilBasedProductRangesView", qc_product)
		mass_field.Update(qc_product.SG)

		sg_field.Update(qc_product.SG)
		density_field.Update(qc_product.Density)
	}

	return OilBasedProductRangesView{group_panel,
		DerivedMassRangesView{&mass_field,
			&sg_field,
			&density_field},
		update}

}
