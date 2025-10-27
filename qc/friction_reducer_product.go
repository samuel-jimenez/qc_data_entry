package qc

import (
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/qc_ui"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type FrictionReducerProduct struct {
	product.BaseProduct
	sg          float64
	string_test int64
	viscosity   int64
}

func (fr_product FrictionReducerProduct) toProduct() product.MeasuredProduct {
	return product.MeasuredProduct{
		BaseProduct: fr_product.Base(),
		PH:          nullable.NewNullFloat64(0, false),
		SG:          nullable.NewNullFloat64(fr_product.sg, true),
		Density:     nullable.NewNullFloat64(formats.Density_from_sg(fr_product.sg), true),
		String_test: nullable.NewNullInt64(fr_product.string_test),
		Viscosity:   nullable.NewNullInt64(fr_product.viscosity),
	}
}

func newFrictionReducerProduct(base_product product.BaseProduct, viscosity, mass, string_test float64) product.MeasuredProduct {

	sg := formats.SG_from_mass(mass)

	return FrictionReducerProduct{base_product, sg, int64(string_test), int64(viscosity)}.toProduct()

}

func (product *FrictionReducerProduct) Check_data() bool {
	return true
}

type FrictionReducerProductViewer interface {
	*windigo.AutoPanel
	Get(base_product product.BaseProduct, replace_sample_point bool) product.MeasuredProduct
	Clear()
	SetFont(font *windigo.Font)
	Refresh()
}

type FrictionReducerProductView struct {
	*windigo.AutoPanel
	visual_field *views.BoolCheckboxView
	viscosity_field,
	string_field *GUI.NumberEditView
	density_field *qc_ui.MassDataView
	sample_point  string
}

func BuildNewFrictionReducerProductView(parent *windigo.AutoPanel, sample_point string, ranges_panel *FrictionReducerProductRangesView) *FrictionReducerProductView {

	view := new(FrictionReducerProductView)

	view.AutoPanel = windigo.NewGroupAutoPanel(parent)
	view.AutoPanel.SetText(sample_point)

	view.visual_field = views.NewBoolCheckboxView(view.AutoPanel, VISUAL_TEXT)

	view.viscosity_field = GUI.NumberEditView_from_new(view.AutoPanel, formats.VISCOSITY_TEXT)

	view.string_field = GUI.NumberEditView_with_Change_from_new(view.AutoPanel, formats.STRING_TEXT_MINI, ranges_panel.string_field)

	view.density_field = qc_ui.MassDataView_from_new(view.AutoPanel, ranges_panel)
	view.density_field.SetZAfter(view.viscosity_field)

	view.AutoPanel.Dock(view.visual_field, windigo.Top)
	view.AutoPanel.Dock(view.viscosity_field, windigo.Top)
	view.AutoPanel.Dock(view.density_field, windigo.Top)
	view.AutoPanel.Dock(view.string_field, windigo.Top)

	view.sample_point = strings.ToUpper(sample_point)
	return view

}

func (view *FrictionReducerProductView) Get(base_product product.BaseProduct, replace_sample_point bool) product.MeasuredProduct {
	base_product.Visual = view.visual_field.Checked()
	if replace_sample_point {
		base_product.Sample_point = view.sample_point
	}
	return newFrictionReducerProduct(base_product, view.viscosity_field.Get(), view.density_field.Get(), view.string_field.Get())
}

func (view *FrictionReducerProductView) SetFont(font *windigo.Font) {
	view.AutoPanel.SetFont(font)
	view.visual_field.SetFont(font)
	view.viscosity_field.SetFont(font)
	view.density_field.SetFont(font)
	view.string_field.SetFont(font)
}

func (view *FrictionReducerProductView) RefreshSize() {
	view.SetSize(GUI.GROUP_WIDTH, GUI.GROUP_HEIGHT)
	view.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.BTM_SPACER_WIDTH, GUI.BTM_SPACER_HEIGHT)

	view.visual_field.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)
	view.viscosity_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.EDIT_FIELD_HEIGHT)
	view.density_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.DATA_SUBFIELD_WIDTH, GUI.DATA_UNIT_WIDTH, GUI.EDIT_FIELD_HEIGHT)
	view.string_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.EDIT_FIELD_HEIGHT)

}

func (view *FrictionReducerProductView) Clear() {
	view.visual_field.Clear()
	view.viscosity_field.Clear()
	view.density_field.Clear()
	view.string_field.Clear()
}
