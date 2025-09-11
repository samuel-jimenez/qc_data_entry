package qc

import (
	"database/sql"
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

	ranges_button, reprint_button, inbound_button *windigo.PushButton
}

func NewTopPanelInternalView(
	QCWindow *QCWindow,
	QC_Product *product.QCProduct,
	product_panel_0_0, product_panel_0_1 *windigo.AutoPanel,
	internal_product_field, customer_field, lot_field, sample_field *GUI.ComboBox,
	ranges_button, reprint_button, inbound_button *windigo.PushButton,
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

	view.ranges_button = ranges_button
	view.reprint_button = reprint_button
	view.inbound_button = inbound_button

	//
	//
	// Dock
	//
	//
	product_panel_0_0.Dock(internal_product_field, windigo.Left)
	product_panel_0_0.Dock(customer_field, windigo.Left)

	product_panel_0_1.Dock(lot_field, windigo.Left)
	product_panel_0_1.Dock(sample_field, windigo.Left)

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
		view.product_data[name] = id

		internal_product_field.AddItem(name)
		return nil
	}, DB.DB_Select_product_info)

	//
	//
	//
	//
	// functionality

	internal_product_field.OnSelectedChange().Bind(func(e *windigo.Event) { view.product_field_pop_data(internal_product_field.GetSelectedItem()) })
	internal_product_field.OnKillFocus().Bind(func(e *windigo.Event) { view.product_field_text_pop_data(internal_product_field.Text()) })

	lot_field.OnSelectedChange().Bind(func(e *windigo.Event) { view.lot_field_pop_data(lot_field.GetSelectedItem()) })
	lot_field.OnKillFocus().Bind(func(e *windigo.Event) { view.lot_field_text_pop_data(lot_field.Text()) })

	customer_field.OnSelectedChange().Bind(func(e *windigo.Event) { view.customer_field_pop_data(customer_field.GetSelectedItem()) })
	customer_field.OnKillFocus().Bind(func(e *windigo.Event) { view.customer_field_text_pop_data(customer_field.Text()) })

	sample_field.OnSelectedChange().Bind(func(e *windigo.Event) { view.sample_field_pop_data(sample_field.GetSelectedItem()) })
	sample_field.OnKillFocus().Bind(func(e *windigo.Event) { view.sample_field_text_pop_data(sample_field.Text()) })

	ranges_button.OnClick().Bind(func(e *windigo.Event) {
		if view.QC_Product.Product_name != "" {
			views.ShowNewQCProductRangesView(view.QC_Product)
			log.Println("debug: ranges_button-product_lot", view.QC_Product)
		}
	})

	reprint_button.OnClick().Bind(func(e *windigo.Event) {
		if view.QC_Product.Lot_number != "" {
			log.Println("debug: reprint_button")
			view.QC_Product.Reprint()
		}
	})

	return view
}

func (view *TopPanelInternalView) SetFont(font *windigo.Font) {

	view.internal_product_field.SetFont(font)
	view.customer_field.SetFont(font)
	view.lot_field.SetFont(font)
	view.sample_field.SetFont(font)

	view.ranges_button.SetFont(font)
	view.reprint_button.SetFont(font)
	view.inbound_button.SetFont(font)

}

func (view *TopPanelInternalView) RefreshSize() {

	view.product_panel_0_0.SetSize(GUI.HPANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.product_panel_0_0.SetMarginTop(GUI.TOP_SPACER_HEIGHT)
	view.product_panel_0_0.SetMarginLeft(GUI.HPANEL_MARGIN)

	view.internal_product_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.customer_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.product_panel_0_1.SetSize(GUI.HPANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.product_panel_0_1.SetMarginTop(GUI.INTER_SPACER_HEIGHT)
	view.product_panel_0_1.SetMarginLeft(GUI.HPANEL_MARGIN)

	view.customer_field.SetMarginLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)
	view.sample_field.SetMarginLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)

	view.lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.sample_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.ranges_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
	view.ranges_button.SetMarginsAll(BUTTON_MARGIN)

	view.reprint_button.SetMarginsAll(BUTTON_MARGIN)
	view.reprint_button.SetMarginLeft(GUI.REPRINT_BUTTON_MARGIN_L)
	view.reprint_button.SetSize(GUI.REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.inbound_button.SetMarginsAll(BUTTON_MARGIN)
	view.inbound_button.SetSize(GUI.REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

}

func (view *TopPanelInternalView) SetTitle(title string) {
	if view.mainWindow == nil {
		return
	}
	view.mainWindow.SetText(title)
}

func (view *TopPanelInternalView) Show() {

	view.product_panel_0_0.Show()
	view.product_panel_0_1.Show()
	view.ranges_button.Show()
	view.reprint_button.Show()
	view.inbound_button.Show()

}

func (view *TopPanelInternalView) Hide() {

	view.product_panel_0_0.Hide()
	view.product_panel_0_1.Hide()
	view.ranges_button.Hide()
	view.reprint_button.Hide()
	view.inbound_button.Hide()

}

func (view *TopPanelInternalView) PopQRData(product QR.QRJson) {
	view.product_field_text_pop_data(product.Product_type)
	view.lot_field_text_pop_data(product.Lot_number)
}

// TODO TopPanelInternalViewAlertProduct replace with proper check
func (view *TopPanelInternalView) AlertProduct() {
	view.internal_product_field.Alert()
}

func (view *TopPanelInternalView) product_field_pop_data(str string) {
	log.Println("Warn: Debug: product_field_pop_data product_id", view.QC_Product.Product_id)
	log.Println("Warn: Debug: product_field_pop_data product_lot_id", view.QC_Product.Product_Lot_id)

	view.internal_product_field.Ok()

	// if product_lot.product_id != product_lot.insel_product_id(str) {
	old_product_id := view.QC_Product.Product_id
	view.QC_Product.Product_name = str
	view.QC_Product.Insel_product_self()

	if view.QC_Product.Product_id != old_product_id {
		view.SetTitle(view.QC_Product.Product_name)
		view.QC_Product.ResetQC()
		view.QC_Product.Blend = nil

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

func (view *TopPanelInternalView) product_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	view.internal_product_field.SetText(formatted_text)
	if view.internal_product_field.Text() != "" {
		view.product_field_pop_data(view.internal_product_field.Text())
		log.Println("Debug: product_field_text_pop_data", view.QC_Product)
	} else {
		view.QC_Product.Product_id = DB.INVALID_ID
		view.internal_product_field.Error()

	}
}

func (view *TopPanelInternalView) lot_field_pop_data(str string) {
	view.QC_Product.Update_lot(str, view.customer_field.Text())
	view.SetTitle(str)
}
func (view *TopPanelInternalView) lot_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	view.lot_field.SetText(formatted_text)

	view.lot_field_pop_data(formatted_text)
}

func (view *TopPanelInternalView) customer_field_pop_data(str string) {
	view.QC_Product.Update_lot(view.lot_field.Text(), str)
}
func (view *TopPanelInternalView) customer_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	view.customer_field.SetText(formatted_text)

	view.customer_field_pop_data(formatted_text)
}

func (view *TopPanelInternalView) sample_field_pop_data(str string) {
	view.QC_Product.Sample_point = str
}
func (view *TopPanelInternalView) sample_field_text_pop_data(str string) {
	formatted_text := strings.ToUpper(strings.TrimSpace(str))
	view.sample_field.SetText(formatted_text)
	view.sample_field_pop_data(formatted_text)
}

func (view *TopPanelInternalView) SetCurrentTab(i int) {
	view.mainWindow.SetCurrentTab(i)
}
