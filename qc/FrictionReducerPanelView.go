package qc

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/util/math"
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

	button_dock_totes, button_dock_cars *GUI.ButtonDock

	product_panel *TopPanelView
}

// create table product_line (product_id integer not null primary key, product_name text);
func Show_fr(parent *windigo.AutoPanel, product_panel *TopPanelView) *FrictionReducerPanelView {
	DELTA_DIFF_VISCO := 200.

	top_text := "Top"
	// bottom_text := "Bottom"
	bottom_text := "Btm"

	view := new(FrictionReducerPanelView)
	view.product_panel = product_panel

	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.ranges_panel = BuildNewFrictionReducerProductRangesView(view.AutoPanel, view.product_panel.QC_Product)

	view.top_group = BuildNewFrictionReducerProductView(view.AutoPanel, top_text, view.ranges_panel)
	view.bottom_group = BuildNewFrictionReducerProductView(view.AutoPanel, bottom_text, view.ranges_panel)

	view.top_group.viscosity_field.Entangle(view.bottom_group.viscosity_field, view.ranges_panel.viscosity_field, DELTA_DIFF_VISCO)

	view.button_dock_totes = GUI.NewMarginalButtonDock(parent, SUBMIT_CLEAR_BTN, []int{40, 0}, []func(){view.tote_cb, view.Clear})
	view.button_dock_cars = GUI.NewMarginalButtonDock(parent, SUBMIT_CLEAR_BTN, []int{40, 0}, []func(){view.submit_cb, view.Clear})
	view.button_dock_totes.Hide()

	view.AutoPanel.Dock(view.top_group, windigo.Left)
	view.AutoPanel.Dock(view.bottom_group, windigo.Left)
	view.AutoPanel.Dock(view.ranges_panel, windigo.Right)

	parent.Dock(view.AutoPanel, windigo.Top)
	parent.Dock(view.button_dock_totes, windigo.Top)
	parent.Dock(view.button_dock_cars, windigo.Top)

	return view
}

func (view *FrictionReducerPanelView) SetFont(font *windigo.Font) {
	view.top_group.SetFont(font)
	view.bottom_group.SetFont(font)
	view.ranges_panel.SetFont(font)
	view.button_dock_totes.SetFont(font)
	view.button_dock_cars.SetFont(font)
}

func (view *FrictionReducerPanelView) RefreshSize() {
	view.SetSize(GUI.OFF_AXIS, GUI.GROUP_HEIGHT)
	view.SetMargins(GUI.GROUP_MARGIN, GUI.GROUP_MARGIN, 0, 0)

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
		view.button_dock_totes.Hide()
	} else {
		view.bottom_group.Hide()
		view.button_dock_cars.Hide()
		// // TODO
		// if qc_product.Container_type == product.CONTAINER_TOTE {
		view.button_dock_totes.Show()
		// } else { // product.CONTAINER_SAMPLE
		// 					NO COA
		// no storage
		// TODO blend013 ensurethis works with testing blends

		// }
	}
}

func (view *FrictionReducerPanelView) Clear() {
	view.top_group.Clear()
	view.bottom_group.Clear()
	view.ranges_panel.Clear()
}

func (view *FrictionReducerPanelView) submit_cb() {
	base_product := view.product_panel.BaseProduct()
	log.Println("DEBUG: FrictionReducerPanelView.submit_cb base_product", base_product)

	top_product := view.top_group.Get(base_product, true)
	bottom_product := view.bottom_group.Get(base_product, true)
	check_dual_data(top_product, bottom_product)
}

func (view *FrictionReducerPanelView) tote_cb() {
	base_product := view.product_panel.BaseProduct()
	log.Println("DEBUG: FrictionReducerPanelView.submit_cb.tote base_product", base_product)

	top_product := view.top_group.Get(base_product, false)
	if top_product.Check_data() {
		top_product.Save()
		err := top_product.Output()
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "top_product.Output", err)
		}
		// TODO view.component_panel.saVE

	}
}

// TODO
func check_dual_data(top_product, bottom_product product.MeasuredProduct) {
	// DELTA_DIFF_VISCO := 200
	var DELTA_DIFF_VISCO int64 = 200 // go sucks

	if math.Abs(top_product.Viscosity.Diff(bottom_product.Viscosity)) <= DELTA_DIFF_VISCO &&
		top_product.Check_data() && bottom_product.Check_data() {

		log.Println("debug: Check_data", top_product)
		// top_product.Save()
		// bottom_product.Save()

		// TODO blend013 ensurethis works with testing blends
		product.Store(top_product, bottom_product)

		// qc_product.Store(9)

		err := top_product.Printout()
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "top_product.Printout", err)
		}

		log.Println("debug: Check_data", bottom_product)
		err = bottom_product.Printout()
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "bottom_product.Printout", err)
		}
		// TODO find closest: RMS?
		err = bottom_product.Output_sample()
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "bottom_product.Output_sample", err)
		}
		err = bottom_product.Save_xl()
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "bottom_product.Save_xl", err)
		}

		// * Check storage
		bottom_product.CheckStorage()

	} else { // TODO show confirm box
		log.Println("ERROR: check_dual_data-Viscosity", top_product.Lot_number, top_product.Product_name, top_product.Viscosity, bottom_product.Viscosity)
	}
}
