package view

import (
	"fmt"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/windigo"
)

/*
 * QCProductRangesView
 *
 */

func ShowNewQCProductRangesView(qc_product *product.QCProduct) {

	rangeWindow := windigo.NewForm(nil)
	var WindowText string
	// rangeWindow.SetTranslucentBackground()

	if qc_product.Product_name_customer != "" {
		WindowText = fmt.Sprintf("%s (%s)", qc_product.Product_name, qc_product.Product_name_customer)
	} else {
		WindowText = qc_product.Product_name
	}
	rangeWindow.SetSize(GUI.RANGES_WINDOW_WIDTH,
		GUI.RANGES_WINDOW_HEIGHT) // (width, height)
	rangeWindow.SetText(WindowText)

	dock := windigo.NewSimpleDock(rangeWindow)
	dock.SetPaddingsAll(GUI.RANGES_WINDOW_PADDING)
	dock.SetPaddingTop(GUI.RANGES_PADDING)

	prod_label := windigo.NewLabel(rangeWindow)
	prod_label.SetText(WindowText)
	prod_label.SetSize(GUI.OFF_AXIS, GUI.RANGES_FIELD_SMALL_HEIGHT)

	radio_dock := product.BuildNewDiscreteView(rangeWindow, "Type", qc_product.Product_type, []string{"Water Based", "Oil Based", "Friction Reducer"})
	radio_dock.SetSize(GUI.OFF_AXIS, GUI.DISCRETE_FIELD_HEIGHT)
	radio_dock.SetItemSize(GUI.PRODUCT_TYPE_WIDTH)
	radio_dock.SetPaddingsAll(GUI.GROUPBOX_CUSHION)

	coa_field := windigo.NewCheckBox(rangeWindow)
	coa_field.SetText("Save to COA published ranges")
	coa_field.SetMarginsAll(GUI.ERROR_MARGIN)
	coa_field.SetSize(GUI.OFF_AXIS, GUI.RANGES_FIELD_SMALL_HEIGHT)

	appearance_dock := product.BuildNewProductAppearanceView(rangeWindow, "Appearance", qc_product.Appearance)

	labels := GUI.NewTextDock(rangeWindow, []string{"", "Min", "Target", "Max"})
	labels.SetMarginsAll(GUI.RANGES_PADDING)
	labels.SetDockSize(GUI.RANGES_FIELD_WIDTH, GUI.RANGES_FIELD_SMALL_HEIGHT)
	//TODO center
	//TODO layout split n

	ph_dock := BuildNewRangeView(rangeWindow, "pH", qc_product.PH, formats.Format_ranges_ph)
	ph_dock.SetLabeledSize(GUI.LABEL_WIDTH, GUI.RANGES_FIELD_WIDTH, GUI.RANGES_FIELD_HEIGHT)

	sg_dock := BuildNewRangeView(rangeWindow, "Specific Gravity", qc_product.SG, formats.Format_ranges_sg)
	sg_dock.SetLabeledSize(GUI.LABEL_WIDTH, GUI.RANGES_FIELD_WIDTH, GUI.RANGES_FIELD_HEIGHT)

	density_dock := BuildNewRangeView(rangeWindow, "Density", qc_product.Density, formats.Format_ranges_density)
	density_dock.SetLabeledSize(GUI.LABEL_WIDTH, GUI.RANGES_FIELD_WIDTH, GUI.RANGES_FIELD_HEIGHT)

	string_dock := BuildNewRangeView(rangeWindow, "String Test \n\t at 0.5gpt", qc_product.String_test, formats.Format_ranges_string_test)
	string_dock.SetLabeledSize(GUI.LABEL_WIDTH, GUI.RANGES_FIELD_WIDTH, GUI.RANGES_FIELD_HEIGHT)
	//TODO store string_amt "at 0.5gpt"

	visco_dock := BuildNewRangeView(rangeWindow, "Viscosity", qc_product.Viscosity, formats.Format_ranges_viscosity)
	visco_dock.SetLabeledSize(GUI.LABEL_WIDTH, GUI.RANGES_FIELD_WIDTH, GUI.RANGES_FIELD_HEIGHT)

	exit := func() {
		rangeWindow.Close()
		windigo.Exit()
	}
	save := func() {
		qc_product.Edit(
			radio_dock.Get(),
			appearance_dock.Get(),
			ph_dock.Get(),
			sg_dock.Get(),
			density_dock.Get(),
			string_dock.Get(),
			visco_dock.Get(),
		)

		qc_product.Upsert()
		threads.Show_status("QC Data Updated")
		qc_product.Update()
		exit()
	}
	save_coa := func() {
		var coa_product product.QCProduct
		coa_product.Product = qc_product.Product
		coa_product.Edit(
			radio_dock.Get(),
			appearance_dock.Get(),
			ph_dock.Get(),
			sg_dock.Get(),
			density_dock.Get(),
			string_dock.Get(),
			visco_dock.Get(),
		)

		coa_product.Upsert_coa()
		threads.Show_status("COA Data Updated")
		qc_product.Select_product_details()
		qc_product.Update()
		exit()
	}

	try_save := func() {
		if radio_dock.Get().Valid {
			radio_dock.Ok()
			if coa_field.Checked() {
				save_coa()
			} else {
				save()
			}
		} else {
			radio_dock.Error()
		}
	}

	load_coa := func() {
		qc_product.Reset()
		qc_product.Select_product_coa_details()
		appearance_dock.Set(qc_product.Appearance)
		ph_dock.Set(qc_product.PH)
		sg_dock.Set(qc_product.SG)
		density_dock.Set(qc_product.Density)
		string_dock.Set(qc_product.String_test)
		visco_dock.Set(qc_product.Viscosity)
	}

	button_dock := GUI.NewButtonDock(rangeWindow, []string{"OK", "Cancel", "Load CoA Data"}, []func(){try_save, exit, load_coa})
	button_dock.SetDockSize(GUI.RANGES_BUTTON_WIDTH, GUI.RANGES_BUTTON_HEIGHT)
	button_dock.SetMarginLeft(GUI.RANGES_PADDING)

	dock.Dock(prod_label, windigo.Top)
	dock.Dock(radio_dock, windigo.Top)
	dock.Dock(coa_field, windigo.Top)
	dock.Dock(appearance_dock, windigo.Top)

	dock.Dock(labels, windigo.Top)
	dock.Dock(ph_dock, windigo.Top)
	dock.Dock(sg_dock, windigo.Top)
	dock.Dock(density_dock, windigo.Top)
	dock.Dock(string_dock, windigo.Top)
	dock.Dock(visco_dock, windigo.Top)
	dock.Dock(button_dock, windigo.Top)

	rangeWindow.Center()
	rangeWindow.Show()
	rangeWindow.OnClose().Bind(
		func(arg *windigo.Event) {
			exit()
		})
	rangeWindow.RunMainLoop() // Must call to start event loop.
}
