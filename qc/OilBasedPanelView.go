package qc

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
	Clear()
	submit_data()
	log_data()
}

type OilBasedPanelView struct {
	*windigo.AutoPanel
	group_panel  *OilBasedProductView
	ranges_panel *OilBasedProductRangesView
	button_dock  *GUI.ButtonDock

	product_panel *TopPanelView
}

func Show_oil_based(parent *windigo.AutoPanel, product_panel *TopPanelView) *OilBasedPanelView {

	view := new(OilBasedPanelView)

	view.AutoPanel = windigo.NewAutoPanel(parent)
	view.product_panel = product_panel

	view.ranges_panel = BuildNewOilBasedProductRangesView(view.AutoPanel, view.product_panel.QC_Product)
	view.group_panel = newOilBasedProductView(view.AutoPanel, view.ranges_panel)

	view.button_dock = GUI.NewMarginalButtonDock(parent, SUBMIT_CLEAR_LOG_BTN, []int{40, 0, 10}, []func(){view.submit_data, view.Clear, view.log_data})

	view.AutoPanel.Dock(view.group_panel, windigo.Left)
	view.AutoPanel.Dock(view.ranges_panel, windigo.Right)

	parent.Dock(view.AutoPanel, windigo.Top)
	parent.Dock(view.button_dock, windigo.Top)

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

func (view *OilBasedPanelView) Clear() {
	view.group_panel.Clear()
	view.ranges_panel.Clear()
}

func (view *OilBasedPanelView) submit_data() {
	measured_product := view.group_panel.Get(view.product_panel.BaseProduct())
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

func (view *OilBasedPanelView) log_data() {
	measured_product := view.group_panel.Get(view.product_panel.BaseProduct())
	if measured_product.Check_data() {
		log.Println("ob log data", measured_product)
		// measured_product.Save()
		product.Store(measured_product)
		measured_product.Export_json()
		// * Check storage
		measured_product.CheckStorage()
	}
}
