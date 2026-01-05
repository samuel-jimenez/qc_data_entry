package qc_ui

import (
	"database/sql"
	"errors"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/util"
	"github.com/samuel-jimenez/windigo"
)

func Check_dupe_data(parent windigo.Controller, err error, print_p bool, products ...*product.MeasuredProduct) {
	if errors.Is(err, product.DupeSampleError) {
		DupeMeasurementView := DupeMeasurementView_from_products(parent, products)
		DupeMeasurementView.SetModal(false)
		DupeMeasurementView.RefreshSize()
		DupeMeasurementView.Show()
	}
}

func Check_dupe_sample_ids(parent windigo.Controller, err error, old_QC_id, new_QC_id int64, action func()) {
	if old_QC_id != new_QC_id && errors.Is(err, product.DupeSampleError) {
		DupeMeasurementView := DupeMeasurementView_from_ids(parent, old_QC_id, new_QC_id, action)
		DupeMeasurementView.SetModal(false)
		DupeMeasurementView.RefreshSize()
		DupeMeasurementView.Show()
	}
}

func Check_dupe_sample_lots(parent windigo.Controller, err error, old_Lot_id, new_Lot_id int64, action func()) {
	if old_Lot_id != new_Lot_id && errors.Is(err, DB.DupeSampleError) { // todo combine`	errors

		DupeMeasurementView := DupeMeasurementView_from_lots(parent, old_Lot_id, new_Lot_id, action)
		DupeMeasurementView.SetModal(false)
		DupeMeasurementView.RefreshSize()
		DupeMeasurementView.Show()
	}
}

/*
 * DupeMeasurementView
 *
 */
type DupeMeasurementView struct {
	*windigo.AutoDialog

	products []*product.MeasuredProduct

	top_message, prompt_message *windigo.Label

	prod_panels []*MeasuredProductDiffView

	prod_bar *windigo.TabView

	button_panel *windigo.AutoPanel
	buttons      []*windigo.PushButton
}

func DupeMeasurementView_from_new(parent windigo.Controller,
) *DupeMeasurementView {
	window_title := "Duplicate Measurement"
	view := new(DupeMeasurementView)

	view.AutoDialog = windigo.AutoDialog_from_new(parent)
	view.SetText(window_title)
	view.top_message = windigo.NewLabel(view)
	view.top_message.SetText("Measurement exists:")
	view.top_message.SetSize(100, 50)
	view.top_message.SetStyleLeft()

	// CardView
	view.prod_bar = windigo.CardView_from_new(
		view,
		windigo.Pen_from_Color(windigo.Black),
		windigo.Pen_from_Color(windigo.Black),
	)

	view.prompt_message = windigo.NewLabel(view)
	view.prompt_message.SetText("Check lot number and sample point or replace existing measurement.")
	view.prompt_message.SetSize(100, 80)
	view.prompt_message.SetStyleLeft()

	view.button_panel = windigo.NewAutoPanel(view)

	view.Dock(view.top_message, windigo.Top)
	view.Dock(view.prod_bar, windigo.Top)
	// TODO ;abels
	view.Dock(view.prod_bar.Panels(), windigo.Top)
	view.Dock(view.prompt_message, windigo.Top)
	view.Dock(view.button_panel, windigo.Bottom)

	accept_button := windigo.NewPushButton(view.button_panel)
	accept_button.SetText("OK")
	view.buttons = append(view.buttons, accept_button)

	view.button_panel.Dock(accept_button, windigo.Left)

	// event handling
	accept_button.OnClick().Bind(view.EventClose)
	view.SetButtons(accept_button, nil)

	return view
}

func DupeMeasurementView_from_products(parent windigo.Controller,
	Measured_Products []*product.MeasuredProduct,
) *DupeMeasurementView {
	view := DupeMeasurementView_from_new(parent)
	view.products = Measured_Products

	for _, Measured_Product := range Measured_Products {
		prod_panel := MeasuredProductDiffView_from_combined(view.prod_bar.Panels(), Measured_Product)
		view.prod_bar.AttachPanel(Measured_Product.Sample_point, prod_panel)
		view.prod_panels = append(view.prod_panels, prod_panel)
	}
	view.RefreshSize()

	resample_button := windigo.NewPushButton(view.button_panel)
	resample_button.SetText("Replace")
	view.buttons = append(view.buttons, resample_button)

	view.button_panel.Dock(resample_button, windigo.Left)

	// event handling
	resample_button.OnClick().Bind(view.resample_button_OnClick)

	return view
}

func DupeMeasurementView_from_ids(parent windigo.Controller,
	old_QC_id, new_QC_id int64,
	action func(),
) *DupeMeasurementView {
	view := DupeMeasurementView_from_new(parent)

	// 			view.products = Measured_Products
	//
	//
	// 			for _, Measured_Product := range Measured_Products {
	// 				prod_panel := MeasuredProductDiffView_from_ids(view.prod_bar.Panels(), Measured_Product)
	// 				view.prod_bar.AttachPanel(Measured_Product.Sample_point, prod_panel)
	// 				view.prod_panels = append(view.prod_panels, prod_panel)
	// 			}

	prod_panel, Sample_point := MeasuredProductDiffView_from_ids(view.prod_bar.Panels(), old_QC_id, new_QC_id)
	view.prod_bar.AttachPanel(Sample_point, prod_panel)
	view.prod_panels = append(view.prod_panels, prod_panel)
	view.RefreshSize()

	overwrite_button := windigo.NewPushButton(view.button_panel)
	overwrite_button.SetText("Overwrite")
	view.buttons = append(view.buttons, overwrite_button)

	view.button_panel.Dock(overwrite_button, windigo.Left)

	// event handling
	overwrite_button.OnClick().Bind(func(*windigo.Event) {
		DeleteSample(old_QC_id)
		action()
		view.Close()
	})

	return view
}

func DupeMeasurementView_from_lots(parent windigo.Controller,
	old_Lot_id, new_Lot_id int64,
	action func(),
) *DupeMeasurementView {
	view := DupeMeasurementView_from_new(parent)

	proc_name := "DupeMeasurementView_from_lots"
	var sample_pts []string
	var old_ids []int64
	DB.Forall_exit(proc_name,
		util.NOOP,
		func(row *sql.Rows) (err error) {
			var sample_pt string
			err = row.Scan(
				&sample_pt,
			)
			sample_pts = append(sample_pts, sample_pt)
			return
		},
		DB.DB_Select_Product__qc_samples_dupe_pts, old_Lot_id, new_Lot_id)

	for _, Sample_point := range sample_pts {
		old_QC_id, _ := product.CheckDupes(
			old_Lot_id, Sample_point,
		)
		new_QC_id, _ := product.CheckDupes(
			new_Lot_id, Sample_point,
		)

		prod_panel, Sample_point := MeasuredProductDiffView_from_ids(view.prod_bar.Panels(), old_QC_id, new_QC_id)
		view.prod_bar.AttachPanel(Sample_point, prod_panel)
		view.prod_panels = append(view.prod_panels, prod_panel)
		old_ids = append(old_ids, old_QC_id)
	}

	view.RefreshSize()

	overwrite_button := windigo.NewPushButton(view.button_panel)
	overwrite_button.SetText("Overwrite")
	view.buttons = append(view.buttons, overwrite_button)

	view.button_panel.Dock(overwrite_button, windigo.Left)

	// event handling
	overwrite_button.OnClick().Bind(func(*windigo.Event) {
		for _, old_QC_id := range old_ids {
			DeleteSample(old_QC_id)
		}
		action()
		view.Close()
	})

	return view
}

func (view *DupeMeasurementView) SetSize(w, h int) {
	view.Dialog.SetSize(w+GUI.WINDOW_FUDGE_MARGIN_W, h+GUI.WINDOW_FUDGE_MARGIN_H)
}

func (view *DupeMeasurementView) SetFont(font *windigo.Font) {
	view.prod_bar.SetFont(font)
	for _, prod_panel := range view.prod_panels {
		prod_panel.SetFont(font)
	}
}

func (view *DupeMeasurementView) RefreshSize() {
	view.SetSize(2*GUI.GROUP_WIDTH, GUI.GROUP_HEIGHT+2*GUI.RANGES_BUTTON_HEIGHT+2*GUI.PRODUCT_FIELD_HEIGHT)

	view.button_panel.SetSize(2*GUI.GROUP_WIDTH, GUI.RANGES_BUTTON_HEIGHT)
	for _, button := range view.buttons {
		button.SetSize(GUI.RANGES_BUTTON_WIDTH, GUI.RANGES_BUTTON_HEIGHT)
	}

	view.prod_bar.SetSize(GUI.GROUP_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.prod_bar.Panels().SetSize(GUI.OFF_AXIS, GUI.GROUP_HEIGHT)

	view.prod_bar.SetPanelSize(2 * GUI.GROUP_WIDTH)
	view.prod_bar.SetPanelMargin(50)

	for _, prod_panel := range view.prod_panels {
		prod_panel.RefreshSize()
	}
}

func (view *DupeMeasurementView) resample_button_OnClick(*windigo.Event) {
	product.Update(view.products...)
	view.Close()
}
