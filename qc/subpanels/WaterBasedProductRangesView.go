package subpanels

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type WaterBasedProductRangesView struct {
	*windigo.AutoPanel
	ph_field,
	sg_field *views.RangeROView

	Update      func(qc_product *product.QCProduct)
	SetFont     func(font *windigo.Font)
	RefreshSize func()
}

func (data_view WaterBasedProductRangesView) Clear() {
	data_view.ph_field.Clear()
	data_view.sg_field.Clear()
}

func BuildNewWaterBasedProductRangesView(parent *windigo.AutoPanel, qc_product *product.QCProduct) WaterBasedProductRangesView {

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

	update := func(qc_product *product.QCProduct) {

		log.Println("Debug:BuildNewWaterBasedProductRangesView  update", qc_product)
		visual_field.Update(qc_product.Appearance)
		sg_field.Update(qc_product.SG)
		ph_field.Update(qc_product.PH)
	}
	setFont := func(font *windigo.Font) {
		visual_field.SetFont(font) //?TODO
		sg_field.SetFont(font)
		ph_field.SetFont(font)
	}
	refresh := func() {
		group_panel.SetSize(GUI.DATA_FIELD_WIDTH, GUI.GROUP_HEIGHT)
		group_panel.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.BTM_SPACER_WIDTH, GUI.BTM_SPACER_HEIGHT)

		visual_field.RefreshSize()
		sg_field.RefreshSize()
		ph_field.RefreshSize()
	}

	return WaterBasedProductRangesView{group_panel, &ph_field, &sg_field, update, setFont, refresh}

}
