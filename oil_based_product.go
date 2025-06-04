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
	have_visual bool, mass float64) Product {
	base_product.Visual = have_visual
	sg := sg_from_mass(mass)

	return OilBasedProduct{base_product, sg}.toProduct()

}

func (product OilBasedProduct) check_data() bool {
	return true
}

func show_oil_based(parent windigo.AutoPanel, qc_product QCProduct, create_new_product_cb func() BaseProduct) func(qc_product QCProduct) {

	group_width := GROUP_WIDTH
	group_height := GROUP_HEIGHT
	group_margin := GROUP_MARGIN

	label_width := LABEL_WIDTH
	field_width := DATA_FIELD_WIDTH
	field_height := FIELD_HEIGHT

	button_width := BUTTON_WIDTH
	button_height := BUTTON_HEIGHT

	bottom_spacer_height := BUTTON_SPACER_HEIGHT

	visual_text := "Visual Inspection"
	mass_text := "Mass"

	group_panel := windigo.NewAutoPanel(parent)
	group_panel.SetSize(group_width, group_height)

	group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, BTM_SPACER_WIDTH, BTM_SPACER_HEIGHT)
	group_panel.SetMargins(group_margin, group_margin, 0, 0)

	ranges_panel := BuildNewOilBasedProductRangesView(parent, qc_product, RANGE_WIDTH, group_height)
	ranges_panel.SetMarginTop(group_margin)

	visual_field := NewBoolCheckboxView(group_panel, OFF_AXIS, field_height, visual_text)

	mass_field := NewMassDataView(group_panel, label_width, field_width, field_height, mass_text, ranges_panel)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(mass_field, windigo.Top)

	submit_cb := func() {
		product := newOilBasedProduct(create_new_product_cb(), visual_field.Get(), mass_field.Get())
		if product.check_data() {
			log.Println("data", product)
			product.save()
			product.output()
		}
	}

	clear_cb := func() {
		visual_field.Clear()
		mass_field.Clear()
		ranges_panel.Clear()

	}

	log_cb := func() {
		product := newOilBasedProduct(create_new_product_cb(), visual_field.Get(), mass_field.Get())
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
	MassRangesView

	Update func(qc_product QCProduct)
}

func (data_view OilBasedProductRangesView) Clear() {
	data_view.MassRangesView.Clear()
}

func BuildNewOilBasedProductRangesView(parent windigo.AutoPanel, qc_product QCProduct, group_width, group_height int) OilBasedProductRangesView {

	visual_text := "Visual Inspection"
	mass_text := "Mass"

	sg_text := "Specific Gravity"
	density_text := "Density"

	group_panel := windigo.NewAutoPanel(parent)
	group_panel.SetSize(group_width, group_height)
	group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, BTM_SPACER_WIDTH, BTM_SPACER_HEIGHT)

	visual_field := BuildNewProductAppearanceROView(group_panel, visual_text, qc_product.Appearance)

	mass_field := BuildNewRangeROViewMap(group_panel, mass_text, qc_product.SG, format_mass, mass_from_sg)
	sg_field := BuildNewRangeROView(group_panel, sg_text, qc_product.SG, format_ranges_sg)
	density_field := BuildNewRangeROView(group_panel, density_text, qc_product.Density, format_ranges_density)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(mass_field, windigo.Top)
	group_panel.Dock(density_field, windigo.Bottom)
	group_panel.Dock(sg_field, windigo.Bottom)

	update := func(qc_product QCProduct) {
		visual_field.Update(qc_product.Appearance)
		mass_field.Update(qc_product.SG)
		sg_field.Update(qc_product.SG)
		density_field.Update(qc_product.Density)
	}

	return OilBasedProductRangesView{group_panel,
		MassRangesView{&mass_field,
			&sg_field,
			&density_field},
		update}

}
