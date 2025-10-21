package qc

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type WaterBasedProductRangesViewer interface {
	SetFont(font *windigo.Font)
	RefreshSize()
	Update(*product.QCProduct)
	Clear()
}

type WaterBasedProductRangesView struct {
	*windigo.AutoPanel
	visual_field *product.ProductAppearanceROView
	ph_field,
	sg_field *views.RangeROView
}

func BuildNewWaterBasedProductRangesView(parent *windigo.AutoPanel, qc_product *product.QCProduct) *WaterBasedProductRangesView {

	view := new(WaterBasedProductRangesView)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.visual_field = product.BuildNewProductAppearanceROView(view.AutoPanel, VISUAL_TEXT, qc_product.Appearance)
	view.ph_field = views.RangeROView_from_new(view.AutoPanel, formats.PH_TEXT, qc_product.PH, formats.Format_ranges_ph)
	view.sg_field = views.RangeROView_from_new(view.AutoPanel, formats.SG_TEXT, qc_product.SG, formats.Format_ranges_sg)

	view.AutoPanel.Dock(view.visual_field, windigo.Top)
	view.AutoPanel.Dock(view.ph_field, windigo.Top)
	view.AutoPanel.Dock(view.sg_field, windigo.Top)

	return view

}

func (view *WaterBasedProductRangesView) SetFont(font *windigo.Font) {
	view.visual_field.SetFont(font)
	view.sg_field.SetFont(font)
	view.ph_field.SetFont(font)
}

func (view *WaterBasedProductRangesView) RefreshSize() {
	view.SetSize(GUI.DATA_FIELD_WIDTH, GUI.GROUP_HEIGHT)
	view.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.BTM_SPACER_WIDTH, GUI.BTM_SPACER_HEIGHT)
	view.SetMarginTop(GUI.GROUP_MARGIN)

	view.visual_field.RefreshSize()
	view.sg_field.RefreshSize()
	view.ph_field.RefreshSize()
}

func (view *WaterBasedProductRangesView) Update(qc_product *product.QCProduct) {

	log.Println("Debug:BuildNewWaterBasedProductRangesView  update", qc_product)
	view.visual_field.Update(qc_product.Appearance)
	view.sg_field.Update(qc_product.SG)
	view.ph_field.Update(qc_product.PH)
}

func (view *WaterBasedProductRangesView) Clear() {
	view.ph_field.Clear()
	view.sg_field.Clear()
}
