package qc

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/qc_ui"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type FrictionReducerPanelViewer interface {
	SetFont(font *windigo.Font)
	RefreshSize()
	Update(*product.QCProduct)
	ChangeContainer(*product.QCProduct)
	Clear()
	submit_cb()
	tote_cb()
}

type FrictionReducerPanelView struct {
	*windigo.AutoPanel
	top_group, bottom_group *FrictionReducerProductView
	ranges_panel            *FrictionReducerProductRangesView

	button_dock_samples, button_dock_totes, button_dock_cars *GUI.ButtonDock

	product_panel *TopPanelView
}

// create table product_line (product_id integer not null primary key, product_name text);
func FrictionReducerPanelView_from_new(parent *windigo.AutoPanel, product_panel *TopPanelView) *FrictionReducerPanelView {
	DELTA_DIFF_VISCO := 200.

	top_text := "Top"
	// bottom_text := "Bottom"
	bottom_text := "Btm"

	view := new(FrictionReducerPanelView)
	view.product_panel = product_panel

	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.ranges_panel = FrictionReducerProductRangesView_from_new(view.AutoPanel, view.product_panel.QC_Product)

	view.top_group = FrictionReducerProductView_from_new(view.AutoPanel, top_text, view.ranges_panel)
	view.bottom_group = FrictionReducerProductView_from_new(view.AutoPanel, bottom_text, view.ranges_panel)

	view.top_group.viscosity_field.Entangle(view.bottom_group.viscosity_field, view.ranges_panel.viscosity_field, DELTA_DIFF_VISCO)
	view.top_group.Interleave(view.bottom_group)

	view.button_dock_samples = GUI.NewMarginalButtonDock(parent, SUBMIT_CLEAR_LOG_BTN, []int{40, 0, 10}, []func(){view.submit_sample, view.Clear, view.log_sample})
	view.button_dock_totes = GUI.NewMarginalButtonDock(parent, SUBMIT_CLEAR_LOG_BTN, []int{40, 0, 10}, []func(){view.submit_tote, view.Clear, view.log_tote})
	view.button_dock_cars = GUI.NewMarginalButtonDock(parent, SUBMIT_CLEAR_LOG_BTN, []int{40, 0, 10}, []func(){view.submit_car, view.Clear, view.log_car})

	view.button_dock_samples.Hide()
	view.button_dock_totes.Hide()

	view.AutoPanel.Dock(view.top_group, windigo.Left)
	view.AutoPanel.Dock(view.bottom_group, windigo.Left)
	view.AutoPanel.Dock(view.ranges_panel, windigo.Right)

	parent.Dock(view.AutoPanel, windigo.Top)
	parent.Dock(view.button_dock_samples, windigo.Top)
	parent.Dock(view.button_dock_totes, windigo.Top)
	parent.Dock(view.button_dock_cars, windigo.Top)

	// view.AddShortcut(windigo.Shortcut{Key: windigo.KeyReturn}, view.submit_car)
	// TODO bool
	// TODO split out button_dock

	return view
}

func (view *FrictionReducerPanelView) SetFont(font *windigo.Font) {
	view.top_group.SetFont(font)
	view.bottom_group.SetFont(font)
	view.ranges_panel.SetFont(font)
	view.button_dock_samples.SetFont(font)
	view.button_dock_totes.SetFont(font)
	view.button_dock_cars.SetFont(font)
}

func (view *FrictionReducerPanelView) RefreshSize() {
	view.SetSize(GUI.OFF_AXIS, GUI.GROUP_HEIGHT)
	view.SetMargins(GUI.GROUP_MARGIN, GUI.GROUP_MARGIN, 0, 0)

	view.button_dock_samples.SetDockSize(GUI.BUTTON_WIDTH, GUI.BUTTON_HEIGHT)
	view.button_dock_totes.SetDockSize(GUI.BUTTON_WIDTH, GUI.BUTTON_HEIGHT)
	view.button_dock_cars.SetDockSize(GUI.BUTTON_WIDTH, GUI.BUTTON_HEIGHT)

	view.top_group.RefreshSize()
	view.bottom_group.RefreshSize()

	view.ranges_panel.RefreshSize()
}

func (view *FrictionReducerPanelView) Update(QC_Product *product.QCProduct) {
	view.ranges_panel.Update(QC_Product)
}

func (view *FrictionReducerPanelView) ChangeContainer(qc_product *product.QCProduct) {
	if qc_product.Container_type == product.CONTAINER_RAILCAR {
		view.bottom_group.Show()
		view.button_dock_cars.Show()
		view.button_dock_samples.Hide()
		view.button_dock_totes.Hide()
	} else {
		view.bottom_group.Hide()
		view.button_dock_cars.Hide()
		if qc_product.Container_type == product.CONTAINER_TOTE {
			view.button_dock_totes.Show()
			view.button_dock_samples.Hide()
		} else { // product.CONTAINER_SAMPLE
			view.button_dock_samples.Show()
			view.button_dock_totes.Hide()
		}
	}
}

func (view *FrictionReducerPanelView) Clear() {
	view.top_group.Clear()
	view.bottom_group.Clear()
	view.ranges_panel.Clear()
}

func (view *FrictionReducerPanelView) Focus() {
	view.top_group.FocusVisual()
}

func (view *FrictionReducerPanelView) submit_car() {
	view.proc_car(true)
}

func (view *FrictionReducerPanelView) log_car() {
	view.proc_car(false)
}

func (view *FrictionReducerPanelView) submit_sample() {
	view.proc_single(false, true)
}

func (view *FrictionReducerPanelView) log_sample() {
	view.proc_single(false, false)
}

func (view *FrictionReducerPanelView) submit_tote() {
	view.proc_single(true, true)
}

func (view *FrictionReducerPanelView) log_tote() {
	view.proc_single(true, false)
}

func (view *FrictionReducerPanelView) proc_car(print_p bool) {
	base_product := view.product_panel.BaseProduct()
	top_product := view.top_group.Get(base_product, true)
	bottom_product := view.bottom_group.Get(base_product, true)
	valid, err := product.Check_dual_data(top_product, bottom_product, print_p)
	if valid {
		view.product_panel.SetMeasuredProduct(top_product, bottom_product)
	} else {
		// TODO show confirm box?
		log.Println("ERROR: check_dual_data-Viscosity", top_product.Lot_number, top_product.Product_name, top_product.Viscosity, bottom_product.Viscosity)
		qc_ui.Check_dupe_data(view, err, print_p, top_product, bottom_product)
	}
}

func (view *FrictionReducerPanelView) proc_single(store_p, print_p bool) {
	measured_product := view.top_group.Get(view.product_panel.BaseProduct(), false)
	valid, err := product.Check_single_data(measured_product, store_p, print_p)
	if valid {
		view.product_panel.SetMeasuredProduct(measured_product)
	}
	qc_ui.Check_dupe_data(view, err, print_p, measured_product)
}
