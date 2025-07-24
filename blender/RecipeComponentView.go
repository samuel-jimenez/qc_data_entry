package blender

import (
	"log"
	"strconv"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/*
 * NumbEditViewer
 *
 */
type NumbEditViewer interface {
	Get() float64
	Set(float64)
}

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
	edit_field := new(NumbEditView)
	edit_field.Edit = windigo.NewEdit(parent)
	return edit_field
}

func (control *NumbEditView) Get() float64 {
	val, _ := strconv.ParseFloat(strings.TrimSpace(control.Text()), 64)
	return val
}

func (control *NumbEditView) Set(val float64) {
	control.SetText(strconv.FormatFloat(val, 'f', 2, 64))
}

/*
 * RecipeComponentViewer
 *
 */
type RecipeComponentViewer interface {
	Get() *RecipeComponent
	Update(*RecipeComponent)
	Update_component_types(component_types_list []string)
	SetFont(font *windigo.Font)
	MoveUp()
	MoveDown()
}

/*
 * RecipeComponentView
 *
 */
type RecipeComponentView struct {
	*windigo.AutoPanel
	RecipeComponent      *RecipeComponent
	component_field      *GUI.SearchBox
	amount_field         *NumbEditView
	component_types_data map[string]int64

	// Component_name   string
	// Component_amount float64
	// Component_id int64
	// Add_order    int64
}

func NewRecipeComponentView(parent *RecipeView) *RecipeComponentView {
	DEL_BUTTON_WIDTH := 20

	view := new(RecipeComponentView)
	view.RecipeComponent = NewRecipeComponent()
	view.RecipeComponent.Add_order = len(parent.Components)
	// view.RecipeComponent.
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
	view.component_field.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.amount_field.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	component_del_button.SetSize(DEL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.AutoPanel.Dock(view.component_field, windigo.Left)
	view.AutoPanel.Dock(view.amount_field, windigo.Left)
	// view.AutoPanel.Dock(order_field, windigo.Left)
	view.AutoPanel.Dock(component_del_button, windigo.Left)

	return view
}

func (view *RecipeComponentView) Get() *RecipeComponent {
	view.RecipeComponent.Component_name = view.component_field.Text()
	view.RecipeComponent.Component_type_id = view.component_types_data[view.RecipeComponent.Component_name]
	if view.RecipeComponent.Component_type_id == DB.INVALID_ID {
		//TODO make error
		return nil
	} //TODO check for zero amount?
	view.RecipeComponent.Component_amount = view.amount_field.Get()
	log.Println("DEBUG: RecipeComponentView update_component_types", view.RecipeComponent, view.component_field.GetSelectedItem(), view.component_field.SelectedItem(), view.component_types_data[view.component_field.Text()])

	return view.RecipeComponent
}

func (view *RecipeComponentView) Update(RecipeComponent *RecipeComponent) {
	log.Println("DEBUG: RecipeComponentView Update", RecipeComponent)

	if view.RecipeComponent == RecipeComponent || RecipeComponent == nil {
		log.Println("Warn: RecipeComponentView Update return", RecipeComponent, view.RecipeComponent)
		return
	}

	view.RecipeComponent = RecipeComponent

	view.component_field.SetText(RecipeComponent.Component_name)
	view.amount_field.Set(RecipeComponent.Component_amount)

}

func (view *RecipeComponentView) Update_component_types(component_types_list []string) {
	text := view.component_field.Text()
	log.Println("DEBUG: RecipeComponentView update_component_types", text)
	view.component_field.Update(component_types_list)
	view.component_field.SetText(text)
	// c_view := windigo.NewLabel(component_panel)
	// c_view.SetText(component.Component_name)
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
