package qc

import (
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type WaterBasedProduct struct {
	product.QCProduct
	sg float64
	ph float64
}

func (wb_product WaterBasedProduct) toProduct() *product.MeasuredProduct {
	PH := nullable.NewNullFloat64(wb_product.ph, true)
	NULL := nullable.NewNullFloat64(0, false)
	if strings.Contains(wb_product.Product_name, "BIONIX") && wb_product.ph == 0 {
		PH = NULL
	}
	return &product.MeasuredProduct{
		QCProduct:   wb_product.Base(),
		PH:          PH,
		SG:          nullable.NewNullFloat64(wb_product.sg, true),
		Density:     NULL,
		String_test: nullable.NullInt64Default(),
		Viscosity:   nullable.NullInt64Default(),
	}

	// TODO Option?
}

func MeasuredProduct_from_WaterBasedProductView(base_product product.QCProduct, have_visual bool, sg, ph float64) *product.MeasuredProduct {
	base_product.Visual = have_visual

	return WaterBasedProduct{base_product, sg, ph}.toProduct()
}

func (product WaterBasedProduct) Check_data() bool {
	return true
}

type WaterBasedProductViewer interface {
	Get(base_product product.BaseProduct) product.MeasuredProduct
	Clear()
	SetFont(font *windigo.Font)
	RefreshSize()
}

type WaterBasedProductView struct {
	*windigo.AutoPanel
	visual_field *views.BoolCheckboxView
	ph_field,
	sg_field *GUI.NumberEditView
}

func WaterBasedProductView_from_new(parent *windigo.AutoPanel, ranges_panel *WaterBasedProductRangesView) *WaterBasedProductView {
	view := new(WaterBasedProductView)

	group_panel := windigo.NewAutoPanel(parent)

	visual_field := views.NewBoolCheckboxView(group_panel, VISUAL_TEXT)

	ph_field := GUI.NumberEditView_with_Change_from_new(group_panel, formats.PH_TEXT, ranges_panel.ph_field)
	sg_field := GUI.NumberEditView_with_PointlessChange_from_new(group_panel, formats.SG_TEXT, ranges_panel.sg_field)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(ph_field, windigo.Top)
	group_panel.Dock(sg_field, windigo.Top)

	view.AutoPanel = group_panel
	view.visual_field = visual_field
	view.ph_field = ph_field
	view.sg_field = sg_field

	view.AddShortcuts()

	return view
}

func (view *WaterBasedProductView) AddShortcuts() {
	num_back_shortcut := windigo.Shortcut{Key: windigo.KeyDivide}
	num_fwd_shortcut := windigo.Shortcut{Key: windigo.KeyMultiply}
	num_prev_shortcut := windigo.Shortcut{Key: windigo.KeySubtract}
	num_next_shortcut := windigo.Shortcut{Key: windigo.KeyAdd}

	kb_prev_shortcut := windigo.Shortcut{Key: windigo.KeyW}
	kb_back_shortcut := windigo.Shortcut{Key: windigo.KeyA}
	kb_next_shortcut := windigo.Shortcut{Key: windigo.KeyS}
	kb_fwd_shortcut := windigo.Shortcut{Key: windigo.KeyD}

	view.visual_field.AddShortcut(num_prev_shortcut, view.FocusVisual)
	view.visual_field.AddShortcut(num_fwd_shortcut, view.FocusVisual)
	view.visual_field.AddShortcut(num_back_shortcut, view.FocusVisual)
	view.visual_field.AddShortcut(num_next_shortcut, view.FocusPH)
	view.visual_field.AddShortcut(kb_prev_shortcut, view.FocusVisual)
	view.visual_field.AddShortcut(kb_fwd_shortcut, view.FocusVisual)
	view.visual_field.AddShortcut(kb_back_shortcut, view.FocusVisual)
	view.visual_field.AddShortcut(kb_next_shortcut, view.FocusPH)

	view.ph_field.AddShortcut(num_prev_shortcut, view.FocusVisual)
	view.ph_field.AddShortcut(num_fwd_shortcut, view.FocusPH)
	view.ph_field.AddShortcut(num_back_shortcut, view.FocusPH)
	view.ph_field.AddShortcut(num_next_shortcut, view.FocusSG)
	view.ph_field.AddShortcut(kb_prev_shortcut, view.FocusVisual)
	view.ph_field.AddShortcut(kb_fwd_shortcut, view.FocusPH)
	view.ph_field.AddShortcut(kb_back_shortcut, view.FocusPH)
	view.ph_field.AddShortcut(kb_next_shortcut, view.FocusSG)

	view.sg_field.AddShortcut(num_prev_shortcut, view.FocusPH)
	view.sg_field.AddShortcut(num_fwd_shortcut, view.FocusSG)
	view.sg_field.AddShortcut(num_back_shortcut, view.FocusSG)
	view.sg_field.AddShortcut(num_next_shortcut, view.FocusSG)
	view.sg_field.AddShortcut(kb_prev_shortcut, view.FocusPH)
	view.sg_field.AddShortcut(kb_fwd_shortcut, view.FocusSG)
	view.sg_field.AddShortcut(kb_back_shortcut, view.FocusSG)
	view.sg_field.AddShortcut(kb_next_shortcut, view.FocusSG)
}

func (view *WaterBasedProductView) Get(base_product product.QCProduct) *product.MeasuredProduct {
	return MeasuredProduct_from_WaterBasedProductView(base_product, view.visual_field.Get(), view.sg_field.Get(), view.ph_field.Get())
}

func (view *WaterBasedProductView) Clear() {
	view.visual_field.Clear()
	view.sg_field.Clear()
	view.ph_field.Clear()
}

func (view *WaterBasedProductView) SetFont(font *windigo.Font) {
	view.AutoPanel.SetFont(font)
	view.visual_field.SetFont(font)
	view.sg_field.SetFont(font)
	view.ph_field.SetFont(font)
}

func (view *WaterBasedProductView) RefreshSize() {
	view.SetSize(GUI.GROUP_WIDTH, GUI.GROUP_HEIGHT)
	view.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.BTM_SPACER_WIDTH, GUI.BTM_SPACER_HEIGHT)

	view.visual_field.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)
	view.sg_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.EDIT_FIELD_HEIGHT)
	view.ph_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.EDIT_FIELD_HEIGHT)
}

func (view *WaterBasedProductView) FocusVisual() bool {
	view.visual_field.SetFocus()
	return true
}

func (view *WaterBasedProductView) FocusPH() bool {
	view.ph_field.SetFocus()
	return true
}

func (view *WaterBasedProductView) FocusSG() bool {
	view.sg_field.SetFocus()
	return true
}
