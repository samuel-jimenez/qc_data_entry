package blender

import (
	"log"

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

	height := FIELD_HEIGHT
	delta_height := FIELD_HEIGHT
		height := FIELD_HEIGHT
	delta_height := FIELD_HEIGHT

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
	height := height
	delta_height := height
// 	view.SetSize(label_width+control_width, height)
// 	for _, component := range view.Components {
// 		component.SetSize(label_width+control_width, height)
// 	}
// }

// func (view *RecipeView) RefreshSize() {
//
// 	panel.SetSize(GUI.OFF_AXIS, GROUP_HEIGHT)
// 	panel.SetMargins(GROUP_MARGIN, GROUP_MARGIN, 0, 0)
//
// 	group_panel.SetSize(GROUP_WIDTH, GROUP_HEIGHT)
// 	group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, BTM_SPACER_WIDTH, BTM_SPACER_HEIGHT)
//
// 	visual_field.SetSize(GUI.OFF_AXIS, FIELD_HEIGHT)
//
// 	mass_field.SetLabeledSize(GUI.LABEL_WIDTH, DATA_FIELD_WIDTH, DATA_SUBFIELD_WIDTH, DATA_UNIT_WIDTH, FIELD_HEIGHT)
//
// 	button_dock.SetDockSize(BUTTON_WIDTH, BUTTON_HEIGHT)
//
// 	ranges_panel.SetMarginTop(GROUP_MARGIN)
// 	ranges_panel.Refresh()
//
// }
