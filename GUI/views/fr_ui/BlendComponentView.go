package fr_ui

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/windigo"
)

/*
 * BlendComponentViewer
 *
 */
type BlendComponentViewer interface {
	windigo.Controller
	Get() *blender.BlendComponent
	Update_component_types()
	SetAmount(amount float64)
	SetFont(font *windigo.Font)
	RefreshSize()
}

/*
 * BlendComponentView
 *
 */
type BlendComponentView struct {
	views.QCBlendComponentView
	// InboundLotView *views.InboundLotView
	// SG_field *views.NumbEditView
	Density_field *views.NumbEditView
	Gallons_field *views.NumbEditView
	// parent        *BlendView
}

func NewBaseBlendComponentView(parent *BlendView, recipeComponent *blender.RecipeComponent) *BlendComponentView {

	view := new(BlendComponentView)
	view.QCBlendComponentView = *views.New_Bare_QCBlendComponentView_from_RecipeComponent_com(parent, recipeComponent)

	// view.SG_field = views.NewNumbEditView(view.AutoPanel)
	// view.AutoPanel.Dock(view.SG_field, windigo.Left)

	view.Density_field = views.NewNumbEditView(view.AutoPanel)
	view.AutoPanel.Dock(view.Density_field, windigo.Left)

	view.Gallons_field = views.NewNumbEditView(view.AutoPanel)
	view.AutoPanel.Dock(view.Gallons_field, windigo.Left)

	// RefreshSize
	view.RefreshSize()

	// DEL_BUTTON_WIDTH := 20
	// lot_add_button := windigo.NewPushButton(view.AutoPanel)
	// lot_add_button.SetText("+")
	// // TODO add lot
	// lot_add_button.OnClick().Bind(func(e *windigo.Event) {
	// 	if view.InboundLotView != nil {
	// 		return
	// 	}
	// 	view.InboundLotView = views.NewInboundLotView(view, recipeComponent)
	// 	view.InboundLotView.RefreshSize()
	// 	view.AutoPanel.Dock(view.InboundLotView, windigo.Left)
	// 	view.InboundLotView.OnClose().Bind(func(e *windigo.Event) {
	// 		view.InboundLotView = nil
	// 		view.Update_component_types()
	// 	})
	// })
	// view.Controls = append(view.Controls, lot_add_button)
	// lot_add_button.SetSize(DEL_BUTTON_WIDTH, GUI.OFF_AXIS)
	// view.AutoPanel.Dock(lot_add_button, windigo.Left)

	view.Density_field.Set(9.) // TODO find a good default

	return view

}

func NewDummyBlendComponentView(parent *BlendView, name string) *BlendComponentView {

	Recipe_Component := blender.NewRecipeComponent()
	Recipe_Component.Component_name = name
	Recipe_Component.Component_amount = 1

	view := NewBaseBlendComponentView(parent, Recipe_Component)

	view.Amount_required_field.SetEnabled(false)
	view.Component_field.SetEnabled(false)
	view.Gallons_field.SetEnabled(false)

	return view

}

func NewHeelBlendComponentView(parent *BlendView) *BlendComponentView {

	view := NewDummyBlendComponentView(parent, "HEEL")

	view.Density_field.OnChange().Bind(func(e *windigo.Event) {
		volume := view.Gallons_field.Get()
		parent.SetHeelVolume(volume)
	})
	return view

}

func NewTotalBlendComponentView(parent *BlendView) *BlendComponentView {

	view := NewDummyBlendComponentView(parent, "Total")

	view.Density_field.SetEnabled(false)
	return view

}

func NewBlendComponentView(parent *BlendView, recipeComponent *blender.RecipeComponent) *BlendComponentView {

	view := NewBaseBlendComponentView(parent, recipeComponent)

	view.Amount_required_field.OnChange().Bind(func(e *windigo.Event) {
		parent.SetAmount(view.Amount_required_field.Get() / view.Recipe_Component.Component_amount)
	})
	// view.Component_field.Hide()

	// grab SG
	view.Component_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		proc_name := "BlendComponentView_Component_field"
		lot := view.Component_types_data[view.Component_field.GetSelectedItem()]
		var SG float64
		if err := DB.Select_ErrNoRows(
			proc_name,
			DB.DB_Select_component_type_density.QueryRow(
				lot.Lot_id,
				lot.Inboundp,
			),
			&SG,
		); err != nil {
			return
			// TODO reset to default?
		}
		log.Println("DEBUG:", proc_name, "SG", SG)
		// view.SG_field.Set(SG)
		view.Density_field.Set(formats.Density_from_sg(SG))
		view.Density_field.OnChange().Fire(nil)
	})
	view.Density_field.OnChange().Bind(func(e *windigo.Event) {
		view.Gallons_field.SetInt(view.Amount_required_field.Get() / view.Density_field.Get())
	})
	view.Gallons_field.OnChange().Bind(func(e *windigo.Event) {
		// view.Density_field.Set(view.Amount_required_field.Get() / view.Gallons_field.Get())
		view.Amount_required_field.SetInt(view.Gallons_field.Get() * view.Density_field.Get())
		view.Amount_required_field.OnChange().Fire(nil)
	})

	return view
}

func (view *BlendComponentView) Get() *blender.BlendComponent {

	BlendComponent := view.QCBlendComponentView.Get()

	if BlendComponent == nil {
		//TODO make error
		return nil
	}
	//TODO check for zero amount?
	BlendComponent.Component_amount = view.Amount_required_field.Get()
	log.Println("DEBUG: BlendComponentView update_component_types", BlendComponent, view.Component_field.GetSelectedItem(), view.Component_field.SelectedItem(), view.Component_types_data[view.Component_field.Text()])

	return BlendComponent
}

func (view *BlendComponentView) SetAmount(amount float64) float64 {
	Amount := amount * view.Recipe_Component.Component_amount
	Gallons := Amount / view.Density_field.Get()
	view.Amount_required_field.SetInt(Amount)
	view.Gallons_field.SetInt(Gallons)
	return Gallons
}

func (view *BlendComponentView) SetVolume(volume float64) float64 {
	mass := volume * view.Density_field.Get()
	view.Amount_required_field.SetInt(mass)
	view.Gallons_field.SetInt(volume)
	return mass
}

func (view *BlendComponentView) SetFont(font *windigo.Font) {
	view.QCBlendComponentView.SetFont(font)
	view.Density_field.SetFont(font)
	view.Gallons_field.SetFont(font)
}

func (view *BlendComponentView) RefreshSize() {
	view.QCBlendComponentView.RefreshSize()
	// view.SG_field.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Density_field.SetSize(GUI.DATA_SUBFIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Gallons_field.SetSize(GUI.DATA_SUBFIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
}
