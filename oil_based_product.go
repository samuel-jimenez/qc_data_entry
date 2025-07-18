package main

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/view"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type OilBasedProduct struct {
	product.BaseProduct
	sg float64
}

func (ob_product OilBasedProduct) toProduct() product.Product {
	return product.Product{
		BaseProduct: ob_product.Base(),
		PH:          nullable.NewNullFloat64(0, false),
		SG:          nullable.NewNullFloat64(ob_product.sg, true),
		Density:     nullable.NewNullFloat64(0, false),
		String_test: nullable.NewNullFloat64(0, false),
		Viscosity:   nullable.NewNullFloat64(0, false),
	}

	//TODO Option?
}

func newOilBasedProduct(base_product product.BaseProduct,
	have_visual bool, mass float64) product.Product {
	base_product.Visual = have_visual
	sg := formats.SG_from_mass(mass)

	return OilBasedProduct{base_product, sg}.toProduct()

}

func (product OilBasedProduct) Check_data() bool {
	return true
}

type OilBasedPanelView struct {
	Update  func(qc_product *product.QCProduct)
	SetFont func(font *windigo.Font)
	Refresh func()
}

func show_oil_based(parent *windigo.AutoPanel, qc_product *product.QCProduct, create_new_product_cb func() product.BaseProduct) *OilBasedPanelView {

	visual_text := "Visual Inspection"

	panel := windigo.NewAutoPanel(parent)

	group_panel := windigo.NewAutoPanel(panel)

	ranges_panel := BuildNewOilBasedProductRangesView(panel, qc_product)

	visual_field := NewBoolCheckboxView(group_panel, visual_text)

	mass_field := NewMassDataView(group_panel, ranges_panel)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(mass_field, windigo.Top)

	submit_cb := func() {
		product := newOilBasedProduct(create_new_product_cb(), visual_field.Get(), mass_field.Get())
		if product.Check_data() {
			log.Println("data", product)
			product.Save()
			err := product.Output()
			if err != nil {
				log.Printf("Error: %q: %s\n", err, "OilBasedProduct.Output")
			}
		}
	}

	clear_cb := func() {
		visual_field.Clear()
		mass_field.Clear()
		ranges_panel.Clear()

	}

	log_cb := func() {
		product := newOilBasedProduct(create_new_product_cb(), visual_field.Get(), mass_field.Get())
		if product.Check_data() {
			log.Println("data", product)
			product.Save()
			product.Export_json()
		}
	}

	button_dock := GUI.NewMarginalButtonDock(parent, []string{"Submit", "Clear", "Log"}, []int{40, 0, 10}, []func(){submit_cb, clear_cb, log_cb})

	panel.Dock(group_panel, windigo.Left)
	panel.Dock(ranges_panel, windigo.Right)

	parent.Dock(panel, windigo.Top)
	parent.Dock(button_dock, windigo.Top)

	setFont := func(font *windigo.Font) {
		ranges_panel.SetFont(font)
		visual_field.SetFont(font)
		mass_field.SetFont(font)
		button_dock.SetFont(font)
	}
	refresh := func() {

		panel.SetSize(GUI.OFF_AXIS, GROUP_HEIGHT)
		panel.SetMargins(GROUP_MARGIN, GROUP_MARGIN, 0, 0)

		group_panel.SetSize(GROUP_WIDTH, GROUP_HEIGHT)
		group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, BTM_SPACER_WIDTH, BTM_SPACER_HEIGHT)

		visual_field.SetSize(GUI.OFF_AXIS, FIELD_HEIGHT)

		mass_field.SetLabeledSize(GUI.LABEL_WIDTH, DATA_FIELD_WIDTH, DATA_SUBFIELD_WIDTH, DATA_UNIT_WIDTH, FIELD_HEIGHT)

		button_dock.SetDockSize(BUTTON_WIDTH, BUTTON_HEIGHT)

		ranges_panel.SetMarginTop(GROUP_MARGIN)
		ranges_panel.Refresh()

	}

	return &OilBasedPanelView{ranges_panel.Update, setFont, refresh}

}

type OilBasedProductRangesView struct {
	*windigo.AutoPanel
	*MassRangesView

	Update  func(qc_product *product.QCProduct)
	SetFont func(font *windigo.Font)
	Refresh func()
}

func (data_view OilBasedProductRangesView) Clear() {
	data_view.MassRangesView.Clear()
}

func BuildNewOilBasedProductRangesView(parent *windigo.AutoPanel, qc_product *product.QCProduct) OilBasedProductRangesView {

	visual_text := "Visual Inspection"
	mass_text := "Mass"

	sg_text := "Specific Gravity"
	density_text := "Density"

	group_panel := windigo.NewAutoPanel(parent)

	visual_field := product.BuildNewProductAppearanceROView(group_panel, visual_text, qc_product.Appearance)

	mass_field := view.BuildNewRangeROViewMap(group_panel, mass_text, qc_product.SG, formats.Format_mass, formats.Mass_from_sg)
	sg_field := view.BuildNewRangeROView(group_panel, sg_text, qc_product.SG, formats.Format_ranges_sg)
	density_field := view.BuildNewRangeROView(group_panel, density_text, qc_product.Density, formats.Format_ranges_density)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(mass_field, windigo.Top)
	group_panel.Dock(density_field, windigo.Bottom)
	group_panel.Dock(sg_field, windigo.Bottom)

	update := func(qc_product *product.QCProduct) {
		visual_field.Update(qc_product.Appearance)
		mass_field.Update(qc_product.SG)
		sg_field.Update(qc_product.SG)
		density_field.Update(qc_product.Density)
	}
	setFont := func(font *windigo.Font) {
		visual_field.SetFont(font)
		mass_field.SetFont(font)
		sg_field.SetFont(font)
		density_field.SetFont(font)
	}

	refresh := func() {
		group_panel.SetSize(RANGE_WIDTH, GROUP_HEIGHT)
		group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, BTM_SPACER_WIDTH, BTM_SPACER_HEIGHT)
		visual_field.Refresh()
		mass_field.Refresh()
		sg_field.Refresh()
		density_field.Refresh()
	}

	return OilBasedProductRangesView{group_panel,
		&MassRangesView{&mass_field,
			&sg_field,
			&density_field},
		update, setFont, refresh}

}
