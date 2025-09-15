package qc

import (
	"database/sql"
	"log"
	"maps"
	"slices"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender/blendbound"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

/*
 * TopPanelInboundViewer
 *
 */
type TopPanelInboundViewer interface {
	windigo.Controller
	SetFont(font *windigo.Font)
	RefreshSize()
	SetTitle(title string)
	// BaseProduct() product.BaseProduct

	Show()
	Hide()

	RefreshLots()
	RemoveInbound() *blendbound.InboundLot

	inbound_lot_field_OnSelectedChange(e *windigo.Event)
	inbound_container_field_OnSelectedChange(e *windigo.Event)
	testing_lot_field_OnSelectedChange(e *windigo.Event)
	sample_button_OnClick(e *windigo.Event)
	release_button_OnClick(e *windigo.Event)
	today_button_OnClick(e *windigo.Event)
}

/*
 * TopPanelInboundView
 *
 */
type TopPanelInboundView struct {
	*windigo.AutoPanel
	Inbound_Lot                              *blendbound.InboundLot
	QC_Product                               *product.QCProduct
	inbound_test                             string
	inbound_lot_data, inbound_container_data map[string]*blendbound.InboundLot

	mainWindow *windigo.Form

	product_panel_1_0, product_panel_1_1 *windigo.AutoPanel

	testing_lot_field, //inbound_product_field,
	inbound_lot_field, inbound_container_field *GUI.ComboBox

	inbound_product_field *windigo.LabeledEdit
	tester_field          *GUI.SearchBox

	sample_button, release_button, today_button, internal_button *windigo.PushButton
}

func NewTopPanelInboundView(
	mainWindow *windigo.Form,
	QC_Product *product.QCProduct,
	product_panel_1_0, product_panel_1_1 *windigo.AutoPanel,
	testing_lot_field, inbound_lot_field, inbound_container_field *GUI.ComboBox, inbound_product_field *windigo.LabeledEdit,
	sample_button, release_button, today_button, internal_button *windigo.PushButton,
) *TopPanelInboundView {

	view := new(TopPanelInboundView)

	// build object
	view.QC_Product = QC_Product
	view.inbound_test = "BSQL%"
	view.inbound_lot_data = make(map[string]*blendbound.InboundLot)
	view.inbound_container_data = make(map[string]*blendbound.InboundLot)

	view.mainWindow = mainWindow

	view.product_panel_1_0 = product_panel_1_0
	view.product_panel_1_1 = product_panel_1_1

	view.inbound_lot_field = inbound_lot_field
	view.testing_lot_field = testing_lot_field
	view.inbound_product_field = inbound_product_field
	view.inbound_container_field = inbound_container_field

	view.sample_button = sample_button
	view.release_button = release_button
	view.today_button = today_button
	view.internal_button = internal_button

	//
	// Dock
	product_panel_1_0.Dock(testing_lot_field, windigo.Left)
	product_panel_1_0.Dock(inbound_product_field, windigo.Left)

	product_panel_1_1.Dock(inbound_lot_field, windigo.Left)
	product_panel_1_1.Dock(inbound_container_field, windigo.Left)

	//
	// combobox
	view.RefreshLots()

	proc_name := "TopPanelInboundView.FillInbound"
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
			view.inbound_lot_data[Inbound.Lot_number] = Inbound
			view.inbound_container_data[Inbound.Container_name] = Inbound
			return nil
		},
		DB.DB_Select_inbound_lot_status, blendbound.Status_AVAILABLE)

	for _, Container_name := range slices.Sorted(maps.Keys(view.inbound_container_data)) {
		inbound_container_field.AddItem(Container_name)
	}

	//
	// functionality
	inbound_lot_field.OnSelectedChange().Bind(view.inbound_lot_field_OnSelectedChange)
	inbound_container_field.OnSelectedChange().Bind(view.inbound_container_field_OnSelectedChange)
	testing_lot_field.OnSelectedChange().Bind(view.testing_lot_field_OnSelectedChange)

	sample_button.OnClick().Bind(view.sample_button_OnClick)
	release_button.OnClick().Bind(view.release_button_OnClick)
	today_button.OnClick().Bind(view.today_button_OnClick)

	return view
}

func (view *TopPanelInboundView) SetFont(font *windigo.Font) {

	view.inbound_lot_field.SetFont(font)
	view.testing_lot_field.SetFont(font)
	view.inbound_product_field.SetFont(font)
	view.inbound_container_field.SetFont(font)

	view.sample_button.SetFont(font)
	view.release_button.SetFont(font)
	view.today_button.SetFont(font)
	view.internal_button.SetFont(font)
}

func (view *TopPanelInboundView) RefreshSize() {
	TODAY_BUTTON_MARGIN_L := 2*GUI.LABEL_WIDTH + GUI.PRODUCT_FIELD_WIDTH + GUI.TOP_PANEL_INTER_SPACER_WIDTH - GUI.SMOL_BUTTON_WIDTH - GUI.REPRINT_BUTTON_WIDTH - 2*BUTTON_MARGIN

	view.product_panel_1_0.SetSize(GUI.HPANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.product_panel_1_0.SetMarginTop(GUI.TOP_SPACER_HEIGHT)
	view.product_panel_1_0.SetMarginLeft(GUI.HPANEL_MARGIN)

	view.testing_lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.inbound_product_field.SetMarginLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)
	view.inbound_product_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.product_panel_1_1.SetSize(GUI.HPANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.product_panel_1_1.SetMarginTop(GUI.INTER_SPACER_HEIGHT)
	view.product_panel_1_1.SetMarginLeft(GUI.HPANEL_MARGIN)

	view.inbound_lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.inbound_container_field.SetMarginLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)
	view.inbound_container_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.sample_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
	view.sample_button.SetMarginsAll(BUTTON_MARGIN)

	view.release_button.SetMarginsAll(BUTTON_MARGIN)
	view.release_button.SetSize(GUI.REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.today_button.SetMarginsAll(BUTTON_MARGIN)
	view.today_button.SetMarginLeft(TODAY_BUTTON_MARGIN_L)
	view.today_button.SetSize(GUI.REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.internal_button.SetMarginsAll(BUTTON_MARGIN)
	view.internal_button.SetSize(GUI.REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

}

func (view *TopPanelInboundView) SetTitle(title string) {
	if view.mainWindow == nil {
		return
	}
	view.mainWindow.SetText(title)
}

func (view *TopPanelInboundView) Show() {
	view.product_panel_1_0.Show()
	view.product_panel_1_1.Show()
	view.sample_button.Show()
	view.release_button.Show()
	view.today_button.Show()
	view.internal_button.Show()

}

func (view *TopPanelInboundView) Hide() {
	view.product_panel_1_0.Hide()
	view.product_panel_1_1.Hide()
	view.sample_button.Hide()
	view.release_button.Hide()
	view.today_button.Hide()
	view.internal_button.Hide()

}

func (view *TopPanelInboundView) RefreshLots() {
	GUI.Fill_combobox_from_query(view.testing_lot_field,
		DB.DB_Select_lot_list_for_name_status, view.inbound_test, product.Status_BLENDED)
	view.testing_lot_field.SetSelectedItem(-1)
	view.inbound_lot_field.SetSelectedItem(-1)
	view.inbound_container_field.SetSelectedItem(-1)
	view.inbound_product_field.SetText("")
}

func (view *TopPanelInboundView) RemoveInbound() *blendbound.InboundLot {
	if view.Inbound_Lot == nil {
		return nil
	}

	Inbound__Lot_ := view.Inbound_Lot
	view.Inbound_Lot = nil

	delete(view.inbound_container_data, Inbound__Lot_.Container_name)
	delete(view.inbound_lot_data, Inbound__Lot_.Lot_number)

	return Inbound__Lot_
}

func (view *TopPanelInboundView) inbound_lot_field_OnSelectedChange(e *windigo.Event) {
	view.testing_lot_field.SetSelectedItem(-1)
	view.Inbound_Lot = view.inbound_lot_data[view.inbound_lot_field.GetSelectedItem()]
	if view.Inbound_Lot == nil {
		return
	}

	view.inbound_product_field.SetText(view.Inbound_Lot.Product_name)
	view.inbound_container_field.SetText(view.Inbound_Lot.Container_name)
}

func (view *TopPanelInboundView) inbound_container_field_OnSelectedChange(e *windigo.Event) {
	view.testing_lot_field.SetSelectedItem(-1)
	view.Inbound_Lot = view.inbound_container_data[view.inbound_container_field.GetSelectedItem()]
	if view.Inbound_Lot == nil {
		return
	}

	view.inbound_product_field.SetText(view.Inbound_Lot.Product_name)
	view.inbound_lot_field.SetText(view.Inbound_Lot.Lot_number)
}

func (view *TopPanelInboundView) testing_lot_field_OnSelectedChange(e *windigo.Event) {
	view.Inbound_Lot = nil
	view.QC_Product.Update_testing_lot(view.testing_lot_field.GetSelectedItem())
	view.inbound_product_field.SetText(view.QC_Product.Product_name)

	blend := view.QC_Product.Blend
	if blend == nil || len(blend.Components) < 1 {
		log.Println("Err: testing_lot_field.OnSelectedChange", view.QC_Product, blend)
		return
	}

	component := blend.Components[0]
	view.inbound_lot_field.SetText(component.Lot_name)
	view.inbound_container_field.SetText(component.Container_name)
	// inbound_product_field.SetText(component.Component_name)
	view.SetTitle(view.QC_Product.Lot_number)

	// update component list
	view.QC_Product.ResetQC()
	view.QC_Product.Container_type = product.CONTAINER_SAMPLE
	view.QC_Product.Update()
}

func (view *TopPanelInboundView) sample_button_OnClick(e *windigo.Event) {
	if view.Inbound_Lot == nil {
		return
	}

	Inbound__Lot_ := view.RemoveInbound()

	Inbound__Lot_.Sample()
	Inbound__Lot_.Quality_test()
	view.RefreshLots()

}

func (view *TopPanelInboundView) release_button_OnClick(e *windigo.Event) {

	proc_name := "TopPanelInboundView.release_button_OnClick"

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

	Inbound__Lot_ := view.RemoveInbound()
	log.Printf("Info: Releasing %s: %s - %s\n", Inbound__Lot_.Container_name, Inbound__Lot_.Product_name, Inbound__Lot_.Lot_number)

	Inbound__Lot_.Update_status(blendbound.Status_UNAVAILABLE)

	log.Println("TRACE: DEBUG: UPdate componnent DB_Update_lot_list__component_status", proc_name, Inbound__Lot_.Lot_number, Inbound__Lot_)

	if err := product.Release_testing_lot(Inbound__Lot_.Lot_number); err != nil {
		log.Println("error[]%S]:", proc_name, err)
		return
	}
	view.RefreshLots()
}

func (view *TopPanelInboundView) today_button_OnClick(e *windigo.Event) {
	if view.QC_Product == nil {
		return
	}
	view.QC_Product.Update_testing_tttToday_Junior()
	view.RefreshLots()
}
