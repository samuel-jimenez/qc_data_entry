package subpanels

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type OilBasedPanelViewer interface {
	SetFont(font *windigo.Font)
	RefreshSize()
	Update(*product.QCProduct)
}

type OilBasedPanelView struct {
	*windigo.AutoPanel
	group_panel  *OilBasedProductView
	ranges_panel *OilBasedProductRangesView
	button_dock  *GUI.ButtonDock
}

func Show_oil_based(parent *windigo.AutoPanel, qc_product *product.QCProduct, create_new_product_cb func() product.BaseProduct) *OilBasedPanelView {

	view := new(OilBasedPanelView)

	panel := windigo.NewAutoPanel(parent)

	ranges_panel := BuildNewOilBasedProductRangesView(panel, qc_product)
	group_panel := newOilBasedProductView(panel, ranges_panel)

	submit_cb := func() {
		measured_product := group_panel.Get(create_new_product_cb())
		if measured_product.Check_data() {
			log.Println("ob sub data", measured_product)
			// TODO blend013 ensurethis works with testing blends
			// measured_product.Save()
			product.Store(measured_product)
			err := measured_product.Output()
			if err != nil {
				log.Printf("Error: [%s]: %q\n", "OilBasedProduct.Output", err)
			}
			// * Check storage
			measured_product.CheckStorage()
		}
	}

	log_cb := func() {
		measured_product := group_panel.Get(create_new_product_cb())
		if measured_product.Check_data() {
			log.Println("ob log data", measured_product)
			// measured_product.Save()
			product.Store(measured_product)
			measured_product.Export_json()
			// * Check storage
			measured_product.CheckStorage()
		}
	}

	button_dock := GUI.NewMarginalButtonDock(parent, []string{"Submit", "Clear", "Log"}, []int{40, 0, 10}, []func(){submit_cb, group_panel.Clear, log_cb})

	panel.Dock(group_panel, windigo.Left)
	panel.Dock(ranges_panel, windigo.Right)

	parent.Dock(panel, windigo.Top)
	parent.Dock(button_dock, windigo.Top)

	view.AutoPanel = panel
	view.group_panel = group_panel
	view.ranges_panel = ranges_panel
	view.button_dock = button_dock

	return view

}

func (view *OilBasedPanelView) SetFont(font *windigo.Font) {
	view.group_panel.SetFont(font)
	view.ranges_panel.SetFont(font)
	view.button_dock.SetFont(font)
}

func (view *OilBasedPanelView) RefreshSize() {

	view.SetSize(GUI.OFF_AXIS, GUI.GROUP_HEIGHT)
	view.SetMargins(GUI.GROUP_MARGIN, GUI.GROUP_MARGIN, 0, 0)

	view.group_panel.RefreshSize()

	view.ranges_panel.RefreshSize()

	view.button_dock.SetDockSize(GUI.BUTTON_WIDTH, GUI.BUTTON_HEIGHT)

}

func (view *OilBasedPanelView) Update(qc_product *product.QCProduct) {
	view.ranges_panel.Update(qc_product)
}
