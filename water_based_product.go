package main

import (
	"log"

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

func newWaterBasedProduct(base_product BaseProduct, visual_field *windigo.CheckBox, sg, ph float64) Product {

	base_product.Visual = visual_field.Checked()

	return WaterBasedProduct{base_product, sg, ph}.toProduct()

}

func (product WaterBasedProduct) check_data() bool {
	return true
}

func show_water_based(parent windigo.AutoPanel, qc_product QCProduct, create_new_product_cb func() BaseProduct) func(qc_product QCProduct) {

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
	sg_text := "SG"
	ph_text := "pH"

	group_panel := windigo.NewAutoPanel(parent)
	group_panel.SetSize(group_width, group_height)

	group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, BTM_SPACER_WIDTH, BTM_SPACER_HEIGHT)
	group_panel.SetMargins(group_margin, group_margin, 0, 0)

	ranges_panel := BuildNewWaterBasedProductRangesView(parent, qc_product, RANGE_WIDTH, group_height)
	ranges_panel.SetMarginTop(group_margin)

	visual_field := windigo.NewCheckBox(group_panel)
	visual_field.SetText(visual_text)
	visual_field.SetMarginsAll(ERROR_MARGIN)

	sg_field := BuildNewNumberEditView(group_panel, label_width, field_width, field_height, sg_text, ranges_panel.sg_field)
	ph_field := BuildNewNumberEditView(group_panel, label_width, field_width, field_height, ph_text, ranges_panel.ph_field)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(sg_field, windigo.Top)
	group_panel.Dock(ph_field, windigo.Top)

	submit_cb := func() {
		product := newWaterBasedProduct(create_new_product_cb(), visual_field, sg_field.Get(), ph_field.Get())
		if product.check_data() {
			log.Println("data", product)
			product.save()
			product.output()
		}
	}

	clear_cb := func() {
		visual_field.SetChecked(false)
		sg_field.Clear()
		ph_field.Clear()
		ranges_panel.Clear()
	}

	log_cb := func() {
		product := newWaterBasedProduct(create_new_product_cb(), visual_field, sg_field.Get(), ph_field.Get())
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

type WaterBasedProductRangesView struct {
	windigo.AutoPanel
	ph_field,
	sg_field *RangeROView

	Update func(qc_product QCProduct)
}

func (data_view WaterBasedProductRangesView) Clear() {
	data_view.ph_field.Clear()
	data_view.sg_field.Clear()
}

func BuildNewWaterBasedProductRangesView(parent windigo.AutoPanel, qc_product QCProduct, group_width, group_height int) WaterBasedProductRangesView {

	visual_text := "Visual Inspection"
	sg_text := "SG"
	ph_text := "pH"

	group_panel := windigo.NewAutoPanel(parent)
	group_panel.SetSize(group_width, group_height)
	group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, BTM_SPACER_WIDTH, BTM_SPACER_HEIGHT)

	// visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)
	visual_field := BuildNewProductAppearanceROView(group_panel, visual_text, qc_product.Appearance)

	sg_field := BuildNewRangeROView(group_panel, sg_text, qc_product.SG, format_ranges_sg)
	ph_field := BuildNewRangeROView(group_panel, ph_text, qc_product.PH, format_ranges_ph)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(sg_field, windigo.Top)
	group_panel.Dock(ph_field, windigo.Top)

	update := func(qc_product QCProduct) {
		visual_field.Update(qc_product.Appearance)
		sg_field.Update(qc_product.SG)
		ph_field.Update(qc_product.PH)
	}

	return WaterBasedProductRangesView{group_panel, &ph_field, &sg_field, update}

}
