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

func show_water_based(parent windigo.AutoPanel, qc_product QCProduct, create_new_product_cb func() BaseProduct) func(qc_product QCProduct) {

	group_width := GROUP_WIDTH
	group_height := GROUP_HEIGHT
	group_margin := GROUP_MARGIN

	label_width := LABEL_WIDTH
	field_width := FIELD_WIDTH
	field_height := FIELD_HEIGHT

	button_width := BUTTON_WIDTH
	button_height := BUTTON_HEIGHT

	top_spacer_width := TOP_SPACER_WIDTH
	top_spacer_height := TOP_SPACER_HEIGHT
	inter_spacer_height := INTER_SPACER_HEIGHT
	bottom_spacer_height := BUTTON_SPACER_HEIGHT

	visual_text := "Visual Inspection"
	sg_text := "SG"
	ph_text := "pH"

	group_panel := windigo.NewAutoPanel(parent)
	group_panel.SetSize(group_width, group_height)

	group_panel.SetPaddings(top_spacer_height, inter_spacer_height, inter_spacer_height, top_spacer_width)
	group_panel.SetMargins(group_margin, 0, 0, group_margin)

	ranges_panel := BuildNewWaterBasedProductRangesView(parent, qc_product, group_width, group_height)
	ranges_panel.SetMarginTop(group_margin)

	visual_field := windigo.NewCheckBox(group_panel)
	visual_field.SetText(visual_text)

	sg_field := show_number_edit(group_panel, label_width, field_width, field_height, sg_text, ranges_panel.sg_field)
	ph_field := show_number_edit(group_panel, label_width, field_width, field_height, ph_text, ranges_panel.ph_field)

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
	parent.Dock(ranges_panel, windigo.Right)

	return ranges_panel.Update

}

type WaterBasedProductRangesView struct {
	windigo.AutoPanel
	ph_field,
	sg_field *RangeROView

	Update func(qc_product QCProduct)
}

func BuildNewWaterBasedProductRangesView(parent windigo.AutoPanel, qc_product QCProduct, group_width, group_height int) WaterBasedProductRangesView {

	top_spacer_width := TOP_SPACER_WIDTH
	top_spacer_height := TOP_SPACER_HEIGHT
	inter_spacer_height := INTER_SPACER_HEIGHT

	visual_text := "Visual Inspection"
	sg_text := "SG"
	ph_text := "pH"

	group_panel := windigo.NewAutoPanel(parent)
	group_panel.SetSize(group_width, group_height)
	group_panel.SetPaddings(top_spacer_height, inter_spacer_height, inter_spacer_height, top_spacer_width)

	// visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)

	visual_field := windigo.NewLabel(group_panel)
	visual_field.SetText(visual_text)

	sg_field := BuildNewRangeROView(group_panel, sg_text, qc_product.SG, format_ranges_sg)
	ph_field := BuildNewRangeROView(group_panel, ph_text, qc_product.PH, format_ranges_ph)

	sg_field.SetMarginTop(inter_spacer_height)
	ph_field.SetMarginTop(inter_spacer_height)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(sg_field, windigo.Top)
	group_panel.Dock(ph_field, windigo.Top)

	update := func(qc_product QCProduct) {
		log.Println("update BuildNewWaterBasedProductRangesView", qc_product)

		sg_field.Update(qc_product.SG)
		ph_field.Update(qc_product.PH)
	}

	return WaterBasedProductRangesView{group_panel, &ph_field, &sg_field, update}

}
