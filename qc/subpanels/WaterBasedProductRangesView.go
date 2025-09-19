package subpanels

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

	visual_text := "Visual Inspection"
	sg_text := "SG"
	ph_text := "pH"

	group_panel := windigo.NewAutoPanel(parent)

	// visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)
	visual_field := product.BuildNewProductAppearanceROView(group_panel, visual_text, qc_product.Appearance)

	ph_field := views.BuildNewRangeROView(group_panel, ph_text, qc_product.PH, formats.Format_ranges_ph)
	sg_field := views.BuildNewRangeROView(group_panel, sg_text, qc_product.SG, formats.Format_ranges_sg)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(ph_field, windigo.Top)
	group_panel.Dock(sg_field, windigo.Top)

	return &WaterBasedProductRangesView{group_panel, &visual_field, &ph_field, &sg_field}

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
