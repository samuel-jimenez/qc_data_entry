package blender

import (
	"log"
	"strconv"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/*
* NumbEditView
* 	// cf NumberEditView

*
 */
type NumbEditView struct {
	// GUI.ErrableView
	*windigo.Edit
}

func NewNumbEditView(parent windigo.Controller) *NumbEditView {
	// edit_field := new(NumbEditView)
	// edit_field.Edit = windigo.NewEdit(parent)
	// return edit_field
	return &NumbEditView{windigo.NewEdit(parent)}

}

func (control *NumbEditView) Get() float64 {
	val, _ := strconv.ParseFloat(strings.TrimSpace(control.Text()), 64)
	return val
}

/*
 * RecipeComponentView
 *
 */
type RecipeComponentView struct {
	*windigo.AutoPanel
	RecipeComponent *RecipeComponent
	// Component_name   string
	// Component_amount float64
	// Component_id int64
	// Add_order    int64
	Get                    func() *RecipeComponent
	Update                 func(*RecipeComponent)
	Update_component_types func(component_types_list []string)
}

func NewRecipeComponentView(parent *RecipeView) *RecipeComponentView {
	DEL_BUTTON_WIDTH := 20

	view := new(RecipeComponentView)
	view.RecipeComponent = NewRecipeComponent()
	view.RecipeComponent.Add_order = len(parent.Components)
	// view.RecipeComponent.
	view.AutoPanel = windigo.NewAutoPanel(parent)

	// component_field := GUI.NewSearchBoxWithLabels(view.AutoPanel, parent.component_types_list)
	component_field := GUI.NewListSearchBoxWithLabels(view.AutoPanel, parent.component_types_list)

	// cf NumberEditView
	// amount_field := windigo.NewEdit(view.AutoPanel)
	amount_field := NewNumbEditView(view.AutoPanel)

	// order_field := windigo.NewEdit(view.AutoPanel) //todo
	// order_field := NewNumbEditView(view.AutoPanel) //todo

	component_del_button := windigo.NewPushButton(view.AutoPanel)
	component_del_button.SetText("-")

	view.AutoPanel.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	component_field.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	amount_field.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	component_del_button.SetSize(DEL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.AutoPanel.Dock(component_field, windigo.Left)
	view.AutoPanel.Dock(amount_field, windigo.Left)
	// view.AutoPanel.Dock(order_field, windigo.Left)
	view.AutoPanel.Dock(component_del_button, windigo.Left)

	view.Get = func() *RecipeComponent {
		view.RecipeComponent.Component_name = component_field.Text()
		view.RecipeComponent.Component_amount = amount_field.Get()
		// view.RecipeComponent.Add_order = order_field.Get() //TODO
		// view.RecipeComponent.Component_id = component_field.Text()
		log.Println("DEBUG: RecipeComponentView update_component_types", view.RecipeComponent)

		return view.RecipeComponent

	}

	view.Update = func(RecipeComponent *RecipeComponent) {
		log.Println("DEBUG: RecipeComponentView Update", RecipeComponent)

		if view.RecipeComponent == RecipeComponent || RecipeComponent == nil {
			log.Println("DEBUG: RecipeComponentView Update return", RecipeComponent, view.RecipeComponent)
			return
		}

		view.RecipeComponent = RecipeComponent

		component_field.SetText(RecipeComponent.Component_name)
	}

	view.Update_component_types = func(component_types_list []string) {
		text := component_field.Text()
		log.Println("DEBUG: RecipeComponentView update_component_types", text)
		component_field.Update(component_types_list)
		component_field.SetText(text)
		// c_view := windigo.NewLabel(component_panel)
		// c_view.SetText(component.Component_name)
	}

	return view
}

func (view *RecipeComponentView) SetFont(font *windigo.Font) {
	// for _, component := range view.Components {
	// 	component.SetFont(font)
	// }
}

// func (view *RecipeComponentView) SetLabeledSize(label_width, control_width, height int) {
// 	view.SetSize(label_width+control_width, height)
// }

//??TODO

func (view *RecipeComponentView) MoveUp() {
}

func (view *RecipeComponentView) MoveDown() {
}

/*
func (view *RecipeComponentView) SetLabeledSize(label_width, control_width, height int) {
	view.SetSize(label_width+control_width, height)

	view.Label().SetSize(label_width, height)

}*/
