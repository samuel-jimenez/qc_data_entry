package subpanels

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type OilBasedProductRangesViewer interface {
	Update(*product.QCProduct)
	Clear()
	SetFont(font *windigo.Font)
	RefreshSize()
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

func (view *OilBasedProductRangesView) RefreshSize() {
	view.SetSize(GUI.DATA_FIELD_WIDTH, GUI.GROUP_HEIGHT)
	view.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.BTM_SPACER_WIDTH, GUI.BTM_SPACER_HEIGHT)

	view.visual_field.RefreshSize()
	view.Mass_field.RefreshSize()
	view.SG_field.RefreshSize()
	view.Density_field.RefreshSize()
}
