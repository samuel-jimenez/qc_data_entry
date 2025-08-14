package views

import (
	"database/sql"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/windigo"
)

type RecipeViewer interface {
	windigo.Pane
	Get() *blender.ProductRecipe
	Update(recipe *blender.ProductRecipe)
	Update_products(product_list []string, product_data map[string]int64)
	Update_component_types(component_types_list []string, component_types_data map[string]int64)
	AddComponent()
	SetFont(font *windigo.Font)
	RefreshSize()
}
type RecipeView struct {
	*windigo.AutoPanel
	Recipe     *blender.ProductRecipe
	panel      *windigo.AutoPanel
	Components []*RecipeComponentView
	product_list,
	component_types_list []string
	Product_data,
	component_types_data map[string]int64

	controls            []windigo.Controller
	component_add_panel *ComponentTypeView
}

func NewRecipeView(parent *RecipeProductView) *RecipeView {

	view := new(RecipeView)
	view.AutoPanel = windigo.NewAutoPanel(parent)
	view.component_types_data = make(map[string]int64)
	view.Product_data = parent.Product_data

	view.panel = windigo.NewAutoPanel(view.AutoPanel)

	component_add_button := windigo.NewPushButton(view.panel)
	component_add_button.SetText("+")
	component_add_button.SetSize(ADD_BUTTON_WIDTH, GUI.OFF_AXIS)

	component_new_button := windigo.NewPushButton(view.panel)
	component_new_button.SetText("New")
	component_new_button.SetSize(NEW_BUTTON_WIDTH, GUI.OFF_AXIS)

	//TODO move to view
	view.component_add_panel = NewComponentTypeView(view)
	view.component_add_panel.Hide()

	view.controls = append(view.controls, component_add_button)
	view.controls = append(view.controls, component_new_button)
	view.controls = append(view.controls, view.component_add_panel)

	// TODO/ RecipeHeader from  RecipeComponent Component Amount

	// Dock

	view.panel.Dock(component_add_button, windigo.Left)
	view.panel.Dock(component_new_button, windigo.Left)
	view.AutoPanel.Dock(view.panel, windigo.Top)

	view.AutoPanel.Dock(view.component_add_panel, windigo.Top)

	//event handling
	component_add_button.OnClick().Bind(func(e *windigo.Event) {
		if view.Recipe != nil {
			view.AddComponent()
		}
	})

	component_new_button.OnClick().Bind(func(e *windigo.Event) {
		view.component_add_panel.Show()
	})

	view.Update_component_types()

	return view
}

func (view *RecipeView) Get() *blender.ProductRecipe {
	if view.Recipe == nil {
		return nil
	}
	//TODO
	view.Recipe.Components = nil
	total := 0.0
	for i, component := range view.Components {
		recipeComponent := component.Get()
		log.Println("DEBUG: RecipeView Get", recipeComponent)
		if recipeComponent != nil {
			recipeComponent.Add_order = i
			log.Println("DEBUG: RecipeView Get", recipeComponent)
			view.Recipe.AddComponent(*recipeComponent)
			total += recipeComponent.Component_amount
		}
	}
	log.Println("DEBUG: RecipeView Get total", total)

	view.Recipe.NormalizeComponents(total)
	view.Recipe.SaveComponents()
	return view.Recipe

}

func (view *RecipeView) Update(recipe *blender.ProductRecipe) {
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
		component_view := view.AddComponent()
		component_view.Update(&component)
	}

}

func (view *RecipeView) Update_products(product_list []string, product_data map[string]int64) {
	view.Product_data = product_data
	view.component_add_panel.Update_products(product_list)
}
func (view *RecipeView) Update_component_types() {
	proc_name := "RecipeView.Update_component_types"
	DB.Forall_err(proc_name,
		func() { view.component_types_list = nil },
		func(row *sql.Rows) error {
			var (
				id   int64
				name string
			)

			if err := row.Scan(
				&id, &name,
			); err != nil {
				return err
			}

			view.component_types_data[name] = id
			log.Println("DEBUG: update_component_types nsme", id, name)

			view.component_types_list = append(view.component_types_list, name)

			return nil
		},
		DB.DB_Select_all_component_types)
	log.Println("DEBUG: update_component_types", view.component_types_list)

	if view.Recipe == nil {
		return
	}
	for _, component := range view.Components {
		log.Println("DEBUG: RecipeView update_component_types", component)
		component.Update_component_types(view.component_types_list)
	}
}

func (view *RecipeView) AddComponent() *RecipeComponentView {
	if view.Recipe == nil {
		return nil
	}

	height := view.ClientHeight()

	component_view := NewRecipeComponentView(view)
	view.Components = append(view.Components, component_view)
	view.AutoPanel.Dock(component_view, windigo.Top)

	view.SetSize(view.ClientWidth(), height+component_view.ClientHeight())

	return component_view
}

func (view *RecipeView) SetFont(font *windigo.Font) {
	for _, control := range view.controls {
		control.SetFont(font)
	}
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
	height := 2 * GUI.PRODUCT_FIELD_HEIGHT
	delta_height := GUI.PRODUCT_FIELD_HEIGHT
	width := view.ClientWidth()

	log.Println("Crit: DEBUG: RecipeView RefreshSize", height, GUI.PRODUCT_FIELD_HEIGHT, len(view.Components))
	view.SetSize(GUI.OFF_AXIS, height+len(view.Components)*delta_height)
	view.panel.SetSize(GUI.TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.component_add_panel.RefreshSize()

	for _, component := range view.Components {

		component.SetSize(width, delta_height)
		// component.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, delta_height)
	}

}
