package views

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
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
	QCBlendComponentView
	amount_field *NumbEditView
}

// TODO
func NewBlendComponentView(parent *BlendView, recipeComponent *blender.RecipeComponent) *BlendComponentView {

	// //TODO
	// view := new(BlendComponentView)
	// view.QCBlendComponentView = *NewQCBlendComponentView(parent,recipeComponent)
	//
	//
	//
	// // cf NumberEditView
	// // view.amount_field = windigo.NewEdit(view.AutoPanel)
	// view.amount_field = NewNumbEditView(view.AutoPanel) //TODO before

	DEL_BUTTON_WIDTH := 20

	view := new(BlendComponentView)
	view.recipeComponent = recipeComponent

	// view.BlendComponent = NewBlendComponent()
	// view.BlendComponent.Add_order = len(parent.Components)
	// view.BlendComponent.

	view.component_types_data = make(map[string]blender.BlendComponent)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	// view.component_field = GUI.NewSearchBoxWithLabels(view.AutoPanel, parent.component_types_list)
	// view.component_field = GUI.NewListSearchBoxWithLabels(view.AutoPanel, view.component_types_list)
	view.component_field = GUI.NewListSearchBox(view.AutoPanel)

	// cf NumberEditView
	// view.amount_field = windigo.NewEdit(view.AutoPanel)
	view.amount_field = NewNumbEditView(view.AutoPanel)

	//??TODO make memnber?
	var InboundLotView *InboundLotView

	lot_add_button := windigo.NewPushButton(view.AutoPanel)
	lot_add_button.SetText("+")
	// TODO add lot
	lot_add_button.OnClick().Bind(func(e *windigo.Event) {
		if InboundLotView != nil {
			return
		}
		InboundLotView = NewInboundLotView(view, recipeComponent)
		InboundLotView.RefreshSize()
		view.AutoPanel.Dock(InboundLotView, windigo.Left)
		InboundLotView.OnClose().Bind(func(e *windigo.Event) {
			InboundLotView = nil
			view.Update_component_types()
		})
	})

	view.controls = append(view.controls, lot_add_button)

	view.AutoPanel.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.component_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.amount_field.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	lot_add_button.SetSize(DEL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.AutoPanel.Dock(view.component_field, windigo.Left)
	view.AutoPanel.Dock(view.amount_field, windigo.Left)
	// view.AutoPanel.Dock(order_field, windigo.Left)
	view.AutoPanel.Dock(lot_add_button, windigo.Left)

	view.component_field.Label().SetText(recipeComponent.Component_name)
	// view.component_field.SetText(RecipeComponent.Component_name)
	view.amount_field.Set(recipeComponent.Component_amount)

	view.Update_component_types()

	view.amount_field.OnChange().Bind(func(e *windigo.Event) {
		parent.SetAmount(view.amount_field.Get() / view.recipeComponent.Component_amount)
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
	BlendComponent.Component_amount = view.amount_field.Get()
	log.Println("DEBUG: BlendComponentView update_component_types", BlendComponent, view.component_field.GetSelectedItem(), view.component_field.SelectedItem(), view.component_types_data[view.component_field.Text()])

	return BlendComponent
}

func (view *BlendComponentView) SetAmount(amount float64) {
	view.amount_field.Set(amount * view.recipeComponent.Component_amount)
}

/*
 * //TODO

func (view *BlendComponentView) SetFont(font *windigo.Font) {
	// for _, component := range view.Components {
	// 	component.SetFont(font)
	// }
}

// func (view *BlendComponentView) SetLabeledSize(label_width, control_width, height int) {
// 	view.SetSize(label_width+control_width, height)
// }


func (view *BlendComponentView) SetLabeledSize(label_width, control_width, height int) {
	view.SetSize(label_width+control_width, height)

	view.Label().SetSize(label_width, height)

}*/

func (view *BlendComponentView) SetFont(font *windigo.Font) {
	view.QCBlendComponentView.SetFont(font)
	view.amount_field.SetFont(font)
}

func (view *BlendComponentView) RefreshSize() {
	view.QCBlendComponentView.RefreshSize()
	view.amount_field.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
}
