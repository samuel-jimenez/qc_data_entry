package main

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/windigo"
)

type WaterBasedProduct struct {
	BaseProduct
	sg float64
	ph float64
}

func (product WaterBasedProduct) toProduct() Product {
	return Product{product.toBaseProduct(), nullable.NewNullFloat64(product.sg, true), nullable.NewNullFloat64(product.ph, true), nullable.NewNullFloat64(0, false), nullable.NewNullFloat64(0, false), nullable.NewNullFloat64(0, false)}

	//TODO Option?
}

func newWaterBasedProduct(base_product BaseProduct, have_visual bool, sg, ph float64) Product {

	base_product.Visual = have_visual

	return WaterBasedProduct{base_product, sg, ph}.toProduct()

}

func (product WaterBasedProduct) check_data() bool {
	return true
}

type WaterBasedPanelView struct {
	Update  func(qc_product QCProduct)
	SetFont func(font *windigo.Font)
	Refresh func()
}

func show_water_based(parent *windigo.AutoPanel, qc_product QCProduct, create_new_product_cb func() BaseProduct) *WaterBasedPanelView {

	visual_text := "Visual Inspection"
	sg_text := "SG"
	ph_text := "pH"

	panel := windigo.NewAutoPanel(parent)

	group_panel := windigo.NewAutoPanel(panel)

	ranges_panel := BuildNewWaterBasedProductRangesView(panel, qc_product)

	visual_field := NewBoolCheckboxView(group_panel, visual_text)

	sg_field := NewNumberEditViewWithChange(group_panel, sg_text, ranges_panel.sg_field)
	ph_field := NewNumberEditViewWithChange(group_panel, ph_text, ranges_panel.ph_field)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(sg_field, windigo.Top)
	group_panel.Dock(ph_field, windigo.Top)

	submit_cb := func() {
		product := newWaterBasedProduct(create_new_product_cb(), visual_field.Get(), sg_field.Get(), ph_field.Get())
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
		product := newWaterBasedProduct(create_new_product_cb(), visual_field.Get(), sg_field.Get(), ph_field.Get())
		if product.check_data() {
			log.Println("data", product)
			product.save()
			product.export_json()
		}
	}

	button_dock := NewMarginalButtonDock(parent, []string{"Submit", "Clear", "Log"}, []int{40, 0, 10}, []func(){submit_cb, clear_cb, log_cb})

	panel.Dock(group_panel, windigo.Left)
	panel.Dock(ranges_panel, windigo.Right)

	parent.Dock(panel, windigo.Top)
	parent.Dock(button_dock, windigo.Top)

	setFont := func(font *windigo.Font) {
		ranges_panel.SetFont(font)
		visual_field.SetFont(font)
		sg_field.SetFont(font)
		ph_field.SetFont(font)
		button_dock.SetFont(font)
	}
	refresh := func() {

		panel.SetSize(OFF_AXIS, GROUP_HEIGHT)
		panel.SetMargins(GROUP_MARGIN, GROUP_MARGIN, 0, 0)

		group_panel.SetSize(GROUP_WIDTH, GROUP_HEIGHT)
		group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, BTM_SPACER_WIDTH, BTM_SPACER_HEIGHT)

		ranges_panel.SetSize(RANGE_WIDTH, GROUP_HEIGHT)
		ranges_panel.SetMarginTop(GROUP_MARGIN)

		visual_field.SetSize(OFF_AXIS, FIELD_HEIGHT)

		sg_field.SetLabeledSize(LABEL_WIDTH, DATA_FIELD_WIDTH, FIELD_HEIGHT)
		ph_field.SetLabeledSize(LABEL_WIDTH, DATA_FIELD_WIDTH, FIELD_HEIGHT)

		button_dock.SetDockSize(BUTTON_WIDTH, BUTTON_HEIGHT)

		ranges_panel.Refresh()
	}

	return &WaterBasedPanelView{ranges_panel.Update, setFont, refresh}

}

type WaterBasedProductRangesView struct {
	*windigo.AutoPanel
	ph_field,
	sg_field *RangeROView

	Update  func(qc_product QCProduct)
	SetFont func(font *windigo.Font)
	Refresh func()
}

func (data_view WaterBasedProductRangesView) Clear() {
	data_view.ph_field.Clear()
	data_view.sg_field.Clear()
}

func BuildNewWaterBasedProductRangesView(parent *windigo.AutoPanel, qc_product QCProduct) WaterBasedProductRangesView {

	visual_text := "Visual Inspection"
	sg_text := "SG"
	ph_text := "pH"

	group_panel := windigo.NewAutoPanel(parent)

	// visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)
	visual_field := BuildNewProductAppearanceROView(group_panel, visual_text, qc_product.Appearance)

	sg_field := BuildNewRangeROView(group_panel, sg_text, qc_product.SG, formats.Format_ranges_sg)
	ph_field := BuildNewRangeROView(group_panel, ph_text, qc_product.PH, formats.Format_ranges_ph)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(sg_field, windigo.Top)
	group_panel.Dock(ph_field, windigo.Top)

	update := func(qc_product QCProduct) {
		visual_field.Update(qc_product.Appearance)
		sg_field.Update(qc_product.SG)
		ph_field.Update(qc_product.PH)
	}
	setFont := func(font *windigo.Font) {
		visual_field.SetFont(font) //?TODO
		sg_field.SetFont(font)
		ph_field.SetFont(font)
	}
	refresh := func() {
		group_panel.SetSize(RANGE_WIDTH, GROUP_HEIGHT)
		group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, BTM_SPACER_WIDTH, BTM_SPACER_HEIGHT)

		visual_field.Refresh()
		sg_field.Refresh()
		ph_field.Refresh()
	}

	return WaterBasedProductRangesView{group_panel, &ph_field, &sg_field, update, setFont, refresh}

}
