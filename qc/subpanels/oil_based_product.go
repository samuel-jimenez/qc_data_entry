package subpanels

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
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
		String_test: nullable.NullInt64Default(),
		Viscosity:   nullable.NullInt64Default(),
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

type OilBasedProductViewer interface {
	*windigo.AutoPanel
	Get(base_product product.BaseProduct) product.Product
	Clear()
	SetFont(font *windigo.Font)
	Refresh()
}

type OilBasedProductView struct {
	*windigo.AutoPanel
	visual_field  *views.BoolCheckboxView
	density_field *views.MassDataView
	ranges_panel  *OilBasedProductRangesView
}

func newOilBasedProductView(parent *windigo.AutoPanel, ranges_panel *OilBasedProductRangesView) *OilBasedProductView {

	visual_text := "Visual Inspection"
	mass_text := "Mass"

	view := new(OilBasedProductView)

	panel := windigo.NewAutoPanel(parent)

	visual_field := views.NewBoolCheckboxView(panel, visual_text)
	mass_field := views.NewNumberEditView(panel, mass_text)

	density_field := views.NewMassDataView(panel, ranges_panel, mass_field)

	panel.Dock(visual_field, windigo.Top)
	panel.Dock(density_field, windigo.Top)

	view.AutoPanel = panel
	view.visual_field = visual_field
	view.density_field = density_field
	view.ranges_panel = ranges_panel

	return view
}

func (view *OilBasedProductView) Get(base_product product.BaseProduct) product.Product {
	// base_product.Visual = view.visual_field.Checked()
	return newOilBasedProduct(base_product, view.visual_field.Get(), view.density_field.Get())
}

func (view *OilBasedProductView) Clear() {
	view.visual_field.Clear()

	view.density_field.Clear()

	view.ranges_panel.Clear()

}

func (view *OilBasedProductView) SetFont(font *windigo.Font) {
	view.AutoPanel.SetFont(font)
	view.visual_field.SetFont(font)
	view.density_field.SetFont(font)
}

func (view *OilBasedProductView) Refresh() {
	view.SetSize(GUI.GROUP_WIDTH, GUI.GROUP_HEIGHT)
	view.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.BTM_SPACER_WIDTH, GUI.BTM_SPACER_HEIGHT)

	view.visual_field.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)
	view.density_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.DATA_SUBFIELD_WIDTH, GUI.DATA_UNIT_WIDTH, GUI.EDIT_FIELD_HEIGHT)

	// view.ranges_panel.SetMarginTop(GUI.GROUP_MARGIN)
	// view.ranges_panel.Refresh()

}

type OilBasedPanelView struct {
	*windigo.AutoPanel
	group_panel  *OilBasedProductView
	ranges_panel *OilBasedProductRangesView
	button_dock  *GUI.ButtonDock

	// Update func(qc_product *product.QCProduct)
	// SetFont func(font *windigo.Font)
	// Refresh func()
}

func Show_oil_based(parent *windigo.AutoPanel, qc_product *product.QCProduct, create_new_product_cb func() product.BaseProduct) *OilBasedPanelView {

	view := new(OilBasedPanelView)

	panel := windigo.NewAutoPanel(parent)

	ranges_panel := BuildNewOilBasedProductRangesView(panel, qc_product)
	group_panel := newOilBasedProductView(panel, ranges_panel)

	submit_cb := func() {
		measured_product := group_panel.Get(create_new_product_cb())
		if measured_product.Check_data() {
			log.Println("ob sub data", measured_product)
			// TODO blend013 ensurethis works with testing blends
			// measured_product.Save()
			product.Store(measured_product)
			err := measured_product.Output()
			if err != nil {
				log.Printf("Error: [%s]: %q\n", "OilBasedProduct.Output", err)
			}
			// * Check storage
			measured_product.CheckStorage()
		}
	}

	log_cb := func() {
		measured_product := group_panel.Get(create_new_product_cb())
		if measured_product.Check_data() {
			log.Println("ob log data", measured_product)
			measured_product.Save()
			measured_product.Export_json()
		}
	}

	button_dock := GUI.NewMarginalButtonDock(parent, []string{"Submit", "Clear", "Log"}, []int{40, 0, 10}, []func(){submit_cb, group_panel.Clear, log_cb})

	panel.Dock(group_panel, windigo.Left)
	panel.Dock(ranges_panel, windigo.Right)

	parent.Dock(panel, windigo.Top)
	parent.Dock(button_dock, windigo.Top)

	view.AutoPanel = panel
	view.group_panel = group_panel
	view.ranges_panel = ranges_panel
	view.button_dock = button_dock

	return view

}

func (view *OilBasedPanelView) SetFont(font *windigo.Font) {
	view.group_panel.SetFont(font)
	view.ranges_panel.SetFont(font)
	view.button_dock.SetFont(font)
}

func (view *OilBasedPanelView) Refresh() {

	view.SetSize(GUI.OFF_AXIS, GUI.GROUP_HEIGHT)
	view.SetMargins(GUI.GROUP_MARGIN, GUI.GROUP_MARGIN, 0, 0)

	view.group_panel.Refresh()

	view.ranges_panel.SetMarginTop(GUI.GROUP_MARGIN)
	view.ranges_panel.Refresh()

	view.button_dock.SetDockSize(GUI.BUTTON_WIDTH, GUI.BUTTON_HEIGHT)

}

func (view *OilBasedPanelView) Update(qc_product *product.QCProduct) {
	view.ranges_panel.Update(qc_product)
}

type OilBasedProductRangesView struct {
	*windigo.AutoPanel
	*views.MassRangesView

	visual_field *product.ProductAppearanceROView
}

func BuildNewOilBasedProductRangesView(parent *windigo.AutoPanel, qc_product *product.QCProduct) *OilBasedProductRangesView {

	visual_text := "Visual Inspection"
	mass_text := "Mass"

	sg_text := "Specific Gravity"
	density_text := "Density"

	group_panel := windigo.NewAutoPanel(parent)

	visual_field := product.BuildNewProductAppearanceROView(group_panel, visual_text, qc_product.Appearance)

	mass_field := views.BuildNewRangeROViewMap(group_panel, mass_text, qc_product.SG, formats.Format_mass, formats.Mass_from_sg)
	sg_field := views.BuildNewRangeROView(group_panel, sg_text, qc_product.SG, formats.Format_ranges_sg)
	density_field := views.BuildNewRangeROView(group_panel, density_text, qc_product.Density, formats.Format_ranges_density)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(mass_field, windigo.Top)
	group_panel.Dock(density_field, windigo.Bottom)
	group_panel.Dock(sg_field, windigo.Bottom)

	return &OilBasedProductRangesView{group_panel,
		&views.MassRangesView{Mass_field: &mass_field,
			SG_field:      &sg_field,
			Density_field: &density_field},
		&visual_field,
	}
}

func (view *OilBasedProductRangesView) Update(qc_product *product.QCProduct) {
	view.visual_field.Update(qc_product.Appearance)

	// view.MassRangesView.Update(qc_product)
	view.Mass_field.Update(qc_product.SG)
	view.SG_field.Update(qc_product.SG)
	view.Density_field.Update(qc_product.Density)
}

func (data_view *OilBasedProductRangesView) Clear() {
	data_view.MassRangesView.Clear()
}

func (view *OilBasedProductRangesView) SetFont(font *windigo.Font) {
	view.visual_field.SetFont(font)
	view.Mass_field.SetFont(font)
	view.SG_field.SetFont(font)
	view.Density_field.SetFont(font)
}

func (view *OilBasedProductRangesView) Refresh() {
	view.SetSize(GUI.DATA_FIELD_WIDTH, GUI.GROUP_HEIGHT)
	// view.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.BTM_SPACER_WIDTH, GUI.BTM_SPACER_HEIGHT)
	view.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.RANGES_RO_PADDING, GUI.BTM_SPACER_HEIGHT)
	view.Mass_field.Refresh()
	view.SG_field.Refresh()
	view.Density_field.Refresh()
}
