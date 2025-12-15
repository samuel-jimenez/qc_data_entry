package qc_ui

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

/*
 * MeasuredProductView
 *
 */
type MeasuredProductView struct {
	*windigo.AutoPanel
	QC_id        int64
	Sample_point string
	Checked      *windigo.CheckBox
	labels       []*windigo.Label
	ROViews      []*NullFloat64ROView
}

func makelabel(parent windigo.AutoDock, text string, labels []*windigo.Label) []*windigo.Label {
	label := windigo.NewLabel(parent)
	label.SetText(text)
	parent.Dock(label, windigo.Top)

	return append(labels, label)
}

// , field_data nullable.NullFloat64, format func(float64) string
func makeFloat(parent windigo.AutoDock, label *NullFloat64ROView, labels []*NullFloat64ROView) []*NullFloat64ROView {
	// label := NullFloat64ROView_from_new(parent, Measured_Product.PH, formats.Format_ranges_ph)
	parent.Dock(label, windigo.Top)

	return append(labels, label)
}

// func MeasuredProductView_from_AutoPanel(AutoPanel *windigo.AutoPanel, Measured_Product *product.MeasuredProduct) *MeasuredProductView {
func MeasuredProductView_from_new(parent windigo.Controller, Measured_Product *product.MeasuredProduct) *MeasuredProductView {
	view := new(MeasuredProductView)
	view.QC_id = Measured_Product.QC_id
	view.Sample_point = Measured_Product.Sample_point

	// view.AutoPanel = AutoPanel
	view.AutoPanel = windigo.NewAutoPanel(parent)
	// view.SetAndClearStyleBits(w32.WS_BORDER, 0)
	// view.SetPaddingsAll(5)
	view.SetPaddingsAll(5)
	// view.SetMarginsAll(5)

	view.Checked = windigo.NewCheckBox(view)
	view.Checked.SetChecked(true)
	view.Checked.SetText("Selected")
	view.Dock(view.Checked, windigo.Top)

	view.labels = makelabel(view, Measured_Product.Lot_number, view.labels)
	view.labels = makelabel(view, Measured_Product.Product_name, view.labels)
	view.labels = makelabel(view, Measured_Product.Sample_point, view.labels)

	view.ROViews = makeFloat(view, NullFloat64ROView_from_new(view, Measured_Product.PH, formats.Format_ranges_ph), view.ROViews)
	view.ROViews = makeFloat(view, NullFloat64ROView_from_new(view, Measured_Product.SG, formats.Format_ranges_sg), view.ROViews)
	// SG_label := NullFloat64ROView_from_new(view, Measured_Product.QCProduct.SG.Min, formats.Format_ranges_sg),view.ROViews) //TESTING
	view.ROViews = makeFloat(view, NullFloat64ROView_from_new(view, Measured_Product.Density, formats.Format_ranges_density), view.ROViews)

	view.labels = makelabel(view, Measured_Product.Viscosity.String(), view.labels)

	view.labels = makelabel(view, Measured_Product.String_test.String(), view.labels)

	// NullFloat64ROView_from_new(view, Measured_Product.Viscosity, formats.Format_ranges_viscosity)

	// 	view.viscosity_field = qc_ui.RangeROView_from_new(view.AutoPanel, formats.VISCOSITY_TEXT, qc_product.Viscosity, formats.Format_ranges_viscosity)
	//
	// 	view.string_field = qc_ui.RangeROView_from_new(view.AutoPanel, formats.STRING_TEXT_MINI, qc_product.SG, formats.Format_ranges_string_test)
	//
	// 	view.Mass_field = qc_ui.RangeROViewMap_from_new(view.AutoPanel, formats.MASS_TEXT, qc_product.SG, formats.Format_mass, formats.Mass_from_sg)
	// 	view.SG_field = qc_ui.RangeROView_from_new(view.AutoPanel, formats.SG_TEXT, qc_product.SG, formats.Format_ranges_sg)
	// 	view.Density_field = qc_ui.RangeROView_from_new(view.AutoPanel, formats.DENSITY_TEXT, qc_product.Density, formats.Format_ranges_density)

	/*
	 *
	 *
	 *
	 * label := windigo.NewLabel(view)
	 * label.SetText(Measured_Product.Sample_point)
	 * view.Dock(label, windigo.Top)
	 *
	 */

	// view.prod_panel = NewClockTimerView(view)
	// view.Dock(Lot_number_label, windigo.Top)
	// view.Dock(Product_name_label, windigo.Top)
	// view.Dock(Sample_point_label, windigo.Top)
	//
	// view.Dock(label, windigo.Top)
	// view.Dock(SG_label, windigo.Top)
	// view.Dock(Density_label, windigo.Top)
	// view.Dock(String_test_label, windigo.Top)
	// view.Dock(Viscosity_label, windigo.Top)

	// view.Dock(view.prod_panel, windigo.Fill)
	// view.Dock(parent, windigo.Fill)
	// parent.Dock(view, windigo.Fill)

	// view.RefreshSize()
	return view
}

func (view *MeasuredProductView) SetFont(font *windigo.Font) {
	for _, label := range view.labels {
		label.SetFont(font)
	}
	for _, ROView := range view.ROViews {
		ROView.SetFont(font)
	}
}

func (view *MeasuredProductView) RefreshSize() {
	view.SetSize(GUI.GROUP_WIDTH, GUI.GROUP_HEIGHT)
	view.SetPaddingsAll(8)

	for _, label := range view.labels {
		label.SetSize(GUI.GROUP_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	}
	for _, ROView := range view.ROViews {
		ROView.SetSize(GUI.GROUP_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	}

	// view.SetSize(GUI.CLOCK_WIDTH, qc.TOP_SUBPANEL_HEIGHT)
	// 	num_rows := 3
	//
	// 	TOP_SUBPANEL_HEIGHT = GUI.TOP_SPACER_HEIGHT + num_rows*(GUI.PRODUCT_FIELD_HEIGHT+GUI.INTER_SPACER_HEIGHT) + GUI.BTM_SPACER_HEIGHT
	// view.SetSize(GUI.CLOCK_WIDTH, GUI.TOP_SPACER_HEIGHT+3*(GUI.PRODUCT_FIELD_HEIGHT+GUI.INTER_SPACER_HEIGHT)+GUI.BTM_SPACER_HEIGHT)
	// log.Println("WhupsView", GUI.CLOCK_WIDTH, GUI.TOP_SPACER_HEIGHT+3*(GUI.PRODUCT_FIELD_HEIGHT+GUI.INTER_SPACER_HEIGHT)+GUI.BTM_SPACER_HEIGHT)
	// // view.prod_panel.RefreshSize()
	// for _, prod_panel := range view.prod_panels {
	// 	prod_panel.SetSize(500, 600)
	// 	// prod_panel.RefreshSize()
	// }
	// log.Println("WhupsView-clock_panel", view.prod_panels.Width(), view.prod_panels.Height())
}
