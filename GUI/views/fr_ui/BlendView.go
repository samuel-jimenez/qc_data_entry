package fr_ui

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/windigo"
)

/*
 * BlendViewer
 *
 */
type BlendViewer interface {
	views.QCBlendViewer
	// windigo.Pane
	// Get() *blender.ProductBlend
	// UpdateRecipe(blend *blender.ProductBlend)
	SetTotal(float64)
	SetProperTotal(float64)
	SetHeel(float64)
	SetAmount(float64)
	SetProperAmount(float64)
	// SetFont(font *windigo.Font)
	// RefreshSize()
}

// RecipeHeader
type BlendHeader struct {
	*windigo.AutoPanel
	Component_Label, Weight_Label, Lot_Label,
	Density_Label, Gallons_Label *windigo.Label
}

func Header_units(title, units string) string {
	return title + "\n(" + units + ")"
}

func NewBlend_Header(parent windigo.Controller) *BlendHeader {
	view := new(BlendHeader)

	view.AutoPanel = windigo.NewAutoPanel(parent)
	// 	TODO!!! RecipeHeader
	// SG_UNITS        = "g/mL"
	view.Component_Label = windigo.NewLabel(view.AutoPanel)
	view.Component_Label.SetText("Component")
	view.AutoPanel.Dock(view.Component_Label, windigo.Left)

	view.Weight_Label = windigo.NewLabel(view.AutoPanel)
	view.Weight_Label.SetText(Header_units("Weight", "lb"))
	view.AutoPanel.Dock(view.Weight_Label, windigo.Left)

	view.Lot_Label = windigo.NewLabel(view.AutoPanel)
	view.Lot_Label.SetText("Component Lot")
	view.AutoPanel.Dock(view.Lot_Label, windigo.Left)

	view.Density_Label = windigo.NewLabel(view.AutoPanel)
	view.Density_Label.SetText(Header_units("Density", formats.DENSITY_UNITS))
	view.AutoPanel.Dock(view.Density_Label, windigo.Left)

	view.Gallons_Label = windigo.NewLabel(view.AutoPanel)
	view.Gallons_Label.SetText(Header_units("Volume", "gal"))
	view.AutoPanel.Dock(view.Gallons_Label, windigo.Left)

	return view
}

func (view *BlendHeader) SetFont(font *windigo.Font) {
	view.Component_Label.SetFont(font)
	view.Weight_Label.SetFont(font)
	view.Lot_Label.SetFont(font)
	view.Density_Label.SetFont(font)
	view.Gallons_Label.SetFont(font)
}

func (view *BlendHeader) RefreshSize() {
	view.SetSize(GUI.OFF_AXIS, 2*GUI.PRODUCT_FIELD_HEIGHT)
	view.Component_Label.SetSize(GUI.SOURCES_LABEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Weight_Label.SetSize(GUI.DATA_SUBFIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Lot_Label.SetSize(GUI.SOURCES_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Density_Label.SetSize(GUI.DATA_SUBFIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Gallons_Label.SetSize(GUI.SOURCES_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
}

//TODO combine QCBlendView
/*
 * BlendView
 *
 */
type BlendView struct {
	views.QCBlendView
	total_field, heel_field, amount_field *views.NumberEditView
	headers                               *BlendHeader
	Heel_Component                        *BlendComponentView
}

func NewBlendView(parent windigo.Controller) *BlendView {
	view := new(BlendView)
	view.QCBlendView = *views.NewQCBlendView(parent)
	//.TODO  strap

	total_text := "Total"
	heel_text := "Heel"

	amount_text := "Blend"
	// amount_text := "Produced"
	// amount_text := "Amount"
	// amount_text := "Quantity"

	view.total_field = views.NewNumberEditView(view.Panel, total_text)
	view.heel_field = views.NewNumberEditView(view.Panel, heel_text)
	view.amount_field = views.NewNumberEditView(view.Panel, amount_text)
	view.headers = NewBlend_Header(view.Panel)

	view.Heel_Component = NewHeelBlendComponentView(view)

	// view.Panel.D	ock(view.amount_field, windigo.Left)
	view.Panel.Dock(view.total_field, windigo.TopLeft)
	view.Panel.Dock(view.heel_field, windigo.TopLeft)
	view.Panel.Dock(view.amount_field, windigo.TopLeft)
	view.Panel.Dock(view.headers, windigo.Top)
	view.Dock(view.Heel_Component, windigo.Top)

	view.total_field.OnChange().Bind(func(e *windigo.Event) {
		view.SetTotal(view.total_field.Get())
	})
	view.heel_field.OnChange().Bind(func(e *windigo.Event) {
		view.SetHeel(view.heel_field.Get())
	})
	view.amount_field.OnChange().Bind(func(e *windigo.Event) {
		view.SetAmount(view.amount_field.Get())
	})

	amount := 100.
	view.heel_field.SetInt(0)
	view.SetProperTotal(amount)
	view.SetProperAmount(amount)

	return view
}

func (view *BlendView) UpdateRecipe(recipe *blender.ProductRecipe) {
	if view.Recipe == recipe {
		return
	}

	for _, component := range view.Components {
		component.Close()
	}
	view.Components = nil
	view.Recipe = recipe
	log.Println("BlendView Update", view.Recipe)
	if view.Recipe == nil {
		return
	}

	amount := view.amount_field.Get()

	height := view.ClientHeight()
	delta_height := GUI.PRODUCT_FIELD_HEIGHT
	width := view.ClientWidth()

	for _, component := range view.Recipe.Components {
		component_view := NewBlendComponentView(view, &component)
		view.Components = append(view.Components, component_view)
		view.AutoPanel.Dock(component_view, windigo.Top)
		height += delta_height
		component_view.SetAmount(amount)
	}
	view.SetSize(width, height)
}

func (view *BlendView) SetTotal(total float64) {
	view.SetProperTotal(total)
	heel := view.heel_field.Get()
	view.SetProperAmount(total - heel)
}

func (view *BlendView) SetProperTotal(total float64) {
	view.total_field.SetInt(total)
}

func (view *BlendView) SetHeel(heel float64) {
	view.heel_field.SetInt(heel)
	view.Heel_Component.SetAmount(heel)
	total := view.total_field.Get()
	view.SetProperAmount(total - heel)
}

func (view *BlendView) SetAmount(amount float64) {
	heel := view.heel_field.Get()
	view.SetProperTotal(amount + heel)
	view.SetProperAmount(amount)
}

func (view *BlendView) SetProperAmount(amount float64) {
	view.amount_field.SetInt(amount)
	for _, component := range view.Components {
		component.(*BlendComponentView).SetAmount(amount)
	}
}

func (view *BlendView) SetFont(font *windigo.Font) {
	view.total_field.SetFont(font)
	view.heel_field.SetFont(font)
	view.amount_field.SetFont(font)
	view.headers.SetFont(font)
	view.Heel_Component.SetFont(font)

	for _, component := range view.Components {
		component.SetFont(font)
	}
}

func (view *BlendView) RefreshSize() {
	height := 5 * GUI.PRODUCT_FIELD_HEIGHT
	delta_height := GUI.PRODUCT_FIELD_HEIGHT

	view.SetSize(GUI.OFF_AXIS, height+len(view.Components)*delta_height)
	view.Panel.SetSize(GUI.TOP_PANEL_WIDTH, height)
	view.total_field.SetLabeledSize(GUI.LABEL_WIDTH-GUI.ERROR_MARGIN, GUI.PRODUCT_FIELD_WIDTH, delta_height)
	view.heel_field.SetLabeledSize(GUI.LABEL_WIDTH-GUI.ERROR_MARGIN, GUI.PRODUCT_FIELD_WIDTH, delta_height)
	view.amount_field.SetLabeledSize(GUI.LABEL_WIDTH-GUI.ERROR_MARGIN, GUI.PRODUCT_FIELD_WIDTH, delta_height)
	view.headers.RefreshSize()
	view.Heel_Component.RefreshSize()

	for _, component := range view.Components {
		component.RefreshSize()
	}

}
