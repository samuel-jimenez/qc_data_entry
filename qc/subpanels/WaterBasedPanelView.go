package subpanels

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type WaterBasedPanelViewer interface {
	SetFont(font *windigo.Font)
	RefreshSize()
	Update(*product.QCProduct)
}

type WaterBasedPanelView struct {
	Update      func(qc_product *product.QCProduct)
	SetFont     func(font *windigo.Font)
	RefreshSize func()
}

func Show_water_based(parent *windigo.AutoPanel, qc_product *product.QCProduct, create_new_product_cb func() product.BaseProduct) *WaterBasedPanelView {

	panel := windigo.NewAutoPanel(parent)

	ranges_panel := BuildNewWaterBasedProductRangesView(panel, qc_product)
	group_panel := newWaterBasedProductView(panel, ranges_panel)

	submit_cb := func() {

		measured_product := group_panel.Get(create_new_product_cb())
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
		group_panel.Clear()
		ranges_panel.Clear()
	}

	log_cb := func() {
		measured_product := group_panel.Get(create_new_product_cb())
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

		button_dock.SetFont(font)
	}
	refresh := func() {

		panel.SetSize(GUI.OFF_AXIS, GUI.GROUP_HEIGHT)
		panel.SetMargins(GUI.GROUP_MARGIN, GUI.GROUP_MARGIN, 0, 0)

		button_dock.SetDockSize(GUI.BUTTON_WIDTH, GUI.BUTTON_HEIGHT)

		group_panel.RefreshSize()
		ranges_panel.RefreshSize()
	}

	return &WaterBasedPanelView{ranges_panel.Update, setFont, refresh}

}
