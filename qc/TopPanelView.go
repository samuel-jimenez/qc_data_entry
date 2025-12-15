package qc

import (
	"log"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/qc_ui"
	"github.com/samuel-jimenez/qc_data_entry/QR"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

var (
	MODE_INTERNAL = 0
	MODE_INBOUND  = 1
)

/*
 * TopPanelViewer
 *
 */
type TopPanelViewer interface {
	windigo.Controller
	SetFont(font *windigo.Font)
	RefreshSize()
	SetTitle(title string)
	BaseProduct() product.BaseProduct

	GoInbound()
	GoInternal()
	PopQRData(product QR.QRJson)

	tester_field_pop_data(str string)
	tester_field_text_pop_data(str string)

	UpdateProduct(QC_Product *product.QCProduct)
	ChangeContainer(qc_product *product.QCProduct)
}

/*
 * TopPanelView
 *
 */
type TopPanelView struct {
	*windigo.AutoPanel
	QC_Product       *product.QCProduct
	Measured_Product []*product.MeasuredProduct

	mainWindow           *QCWindow
	TopPanelInternalView *TopPanelInternalView
	TopPanelInboundView  *TopPanelInboundView

	whups_button *windigo.PushButton

	product_panel_0,
	product_panel_0_2 *windigo.AutoPanel

	tester_field *GUI.SearchBox

	container_field *product.DiscreteView
	clock_panel     *views.ClockTimerView

	mode int
}

func (view *TopPanelView) SetMeasuredProduct(products ...*product.MeasuredProduct) {
	view.Measured_Product = products
}

// func TopPanelView_from_new(parent windigo.Controller,	QCWindow *QCWindow) *TopPanelView {
func TopPanelView_from_new(parent *QCWindow) *TopPanelView {
	view := new(TopPanelView)
	// view.QCWindow = QCWindow
	view.mainWindow = parent
	view.QC_Product = product.QCProduct_from_new()
	view.QC_Product.SetUpdate(view.mainWindow.UpdateProduct)
	view.mode = MODE_INTERNAL

	product_text := "Product"
	sample_text := "Sample Point" // TODO
	customer_text := "Customer Name"
	tester_text := "Tester"

	ranges_text := "Ranges"
	inventory_text := "Inventory"
	reprint_text := "Reprint"
	inbound_text := "Inbound"
	sample_button_text := "Sample"
	whups_text := "Whups"

	release_button_text := "Release"
	today_button_text := "Today"

	inbound_lot_text := "Inbound Lot"
	internal_text := "Internal"
	container_text := "Container"

	product_panel := windigo.NewAutoPanel(parent)

	product_panel_0 := windigo.NewAutoPanel(product_panel)

	// TODO array db_select_all_product

	product_panel_0_0 := windigo.NewAutoPanel(product_panel_0)

	internal_product_field := GUI.ComboBox_from_new(product_panel_0_0, product_text)
	customer_field := GUI.ComboBox_from_new(product_panel_0_0, customer_text)

	product_panel_1_0 := windigo.NewAutoPanel(product_panel_0)
	product_panel_1_0.Hide()

	testing_lot_field := GUI.List_ComboBox_from_new(product_panel_1_0, formats.COL_LABEL_LOT)
	// inbound_product_field := GUI.NewComboBox(product_panel_1_0, product_text)
	// inbound_product_field := windigo.LabeledComboBox_from_new(product_panel_1_0, product_text)
	inbound_product_field := windigo.LabeledEdit_from_new(product_panel_1_0, product_text)
	inbound_product_field.SetReadOnly(true)

	product_panel_0_1 := windigo.NewAutoPanel(product_panel_0)

	lot_field := GUI.ComboBox_from_new(product_panel_0_1, formats.COL_LABEL_LOT)
	sample_field := GUI.ComboBox_from_new(product_panel_0_1, sample_text)

	product_panel_1_1 := windigo.NewAutoPanel(product_panel_0)
	product_panel_1_1.Hide()

	// inbound_lot_field := GUI.NewComboBox(product_panel_1_1, product_text)
	inbound_lot_field := GUI.ComboBox_from_new(product_panel_1_1, inbound_lot_text)
	inbound_container_field := GUI.ComboBox_from_new(product_panel_1_1, container_text)
	// inbound_lot_field container_text

	product_panel_0_2 := windigo.NewAutoPanel(product_panel_0)

	tester_field := GUI.NewLabeledSearchBoxFromQuery(product_panel_0_2, tester_text, DB.DB_Select_all_qc_tester)

	clock_panel := views.NewClockTimerView(product_panel_0)

	ranges_button := windigo.NewPushButton(product_panel)
	ranges_button.SetText(ranges_text)

	sample_button := windigo.NewPushButton(product_panel)
	sample_button.SetText(sample_button_text)
	sample_button.Hide()

	container_field := product.BuildNewDiscreteView_NOUPDATE(product_panel, "Container Type", "Sample", "Tote", "Railcar", "ISO") // bs.container_types

	release_button := windigo.NewPushButton(product_panel)
	release_button.SetText(release_button_text)
	release_button.Hide()

	inventory_button := windigo.NewPushButton(product_panel)
	inventory_button.SetText(inventory_text)

	reprint_button := windigo.NewPushButton(product_panel)
	reprint_button.SetText(reprint_text)

	today_button := windigo.NewPushButton(product_panel)
	today_button.SetText(today_button_text)
	today_button.Hide()

	inbound_button := windigo.NewPushButton(product_panel)
	inbound_button.SetText(inbound_text)

	internal_button := windigo.NewPushButton(product_panel)
	internal_button.SetText(internal_text)
	internal_button.Hide()

	view.whups_button = windigo.NewPushButton(product_panel)
	view.whups_button.SetText(whups_text)

	// build object
	view.AutoPanel = product_panel
	view.product_panel_0 = product_panel_0
	view.TopPanelInternalView = NewTopPanelInternalView(
		view.mainWindow,
		view.QC_Product,
		product_panel_0_0, product_panel_0_1,
		internal_product_field, customer_field, lot_field, sample_field,
		container_field,
		ranges_button, inventory_button, reprint_button, inbound_button)
	view.TopPanelInboundView = NewTopPanelInboundView(
		view.mainWindow.Form,
		view.QC_Product,
		product_panel_1_0, product_panel_1_1,
		testing_lot_field, inbound_lot_field, inbound_container_field, inbound_product_field,
		sample_button, release_button, today_button, internal_button)

	view.product_panel_0_2 = product_panel_0_2
	view.clock_panel = clock_panel

	view.container_field = container_field

	view.tester_field = tester_field

	//
	//
	// Dock
	product_panel_0_2.Dock(tester_field, windigo.Left)

	product_panel_0.Dock(clock_panel, windigo.Right)
	product_panel_0.Dock(product_panel_0_0, windigo.Top)
	product_panel_0.Dock(product_panel_1_0, windigo.Top)
	product_panel_0.Dock(product_panel_0_1, windigo.Top)
	product_panel_0.Dock(product_panel_1_1, windigo.Top)
	product_panel_0.Dock(product_panel_0_2, windigo.Top)

	product_panel.Dock(product_panel_0, windigo.Top)
	product_panel.Dock(sample_button, windigo.Left)
	product_panel.Dock(ranges_button, windigo.Left)
	product_panel.Dock(release_button, windigo.Left)
	product_panel.Dock(container_field, windigo.Left)
	product_panel.Dock(today_button, windigo.Left)
	product_panel.Dock(inventory_button, windigo.Left)
	product_panel.Dock(reprint_button, windigo.Left)
	product_panel.Dock(inbound_button, windigo.Left)
	product_panel.Dock(internal_button, windigo.Left)
	product_panel.Dock(view.whups_button, windigo.Left)

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

	view.AutoPanel.AddShortcut(windigo.Shortcut{Key: windigo.KeyEnter}, parent.FocusTab)

	tester_field.OnSelectedChange().Bind(func(e *windigo.Event) { view.tester_field_pop_data(tester_field.GetSelectedItem()) })
	tester_field.OnKillFocus().Bind(func(e *windigo.Event) { view.tester_field_text_pop_data(tester_field.Text()) })

	inbound_button.OnClick().Bind(func(e *windigo.Event) { view.GoInbound() })

	internal_button.OnClick().Bind(func(e *windigo.Event) { view.GoInternal() })

	view.whups_button.OnClick().Bind(view.whups_button_OnClick)

	// TODO 	clock_panel.OnClick().Bind(view.clock_panel_OnClick)
	// clock_panel.OnLBDbl().Bind(view.clock_panel_OnClick)// worthless
	// clock_panel.OnLBUp().Bind(view.clock_panel_OnClick)
	// clock_panel.OnLBDown().Bind(view.clock_panel_OnClick)

	// clock_panel.SetOnLBUp(view.clock_panel_OnClick)
	clock_panel.SetOnLBDown(view.clock_panel_OnClick)

	container_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		view.QC_Product.Container_type = product.ProductContainerType(container_field.Get().Int32)
		view.ChangeContainer(view.QC_Product)
	})

	return view
}

func (view *TopPanelView) SetFont(font *windigo.Font) {
	view.TopPanelInternalView.SetFont(font)
	view.TopPanelInboundView.SetFont(font)
	view.tester_field.SetFont(font)
	view.container_field.SetFont(font)

	view.whups_button.SetFont(font)
}

func (view *TopPanelView) RefreshSize() {
	view.SetSize(GUI.TOP_PANEL_WIDTH, TOP_PANEL_HEIGHT)

	view.product_panel_0.SetSize(GUI.TOP_PANEL_WIDTH, TOP_SUBPANEL_HEIGHT)
	view.TopPanelInternalView.RefreshSize()
	view.TopPanelInboundView.RefreshSize()

	view.product_panel_0_2.SetSize(GUI.HPANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.product_panel_0_2.SetMarginTop(GUI.INTER_SPACER_HEIGHT)
	view.product_panel_0_2.SetMarginLeft(GUI.HPANEL_MARGIN)

	view.tester_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.clock_panel.RefreshSize()

	view.container_field.SetSize(GUI.DISCRETE_FIELD_WIDTH, GUI.OFF_AXIS)
	view.container_field.SetItemSize(CONTAINER_ITEM_WIDTH)
	view.container_field.SetPaddingsAll(GUI.GROUPBOX_CUSHION)
	view.container_field.SetPaddingLeft(0)

	view.whups_button.SetMarginsAll(BUTTON_MARGIN)
	view.whups_button.SetSize(GUI.REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)
}

func (view *TopPanelView) SetTitle(title string) {
	if view.mainWindow == nil {
		return
	}
	view.mainWindow.SetText(title)
}

func (view *TopPanelView) BaseProduct() product.QCProduct {
	log.Println("Debug: TopPanelView-BaseProduct", view.QC_Product.Base())

	// there has to be better way

	view.QC_Product.Valid = view.QC_Product.Tester.Valid

	// only last test holds dropdown
	if !view.QC_Product.Tester.Valid {
		view.tester_field.Alert()
	}

	if true &&
		view.QC_Product.Product_id == DB.INVALID_ID &&
		// and not internal
		view.QC_Product.Blend == nil {
		view.QC_Product.Valid = false
		// TODO TopPanelInternalViewAlertProduct replace with proper check
		// isInternal && check()
		view.TopPanelInternalView.AlertProduct()

	}

	if view.QC_Product.Valid {
		view.mainWindow.SetBlend(view.QC_Product)
	}
	return view.QC_Product.Base()
}

func (view *TopPanelView) GoInbound() {
	if view.mode == MODE_INBOUND {
		return
	}
	view.mode = MODE_INBOUND

	view.TopPanelInboundView.Show()
	view.TopPanelInternalView.Hide()

	view.QC_Product.Container_type = product.CONTAINER_SAMPLE
	view.ChangeContainer(view.QC_Product)
	view.mainWindow.ComponentsDisable()
}

func (view *TopPanelView) GoInternal() {
	if view.mode == MODE_INTERNAL {
		return
	}
	view.mode = MODE_INTERNAL

	view.TopPanelInternalView.Show()
	view.TopPanelInboundView.Hide()

	view.QC_Product.Blend = nil
	view.mainWindow.UpdateProduct(view.QC_Product)
	view.mainWindow.ComponentsEnable()
}

func (view *TopPanelView) PopQRData(product QR.QRJson) {
	view.GoInternal()
	view.TopPanelInternalView.PopQRData(product)
}

func (view *TopPanelView) tester_field_pop_data(str string) {
	view.QC_Product.SetTester(str)
	if view.QC_Product.Tester.Valid {
		view.tester_field.Ok()
	} else {
		view.tester_field.Error()
	}
}

func (view *TopPanelView) tester_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	view.tester_field.SetText(formatted_text)

	view.tester_field_pop_data(formatted_text)
}

func (view *TopPanelView) UpdateProduct(QC_Product *product.QCProduct) {
	view.container_field.Update(product.DiscreteFromContainer(QC_Product.Container_type))
}

func (view *TopPanelView) ChangeContainer(qc_product *product.QCProduct) {
	view.mainWindow.ChangeContainer(qc_product)
}

func (view *TopPanelView) clock_panel_OnClick(*windigo.Event) {
	log.Println("clock_panel_OnClick") // ?? TODO
	ClockPopoutView := views.ClockPopoutView_from_new(view)
	ClockPopoutView.SetModal(false)
	ClockPopoutView.Show()
	ClockPopoutView.RefreshSize()
}

func (view *TopPanelView) whups_button_OnClick(*windigo.Event) {
	WhupsView := qc_ui.WhupsView_from_new(view, view.Measured_Product)
	WhupsView.SetModal(false)
	WhupsView.RefreshSize()
	WhupsView.Show()
}
