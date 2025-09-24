package qc

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



	view := new(FrictionReducerProductRangesView)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.MassRangesView = &views.MassRangesView{}

	view.visual_field = product.BuildNewProductAppearanceROView(view.AutoPanel, VISUAL_TEXT, qc_product.Appearance)

	view.viscosity_field = views.BuildNewRangeROView(view.AutoPanel, views.VISCOSITY_TEXT, qc_product.Viscosity, formats.Format_ranges_viscosity)

	view.string_field = views.BuildNewRangeROView(view.AutoPanel, STRING_TEXT, qc_product.SG, formats.Format_ranges_string_test)

	view.Mass_field = views.BuildNewRangeROViewMap(view.AutoPanel, views.MASS_TEXT, qc_product.SG, formats.Format_mass, formats.Mass_from_sg)
	view.SG_field = views.BuildNewRangeROView(view.AutoPanel, views.SG_TEXT, qc_product.SG, formats.Format_ranges_sg)
	view.Density_field = views.BuildNewRangeROView(view.AutoPanel, views.DENSITY_TEXT, qc_product.Density, formats.Format_ranges_density)

	view.AutoPanel.Dock(view.visual_field, windigo.Top)
	view.AutoPanel.Dock(view.viscosity_field, windigo.Top)
	view.AutoPanel.Dock(view.Mass_field, windigo.Top)
	view.AutoPanel.Dock(view.string_field, windigo.Top)
	view.AutoPanel.Dock(view.Density_field, windigo.Bottom)
	view.AutoPanel.Dock(view.SG_field, windigo.Bottom)

	return view

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
	view.MassRangesView.Update(qc_product)
}

func (view *FrictionReducerProductRangesView) Clear() {
	view.MassRangesView.Clear()
	view.viscosity_field.Clear()
	view.string_field.Clear()
}
