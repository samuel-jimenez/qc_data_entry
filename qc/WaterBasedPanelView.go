package qc

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/qc_ui"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type WaterBasedPanelViewer interface {
	SetFont(font *windigo.Font)
	RefreshSize()
	Update(*product.QCProduct)
	Clear()
	submit_data()
	log_data()
}

type WaterBasedPanelView struct {
	*windigo.AutoPanel
	group_panel   *WaterBasedProductView
	ranges_panel  *WaterBasedProductRangesView
	button_dock   *GUI.ButtonDock
	product_panel *TopPanelView
}

func WaterBasedPanelView_from_new(parent *windigo.AutoPanel, product_panel *TopPanelView) *WaterBasedPanelView {
	view := new(WaterBasedPanelView)

	view.product_panel = product_panel

	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.ranges_panel = WaterBasedProductRangesView_from_new(view.AutoPanel, view.product_panel.QC_Product)
	view.group_panel = WaterBasedProductView_from_new(view.AutoPanel, view.ranges_panel)

	view.button_dock = GUI.NewMarginalButtonDock(parent, SUBMIT_CLEAR_LOG_BTN, []int{40, 0, 10}, []func(){view.submit_data, view.Clear, view.log_data})

	view.AutoPanel.Dock(view.group_panel, windigo.Left)
	view.AutoPanel.Dock(view.ranges_panel, windigo.Right)

	parent.Dock(view.AutoPanel, windigo.Top)
	parent.Dock(view.button_dock, windigo.Top)

	return view
}

func (view *WaterBasedPanelView) SetFont(font *windigo.Font) {
	view.group_panel.SetFont(font)
	view.ranges_panel.SetFont(font)
	view.button_dock.SetFont(font)
}

func (view *WaterBasedPanelView) RefreshSize() {
	view.SetSize(GUI.OFF_AXIS, GUI.GROUP_HEIGHT)
	view.SetMargins(GUI.GROUP_MARGIN, GUI.GROUP_MARGIN, 0, 0)

	view.group_panel.RefreshSize()

	view.ranges_panel.RefreshSize()

	view.button_dock.SetDockSize(GUI.BUTTON_WIDTH, GUI.BUTTON_HEIGHT)
}

func (view *WaterBasedPanelView) Update(qc_product *product.QCProduct) {
	view.ranges_panel.Update(qc_product)
}

func (view *WaterBasedPanelView) Clear() {
	view.group_panel.Clear()
	view.ranges_panel.Clear()
}

func (view *WaterBasedPanelView) ChangeContainer(qc_product *product.QCProduct) {
}

func (view *WaterBasedPanelView) Focus() {
	view.group_panel.FocusVisual()
}

func (view *WaterBasedPanelView) submit_data() {
	measured_product := view.group_panel.Get(view.product_panel.BaseProduct())
	valid, err := product.Check_single_data(measured_product, true, true)
	if valid {
		view.product_panel.SetMeasuredProduct(measured_product)
	}
	qc_ui.Check_dupe_data(view,err, measured_product)
}

func (view *WaterBasedPanelView) log_data() {
	measured_product := view.group_panel.Get(view.product_panel.BaseProduct())
	valid, err := product.Check_single_data(measured_product, true, false)
	if valid {
		view.product_panel.SetMeasuredProduct(measured_product)
	}
	qc_ui.Check_dupe_data(view,err, measured_product)
}


