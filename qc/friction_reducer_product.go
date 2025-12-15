package qc

import (
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/qc_ui"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

type FrictionReducerProduct struct {
	product.QCProduct
	sg          float64
	string_test int64
	viscosity   int64
}

func (fr_product FrictionReducerProduct) toProduct() *product.MeasuredProduct {
	return &product.MeasuredProduct{
		QCProduct:   fr_product.Base(),
		PH:          nullable.NewNullFloat64(0, false),
		SG:          nullable.NewNullFloat64(fr_product.sg, true),
		Density:     nullable.NewNullFloat64(formats.Density_from_sg(fr_product.sg), true),
		String_test: nullable.NewNullInt64(fr_product.string_test),
		Viscosity:   nullable.NewNullInt64(fr_product.viscosity),
	}
}

func MeasuredProduct_from_FrictionReducerProductView(base_product product.QCProduct, viscosity, mass, string_test float64) *product.MeasuredProduct {
	sg := formats.SG_from_mass(mass)

	return FrictionReducerProduct{base_product, sg, int64(string_test), int64(viscosity)}.toProduct()
}

func (product *FrictionReducerProduct) Check_data() bool {
	return true
}

type FrictionReducerProductViewer interface {
	*windigo.AutoPanel
	Get(base_product product.BaseProduct, replace_sample_point bool) product.MeasuredProduct
	Clear()
	SetFont(font *windigo.Font)
	Refresh()
}

type FrictionReducerProductView struct {
	*windigo.AutoPanel
	visual_field *views.BoolCheckboxView
	viscosity_field,
	string_field *GUI.NumberEditView
	density_field *qc_ui.MassDataView
	sample_point  string
}

func FrictionReducerProductView_from_new(parent *windigo.AutoPanel, sample_point string, ranges_panel *FrictionReducerProductRangesView) *FrictionReducerProductView {
	view := new(FrictionReducerProductView)

	view.AutoPanel = windigo.NewGroupAutoPanel(parent)
	view.AutoPanel.SetText(sample_point)

	view.visual_field = views.NewBoolCheckboxView(view.AutoPanel, VISUAL_TEXT)

	view.viscosity_field = GUI.NumberEditView_from_new(view.AutoPanel, formats.VISCOSITY_TEXT)

	view.string_field = GUI.NumberEditView_with_Change_from_new(view.AutoPanel, formats.STRING_TEXT_MINI, ranges_panel.string_field)

	view.density_field = qc_ui.MassDataView_from_new(view.AutoPanel, ranges_panel)
	view.density_field.SetZAfter(view.viscosity_field)

	view.AutoPanel.Dock(view.visual_field, windigo.Top)
	view.AutoPanel.Dock(view.viscosity_field, windigo.Top)
	view.AutoPanel.Dock(view.density_field, windigo.Top)
	view.AutoPanel.Dock(view.string_field, windigo.Top)

	view.sample_point = strings.ToUpper(sample_point)
	return view
}

func (view *FrictionReducerProductView) Get(base_product product.QCProduct, replace_sample_point bool) *product.MeasuredProduct {
	base_product.Visual = view.visual_field.Checked()
	if replace_sample_point {
		base_product.Sample_point = view.sample_point
	}
	return MeasuredProduct_from_FrictionReducerProductView(base_product, view.viscosity_field.Get(), view.density_field.Get(), view.string_field.Get())
}

func (view *FrictionReducerProductView) SetFont(font *windigo.Font) {
	view.AutoPanel.SetFont(font)
	view.visual_field.SetFont(font)
	view.viscosity_field.SetFont(font)
	view.density_field.SetFont(font)
	view.string_field.SetFont(font)
}

func (view *FrictionReducerProductView) RefreshSize() {
	view.SetSize(GUI.GROUP_WIDTH, GUI.GROUP_HEIGHT)
	view.SetPaddings(GUI.TOP_SPACER_WIDTH, GUI.TOP_SPACER_HEIGHT, GUI.BTM_SPACER_WIDTH, GUI.BTM_SPACER_HEIGHT)

	view.visual_field.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)
	view.viscosity_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.EDIT_FIELD_HEIGHT)
	view.density_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.DATA_SUBFIELD_WIDTH, GUI.DATA_UNIT_WIDTH, GUI.EDIT_FIELD_HEIGHT)
	view.string_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.EDIT_FIELD_HEIGHT)
}

func (view *FrictionReducerProductView) Clear() {
	view.visual_field.Clear()
	view.viscosity_field.Clear()
	view.density_field.Clear()
	view.string_field.Clear()
}

func (view *FrictionReducerProductView) FocusVisual() bool { view.visual_field.SetFocus(); return true }
func (view *FrictionReducerProductView) FocusViscosity() bool {
	view.viscosity_field.SetFocus()
	return true
}

func (view *FrictionReducerProductView) FocusDensity() bool {
	view.density_field.SetFocus()
	return true
}
func (view *FrictionReducerProductView) FocusString() bool { view.string_field.SetFocus(); return true }

func (view *FrictionReducerProductView) Interleave(bottom_group *FrictionReducerProductView) {
	num_back_shortcut := windigo.Shortcut{Key: windigo.KeyDivide}
	num_fwd_shortcut := windigo.Shortcut{Key: windigo.KeyMultiply}
	num_prev_shortcut := windigo.Shortcut{Key: windigo.KeySubtract}
	num_next_shortcut := windigo.Shortcut{Key: windigo.KeyAdd}

	kb_prev_shortcut := windigo.Shortcut{Key: windigo.KeyW}
	kb_back_shortcut := windigo.Shortcut{Key: windigo.KeyA}
	kb_next_shortcut := windigo.Shortcut{Key: windigo.KeyS}
	kb_fwd_shortcut := windigo.Shortcut{Key: windigo.KeyD}

	view.visual_field.AddShortcut(num_prev_shortcut, view.FocusVisual)
	view.visual_field.AddShortcut(num_fwd_shortcut, bottom_group.FocusVisual)
	view.visual_field.AddShortcut(num_back_shortcut, view.FocusVisual)
	view.visual_field.AddShortcut(num_next_shortcut, view.FocusViscosity)
	view.visual_field.AddShortcut(kb_prev_shortcut, view.FocusVisual)
	view.visual_field.AddShortcut(kb_fwd_shortcut, bottom_group.FocusVisual)
	view.visual_field.AddShortcut(kb_back_shortcut, view.FocusVisual)
	view.visual_field.AddShortcut(kb_next_shortcut, view.FocusViscosity)

	bottom_group.visual_field.AddShortcut(num_prev_shortcut, bottom_group.FocusVisual)
	bottom_group.visual_field.AddShortcut(num_fwd_shortcut, bottom_group.FocusVisual)
	bottom_group.visual_field.AddShortcut(num_back_shortcut, view.FocusVisual)
	bottom_group.visual_field.AddShortcut(num_next_shortcut, bottom_group.FocusViscosity)
	bottom_group.visual_field.AddShortcut(kb_prev_shortcut, bottom_group.FocusVisual)
	bottom_group.visual_field.AddShortcut(kb_fwd_shortcut, bottom_group.FocusVisual)
	bottom_group.visual_field.AddShortcut(kb_back_shortcut, view.FocusVisual)
	bottom_group.visual_field.AddShortcut(kb_next_shortcut, bottom_group.FocusViscosity)

	view.viscosity_field.AddShortcut(num_prev_shortcut, view.FocusVisual)
	view.viscosity_field.AddShortcut(num_fwd_shortcut, bottom_group.FocusViscosity)
	view.viscosity_field.AddShortcut(num_back_shortcut, view.FocusViscosity)
	view.viscosity_field.AddShortcut(num_next_shortcut, view.FocusDensity)
	view.viscosity_field.AddShortcut(kb_prev_shortcut, view.FocusVisual)
	view.viscosity_field.AddShortcut(kb_fwd_shortcut, bottom_group.FocusViscosity)
	view.viscosity_field.AddShortcut(kb_back_shortcut, view.FocusViscosity)
	view.viscosity_field.AddShortcut(kb_next_shortcut, view.FocusDensity)

	bottom_group.viscosity_field.AddShortcut(num_prev_shortcut, bottom_group.FocusVisual)
	bottom_group.viscosity_field.AddShortcut(num_fwd_shortcut, bottom_group.FocusViscosity)
	bottom_group.viscosity_field.AddShortcut(num_back_shortcut, view.FocusViscosity)
	bottom_group.viscosity_field.AddShortcut(num_next_shortcut, bottom_group.FocusDensity)
	bottom_group.viscosity_field.AddShortcut(kb_prev_shortcut, bottom_group.FocusVisual)
	bottom_group.viscosity_field.AddShortcut(kb_fwd_shortcut, bottom_group.FocusViscosity)
	bottom_group.viscosity_field.AddShortcut(kb_back_shortcut, view.FocusViscosity)
	bottom_group.viscosity_field.AddShortcut(kb_next_shortcut, bottom_group.FocusDensity)

	view.density_field.AddShortcut(num_prev_shortcut, view.FocusViscosity)
	view.density_field.AddShortcut(num_fwd_shortcut, bottom_group.FocusDensity)
	view.density_field.AddShortcut(num_back_shortcut, view.FocusDensity)
	view.density_field.AddShortcut(num_next_shortcut, view.FocusString)
	view.density_field.AddShortcut(kb_prev_shortcut, view.FocusViscosity)
	view.density_field.AddShortcut(kb_fwd_shortcut, bottom_group.FocusDensity)
	view.density_field.AddShortcut(kb_back_shortcut, view.FocusDensity)
	view.density_field.AddShortcut(kb_next_shortcut, view.FocusString)

	bottom_group.density_field.AddShortcut(num_prev_shortcut, bottom_group.FocusViscosity)
	bottom_group.density_field.AddShortcut(num_fwd_shortcut, bottom_group.FocusDensity)
	bottom_group.density_field.AddShortcut(num_back_shortcut, view.FocusDensity)
	bottom_group.density_field.AddShortcut(num_next_shortcut, bottom_group.FocusString)
	bottom_group.density_field.AddShortcut(kb_prev_shortcut, bottom_group.FocusViscosity)
	bottom_group.density_field.AddShortcut(kb_fwd_shortcut, bottom_group.FocusDensity)
	bottom_group.density_field.AddShortcut(kb_back_shortcut, view.FocusDensity)
	bottom_group.density_field.AddShortcut(kb_next_shortcut, bottom_group.FocusString)

	view.string_field.AddShortcut(num_prev_shortcut, view.FocusDensity)
	view.string_field.AddShortcut(num_fwd_shortcut, bottom_group.FocusString)
	view.string_field.AddShortcut(num_back_shortcut, view.FocusString)
	view.string_field.AddShortcut(num_next_shortcut, view.FocusString)
	view.string_field.AddShortcut(kb_prev_shortcut, view.FocusDensity)
	view.string_field.AddShortcut(kb_fwd_shortcut, bottom_group.FocusString)
	view.string_field.AddShortcut(kb_back_shortcut, view.FocusString)
	view.string_field.AddShortcut(kb_next_shortcut, view.FocusString)

	bottom_group.string_field.AddShortcut(kb_prev_shortcut, bottom_group.FocusDensity)
	bottom_group.string_field.AddShortcut(kb_fwd_shortcut, bottom_group.FocusString)
	bottom_group.string_field.AddShortcut(kb_back_shortcut, view.FocusString)
	bottom_group.string_field.AddShortcut(kb_next_shortcut, bottom_group.FocusString)
	bottom_group.string_field.AddShortcut(num_prev_shortcut, bottom_group.FocusDensity)
	bottom_group.string_field.AddShortcut(num_fwd_shortcut, bottom_group.FocusString)
	bottom_group.string_field.AddShortcut(num_back_shortcut, view.FocusString)
	bottom_group.string_field.AddShortcut(num_next_shortcut, bottom_group.FocusString)
}
