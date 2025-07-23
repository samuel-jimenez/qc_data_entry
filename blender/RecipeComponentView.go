package blender

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

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
	// view.RecipeComponent.
	view.AutoPanel = windigo.NewAutoPanel(parent)

	component_field := GUI.NewSearchBoxWithLabels(view.AutoPanel, parent.component_types_list)

	component_del_button := windigo.NewPushButton(view.AutoPanel)
	component_del_button.SetText("+")
	component_del_button.SetSize(DEL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.AutoPanel.Dock(component_field, windigo.Left)
	view.AutoPanel.Dock(component_del_button, windigo.Left)

	view.Get = func() *RecipeComponent {
		view.RecipeComponent.Component_name = component_field.Text()
		return view.RecipeComponent

	}

	view.Update = func(RecipeComponent *RecipeComponent) {
		if view.RecipeComponent == RecipeComponent {
			return
		}

		view.RecipeComponent = RecipeComponent
		if view.RecipeComponent == nil {
			return
		}

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
