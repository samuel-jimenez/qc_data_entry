package qc

import (
	"database/sql"
	"encoding/json"
	"log"
	"maps"
	"slices"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/QR"
	"github.com/samuel-jimenez/qc_data_entry/blender/blendbound"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/windigo"
)

/*
 * TopPanelViewer
 *
 */
type TopPanelViewer interface {
	windigo.Controller
	GoInbound()
	GoInternal()
	SetFont(font *windigo.Font)
	RefreshSize()
	SetMainWindow(mainWindow *windigo.Form)
	SetTitle(title string)
	BaseProduct() product.BaseProduct
	product_field_pop_data(str string)
	product_field_text_pop_data(str string)
	lot_field_pop_data(str string)
	lot_field_text_pop_data(str string)
	customer_field_pop_data(str string)
	customer_field_text_pop_data(str string)
	sample_field_pop_data(str string)
	sample_field_text_pop_data(str string)
	tester_field_pop_data(str string)
	tester_field_text_pop_data(str string)
	PopQRData(product QR.QRJson)
	SetUpdate(func())

	ChangeContainer(qc_product *product.QCProduct)
	SetCurrentTab(int)
}

/*
 * TopPanelView
 *
 */
type TopPanelView struct {
	*windigo.AutoPanel
	Inbound_Lot *blendbound.InboundLot
	QC_Product  *product.QCProduct

	mainWindow *windigo.Form

	product_panel_0,
	product_panel_0_0, product_panel_1_0,
	product_panel_0_1, product_panel_1_1,
	product_panel_0_2 *windigo.AutoPanel

	internal_product_field, customer_field,
	lot_field, sample_field,
	// tester_field,

	testing_lot_field, //inbound_product_field,
	inbound_lot_field, inbound_container_field *GUI.ComboBox

	inbound_product_field *windigo.LabeledEdit
	tester_field          *GUI.SearchBox

	ranges_button, reprint_button, inbound_button,
	sample_button, release_button, internal_button *windigo.PushButton
	container_field *product.DiscreteView
	clock_panel     *views.ClockTimerView

	controls []windigo.Controller

	ChangeContainer func(qc_product *product.QCProduct)
	SetCurrentTab   func(int)
}

func NewTopPanelView(parent windigo.Controller) *TopPanelView {
	view := new(TopPanelView)
	view.QC_Product = product.NewQCProduct()

	proc_name := "TopPanelView.show_window"

	product_data := make(map[string]int)
	inbound_lot_data := make(map[string]*blendbound.InboundLot)
	inbound_container_data := make(map[string]*blendbound.InboundLot)
	inbound_test := "BSQL%"

	product_text := "Product"
	lot_text := "Lot Number"
	sample_text := "Sample Point"
	customer_text := "Customer Name"
	tester_text := "Tester"

	ranges_text := "Ranges"
	reprint_text := "Reprint"
	inbound_text := "Inbound"
	sample_button_text := "Sample"
	release_button_text := "Release"

	inbound_lot_text := "Inbound Lot"
	internal_text := "Internal"
	container_text := "Container"

	product_panel := windigo.NewAutoPanel(parent)

	product_panel_0 := windigo.NewAutoPanel(product_panel)

	//TODO array db_select_all_product

	product_panel_0_0 := windigo.NewAutoPanel(product_panel_0)

	internal_product_field := GUI.NewComboBox(product_panel_0_0, product_text)
	customer_field := GUI.NewComboBox(product_panel_0_0, customer_text)

	product_panel_1_0 := windigo.NewAutoPanel(product_panel_0)
	product_panel_1_0.Hide()

	testing_lot_field := GUI.NewListComboBox(product_panel_1_0, lot_text)
	// inbound_product_field := GUI.NewComboBox(product_panel_1_0, product_text)
	// inbound_product_field := windigo.NewLabeledComboBox(product_panel_1_0, product_text)
	inbound_product_field := windigo.NewLabeledEdit(product_panel_1_0, product_text)
	inbound_product_field.SetReadOnly(true)

	product_panel_0_1 := windigo.NewAutoPanel(product_panel_0)

	lot_field := GUI.NewComboBox(product_panel_0_1, lot_text)
	sample_field := GUI.NewComboBox(product_panel_0_1, sample_text)

	product_panel_1_1 := windigo.NewAutoPanel(product_panel_0)
	product_panel_1_1.Hide()

	// inbound_lot_field := GUI.NewComboBox(product_panel_1_1, product_text)
	inbound_lot_field := GUI.NewComboBox(product_panel_1_1, inbound_lot_text)
	inbound_container_field := GUI.NewComboBox(product_panel_1_1, container_text)
	// inbound_lot_field container_text

	product_panel_0_2 := windigo.NewAutoPanel(product_panel_0)

	tester_field := GUI.NewLabeledSearchBoxFromQuery(product_panel_0_2, tester_text, DB.DB_Select_all_qc_tester)

	clock_panel := views.NewClockTimerView(product_panel_0)

	ranges_button := windigo.NewPushButton(product_panel)
	ranges_button.SetText(ranges_text)

	sample_button := windigo.NewPushButton(product_panel)
	sample_button.SetText(sample_button_text)
	sample_button.Hide()

	release_button := windigo.NewPushButton(product_panel)
	release_button.SetText(release_button_text)
	release_button.Hide()

	container_field := product.BuildNewDiscreteView(product_panel, "Container Type", view.QC_Product.Container_type, []string{"Sample", "Tote", "Railcar"}) // bs.container_types

	reprint_button := windigo.NewPushButton(product_panel)
	reprint_button.SetText(reprint_text)

	inbound_button := windigo.NewPushButton(product_panel)
	inbound_button.SetText(inbound_text)

	internal_button := windigo.NewPushButton(product_panel)
	internal_button.SetText(internal_text)
	internal_button.Hide()

	// build object
	view.AutoPanel = product_panel
	view.product_panel_0 = product_panel_0
	view.product_panel_0_0 = product_panel_0_0
	view.product_panel_1_0 = product_panel_1_0
	view.product_panel_0_1 = product_panel_0_1
	view.product_panel_1_1 = product_panel_1_1
	view.product_panel_0_2 = product_panel_0_2
	view.clock_panel = clock_panel

	view.internal_product_field = internal_product_field
	view.customer_field = customer_field
	view.lot_field = lot_field
	view.sample_field = sample_field
	view.tester_field = tester_field
	view.inbound_lot_field = inbound_lot_field
	view.testing_lot_field = testing_lot_field
	view.inbound_product_field = inbound_product_field
	view.inbound_container_field = inbound_container_field

	view.ranges_button = ranges_button
	view.sample_button = sample_button
	view.container_field = container_field
	view.reprint_button = reprint_button
	view.release_button = release_button
	view.inbound_button = inbound_button
	view.internal_button = internal_button

	//
	//
	// Dock
	//
	//
	product_panel_0_0.Dock(internal_product_field, windigo.Left)
	product_panel_0_0.Dock(customer_field, windigo.Left)

	product_panel_0_1.Dock(lot_field, windigo.Left)
	product_panel_0_1.Dock(sample_field, windigo.Left)

	product_panel_0_2.Dock(tester_field, windigo.Left)

	product_panel_1_0.Dock(testing_lot_field, windigo.Left)
	product_panel_1_0.Dock(inbound_product_field, windigo.Left)

	product_panel_1_1.Dock(inbound_lot_field, windigo.Left)
	product_panel_1_1.Dock(inbound_container_field, windigo.Left)

	product_panel_0.Dock(clock_panel, windigo.Right)
	product_panel_0.Dock(product_panel_0_0, windigo.Top)
	product_panel_0.Dock(product_panel_1_0, windigo.Top)
	product_panel_0.Dock(product_panel_0_1, windigo.Top)
	product_panel_0.Dock(product_panel_1_1, windigo.Top)
	product_panel_0.Dock(product_panel_0_2, windigo.Top)

	product_panel.Dock(product_panel_0, windigo.Top)
	product_panel.Dock(ranges_button, windigo.Left)
	product_panel.Dock(sample_button, windigo.Left)
	product_panel.Dock(container_field, windigo.Left)
	product_panel.Dock(release_button, windigo.Left)
	product_panel.Dock(reprint_button, windigo.Left)
	product_panel.Dock(inbound_button, windigo.Left)
	product_panel.Dock(internal_button, windigo.Left)

	//
	//
	//
	//
	//
	// combobox
	GUI.Fill_combobox_from_query_rows(internal_product_field, func(row *sql.Rows) error {
		var (
			id                   int
			internal_name        string
			product_moniker_name string
		)
		if err := row.Scan(&id, &internal_name, &product_moniker_name); err != nil {
			return err
		}
		name := product_moniker_name + " " + internal_name
		product_data[name] = id

		internal_product_field.AddItem(name)
		return nil
	}, DB.DB_Select_product_info)

	GUI.Fill_combobox_from_query(testing_lot_field,
		DB.DB_Select_lot_list_for_name_status, inbound_test, product.Status_BLENDED)
	testing_lot_field.SetSelectedItem(-1)

	proc_name = "TopPanelView.FillInbound"
	DB.Forall_err(proc_name,
		func() {
			inbound_container_field.DeleteAllItems()
			inbound_lot_field.DeleteAllItems()
		},
		func(row *sql.Rows) error {
			Inbound, err := blendbound.NewInboundLotFromRow(row)
			if err != nil {
				return err
			}
			inbound_lot_field.AddItem(Inbound.Lot_number)
			inbound_lot_data[Inbound.Lot_number] = Inbound
			inbound_container_data[Inbound.Container_name] = Inbound
			return nil
		},
		DB.DB_Select_inbound_lot_status, blendbound.Status_AVAILABLE)

	for _, Container_name := range slices.Sorted(maps.Keys(inbound_container_data)) {
		inbound_container_field.AddItem(Container_name)
	}

	//
	//
	//
	//
	// functionality

	//TODO product.NewBin
	// func (product Product) GetStorage(numSamples int) int {
	// add button to product.NewBin to account for unlogged samples
	// newbin := func() {
	// 	TODO confirm
	//
	// 	proc_name := "product.NewBin"
	// 	if view.QC_Product.Product_id == DB.INVALID_ID {
	// 		return
	// 	}
	//
	// 	qc_sample_storage_id := view.QC_Product.NewStorageBin()
	//
	// 	if view.QC_Product.QC_id == DB.INVALID_ID {
	// 		return
	// 	}
	//
	// 	DB.Update(proc_name,
	// 		DB.DB_Update_qc_samples_storage,
	// 		view.QC_Product.QC_id, qc_sample_storage_id)
	//
	// }

	internal_product_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		view.product_field_pop_data(internal_product_field.GetSelectedItem())

	})
	internal_product_field.OnKillFocus().Bind(func(e *windigo.Event) {
		view.product_field_text_pop_data(internal_product_field.Text())
	})

	lot_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		view.lot_field_pop_data(lot_field.GetSelectedItem())
	})
	lot_field.OnKillFocus().Bind(func(e *windigo.Event) {
		view.lot_field_text_pop_data(lot_field.Text())
	})

	customer_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		view.customer_field_pop_data(customer_field.GetSelectedItem())
	})
	customer_field.OnKillFocus().Bind(func(e *windigo.Event) {
		view.customer_field_text_pop_data(customer_field.Text())
	})

	sample_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		view.sample_field_pop_data(sample_field.GetSelectedItem())
	})
	sample_field.OnKillFocus().Bind(func(e *windigo.Event) {
		view.sample_field_text_pop_data(sample_field.Text())
	})

	testing_lot_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		view.Inbound_Lot = nil
		view.QC_Product.Update_testing_lot(testing_lot_field.GetSelectedItem())
		inbound_product_field.SetText(view.QC_Product.Product_name)

		blend := view.QC_Product.Blend
		if blend == nil || len(blend.Components) < 1 {
			log.Println("Err: testing_lot_field.OnSelectedChange", view.QC_Product, blend)
			return
		}

		component := blend.Components[0]
		inbound_lot_field.SetText(component.Lot_name)
		inbound_container_field.SetText(component.Container_name)
		// inbound_product_field.SetText(component.Component_name)

		// update component list
		view.QC_Product.ResetQC()
		view.QC_Product.Container_type = product.DiscreteFromInt(CONTAINER_SAMPLE)
		view.QC_Product.Update()

	})

	inbound_lot_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		testing_lot_field.SetSelectedItem(-1)
		view.Inbound_Lot = inbound_lot_data[inbound_lot_field.GetSelectedItem()]
		if view.Inbound_Lot == nil {
			return
		}

		inbound_product_field.SetText(view.Inbound_Lot.Product_name)
		inbound_container_field.SetText(view.Inbound_Lot.Container_name)
	})
	inbound_container_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		testing_lot_field.SetSelectedItem(-1)
		view.Inbound_Lot = inbound_container_data[inbound_container_field.GetSelectedItem()]
		if view.Inbound_Lot == nil {
			return
		}

		inbound_product_field.SetText(view.Inbound_Lot.Product_name)
		inbound_lot_field.SetText(view.Inbound_Lot.Lot_number)
	})

	tester_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		view.tester_field_pop_data(tester_field.GetSelectedItem())
	})
	tester_field.OnKillFocus().Bind(func(e *windigo.Event) {
		view.tester_field_text_pop_data(tester_field.Text())
	})

	ranges_button.OnClick().Bind(func(e *windigo.Event) {
		if view.QC_Product.Product_name != "" {
			views.ShowNewQCProductRangesView(view.QC_Product)
			log.Println("debug: ranges_button-product_lot", view.QC_Product)
		}
	})

	sample_button.OnClick().Bind(func(e *windigo.Event) {
		if view.Inbound_Lot == nil {
			return
		}

		Inbound__Lot_ := view.Inbound_Lot
		view.Inbound_Lot = nil

		Inbound__Lot_.Update_status(blendbound.Status_SAMPLED)
		Inbound__Lot_.Quality_test()
		GUI.Fill_combobox_from_query(testing_lot_field,
			DB.DB_Select_lot_list_for_name_status, inbound_test, product.Status_BLENDED)
		delete(inbound_container_data, Inbound__Lot_.Container_name)
		delete(inbound_lot_data, Inbound__Lot_.Lot_number)
		testing_lot_field.SetSelectedItem(-1)
		inbound_lot_field.SetSelectedItem(-1)
		inbound_container_field.SetSelectedItem(-1)
		inbound_product_field.SetText("")
	})

	release_button.OnClick().Bind(func(e *windigo.Event) {
		// TODO make a method
		if view.Inbound_Lot == nil {
			if view.QC_Product == nil {
				return
			}
			blend := view.QC_Product.Blend
			if blend == nil || len(blend.Components) < 1 {
				log.Println("Err: release_button.OnClick", view.QC_Product, blend)
				return
			}

			view.Inbound_Lot = blendbound.NewInboundLotFromBlendComponent(&blend.Components[0])
			if view.Inbound_Lot == nil {
				return
			}
		}

		Inbound__Lot_ := view.Inbound_Lot
		view.Inbound_Lot = nil
		log.Printf("Info: Releasing %s: %s - %s\n", Inbound__Lot_.Container_name, Inbound__Lot_.Product_name, Inbound__Lot_.Lot_number)

		Inbound__Lot_.Update_status(blendbound.Status_UNAVAILABLE)

		log.Println("TRACE: DEBUG: UPdate componnent DB_Update_lot_list__component_status", proc_name, Inbound__Lot_.Lot_number, Inbound__Lot_)

		if err := product.Release_testing_lot(Inbound__Lot_.Lot_number); err != nil {
			log.Println("error[]%S]:", proc_name, err)
			return
		}
		GUI.Fill_combobox_from_query(testing_lot_field,
			DB.DB_Select_lot_list_for_name_status, inbound_test, product.Status_BLENDED)
		delete(inbound_container_data, Inbound__Lot_.Container_name)
		delete(inbound_lot_data, Inbound__Lot_.Lot_number)
		testing_lot_field.SetSelectedItem(-1)
		inbound_lot_field.SetSelectedItem(-1)
		inbound_container_field.SetSelectedItem(-1)
		inbound_product_field.SetText("")
	})

	reprint_button.OnClick().Bind(func(e *windigo.Event) {
		if view.QC_Product.Lot_number != "" {
			log.Println("debug: reprint_button")
			view.QC_Product.Reprint()
		}
	})

	inbound_button.OnClick().Bind(func(e *windigo.Event) {
		view.GoInbound()
	})

	internal_button.OnClick().Bind(func(e *windigo.Event) {
		view.GoInternal()
	})

	container_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		view.QC_Product.Container_type = container_field.Get()
		view.ChangeContainer(view.QC_Product)
	})

	return view
}

func (view *TopPanelView) GoInbound() {
	view.product_panel_1_0.Show()
	view.product_panel_1_1.Show()
	view.sample_button.Show()
	view.release_button.Show()
	view.internal_button.Show()

	view.product_panel_0_0.Hide()
	view.product_panel_0_1.Hide()
	view.ranges_button.Hide()
	view.reprint_button.Hide()
	view.inbound_button.Hide()

	view.QC_Product.Container_type = product.DiscreteFromInt(CONTAINER_SAMPLE)
	view.ChangeContainer(view.QC_Product)
}

func (view *TopPanelView) GoInternal() {
	view.product_panel_1_0.Hide()
	view.product_panel_1_1.Hide()
	view.sample_button.Hide()
	view.release_button.Hide()
	view.internal_button.Hide()

	view.product_panel_0_0.Show()
	view.product_panel_0_1.Show()
	view.ranges_button.Show()
	view.reprint_button.Show()
	view.inbound_button.Show()
}

func (view *TopPanelView) SetFont(font *windigo.Font) {

	view.internal_product_field.SetFont(font)
	view.customer_field.SetFont(font)
	view.lot_field.SetFont(font)
	view.sample_field.SetFont(font)
	view.tester_field.SetFont(font)
	view.inbound_lot_field.SetFont(font)
	view.testing_lot_field.SetFont(font)
	view.inbound_product_field.SetFont(font)
	view.inbound_container_field.SetFont(font)

	view.ranges_button.SetFont(font)
	view.sample_button.SetFont(font)
	view.container_field.SetFont(font)
	view.reprint_button.SetFont(font)
	view.release_button.SetFont(font)
	view.inbound_button.SetFont(font)
	view.internal_button.SetFont(font)
	for _, control := range view.controls {
		control.SetFont(font)
	}
}

// container_item_width

// func (view *TopPanelView) RefreshSize() {
func (view *TopPanelView) RefreshSize(font_size int) {
	var (
		top_panel_height,

		hpanel_width,

		top_spacer_height,
		top_subpanel_height,

		container_item_width,

		reprint_button_margin_l int
	)

	num_rows := 3

	top_spacer_height = 20
	top_subpanel_height = top_spacer_height + num_rows*(GUI.PRODUCT_FIELD_HEIGHT+INTER_SPACER_HEIGHT) + BTM_SPACER_HEIGHT

	top_panel_height = top_subpanel_height + num_rows*GUI.GROUPBOX_CUSHION + GUI.PRODUCT_FIELD_HEIGHT

	hpanel_width = GUI.TOP_PANEL_WIDTH - GUI.CLOCK_WIDTH

	container_item_width = 6 * font_size

	reprint_button_margin_l = 2*GUI.LABEL_WIDTH + GUI.PRODUCT_FIELD_WIDTH + TOP_PANEL_INTER_SPACER_WIDTH - GUI.SMOL_BUTTON_WIDTH - GUI.DISCRETE_FIELD_WIDTH

	view.SetSize(GUI.TOP_PANEL_WIDTH, top_panel_height)

	view.product_panel_0.SetSize(GUI.TOP_PANEL_WIDTH, top_subpanel_height)

	view.product_panel_0_0.SetSize(hpanel_width, GUI.PRODUCT_FIELD_HEIGHT)
	view.product_panel_0_0.SetMarginTop(top_spacer_height)
	view.product_panel_0_0.SetMarginLeft(GUI.HPANEL_MARGIN)

	view.customer_field.SetMarginLeft(TOP_PANEL_INTER_SPACER_WIDTH)
	view.sample_field.SetMarginLeft(TOP_PANEL_INTER_SPACER_WIDTH)

	view.product_panel_0_1.SetSize(hpanel_width, GUI.PRODUCT_FIELD_HEIGHT)
	view.product_panel_0_1.SetMarginTop(INTER_SPACER_HEIGHT)
	view.product_panel_0_1.SetMarginLeft(GUI.HPANEL_MARGIN)

	view.internal_product_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.customer_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.sample_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.product_panel_1_0.SetSize(hpanel_width, GUI.PRODUCT_FIELD_HEIGHT)
	view.product_panel_1_0.SetMarginTop(top_spacer_height)
	view.product_panel_1_0.SetMarginLeft(GUI.HPANEL_MARGIN)

	view.testing_lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.inbound_product_field.SetMarginLeft(TOP_PANEL_INTER_SPACER_WIDTH)
	view.inbound_product_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.product_panel_1_1.SetSize(hpanel_width, GUI.PRODUCT_FIELD_HEIGHT)
	view.product_panel_1_1.SetMarginTop(INTER_SPACER_HEIGHT)
	view.product_panel_1_1.SetMarginLeft(GUI.HPANEL_MARGIN)

	view.inbound_lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.inbound_container_field.SetMarginLeft(TOP_PANEL_INTER_SPACER_WIDTH)
	view.inbound_container_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.product_panel_0_2.SetSize(hpanel_width, GUI.PRODUCT_FIELD_HEIGHT)
	view.product_panel_0_2.SetMarginTop(INTER_SPACER_HEIGHT)
	view.product_panel_0_2.SetMarginLeft(GUI.HPANEL_MARGIN)

	view.tester_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.clock_panel.RefreshSize()

	view.ranges_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
	view.ranges_button.SetMarginsAll(BUTTON_MARGIN)

	view.sample_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
	view.sample_button.SetMarginsAll(BUTTON_MARGIN)

	view.container_field.SetSize(GUI.DISCRETE_FIELD_WIDTH, GUI.OFF_AXIS)
	view.container_field.SetItemSize(container_item_width)
	view.container_field.SetPaddingsAll(GUI.GROUPBOX_CUSHION)
	view.container_field.SetPaddingLeft(0)

	view.reprint_button.SetMarginsAll(BUTTON_MARGIN)
	view.reprint_button.SetMarginLeft(reprint_button_margin_l)
	view.reprint_button.SetSize(REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.release_button.SetMarginsAll(BUTTON_MARGIN)
	view.release_button.SetMarginLeft(reprint_button_margin_l)
	view.release_button.SetSize(REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.inbound_button.SetMarginsAll(BUTTON_MARGIN)
	view.inbound_button.SetSize(REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.internal_button.SetMarginsAll(BUTTON_MARGIN)
	view.internal_button.SetSize(REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

}

func (view *TopPanelView) SetMainWindow(mainWindow *windigo.Form) {
	view.mainWindow = mainWindow
}

func (view *TopPanelView) SetTitle(title string) {
	if view.mainWindow == nil {
		return
	}
	view.mainWindow.SetText(title)
}

func (view *TopPanelView) BaseProduct() product.BaseProduct {
	log.Println("Debug: TopPanelView-BaseProduct", view.QC_Product.Base())

	// there has to be better way
	// if !product_panel.QC_Product.Tester.Valid{
	// 	tester_field.ShowDropdown(true)
	// 	return nil
	// }
	view.QC_Product.Valid = view.QC_Product.Tester.Valid

	// only last test holds dropdown
	if !view.QC_Product.Tester.Valid {
		view.tester_field.ShowDropdown(true)
		view.tester_field.SetFocus()
	}
	// return &(product_panel.QC_Product.Base())
	return view.QC_Product.Base()
}

func (view *TopPanelView) product_field_pop_data(str string) {
	log.Println("Warn: Debug: product_field_pop_data product_id", view.QC_Product.Product_id)
	log.Println("Warn: Debug: product_field_pop_data product_lot_id", view.QC_Product.Product_Lot_id)

	// if product_lot.product_id != product_lot.insel_product_id(str) {
	old_product_id := view.QC_Product.Product_id
	view.QC_Product.Product_name = str
	view.QC_Product.Insel_product_self()

	if view.QC_Product.Product_id != old_product_id {
		view.SetTitle(view.QC_Product.Product_name)
		view.QC_Product.ResetQC()

		view.QC_Product.Select_product_details()
		view.QC_Product.Update()

		if view.QC_Product.Product_type.Valid {
			view.SetCurrentTab(view.QC_Product.Product_type.Index())
		}

		GUI.Fill_combobox_from_query(view.lot_field, DB.DB_Select_product_lot_product, view.QC_Product.Product_id)
		GUI.Fill_combobox_from_query(view.customer_field, DB.DB_Select_product_customer_info, view.QC_Product.Product_id)
		GUI.Fill_combobox_from_query(view.sample_field, DB.DB_Select_product_sample_points, view.QC_Product.Product_id)

		view.QC_Product.Update_lot(view.lot_field.Text(), view.customer_field.Text())

		view.QC_Product.Sample_point = view.sample_field.Text()

	}
}

func (view *TopPanelView) product_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	view.internal_product_field.SetText(formatted_text)
	if view.internal_product_field.Text() != "" {
		view.product_field_pop_data(view.internal_product_field.Text())
		log.Println("Debug: product_field_text_pop_data", view.QC_Product)
	} else {
		view.QC_Product.Product_id = DB.INVALID_ID
	}
}

func (view *TopPanelView) lot_field_pop_data(str string) {
	view.QC_Product.Update_lot(str, view.customer_field.Text())
	view.SetTitle(str)
}
func (view *TopPanelView) lot_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	view.lot_field.SetText(formatted_text)

	view.lot_field_pop_data(formatted_text)
}

func (view *TopPanelView) customer_field_pop_data(str string) {
	view.QC_Product.Update_lot(view.lot_field.Text(), str)
}
func (view *TopPanelView) customer_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	view.customer_field.SetText(formatted_text)

	view.customer_field_pop_data(formatted_text)
}

func (view *TopPanelView) sample_field_pop_data(str string) {
	view.QC_Product.Sample_point = str
}
func (view *TopPanelView) sample_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	view.sample_field.SetText(formatted_text)
	view.sample_field_pop_data(formatted_text)
}

func (view *TopPanelView) tester_field_pop_data(str string) {
	view.QC_Product.SetTester(str)
}
func (view *TopPanelView) tester_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	view.tester_field.SetText(formatted_text)

	view.tester_field_pop_data(formatted_text)
}
func (view *TopPanelView) PopQRData(product QR.QRJson) {
	view.GoInternal()
	view.product_field_text_pop_data(product.Product_type)
	view.lot_field_text_pop_data(product.Lot_number)
}
func (view *TopPanelView) SetUpdate(update func()) {
	view.QC_Product.Update = update
}

func Show_window() {

	log.Println("Info: Process started")
	// DEBUG
	// log.Println(time.Now().UTC().UnixNano())

	//
	//
	// definitions
	//
	//
	window_title := "QC Data Entry"

	windigo.DefaultFont = windigo.NewFont("MS Shell Dlg 2", GUI.BASE_FONT_SIZE, windigo.FontNormal)

	//
	//
	//
	// build window
	//
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetText(window_title)

	keygrab := windigo.NewEdit(mainWindow)
	keygrab.Hide()

	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := NewTopPanelView(mainWindow)
	product_panel.SetMainWindow(mainWindow)

	tabs := windigo.NewTabView(mainWindow)
	tab_wb := tabs.AddAutoPanel("Water Based")
	tab_oil := tabs.AddAutoPanel("Oil Based")
	tab_fr := tabs.AddAutoPanel("Friction Reducer")

	//
	//
	// Dock
	//
	//

	dock.Dock(product_panel, windigo.Top)
	dock.Dock(tabs, windigo.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), windigo.Fill) // tab panels dock just below tabs and fill area

	threads.Status_bar = windigo.NewStatusBar(mainWindow)
	mainWindow.SetStatusBar(threads.Status_bar)

	//
	//
	//
	//
	// functionality

	panel_water_based := show_water_based(tab_wb, product_panel.QC_Product, product_panel.BaseProduct)
	panel_oil_based := show_oil_based(tab_oil, product_panel.QC_Product, product_panel.BaseProduct)
	panel_fr := show_fr(tab_fr, product_panel.QC_Product, product_panel.BaseProduct)

	product_panel.SetUpdate(func() {
		product_panel.container_field.Update(product_panel.QC_Product.Container_type)
		log.Println("Debug: update new_product_cb", product_panel.QC_Product)
		panel_water_based.Update(product_panel.QC_Product)
		panel_oil_based.Update(product_panel.QC_Product)
		panel_fr.Update(product_panel.QC_Product)
		panel_fr.ChangeContainer(product_panel.QC_Product)
	})
	product_panel.ChangeContainer = panel_fr.ChangeContainer
	product_panel.SetCurrentTab = func(i int) {
		tabs.SetCurrent(i)
	}

	//
	//
	//
	//
	// sizing
	refresh := func(font_size int) {
		refresh_globals(font_size)

		mainWindow.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)
		// product_panel.RefreshSize()
		product_panel.RefreshSize(font_size)

		tabs.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		panel_water_based.Refresh()
		panel_oil_based.Refresh()
		panel_fr.Refresh()

	}

	set_font := func(font_size int) {
		GUI.BASE_FONT_SIZE = font_size

		config.Main_config.Set("font_size", GUI.BASE_FONT_SIZE)
		config.Write_config(config.Main_config)

		old_font := windigo.DefaultFont
		windigo.DefaultFont = windigo.NewFont(old_font.Family(), GUI.BASE_FONT_SIZE, 0)
		old_font.Dispose()

		mainWindow.SetFont(windigo.DefaultFont)
		product_panel.SetFont(windigo.DefaultFont)
		tabs.SetFont(windigo.DefaultFont)

		panel_water_based.SetFont(windigo.DefaultFont)
		panel_oil_based.SetFont(windigo.DefaultFont)
		panel_fr.SetFont(windigo.DefaultFont)
		threads.Status_bar.SetFont(windigo.DefaultFont)
		refresh(font_size)

	}
	set_font(GUI.BASE_FONT_SIZE)

	AddShortcuts(mainWindow, keygrab, set_font, product_panel.PopQRData)

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)
	mainWindow.RunMainLoop() // Must call to start event loop.
}

func AddShortcuts(mainWindow *windigo.Form, keygrab *windigo.Edit, set_font func(int), PopQRData func(QR.QRJson)) {
	// QR keyboard handling

	keygrab.OnSetFocus().Bind(func(e *windigo.Event) {
		keygrab.SetText("{")
		keygrab.SelectText(1, 1)
	})

	mainWindow.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModShift, Key: windigo.KeyOEM4}, // {
		func() bool {
			keygrab.SetText("")
			keygrab.SetFocus()
			return true
		})

	mainWindow.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModShift, Key: windigo.KeyOEM6}, // }
		func() bool {
			var product QR.QRJson

			qr_json := keygrab.Text() + "}"
			log.Println("debug: ReadFromScanner: ", qr_json)
			err := json.Unmarshal([]byte(qr_json), &product)
			if err == nil {
				PopQRData(product)
			} else {
				log.Printf("error: [%s]: %q\n", "qr_json_keygrab", err)
			}
			keygrab.SetText("")
			mainWindow.SetFocus()
			return true
		})

	// Resize handling
	mainWindow.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModControl, Key: windigo.KeyOEMPlus}, // {
		func() bool {
			set_font(GUI.BASE_FONT_SIZE + 1)
			return true
		})

	mainWindow.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModControl, Key: windigo.KeyOEMMinus}, // {
		func() bool {
			set_font(GUI.BASE_FONT_SIZE - 1)
			return true
		})
}

func wndOnClose(arg *windigo.Event) {
	GUI.OKPen.Dispose()
	GUI.ErroredPen.Dispose()
	windigo.Exit()
}
