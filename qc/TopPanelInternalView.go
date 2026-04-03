package qc

import (
	"database/sql"
	"log"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/QR"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

/*
 * TopPanelInternalViewer
 *
 */
type TopPanelInternalViewer interface {
	windigo.Controller
	SetFont(font *windigo.Font)
	RefreshSize()
	SetTitle(title string)
	// BaseProduct() product.BaseProduct

	Show()
	Hide()
	// TODO TopPanelInternalViewAlertProduct replace with proper check
	AlertProduct()
	PopQRData(product QR.QRJson)

	product_field_pop_data(str string)
	product_field_text_pop_data(str string)
	lot_field_pop_data(str string)
	lot_field_text_pop_data(str string)
	customer_field_pop_data(str string)
	customer_field_text_pop_data(str string)
	sample_field_pop_data(str string)
	sample_field_text_pop_data(str string)

	SetCurrentTab(int)

	ranges_button_OnClick(*windigo.Event)
	inventory_button_OnClick(*windigo.Event)
	reprint_button_OnClick(*windigo.Event)
}

/*
 * TopPanelInternalView
 *
 */
type TopPanelInternalView struct {
	*windigo.AutoPanel
	QC_Product   *product.QCProduct
	product_data map[string]int

	mainWindow *QCWindow

	product_panel_0_0, product_panel_0_1 *windigo.AutoPanel

	internal_product_field, customer_field,
	lot_field, sample_field *GUI.ComboBox

	container_field *product.DiscreteView

	ranges_button, today_button, inventory_button, reprint_button, inbound_button *windigo.PushButton
}

func NewTopPanelInternalView(
	QCWindow *QCWindow,
	QC_Product *product.QCProduct,
	product_panel_0_0, product_panel_0_1 *windigo.AutoPanel,
	internal_product_field, customer_field, lot_field, sample_field *GUI.ComboBox,
	container_field *product.DiscreteView,
	ranges_button, today_button, inventory_button, reprint_button, inbound_button *windigo.PushButton,
) *TopPanelInternalView {
	view := new(TopPanelInternalView)

	// build object
	view.mainWindow = QCWindow

	view.product_data = make(map[string]int)
	view.QC_Product = QC_Product

	view.product_panel_0_0 = product_panel_0_0
	view.product_panel_0_1 = product_panel_0_1

	view.internal_product_field = internal_product_field
	view.customer_field = customer_field

	view.lot_field = lot_field
	view.sample_field = sample_field

	view.container_field = container_field

	view.ranges_button = ranges_button
	view.today_button = today_button
	view.inventory_button = inventory_button
	view.reprint_button = reprint_button
	view.inbound_button = inbound_button

	//
	// Dock
	view.product_panel_0_0.Dock(view.internal_product_field, windigo.Left)
	view.product_panel_0_0.Dock(view.customer_field, windigo.Left)

	view.product_panel_0_1.Dock(view.lot_field, windigo.Left)
	view.product_panel_0_1.Dock(view.sample_field, windigo.Left)

	//
	// combobox
	GUI.Fill_combobox_from_query_rows(view.internal_product_field, func(row *sql.Rows) error {
		var (
			id                   int
			internal_name        string
			product_moniker_name string
		)
		if err := row.Scan(&id, &internal_name, &product_moniker_name); err != nil {
			return err
		}
		name := product_moniker_name + " " + internal_name
		view.product_data[name] = id

		view.internal_product_field.AddItem(name)
		return nil
	}, DB.DB_Select_product_info_all)

	//
	// functionality
	view.internal_product_field.OnSelectedChange().Bind(func(e *windigo.Event) { view.product_field_pop_data(internal_product_field.GetSelectedItem()) })
	view.internal_product_field.OnKillFocus().Bind(func(e *windigo.Event) { view.product_field_text_pop_data(internal_product_field.Text()) })

	view.lot_field.OnSelectedChange().Bind(func(e *windigo.Event) { view.lot_field_pop_data(lot_field.GetSelectedItem()) })
	view.lot_field.OnKillFocus().Bind(func(e *windigo.Event) { view.lot_field_text_pop_data(lot_field.Text()) })

	view.customer_field.OnSelectedChange().Bind(func(e *windigo.Event) { view.customer_field_pop_data(customer_field.GetSelectedItem()) })
	view.customer_field.OnKillFocus().Bind(func(e *windigo.Event) { view.customer_field_text_pop_data(customer_field.Text()) })

	view.sample_field.OnSelectedChange().Bind(func(e *windigo.Event) { view.sample_field_pop_data(sample_field.GetSelectedItem()) })
	view.sample_field.OnKillFocus().Bind(func(e *windigo.Event) { view.sample_field_text_pop_data(sample_field.Text()) })

	view.ranges_button.OnClick().Bind(view.ranges_button_OnClick)

	view.today_button.OnClick().Bind(view.today_button_OnClick)

	view.inventory_button.OnClick().Bind(view.inventory_button_OnClick)

	view.reprint_button.OnClick().Bind(view.reprint_button_OnClick)

	return view
}

func (self *TopPanelInternalView) SetFont(font *windigo.Font) {
	self.internal_product_field.SetFont(font)
	self.customer_field.SetFont(font)
	self.lot_field.SetFont(font)
	self.sample_field.SetFont(font)

	self.ranges_button.SetFont(font)
	self.today_button.SetFont(font)
	self.inventory_button.SetFont(font)
	self.reprint_button.SetFont(font)
	self.inbound_button.SetFont(font)
}

func (self *TopPanelInternalView) RefreshSize() {
	self.product_panel_0_0.SetSize(GUI.HPANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	self.product_panel_0_0.SetMarginTop(GUI.TOP_SPACER_HEIGHT)
	self.product_panel_0_0.SetMarginLeft(GUI.HPANEL_MARGIN)

	self.internal_product_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	self.customer_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	self.product_panel_0_1.SetSize(GUI.HPANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	self.product_panel_0_1.SetMarginTop(GUI.INTER_SPACER_HEIGHT)
	self.product_panel_0_1.SetMarginLeft(GUI.HPANEL_MARGIN)

	self.customer_field.SetMarginLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)
	self.sample_field.SetMarginLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)

	self.lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	self.sample_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	self.ranges_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
	self.ranges_button.SetMarginsAll(BUTTON_MARGIN)

	self.inventory_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
	self.inventory_button.SetMarginsAll(BUTTON_MARGIN)

	self.reprint_button.SetMarginsAll(BUTTON_MARGIN)
	self.reprint_button.SetMarginLeft(GUI.REPRINT_BUTTON_MARGIN_L)
	self.reprint_button.SetSize(GUI.REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

	self.inbound_button.SetMarginsAll(BUTTON_MARGIN)
	self.inbound_button.SetSize(GUI.REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)
}

func (self *TopPanelInternalView) SetTitle(title string) {
	if self.mainWindow == nil {
		return
	}
	self.mainWindow.SetText(title)
}

func (self *TopPanelInternalView) Show() {
	self.product_panel_0_0.Show()
	self.product_panel_0_1.Show()
	self.ranges_button.Show()
	self.today_button.Show()
	self.container_field.Show()
	self.reprint_button.Show()
	self.inventory_button.Show()
	self.inbound_button.Show()

	lot := self.lot_field.Text()
	self.internal_product_field.OnSelectedChange().Fire(nil)
	self.lot_field_text_pop_data(lot)
}

func (self *TopPanelInternalView) Hide() {
	self.product_panel_0_0.Hide()
	self.product_panel_0_1.Hide()
	self.ranges_button.Hide()
	self.today_button.Hide()
	self.container_field.Hide()
	self.reprint_button.Hide()
	self.inventory_button.Hide()
	self.inbound_button.Hide()
}

func (self *TopPanelInternalView) PopQRData(product QR.QRJson) {
	self.product_field_text_pop_data(product.Product_type)
	self.lot_field_text_pop_data(product.Lot_number)
}

// TODO TopPanelInternalViewAlertProduct replace with proper check
func (self *TopPanelInternalView) AlertProduct() {
	self.internal_product_field.Alert()
}

func (self *TopPanelInternalView) product_field_pop_data(str string) {
	self.internal_product_field.Ok()

	// if product_lot.product_id != product_lot.insel_product_id(str) {
	old_product_id := self.QC_Product.Product_id
	self.QC_Product.Product_name = str
	self.QC_Product.Insel_product_self()

	if self.QC_Product.Product_id != old_product_id {
		self.SetTitle(self.QC_Product.Product_name)
		self.QC_Product.ResetQC()
		self.QC_Product.Blend = nil

		self.QC_Product.Select_product_details()
		self.QC_Product.Update()

		if self.QC_Product.Product_type.Valid {
			self.SetCurrentTab(self.QC_Product.Product_type.Index())
		}

		GUI.Fill_combobox_from_query(self.lot_field, DB.DB_Select_product_lot_product, self.QC_Product.Product_id)
		GUI.Fill_combobox_from_query(self.customer_field, DB.DB_Select_product_customer_info, self.QC_Product.Product_id)
		GUI.Fill_combobox_from_query(self.sample_field, DB.DB_Select_product_sample_points, self.QC_Product.Product_id)

		self.QC_Product.Update_lot(self.lot_field.Text(), self.customer_field.Text())

		self.QC_Product.Sample_point = self.sample_field.Text()

	}
}

func (self *TopPanelInternalView) product_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	self.internal_product_field.SetText(formatted_text)
	if self.internal_product_field.Text() != "" {
		self.product_field_pop_data(self.internal_product_field.Text())
		log.Println("Debug: product_field_text_pop_data", self.QC_Product)
	} else {
		self.QC_Product.Product_id = DB.INVALID_ID
		self.internal_product_field.Error()

	}
}

func (self *TopPanelInternalView) lot_field_pop_data(str string) {
	self.QC_Product.Update_lot(str, self.customer_field.Text())
	self.SetTitle(str)
}

func (self *TopPanelInternalView) lot_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	self.lot_field.SetText(formatted_text)

	self.lot_field_pop_data(formatted_text)
}

func (self *TopPanelInternalView) customer_field_pop_data(str string) {
	self.QC_Product.Update_lot(self.lot_field.Text(), str)
}

func (self *TopPanelInternalView) customer_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	self.customer_field.SetText(formatted_text)

	self.customer_field_pop_data(formatted_text)
}

func (self *TopPanelInternalView) sample_field_pop_data(str string) {
	self.QC_Product.Sample_point = str
}

func (self *TopPanelInternalView) sample_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	self.sample_field.SetText(formatted_text)
	self.sample_field_pop_data(formatted_text)
}

func (self *TopPanelInternalView) SetCurrentTab(i int) {
	self.mainWindow.SetCurrentTab(i)
}

func (self *TopPanelInternalView) ranges_button_OnClick(*windigo.Event) {
	if self.QC_Product.Product_name != "" {
		views.ShowNewQCProductRangesView(self.QC_Product)
		log.Println("debug: ranges_button-product_lot", self.QC_Product)
	}
}

func (self *TopPanelInternalView) today_button_OnClick(*windigo.Event) {
	if self.QC_Product == nil {
		return
	}

	lot_date := blender.BlendProductLOTS()
	self.lot_field.SetText("BSW" + lot_date)
	self.QC_Product.Product_Lot_id = DB.DEFAULT_LOT_ID
	self.QC_Product.Lot_id = DB.DEFAULT_LOT_ID
}

func (self *TopPanelInternalView) inventory_button_OnClick(*windigo.Event) {
	// views.ShowNewQCProductRangesView()
	InventoryView_from_new().Start()
}

func (self *TopPanelInternalView) reprint_button_OnClick(*windigo.Event) {
	if self.QC_Product.Lot_number != "" {
		log.Println("debug: reprint_button")
		self.QC_Product.Reprint()
	}
}
