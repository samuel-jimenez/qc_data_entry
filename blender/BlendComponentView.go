package blender

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/*
 * BlendComponentViewer
 *
 */
type BlendComponentViewer interface {
	// Get() *BlendComponent
	// Update(*BlendComponent)
	// Update_component_types(component_types_list []string)
	// SetFont(font *windigo.Font)
	// MoveUp()
	// MoveDown()
}

/*
 * BlendComponentView
 *
 */
type BlendComponentView struct {
	*windigo.AutoPanel
	BlendComponent       *BlendComponent
	component_field      *GUI.SearchBox
	amount_field         *NumbEditView
	component_types_data map[string]int64

	// Component_name   string
	// Component_amount float64
	// Component_id int64
	// Add_order    int64
}

/*
func NewBlendComponentView(parent *BlendView) *BlendComponentView {
	DEL_BUTTON_WIDTH := 20

	view := new(BlendComponentView)
	view.BlendComponent = NewBlendComponent()
	view.BlendComponent.Add_order = len(parent.Components)
	// view.BlendComponent.
	view.component_types_data = parent.component_types_data
	view.AutoPanel = windigo.NewAutoPanel(parent)

	// view.component_field = GUI.NewSearchBoxWithLabels(view.AutoPanel, parent.component_types_list)
	view.component_field = GUI.NewListSearchBoxWithLabels(view.AutoPanel, parent.component_types_list)

	// cf NumberEditView
	// view.amount_field = windigo.NewEdit(view.AutoPanel)
	view.amount_field = NewNumbEditView(view.AutoPanel)

	component_del_button := windigo.NewPushButton(view.AutoPanel)
	component_del_button.SetText("-")
	// TODO delete

	view.AutoPanel.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.component_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.amount_field.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	component_del_button.SetSize(DEL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.AutoPanel.Dock(view.component_field, windigo.Left)
	view.AutoPanel.Dock(view.amount_field, windigo.Left)
	// view.AutoPanel.Dock(order_field, windigo.Left)
	view.AutoPanel.Dock(component_del_button, windigo.Left)

	return view
}

func (view *BlendComponentView) Get() *BlendComponent {
	view.BlendComponent.Component_name = view.component_field.Text()
	view.BlendComponent.Component_type_id = view.component_types_data[view.BlendComponent.Component_name]
	if view.BlendComponent.Component_type_id == DB.INVALID_ID {
		//TODO make error
		return nil
	} //TODO check for zero amount?
	view.BlendComponent.Component_amount = view.amount_field.Get()
	log.Println("DEBUG: BlendComponentView update_component_types", view.BlendComponent, view.component_field.GetSelectedItem(), view.component_field.SelectedItem(), view.component_types_data[view.component_field.Text()])

	return view.BlendComponent
}

func (view *BlendComponentView) Update(BlendComponent *BlendComponent) {
	log.Println("DEBUG: BlendComponentView Update", BlendComponent)

	if view.BlendComponent == BlendComponent || BlendComponent == nil {
		log.Println("Warn: BlendComponentView Update return", BlendComponent, view.BlendComponent)
		return
	}

	view.BlendComponent = BlendComponent

	view.component_field.SetText(BlendComponent.Component_name)
	view.amount_field.Set(BlendComponent.Component_amount)

}

func (view *BlendComponentView) Update_component_types(component_types_list []string) {
	text := view.component_field.Text()
	log.Println("DEBUG: BlendComponentView update_component_types", text)
	view.component_field.Update(component_types_list)
	view.component_field.SetText(text)
	// c_view := windigo.NewLabel(component_panel)
	// c_view.SetText(component.Component_name)
}

func (view *BlendComponentView) SetFont(font *windigo.Font) {
	// for _, component := range view.Components {
	// 	component.SetFont(font)
	// }
}

// func (view *BlendComponentView) SetLabeledSize(label_width, control_width, height int) {
// 	view.SetSize(label_width+control_width, height)
// }

//??TODO

func (view *BlendComponentView) MoveUp() {
}

func (view *BlendComponentView) MoveDown() {
}*/

/*
func (view *BlendComponentView) SetLabeledSize(label_width, control_width, height int) {
	view.SetSize(label_width+control_width, height)

	view.Label().SetSize(label_width, height)

}*/
