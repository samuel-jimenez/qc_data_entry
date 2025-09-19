package subpanels

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type FrictionReducerProductRangesViewer interface {
	*windigo.AutoPanel
	*views.MassRangesView

	SetFont(font *windigo.Font)
	RefreshSize()
	Update(*product.QCProduct)
	Clear()
}

type FrictionReducerProductRangesView struct {
	*windigo.AutoPanel
	*views.MassRangesView

	visual_field *product.ProductAppearanceROView

	viscosity_field,
	// mass_field,
	// sg_field,
	// density_field,
	string_field *views.RangeROView
}

func BuildNewFrictionReducerProductRangesView(parent *windigo.AutoPanel, qc_product *product.QCProduct) *FrictionReducerProductRangesView {

	visual_text := "Visual Inspection"
	viscosity_text := "Viscosity"
	mass_text := "Mass"
	string_text := "String"
	sg_text := "Specific Gravity"
	density_text := "Density"

	group_panel := windigo.NewAutoPanel(parent)

	visual_field := product.BuildNewProductAppearanceROView(group_panel, visual_text, qc_product.Appearance)

	viscosity_field := views.BuildNewRangeROView(group_panel, viscosity_text, qc_product.Viscosity, formats.Format_ranges_viscosity)

	string_field := views.BuildNewRangeROView(group_panel, string_text, qc_product.SG, formats.Format_ranges_string_test)

	mass_field := views.BuildNewRangeROViewMap(group_panel, mass_text, qc_product.SG, formats.Format_mass, formats.Mass_from_sg)

	sg_field := views.BuildNewRangeROView(group_panel, sg_text, qc_product.SG, formats.Format_ranges_sg)
	density_field := views.BuildNewRangeROView(group_panel, density_text, qc_product.Density, formats.Format_ranges_density)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(viscosity_field, windigo.Top)
	group_panel.Dock(mass_field, windigo.Top)
	group_panel.Dock(string_field, windigo.Top)
	group_panel.Dock(density_field, windigo.Bottom)
	group_panel.Dock(sg_field, windigo.Bottom)

	return &FrictionReducerProductRangesView{
		group_panel,
		&views.MassRangesView{Mass_field: &mass_field,
			SG_field:      &sg_field,
			Density_field: &density_field},
		&visual_field,
		&viscosity_field,
		&string_field,
	}

}

func (view *FrictionReducerProductRangesView) SetFont(font *windigo.Font) {
	view.visual_field.SetFont(font)
	view.viscosity_field.SetFont(font)
	view.Mass_field.SetFont(font)
	view.string_field.SetFont(font)
	view.SG_field.SetFont(font)
	view.Density_field.SetFont(font)
}

func (view *FrictionReducerProductRangesView) RefreshSize() {
	view.SetSize(GUI.DATA_FIELD_WIDTH, GUI.GROUP_HEIGHT)
	view.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.BTM_SPACER_WIDTH, GUI.BTM_SPACER_HEIGHT)
	view.SetMarginTop(GUI.GROUP_MARGIN)

	view.visual_field.RefreshSize()
	view.viscosity_field.RefreshSize()
	view.Mass_field.RefreshSize()
	view.string_field.RefreshSize()
	view.SG_field.RefreshSize()
	view.Density_field.RefreshSize()
}

func (view *FrictionReducerProductRangesView) Update(qc_product *product.QCProduct) {
	view.visual_field.Update(qc_product.Appearance)
	view.viscosity_field.Update(qc_product.Viscosity)
	view.string_field.Update(qc_product.String_test)

	// view.MassRangesView.Update(qc_product)
	view.Mass_field.Update(qc_product.SG)
	view.SG_field.Update(qc_product.SG)
	view.Density_field.Update(qc_product.Density)
}

func (view *FrictionReducerProductRangesView) Clear() {
	view.MassRangesView.Clear()
	view.viscosity_field.Clear()
	view.string_field.Clear()
}
