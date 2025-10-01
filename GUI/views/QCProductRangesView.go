package views

import (
	"fmt"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/windigo"
)

/*
 * QCProductRangesViewer
 *
 */
type QCProductRangesViewer interface {
	OnExit(arg *windigo.Event)
	Exit()
	Save()
	try_save()
	save_coa()
	load_coa()
}

/*
 * QCProductRangesView
 *
 */
type QCProductRangesView struct {
	*windigo.Form
	qc_product                                              *product.QCProduct
	prod_label                                              *windigo.Label
	radio_dock                                              *product.DiscreteView
	coa_field                                               *windigo.CheckBox
	appearance_dock                                         *product.ProductAppearanceView
	labels                                                  *GUI.TextDock
	ph_dock, sg_dock, density_dock, string_dock, visco_dock RangeView
	button_dock                                             *GUI.ButtonDock
}

func NewQCProductRangesView(parent windigo.Controller, qc_product *product.QCProduct) *QCProductRangesView {

	view := new(QCProductRangesView)
	view.qc_product = qc_product
	view.Form = windigo.NewForm(parent)
	var WindowText string
	// rangeWindow.SetTranslucentBackground()

	if view.qc_product.Product_name_customer != "" {
		WindowText = fmt.Sprintf("%s (%s)", view.qc_product.Product_name, view.qc_product.Product_name_customer)
	} else {
		WindowText = view.qc_product.Product_name
	}
	view.SetSize(GUI.RANGES_WINDOW_WIDTH,
		GUI.RANGES_WINDOW_HEIGHT) // (width, height)
	view.SetText(WindowText)

	dock := windigo.NewSimpleDock(view)
	dock.SetPaddingsAll(GUI.RANGES_WINDOW_PADDING)
	dock.SetPaddingTop(GUI.RANGES_PADDING)

	view.prod_label = windigo.NewLabel(view)
	view.prod_label.SetText(WindowText)
	view.prod_label.SetSize(GUI.OFF_AXIS, GUI.RANGES_FIELD_SMALL_HEIGHT)

	view.radio_dock = product.BuildNewDiscreteView(view, "Type", view.qc_product.Product_type, BLEND_WB, BLEND_OIL, BLEND_FR)
	view.radio_dock.SetSize(GUI.OFF_AXIS, GUI.DISCRETE_FIELD_HEIGHT)
	view.radio_dock.SetItemSize(GUI.PRODUCT_TYPE_WIDTH)
	view.radio_dock.SetPaddingsAll(GUI.GROUPBOX_CUSHION)

	view.coa_field = windigo.NewCheckBox(view)
	view.coa_field.SetText("Save to COA published ranges")
	view.coa_field.SetMarginsAll(GUI.ERROR_MARGIN)
	view.coa_field.SetSize(GUI.OFF_AXIS, GUI.RANGES_FIELD_SMALL_HEIGHT)

	view.appearance_dock = product.BuildNewProductAppearanceView(view, "Appearance", view.qc_product.Appearance)

	view.labels = GUI.NewTextDock(view, "", "Min", "Target", "Max")
	view.labels.SetMarginsAll(GUI.RANGES_PADDING)
	view.labels.SetDockSize(GUI.RANGES_FIELD_WIDTH, GUI.RANGES_FIELD_SMALL_HEIGHT)
	//TODO center
	//TODO layout split n

	view.ph_dock = BuildNewRangeView(view, PH_TEXT, view.qc_product.PH, formats.Format_ranges_ph)
	view.ph_dock.SetLabeledSize(GUI.LABEL_WIDTH, GUI.RANGES_FIELD_WIDTH, GUI.RANGES_FIELD_HEIGHT)

	view.sg_dock = BuildNewRangeView(view, SG_TEXT, view.qc_product.SG, formats.Format_ranges_sg)
	view.sg_dock.SetLabeledSize(GUI.LABEL_WIDTH, GUI.RANGES_FIELD_WIDTH, GUI.RANGES_FIELD_HEIGHT)

	view.density_dock = BuildNewRangeView(view, DENSITY_TEXT, view.qc_product.Density, formats.Format_ranges_density)
	view.density_dock.SetLabeledSize(GUI.LABEL_WIDTH, GUI.RANGES_FIELD_WIDTH, GUI.RANGES_FIELD_HEIGHT)

	view.string_dock = BuildNewRangeView(view, "String Test \n\t at 0.5gpt", view.qc_product.String_test, formats.Format_ranges_string_test)
	view.string_dock.SetLabeledSize(GUI.LABEL_WIDTH, GUI.RANGES_FIELD_WIDTH, GUI.RANGES_FIELD_HEIGHT)
	//TODO store string_amt "at 0.5gpt"

	view.visco_dock = BuildNewRangeView(view, VISCOSITY_TEXT, view.qc_product.Viscosity, formats.Format_ranges_viscosity)
	view.visco_dock.SetLabeledSize(GUI.LABEL_WIDTH, GUI.RANGES_FIELD_WIDTH, GUI.RANGES_FIELD_HEIGHT)

	view.button_dock = GUI.NewButtonDock(view, []string{"OK", "Cancel", "Load CoA Data"}, []func(){view.try_save, view.Exit, view.load_coa})
	view.button_dock.SetDockSize(GUI.RANGES_BUTTON_WIDTH, GUI.RANGES_BUTTON_HEIGHT)
	view.button_dock.SetMarginLeft(GUI.RANGES_PADDING)

	dock.Dock(view.prod_label, windigo.Top)
	dock.Dock(view.radio_dock, windigo.Top)
	dock.Dock(view.coa_field, windigo.Top)
	dock.Dock(view.appearance_dock, windigo.Top)

	dock.Dock(view.labels, windigo.Top)
	dock.Dock(view.ph_dock, windigo.Top)
	dock.Dock(view.sg_dock, windigo.Top)
	dock.Dock(view.density_dock, windigo.Top)
	dock.Dock(view.string_dock, windigo.Top)
	dock.Dock(view.visco_dock, windigo.Top)
	dock.Dock(view.button_dock, windigo.Top)

	return view
}

func ShowNewQCProductRangesView(qc_product *product.QCProduct) {
	view := NewQCProductRangesView(nil, qc_product)

	view.Center()
	view.Show()
	view.OnClose().Bind(view.OnExit)
	view.RunMainLoop()
}

func (view *QCProductRangesView) OnExit(arg *windigo.Event) {
	view.Exit()
}

func (view *QCProductRangesView) Exit() {
	view.Close()
	windigo.Exit()
}

func (view *QCProductRangesView) Save() {
	view.qc_product.Edit(
		view.radio_dock.Get(),
		view.appearance_dock.Get(),
		view.ph_dock.Get(),
		view.sg_dock.Get(),
		view.density_dock.Get(),
		view.string_dock.Get(),
		view.visco_dock.Get(),
	)

	view.qc_product.Upsert()
	threads.Show_status("QC Data Updated")
	view.qc_product.Update()
	view.Exit()
}

func (view *QCProductRangesView) try_save() {
	if !view.radio_dock.Get().Valid {
		view.radio_dock.Error()
		return
	}
	view.radio_dock.Ok()
	if view.coa_field.Checked() {
		view.save_coa()
	} else {
		view.Save()
	}
}

func (view *QCProductRangesView) save_coa() {
	var coa_product product.QCProduct
	coa_product.Product = view.qc_product.Product
	coa_product.Edit(
		view.radio_dock.Get(),
		view.appearance_dock.Get(),
		view.ph_dock.Get(),
		view.sg_dock.Get(),
		view.density_dock.Get(),
		view.string_dock.Get(),
		view.visco_dock.Get(),
	)

	coa_product.Upsert_coa()
	threads.Show_status("COA Data Updated")
	view.qc_product.Select_product_details()
	view.qc_product.Update()
	view.Exit()
}

func (view *QCProductRangesView) load_coa() {
	view.qc_product.ResetQC()
	view.qc_product.Select_product_coa_details()
	view.sg_dock.Set(view.qc_product.SG)
	view.density_dock.Set(view.qc_product.Density)
	view.string_dock.Set(view.qc_product.String_test)
	view.visco_dock.Set(view.qc_product.Viscosity)
}
