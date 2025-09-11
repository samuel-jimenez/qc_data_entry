package qc

import (
	"log"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/QR"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
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
	QC_Product *product.QCProduct

	mainWindow           *QCWindow
	TopPanelInternalView *TopPanelInternalView
	TopPanelInboundView  *TopPanelInboundView

	product_panel_0,
	product_panel_0_2 *windigo.AutoPanel

	tester_field *GUI.SearchBox

	container_field *product.DiscreteView
	clock_panel     *views.ClockTimerView
}

// func NewTopPanelView(parent windigo.Controller,	QCWindow *QCWindow) *TopPanelView {
func NewTopPanelView(parent *QCWindow) *TopPanelView {
	view := new(TopPanelView)
	// view.QCWindow = QCWindow
	view.mainWindow = parent
	view.QC_Product = product.NewQCProduct()
	view.QC_Product.SetUpdate(view.mainWindow.UpdateProduct)

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

	container_field := product.BuildNewDiscreteView_NOUPDATE(product_panel, "Container Type", []string{"Sample", "Tote", "Railcar", "ISO"}) // bs.container_types

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
	view.TopPanelInternalView = NewTopPanelInternalView(
		view.mainWindow,
		view.QC_Product,
		product_panel_0_0, product_panel_0_1,
		internal_product_field, customer_field, lot_field, sample_field,
		ranges_button, reprint_button, inbound_button)
	view.TopPanelInboundView = NewTopPanelInboundView(
		view.mainWindow.Form,
		view.QC_Product,
		product_panel_1_0, product_panel_1_1,
		testing_lot_field, inbound_lot_field, inbound_container_field, inbound_product_field,
		sample_button, release_button, internal_button)

	view.product_panel_0_2 = product_panel_0_2
	view.clock_panel = clock_panel

	view.tester_field = tester_field
	view.container_field = container_field

	//
	//
	// Dock
	//
	//

	product_panel_0_2.Dock(tester_field, windigo.Left)

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

	tester_field.OnSelectedChange().Bind(func(e *windigo.Event) { view.tester_field_pop_data(tester_field.GetSelectedItem()) })
	tester_field.OnKillFocus().Bind(func(e *windigo.Event) { view.tester_field_text_pop_data(tester_field.Text()) })

	inbound_button.OnClick().Bind(func(e *windigo.Event) { view.GoInbound() })

	internal_button.OnClick().Bind(func(e *windigo.Event) { view.GoInternal() })

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

}

// container_item_width

// func (view *TopPanelView) RefreshSize() {
func (view *TopPanelView) RefreshSize(font_size int) {
	var (
		top_panel_height,

		top_subpanel_height,

		container_item_width int
	)

	num_rows := 3

	top_subpanel_height = GUI.TOP_SPACER_HEIGHT + num_rows*(GUI.PRODUCT_FIELD_HEIGHT+GUI.INTER_SPACER_HEIGHT) + GUI.BTM_SPACER_HEIGHT

	top_panel_height = top_subpanel_height + num_rows*GUI.GROUPBOX_CUSHION + GUI.PRODUCT_FIELD_HEIGHT

	container_item_width = 6 * font_size

	view.SetSize(GUI.TOP_PANEL_WIDTH, top_panel_height)

	view.product_panel_0.SetSize(GUI.TOP_PANEL_WIDTH, top_subpanel_height)
	view.TopPanelInternalView.RefreshSize()
	view.TopPanelInboundView.RefreshSize()

	view.product_panel_0_2.SetSize(GUI.HPANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.product_panel_0_2.SetMarginTop(GUI.INTER_SPACER_HEIGHT)
	view.product_panel_0_2.SetMarginLeft(GUI.HPANEL_MARGIN)

	view.tester_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.clock_panel.RefreshSize()

	view.container_field.SetSize(GUI.DISCRETE_FIELD_WIDTH, GUI.OFF_AXIS)
	view.container_field.SetItemSize(container_item_width)
	view.container_field.SetPaddingsAll(GUI.GROUPBOX_CUSHION)
	view.container_field.SetPaddingLeft(0)

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

		//isInternal && check()
		view.TopPanelInternalView.AlertProduct()

	}

	// return &(product_panel.QC_Product.Base())
	return view.QC_Product.Base()
}

func (view *TopPanelView) GoInbound() {
	view.TopPanelInboundView.Show()
	view.TopPanelInternalView.Hide()

	view.QC_Product.Container_type = product.CONTAINER_SAMPLE
	view.ChangeContainer(view.QC_Product)
}

func (view *TopPanelView) GoInternal() {

	view.TopPanelInternalView.Show()
	view.TopPanelInboundView.Hide()

	view.QC_Product.Blend = nil

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
	view.container_field.Update(product.DiscreteFromInt(int(QC_Product.Container_type)))
}

func (view *TopPanelView) ChangeContainer(qc_product *product.QCProduct) {
	view.mainWindow.ChangeContainer(qc_product)
}
