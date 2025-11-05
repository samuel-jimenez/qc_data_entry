package blender_ui

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/blender"
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

//TODO combine QCBlendView
/*
 * BlendView
 *
 */
type BlendView struct {
	views.QCBlendView

	parent                                             *BlendStrappingProductView
	total_field, gross_field, heel_field, amount_field *GUI.NumbestEditView
	headers                                            *BlendHeader
	Heel_Component, Total_Component                    *BlendComponentView
	Heel_Gross, Heel_Net, TareMass                     float64
}

func BlendView_from_new(parent *BlendStrappingProductView) *BlendView {
	view := new(BlendView)
	view.QCBlendView = *views.QCBlendView_from_new(parent)
	view.parent = parent

	// total_text := "Total"
	total_text := "Net"
	gross_text := "Gross"
	heel_text := "Heel"

	amount_text := "To Blend"
	// amount_text := "Produced"
	// amount_text := "Amount"
	// amount_text := "Quantity"

	view.gross_field = GUI.NumbestEditView_from_new(view.Panel, gross_text)
	view.total_field = GUI.NumbestEditView_from_new(view.Panel, total_text)
	view.heel_field = GUI.NumbestEditView_from_new(view.Panel, heel_text)
	view.amount_field = GUI.NumbestEditView_from_new(view.Panel, amount_text)
	view.headers = BlendHeader_from_new(view.Panel)

	view.Heel_Component = NewHeelBlendComponentView(view)
	view.Total_Component = NewTotalBlendComponentView(view)

	// view.Panel.D	ock(view.amount_field, windigo.Left)
	view.Panel.Dock(view.gross_field, windigo.TopLeft)
	view.Panel.Dock(view.total_field, windigo.TopLeft)
	view.Panel.Dock(view.heel_field, windigo.TopLeft)
	view.Panel.Dock(view.amount_field, windigo.TopLeft)
	view.Panel.Dock(view.headers, windigo.Top)
	view.Dock(view.Heel_Component, windigo.Top)
	view.Dock(view.Total_Component, windigo.Bottom)

	view.gross_field.OnChange().Bind(view.gross_field_OnChange)

	view.total_field.OnChange().Bind(view.total_field_OnChange)
	view.heel_field.OnChange().Bind(view.heel_field_OnChange)
	view.amount_field.OnChange().Bind(view.amount_field_OnChange)

	amount := 100.
	view.SetProperGrossTotal(amount)
	view.SetProperNetTotal(amount)
	view.SetHeel(0)

	// WE are strapping this anyway
	view.heel_field.Hide()

	return view
}

func (view *BlendView) SetFont(font *windigo.Font) {
	view.gross_field.SetFont(font)
	view.total_field.SetFont(font)
	view.heel_field.SetFont(font)
	view.amount_field.SetFont(font)
	view.headers.SetFont(font)
	view.Heel_Component.SetFont(font)
	view.Total_Component.SetFont(font)

	for _, component := range view.Components {
		component.SetFont(font)
	}
}

func (view *BlendView) RefreshSize() (width, height int) {
	height = 5 * GUI.PRODUCT_FIELD_HEIGHT
	delta_height := GUI.PRODUCT_FIELD_HEIGHT

	view.Panel.SetSize(GUI.TOP_PANEL_WIDTH, height)

	width, height = GUI.OFF_AXIS, height+(len(view.Components)+2)*delta_height
	view.SetSize(width, height)

	view.gross_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, delta_height)
	view.total_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, delta_height)
	view.heel_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, delta_height)
	view.amount_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, delta_height)
	view.headers.RefreshSize()
	view.Heel_Component.RefreshSize()
	view.Total_Component.RefreshSize()

	for _, component := range view.Components {
		component.RefreshSize()
	}
	return
}

func (view *BlendView) Get() *blender.ProductBlend {
	if view.Recipe == nil {
		return nil
	}
	view.Blend = blender.NewProductBlend_from_Recipe(view.Recipe)
	for _, component := range view.Components {
		blendComponent := component.Get()

		if blendComponent == nil {
			// return nil
			// TODO return once we get everything nailed down?
			continue
		}
		// blendComponent
		view.Blend.AddComponent(*blendComponent)
	}
	view.Blend.Total, view.Blend.Heel, view.Blend.Amount = view.total_field.GetInt(), view.heel_field.GetInt(), view.total_field.GetInt()
	// view.Blend.Total, view.Blend.Heel, view.Blend.Amount = view.total_field.Get(),view.heel_field.Get(),view.total_field.Get()
	return view.Blend
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

	delta_height := GUI.PRODUCT_FIELD_HEIGHT
	height := view.Panel.ClientHeight() + 2*delta_height
	width := view.ClientWidth()

	if view.Recipe == nil {
		view.SetSize(width, height)
		return
	}

	total := view.Recipe.Total_Default.Float64
	if !view.Recipe.Total_Default.Valid {
		total = view.total_field.Get()
	}

	for _, component := range view.Recipe.Components {
		component_view := BlendComponentView_from_new(view, &component)
		view.Components = append(view.Components, component_view)
		view.AutoPanel.Dock(component_view, windigo.Top)
		height += delta_height
	}
	view.SetNetTotal(total)
	view.SetSize(width, height)
}

func (view *BlendView) SetGrossTotal(gross float64) {
	view.SetProperGrossTotal(gross)
	view.SetProperNetTotal(gross - view.TareMass)
	view.SetProperAmount(gross - view.Heel_Gross)
}

func (view *BlendView) SetProperGrossTotal(gross float64) {
	view.gross_field.SetInt(gross)
}

func (view *BlendView) SetNetTotal(total float64) {
	view.SetProperNetTotal(total)
	view.SetProperGrossTotal(total + view.TareMass)
	view.SetProperAmount(total - view.Heel_Net)
}

func (view *BlendView) SetProperNetTotal(total float64) {
	view.total_field.SetInt(total)
	view.Total_Component.SetAmount(total)
}

func (view *BlendView) SetProperHeel(heel float64) {
	view.Heel_Gross = heel
	view.Heel_Net = view.Heel_Gross - view.TareMass
	view.heel_field.SetInt(view.Heel_Gross)
}

func (view *BlendView) SetHeel(heel float64) {
	view.SetProperHeel(heel)
	view.Heel_Component.SetMassDensity(view.Heel_Net)
	total := view.total_field.Get()
	view.SetProperAmount(total - view.Heel_Net)
}

func (view *BlendView) SetTare(tare float64) {
	view.TareMass = tare
	view.Heel_Net = view.Heel_Gross - view.TareMass
	view.Heel_Component.SetMassDensity(view.Heel_Net)

	// view.Heel_Gross = view.Heel_Net + view.TareMass
	// view.Heel_Component.SetMassDensity(view.Heel_Net)

	// view.SetGrossTotal(view.gross_field.Get())
	view.SetNetTotal(view.total_field.Get())
}

func (view *BlendView) SetHeelMassVolume(heel float64) {
	view.SetProperHeel(heel)
	volume := view.Heel_Component.SetAmount(heel)
	total := view.total_field.Get()
	view.SetProperAmount(total - view.Heel_Net)
	view.SetStrap(volume)
}

func (view *BlendView) SetStrap(volume float64) {
	view.parent.SetStrap(volume)
}

func (view *BlendView) GetStrap(volume float64) float64 {
	return view.parent.GetStrap(volume)
}

func (view *BlendView) SetHeelVolumeMass(heel float64) {
	view.SetProperHeel(view.Heel_Component.SetVolume(heel) + view.TareMass)
	total := view.total_field.Get()

	view.SetProperAmount(total - view.Heel_Net)
}

func (view *BlendView) SetHeelVolume(heel float64) {
	view.Heel_Component.SetVolumeDensity(heel)
	view.ReTotal()
}

func (view *BlendView) SetAmount(amount float64) {
	view.SetProperGrossTotal(amount + view.Heel_Gross)
	view.SetProperNetTotal(amount + view.Heel_Net)
	view.SetProperAmount(amount)
}

func (view *BlendView) SetProperAmount(amount float64) {
	view.amount_field.SetInt(amount)
	for _, component := range view.Components {
		component.(*BlendComponentView).SetAmount(amount)
	}
	view.ReTotal()
}

func (view *BlendView) ReTotal() {
	// heelVol could be useful here?
	total_vol := 0.
	total_vol = view.Heel_Component.SetTotal(total_vol)

	for _, component := range view.Components {
		total_vol = component.(*BlendComponentView).SetTotal(total_vol)
	}
	view.Total_Component.SetVolumeDensity(total_vol)
}

func (view *BlendView) gross_field_OnChange(*windigo.Event) {
	view.SetGrossTotal(view.gross_field.Get())
}

func (view *BlendView) total_field_OnChange(*windigo.Event) {
	view.SetNetTotal(view.total_field.Get())
}

func (view *BlendView) heel_field_OnChange(*windigo.Event) {
	view.SetHeelMassVolume(view.heel_field.Get())
}

func (view *BlendView) amount_field_OnChange(*windigo.Event) {
	view.SetAmount(view.amount_field.Get())
}
