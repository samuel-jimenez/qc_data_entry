package subpanels

import (
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
	RefreshSize()
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

func (view *OilBasedProductView) RefreshSize() {
	view.SetSize(GUI.GROUP_WIDTH, GUI.GROUP_HEIGHT)
	view.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.BTM_SPACER_WIDTH, GUI.BTM_SPACER_HEIGHT)

	view.visual_field.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)
	view.density_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.DATA_SUBFIELD_WIDTH, GUI.DATA_UNIT_WIDTH, GUI.EDIT_FIELD_HEIGHT)

}
