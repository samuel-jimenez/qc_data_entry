package blender

import (
	"log"

	"github.com/samuel-jimenez/windigo"
)

type RecipeViewer interface {
	windigo.Pane
	Get() *ProductRecipe
	Update(recipe *ProductRecipe)
	Update_component_types(component_types_list []string)
	AddComponent()
	SetFont(font *windigo.Font)
}

type RecipeView struct {
	*windigo.AutoPanel
	Recipe               *ProductRecipe
	Components           []*RecipeComponentView
	component_types_list []string
	// Product_id int64
	// Recipe_id  int64
}

func (view *RecipeView) Get() *ProductRecipe {
	if view.Recipe == nil {
		return nil
	}
	//TODO
	// view.Recipe.Components = nil
	// for _, component := range view.Components {
	// 	view.Recipe.AddComponent(component.Get())
	// }
	return view.Recipe

}

func (view *RecipeView) Update(recipe *ProductRecipe) {
	if view.Recipe == recipe {
		return
	}

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
		component_view := NewRecipeComponentView(view)
		view.Components = append(view.Components, component_view)
		view.AutoPanel.Dock(component_view, windigo.Top)
		//TODO
		// component_view.Update(component)
	}

}

func (view *RecipeView) Update_component_types(component_types_list []string) {
	view.component_types_list = component_types_list
	if view.Recipe == nil {
		return
	}

	for _, component := range view.Components {
		log.Println("DEBUG: RecipeView update_component_types", component)
		component.Update_component_types(component_types_list)
	}
}

func (view *RecipeView) AddComponent() {
	if view.Recipe != nil {
		// object.Recipe.AddComponent
		//TODO
		// (object.Recipe_id,)
		component_view := NewRecipeComponentView(view)
		view.Components = append(view.Components, component_view)
		view.AutoPanel.Dock(component_view, windigo.Top)
	}
}

func (view *RecipeView) SetFont(font *windigo.Font) {
	for _, component := range view.Components {
		component.SetFont(font)
	}
}

func NewRecipeView(parent windigo.Controller) *RecipeView {
	view := new(RecipeView)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	//TODO normalize amounts

	return view
}
