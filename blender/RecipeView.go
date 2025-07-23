package blender

import (
	"log"

	"github.com/samuel-jimenez/windigo"
)

type RecipeViewer interface {
	Get() *ProductRecipe
	Update(recipe *ProductRecipe)
	Update_component_types(component_types_list []string)
	AddComponent()
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
	// object.Recipe.AddComponent

	// Recipe =SQLFilterDiscrete{key,
	// 	selection_options.Get()}
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
		component_data := NewRecipeComponentView(view)
		view.Components = append(view.Components, component_data)
	}
}

func NewRecipeView(parent windigo.Controller) *RecipeView {
	view := new(RecipeView)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	return view
}
