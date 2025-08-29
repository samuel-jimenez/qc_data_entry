package subpanels

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
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

func (fr_product FrictionReducerProduct) toProduct() product.Product {
	return product.Product{
		BaseProduct: fr_product.Base(),
		PH:          nullable.NewNullFloat64(0, false),
		SG:          nullable.NewNullFloat64(fr_product.sg, true),
		Density:     nullable.NewNullFloat64(formats.Density_from_sg(fr_product.sg), true),
		String_test: nullable.NewNullInt64(fr_product.string_test),
		Viscosity:   nullable.NewNullInt64(fr_product.viscosity),
	}
}

func newFrictionReducerProduct(base_product product.BaseProduct, viscosity, mass, string_test float64) product.Product {

	sg := formats.SG_from_mass(mass)

	return FrictionReducerProduct{base_product, sg, int64(string_test), int64(viscosity)}.toProduct()

}

func (product *FrictionReducerProduct) Check_data() bool {
	return true
}

type FrictionReducerProductViewer interface {
	*windigo.AutoPanel
	Get(base_product product.BaseProduct, replace_sample_point bool) product.Product
	Clear()
	SetFont(font *windigo.Font)
	Refresh()
}

type FrictionReducerProductView struct {
	*windigo.AutoPanel
	visual_field *views.BoolCheckboxView
	viscosity_field,
	string_field *views.NumberEditView
	density_field *views.MassDataView
	ranges_panel  *FrictionReducerProductRangesView
	sample_point  string
}

func BuildNewFrictionReducerProductView(parent *windigo.AutoPanel, sample_point string, ranges_panel *FrictionReducerProductRangesView) *FrictionReducerProductView {

	visual_text := "Visual Inspection"
	viscosity_text := "Viscosity"
	string_text := "String"
	mass_text := "Mass"

	view := new(FrictionReducerProductView)

	group_panel := windigo.NewGroupAutoPanel(parent)
	group_panel.SetText(sample_point)

	visual_field := views.NewBoolCheckboxView(group_panel, visual_text)

	viscosity_field := views.NewNumberEditView(group_panel, viscosity_text)

	mass_field := views.NewNumberEditView(group_panel, mass_text)

	string_field := views.NewNumberEditViewWithChange(group_panel, string_text, ranges_panel.string_field)

	density_field := views.NewMassDataView(group_panel, ranges_panel, mass_field)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(viscosity_field, windigo.Top)
	group_panel.Dock(density_field, windigo.Top)
	group_panel.Dock(string_field, windigo.Top)

	view.AutoPanel = group_panel
	view.visual_field = visual_field
	view.viscosity_field = viscosity_field
	view.string_field = string_field
	view.density_field = density_field
	view.ranges_panel = ranges_panel
	view.sample_point = sample_point
	return view

}

func (view *FrictionReducerProductView) Get(base_product product.BaseProduct, replace_sample_point bool) product.Product {
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

func (view *FrictionReducerProductView) Refresh() {
	view.SetSize(GUI.GROUP_WIDTH, GUI.GROUP_HEIGHT)
	view.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.BTM_SPACER_WIDTH, GUI.BTM_SPACER_HEIGHT)

	view.visual_field.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)
	view.viscosity_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.EDIT_FIELD_HEIGHT)
	view.density_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.DATA_SUBFIELD_WIDTH, GUI.DATA_UNIT_WIDTH, GUI.EDIT_FIELD_HEIGHT)
	view.string_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.EDIT_FIELD_HEIGHT)

}

func (view *FrictionReducerProductView) Clear() {

	view.visual_field.SetChecked(false)
	view.viscosity_field.Clear()

	view.density_field.Clear()

	view.string_field.Clear()

	view.ranges_panel.Clear()

}

/*
func (view FrictionReducerProductView) Clear() {
	view.MassRangesView.Clear()
	view.viscosity_field.Clear()
	view.string_field.Clear()
}*/
