package blender

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

type RecipeViewer interface {
	windigo.Pane
	Get() *ProductRecipe
	Update(recipe *ProductRecipe)
	Update_component_types(component_types_list []string, component_types_data map[string]int)
	AddComponent()
	SetFont(font *windigo.Font)
}

type RecipeView struct {
	*windigo.AutoPanel
	Recipe               *ProductRecipe
	Components           []*RecipeComponentView
	component_types_list []string
	component_types_data map[string]int
	// Product_id int64
	// Recipe_id  int64
}

func NewRecipeView(parent windigo.Controller) *RecipeView {
	view := new(RecipeView)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	//TODO normalize amounts

	return view
}

func (view *RecipeView) Get() *ProductRecipe {
	if view.Recipe == nil {
		return nil
	}
	//TODO
	view.Recipe.Components = nil
	for _, component := range view.Components {
		log.Println("DEBUG: RecipeView Get", component.Get())

		// view.Recipe.AddComponent(component.Get())
	}
	return view.Recipe

}

func (view *RecipeView) Update(recipe *ProductRecipe) {
	if view.Recipe == recipe {
		return
	}

	// height := FIELD_HEIGHT
	// delta_height := FIELD_HEIGHT
	// 	total_height := height
	// delta_height := height
	// 	height := FIELD_HEIGHT
	// delta_height := FIELD_HEIGHT

	for _, component := range view.Components {
		component.Close()
	}
	view.Components = nil
	view.Recipe = recipe
	log.Println("RecipeView Update", view.Recipe)
	if view.Recipe == nil {
		return
	}

	for _, component := range view.Recipe.Components {
		log.Println("DEBUG: RecipeView update_components", component)
		// view.AddComponent()
		//TODO
		component_view := view.AddComponent()
		component_view.Update(component)
	}

}

func (view *RecipeView) Update_component_types(component_types_list []string, component_types_data map[string]int) {
	view.component_types_list = component_types_list
	view.component_types_data = component_types_data
	if view.Recipe == nil {
		return
	}

	for _, component := range view.Components {
		log.Println("DEBUG: RecipeView update_component_types", component)
		component.Update_component_types(component_types_list)
	}
}

func (view *RecipeView) AddComponent() *RecipeComponentView {
	if view.Recipe == nil {
		return nil
	}
	// object.Recipe.AddComponent
	//TODO
	// (object.Recipe_id,)
	component_view := NewRecipeComponentView(view)
	view.Components = append(view.Components, component_view)
	view.AutoPanel.Dock(component_view, windigo.Top)
	return component_view
}

func (view *RecipeView) SetFont(font *windigo.Font) {
	for _, component := range view.Components {
		component.SetFont(font)
	}
}

// //TODO grow
// func (view *RecipeView) SetLabeledSize(label_width, control_width, height int) {
// total_height := height
// delta_height := height
// 	view.SetSize(label_width+control_width, height)
// 	for _, component := range view.Components {
// 		component.SetSize(label_width+control_width, height)
// 	}
// }

func (view *RecipeView) RefreshSize() {
	height := GUI.PRODUCT_FIELD_HEIGHT
	delta_height := GUI.PRODUCT_FIELD_HEIGHT
	width := view.ClientWidth()

	view.SetSize(GUI.OFF_AXIS, height+len(view.Components)*delta_height)
	for _, component := range view.Components {

		component.SetSize(width, delta_height)
		// component.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, delta_height)
	}

}
