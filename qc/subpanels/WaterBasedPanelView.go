package subpanels

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type WaterBasedPanelView struct {
	Update      func(qc_product *product.QCProduct)
	SetFont     func(font *windigo.Font)
	RefreshSize func()
}

func Show_water_based(parent *windigo.AutoPanel, qc_product *product.QCProduct, create_new_product_cb func() product.BaseProduct) *WaterBasedPanelView {

	visual_text := "Visual Inspection"
	sg_text := "SG"
	ph_text := "pH"

	panel := windigo.NewAutoPanel(parent)

	group_panel := windigo.NewAutoPanel(panel)

	ranges_panel := BuildNewWaterBasedProductRangesView(panel, qc_product)

	visual_field := views.NewBoolCheckboxView(group_panel, visual_text)

	ph_field := views.NewNumberEditViewWithChange(group_panel, ph_text, ranges_panel.ph_field)
	sg_field := views.NewNumberEditViewWithChange(group_panel, sg_text, ranges_panel.sg_field)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(ph_field, windigo.Top)
	group_panel.Dock(sg_field, windigo.Top)

	submit_cb := func() {
		measured_product := newWaterBasedProduct(create_new_product_cb(), visual_field.Get(), sg_field.Get(), ph_field.Get())
		if measured_product.Check_data() {
			log.Println("wb sub-data", measured_product)
			// TODO blend013 ensurethis works with testing blends
			// measured_product.Save()
			product.Store(measured_product)
			err := measured_product.Output()
			if err != nil {
				log.Printf("Error: [%s]: %q\n", "WaterBasedProduct.Output", err)
			}
			// * Check storage
			measured_product.CheckStorage()
		}
	}

	clear_cb := func() {
		visual_field.SetChecked(false)
		sg_field.Clear()
		ph_field.Clear()
		ranges_panel.Clear()
	}

	log_cb := func() {
		measured_product := newWaterBasedProduct(create_new_product_cb(), visual_field.Get(), sg_field.Get(), ph_field.Get())
		if measured_product.Check_data() {
			log.Println("wb log-data", measured_product)
			// measured_product.Save()
			product.Store(measured_product)
			measured_product.Export_json()
			// * Check storage
			measured_product.CheckStorage()
		}
	}

	button_dock := GUI.NewMarginalButtonDock(parent, []string{"Submit", "Clear", "Log"}, []int{40, 0, 10}, []func(){submit_cb, clear_cb, log_cb})

	panel.Dock(group_panel, windigo.Left)
	panel.Dock(ranges_panel, windigo.Right)

	parent.Dock(panel, windigo.Top)
	parent.Dock(button_dock, windigo.Top)

	setFont := func(font *windigo.Font) {
		ranges_panel.SetFont(font)
		visual_field.SetFont(font)
		sg_field.SetFont(font)
		ph_field.SetFont(font)
		button_dock.SetFont(font)
	}
	refresh := func() {

		panel.SetSize(GUI.OFF_AXIS, GUI.GROUP_HEIGHT)
		panel.SetMargins(GUI.GROUP_MARGIN, GUI.GROUP_MARGIN, 0, 0)

		group_panel.SetSize(GUI.GROUP_WIDTH, GUI.GROUP_HEIGHT)
		group_panel.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.BTM_SPACER_WIDTH, GUI.BTM_SPACER_HEIGHT)

		ranges_panel.SetSize(GUI.DATA_FIELD_WIDTH, GUI.GROUP_HEIGHT)
		ranges_panel.SetMarginTop(GUI.GROUP_MARGIN)

		visual_field.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)

		sg_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.EDIT_FIELD_HEIGHT)
		ph_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.EDIT_FIELD_HEIGHT)

		button_dock.SetDockSize(GUI.BUTTON_WIDTH, GUI.BUTTON_HEIGHT)

		ranges_panel.RefreshSize()
	}

	return &WaterBasedPanelView{ranges_panel.Update, setFont, refresh}

}
