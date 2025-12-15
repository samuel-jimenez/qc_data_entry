package qc

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/qc_ui"
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

func OilBasedPanelView_from_new(parent *windigo.AutoPanel, product_panel *TopPanelView) *OilBasedPanelView {
	view := new(OilBasedPanelView)

	view.AutoPanel = windigo.NewAutoPanel(parent)
	view.product_panel = product_panel

	view.ranges_panel = OilBasedProductRangesView_from_new(view.AutoPanel, view.product_panel.QC_Product)
	view.group_panel = OilBasedProductView_from_new(view.AutoPanel, view.ranges_panel)

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

func (view *OilBasedPanelView) ChangeContainer(qc_product *product.QCProduct) {
}

func (view *OilBasedPanelView) Focus() {
	view.group_panel.FocusVisual()
}

func (view *OilBasedPanelView) submit_data() {
	measured_product := view.group_panel.Get(view.product_panel.BaseProduct())
	valid, err := product.Check_single_data(measured_product, true, true)
	if valid {
		view.product_panel.SetMeasuredProduct(measured_product)
	}
	qc_ui.Check_dupe_data(view, err, measured_product)
}

func (view *OilBasedPanelView) log_data() {
	measured_product := view.group_panel.Get(view.product_panel.BaseProduct())
	valid, err := product.Check_single_data(measured_product, true, false)
	if valid {
		view.product_panel.SetMeasuredProduct(measured_product)
	}
	qc_ui.Check_dupe_data(view, err, measured_product)
}
