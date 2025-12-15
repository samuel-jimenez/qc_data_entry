package qc_ui

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

const (
	UPDATE_LOT_TAB_NO = iota
	UPDATE_SAMPLE_TAB_NO
	UPDATE_SAMPLE_PT_TAB_NO
)

/*
 * WhupsView
 *
 */
type WhupsView struct {
	*windigo.AutoDialog

	// prod_panel *MeasuredProductView
	prod_panels []*MeasuredProductView

	tab_bar, prod_bar *windigo.TabView

	button_panel *windigo.AutoPanel
	buttons      []*windigo.PushButton

	// lot_field     *GUI.SearchBox
	// lot_field, update_sample_tab, update_lot_tab *GUI.SearchBox
	lot_field, update_lot_tab          *GUI.SearchBox
	sample_field, update_sample_tab    *GUI.SearchBox
	COL_ITEMS_LOT                      []string
	Product_id, Lot_id, Product_Lot_id int64
	Lot_number                         string

	/*
		product_panel, lot_panel *windigo.AutoPanel

		product_field, moniker_field,
		sample_field *GUI.ComboBox*/
}

func WhupsView_from_new(parent windigo.Controller,
	Measured_Products []*product.MeasuredProduct,
) *WhupsView {
	window_title := "Whups"
	view := new(WhupsView)

	if len(Measured_Products) > 0 {
		Measured_Product := Measured_Products[0]
		view.Lot_number = Measured_Product.Lot_number
	}

	sample_string := "Sample"
	if len(Measured_Products) > 1 {
		sample_string = "Samples"
	}

	update_lot_text := "Change Lot Number"
	//  make sure no collision
	// 	if using : format, reorder
	update_sample_text := "Move " + sample_string + " To New Lot"
	update_sample_pt_text := "Change Sample Pt"

	// "Change Lot Product" ??? really?

	// Change Lot Number

	// Move Sample To New Lot

	// Change Sample Pt

	// Change Lot Product

	// 	 typo

	view.AutoDialog = windigo.AutoDialog_from_new(parent)
	view.Dialog.SetText(window_title)

	// 	get QC
	view.lot_field = GUI.NewLabeledListSearchBoxFromQuery(view, formats.COL_LABEL_LOT,
		DB.DB_Select_product_lot_all)

	// view.lot_field.SetTheme("DarkMode_Explorer")
	// view.lot_field.SetTheme("DarkMode_CFD")
	// view.lot_field.SetTheme("Explorer")
	// view.lot_field.SetBGColor(windigo.Olive)

	// view.sample_field = GUI.NewLabeledSearchBoxFromQuery(view, formats.COL_LABEL_SAMPLE_PT,
	// DB.DB_Select_all_sample_points)

	// CardView
	view.prod_bar = windigo.CardView_from_new(
		view,
		// windigo.Pen_SolidColor_from_new(3, windigo.CornflowerBlue),
		windigo.Pen_SolidColor_from_new(3, windigo.MediumSlateBlue),

		// windigo.Pen_SolidColor_from_new(3, windigo.LightCyan),
		// windigo.Pen_SolidColor_from_new(3, windigo.RGB(128, 128, 255)),
		windigo.Pen_from_Color(windigo.Black),
	)

	// highlightPen := windigo.Pen_from_flags(w32.PS_GEOMETRIC, 2, windigo.NewSolidColorBrush(windigo.RGB(128, 128, 255)))
	// highlightPen := windigo.Pen_geometric_from_new(windigo.Style_Solid, 2, windigo.NewSolidColorBrush(windigo.RGB(128, 128, 255)))
	// highlightPen := windigo.Pen_SolidColor_from_new(3, windigo)
	// normyPen := windigo.Pen_SolidColor_from_new(1, windigo.Black)
	/*
		view.prod_bar.SetSelectedBorder(highlightPen)
		view.prod_bar.SetUnselectedBorder(normyPen)
		view.prod_bar.SetShowOnePanel(false)*/

	// tabs
	view.tab_bar = windigo.TabView_from_new(view)
	view.update_lot_tab = GUI.NewLabeledSearchBox(view.tab_bar.AddAutoPanel(update_lot_text), formats.COL_LABEL_LOT)
	view.update_lot_tab.Update(view.lot_field.Entries())
	view.tab_bar.AttachPanel(update_sample_text, view.update_lot_tab.Parent().(*windigo.AutoPanel))
	view.update_sample_tab = GUI.NewLabeledSearchBoxFromQuery(view.tab_bar.AddAutoPanel(update_sample_pt_text), formats.COL_LABEL_SAMPLE_PT,
		DB.DB_Select_all_sample_points)

	view.button_panel = windigo.NewAutoPanel(view)

	view.Dock(view.lot_field, windigo.Top)

	view.Dock(view.prod_bar, windigo.Top)
	// TODO ;abels
	view.Dock(view.prod_bar.Panels(), windigo.Top)

	// view.Dock(lot_field_AutoPanel, windigo.Top)
	view.Dock(view.button_panel, windigo.Bottom)
	view.Dock(view.tab_bar.Panels(), windigo.Bottom)
	view.Dock(view.tab_bar, windigo.Bottom|windigo.Overflow_Expand)

	view.Update(Measured_Products)

	accept_button := windigo.NewPushButton(view.button_panel)
	accept_button.SetText("OK")
	view.buttons = append(view.buttons, accept_button)
	// accept_button.SetTheme("DarkMode_Explorer")

	cancel_button := windigo.NewPushButton(view.button_panel)
	cancel_button.SetText("Cancel")
	view.buttons = append(view.buttons, cancel_button)
	// cancel_button.SetTheme("Explorer")
	// cancel_button.SetBGColor(windigo.Olive)

	view.button_panel.Dock(accept_button, windigo.Left)
	view.button_panel.Dock(cancel_button, windigo.Left)

	// event handling
	accept_button.OnClick().Bind(view.accept_button_OnClick)
	cancel_button.OnClick().Bind(view.EventClose)
	view.SetButtons(accept_button, cancel_button)

	view.lot_field.OnSelectedChange().Bind(view.lot_field_OnSelectedChange)

	return view
}

func (view *WhupsView) SetSize(w, h int) {
	view.Dialog.SetSize(w+GUI.WINDOW_FUDGE_MARGIN_W, h+GUI.WINDOW_FUDGE_MARGIN_H)
}

func (view *WhupsView) SetFont(font *windigo.Font) {
	view.prod_bar.SetFont(font)
	for _, prod_panel := range view.prod_panels {
		prod_panel.SetFont(font)
	}
}

func (view *WhupsView) RefreshSize() {
	view.SetSize(2*GUI.GROUP_WIDTH, GUI.GROUP_HEIGHT+2*GUI.RANGES_BUTTON_HEIGHT+2*GUI.PRODUCT_FIELD_HEIGHT)

	view.lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.button_panel.SetSize(2*GUI.GROUP_WIDTH, GUI.RANGES_BUTTON_HEIGHT)
	for _, button := range view.buttons {
		button.SetSize(GUI.RANGES_BUTTON_WIDTH, GUI.RANGES_BUTTON_HEIGHT)
	}

	view.tab_bar.SetBaseSize(GUI.GROUP_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.tab_bar.Panels().SetSize(2*GUI.GROUP_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.update_lot_tab.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.update_sample_tab.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.prod_bar.SetSize(GUI.GROUP_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.prod_bar.Panels().SetSize(GUI.OFF_AXIS, GUI.GROUP_HEIGHT)

	view.prod_bar.SetPanelSize(GUI.GROUP_WIDTH)
	view.prod_bar.SetPanelMargin(50)

	for _, prod_panel := range view.prod_panels {
		prod_panel.RefreshSize()
	}
}

func (view *WhupsView) Update(Measured_Products []*product.MeasuredProduct) {
	sample_string := "Sample"

	if len(Measured_Products) > 1 {
		sample_string = "Samples"
	}
	update_sample_text := "Move " + sample_string + " To New Lot"
	view.tab_bar.SetPanelText(UPDATE_SAMPLE_TAB_NO, update_sample_text)

	for _, prod_panel := range view.prod_panels {

		// for i, prod_panel := range view.prod_panels {
		// log.Println("Measured_Product.len", len(view.prod_panels), i)
		// if i == 0 {
		view.prod_bar.DeletePanel(0)
		// }
		prod_panel.Close()
	}
	view.prod_panels = nil

	for _, Measured_Product := range Measured_Products {
		// view.prod_panel = MeasuredProductView_from_new(view, Measured_Product)
		// view.prod_panels = append(view.prod_panels, MeasuredProductView_from_new(view, Measured_Product))

		// prod_panel := MeasuredProductView_from_new(view, Measured_Product)
		// view.prod_bar.AttachPanel(Measured_Product.Sample_point, prod_panel)
		// view.prod_panels = append(view.prod_panels, prod_panel)
		// view.Dock(prod_panel, windigo.Left)
		// view.prod_panels = append(view.prod_panels, MeasuredProductView_from_AutoPanel(view.prod_bar.AddAutoPanel(Measured_Product.Sample_point), Measured_Product))

		prod_panel := MeasuredProductView_from_new(view.prod_bar.Panels(), Measured_Product)
		view.prod_bar.AttachPanel(Measured_Product.Sample_point, prod_panel)
		view.prod_panels = append(view.prod_panels, prod_panel)

	}
	view.lot_field.SetText(view.Lot_number)
	view.update_lot_tab.SetText(view.Lot_number)

	view.RefreshSize()
}

// TODO extract Lot_Oper_07
// sample package
// make sure no collision
//
// if using ':' format, reorder
func MoveLot(Product_id, Lot_id, Product_Lot_id int64, src_lot, dest_lot string) (other_Lot_id int64, err error) {
	proc_name := "movelot"

	Dest_lot, err, other_Lot_id := DB.DestLot(dest_lot, Product_id, Lot_id, Product_Lot_id)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return
	}
	log.Printf("Modify: [MoveLot] Updating {Lot_number %s} (%d) to %s", src_lot, Lot_id, Dest_lot)

	DB.Update(proc_name,
		DB.DB_Update_lot_list_name,
		Lot_id,
		Dest_lot,
	)

	DB.ShrinkLot(src_lot)
	return
}

func MoveSample(QC_id, Lot_id int64, src_lot, dest_lot string) {
	proc_name := "MoveSample"
	log.Printf("Modify: [WhupsView] Moving {sample %d} from Lot_number %s to Lot_number %s (%d)", QC_id, src_lot, dest_lot, Lot_id)
	DB.Update(proc_name,
		DB.DB_Update_qc_samples_id,
		QC_id,
		Lot_id,
	)
}

func DeleteSample(QC_id int64) {
	proc_name := "DeleteSample"
	log.Printf("Modify: [WhupsView] Deleting {sample %d} ", QC_id)
	DB.Delete(proc_name,
		DB.DB_Delete_qc_samples_id,
		QC_id,
	)
}

func MoveSamplePt(QC_id int64, src_Sample_point, dest_Sample_point string) {
	proc_name := "Update-Sample_point"
	log.Printf("Modify: [WhupsView] Moving {sample %d} from Sample_point %s to Sample_point %s", QC_id, src_Sample_point, dest_Sample_point)
	DB.Insert(proc_name, DB.DB_Insel_sample_point, dest_Sample_point)

	DB.Update(proc_name,
		DB.DB_Update_qc_samples_sample_point,
		QC_id,
		dest_Sample_point,
	)
}

func (view *WhupsView) accept_button_OnClick(*windigo.Event) {
	proc_name := "WhupsView-Get-ID"
	DB.Select_Error(proc_name,
		DB.DB_Select_product_lot__product__lot.QueryRow(view.Lot_number),
		&view.Product_id,
		&view.Lot_id,
		&view.Product_Lot_id,
	)
	switch view.tab_bar.Current() {
	case UPDATE_LOT_TAB_NO:
		src_lot, dest_lot := view.lot_field.GetSelectedItem(), view.update_lot_tab.Text()
		if other_Lot_id, err := MoveLot(
			view.Product_id,
			view.Lot_id,
			view.Product_Lot_id,
			src_lot,
			dest_lot,
		); err != nil {
			Check_dupe_sample_lots(view.Parent(), err, other_Lot_id, view.Lot_id, func() {
				MoveLot(
					view.Product_id,
					view.Lot_id,
					view.Product_Lot_id,
					src_lot,
					dest_lot,
				)
			})
		}
	case UPDATE_SAMPLE_TAB_NO:
		src_lot := view.lot_field.GetSelectedItem()
		Lot_number, Lot_id, _ := DB.Select_product_lot_name(view.update_lot_tab.Text(), view.Product_id)
		for _, prod_panel := range view.prod_panels {
			if prod_panel.Checked.Checked() {
				// Check for dupes
				if old_QC_id, err := product.CheckDupes(
					Lot_id, prod_panel.Sample_point,
				); err != nil {
					Check_dupe_sample_ids(view.Parent(), err, old_QC_id, prod_panel.QC_id, func() {
						MoveSample(prod_panel.QC_id, Lot_id, src_lot, Lot_number)
					})
					continue
				}
				MoveSample(prod_panel.QC_id, Lot_id, src_lot, Lot_number)
			}
		}
	case UPDATE_SAMPLE_PT_TAB_NO:
		prod_panel := view.prod_bar.Panels().At(view.prod_bar.Current()).(*MeasuredProductView)
		Sample_point := view.update_sample_tab.Text()
		// Check for dupes
		if old_QC_id, err := product.CheckDupes(
			view.Lot_id, Sample_point,
		); err != nil {
			Check_dupe_sample_ids(view.Parent(), err, old_QC_id, prod_panel.QC_id, func() {
				MoveSamplePt(prod_panel.QC_id, prod_panel.Sample_point, Sample_point)
			})
			view.Close()
		}
		MoveSamplePt(prod_panel.QC_id, prod_panel.Sample_point, Sample_point)
	}

	view.Close()
}

func (view *WhupsView) lot_field_OnSelectedChange(*windigo.Event) {
	view.Lot_number = view.lot_field.GetSelectedItem()
	view.Update(product.MeasuredProduct_Array_from_SQL_Lot_number(view.Lot_number))
}
