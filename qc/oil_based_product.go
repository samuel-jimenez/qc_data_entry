package qc

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/qc_ui"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type OilBasedProduct struct {
	product.QCProduct
	sg float64
}

func (ob_product OilBasedProduct) toProduct() *product.MeasuredProduct {
	return &product.MeasuredProduct{
		QCProduct:   ob_product.Base(),
		PH:          nullable.NewNullFloat64(0, false),
		SG:          nullable.NewNullFloat64(ob_product.sg, true),
		Density:     nullable.NewNullFloat64(0, false),
		String_test: nullable.NullInt64Default(),
		Viscosity:   nullable.NullInt64Default(),
	}

	// TODO Option?
}

func newOilBasedProduct(base_product product.QCProduct,
	have_visual bool, mass float64,
) *product.MeasuredProduct {
	base_product.Visual = have_visual
	sg := formats.SG_from_mass(mass)

	return OilBasedProduct{base_product, sg}.toProduct()
}

func (product OilBasedProduct) Check_data() bool {
	return true
}

type OilBasedProductViewer interface {
	*windigo.AutoPanel
	Get(base_product product.BaseProduct) product.MeasuredProduct
	Clear()
	SetFont(font *windigo.Font)
	RefreshSize()
}

type OilBasedProductView struct {
	*windigo.AutoPanel
	visual_field  *views.BoolCheckboxView
	density_field *qc_ui.MassDataView
}

func newOilBasedProductView(parent *windigo.AutoPanel, ranges_panel *OilBasedProductRangesView) *OilBasedProductView {
	view := new(OilBasedProductView)

	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.visual_field = views.NewBoolCheckboxView(view.AutoPanel, VISUAL_TEXT)

	view.density_field = qc_ui.MassDataView_from_new(view.AutoPanel, ranges_panel)

	view.AutoPanel.Dock(view.visual_field, windigo.Top)
	view.AutoPanel.Dock(view.density_field, windigo.Top)

	return view
}

func (view *OilBasedProductView) Get(base_product product.QCProduct) *product.MeasuredProduct {
	// base_product.Visual = view.visual_field.Checked()
	return newOilBasedProduct(base_product, view.visual_field.Get(), view.density_field.Get())
}

func (view *OilBasedProductView) Clear() {
	view.visual_field.Clear()
	view.density_field.Clear()
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
