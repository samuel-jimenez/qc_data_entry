package views

import (
	"log"
	"strconv"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/windigo"
)

var (
	ADD_BUTTON_WIDTH = 20
	NEW_BUTTON_WIDTH = 50
)

type RecipeProductViewer interface {
	windigo.Pane
	// Get() *RecipeProduct
	Update(recipe *blender.RecipeProduct)
	Update_component_types(component_types_list []string, component_types_data map[string]int64)
	Update_products(product_list []string, product_data map[string]int64)
	AddComponent()
	SetFont(font *windigo.Font)
	RefreshSize()
}

type RecipeProductView struct {
	*windigo.AutoPanel
	Product                         *blender.RecipeProduct
	Recipe                          *RecipeView
	Product_data                    map[string]int64
	Product_Field, Recipe_sel_field *GUI.ComboBox
	panels                          []*windigo.AutoPanel
	controls                        []windigo.Controller
}

func NewRecipeProductView(parent windigo.Controller) *RecipeProductView {

	product_text := "Product"
	recipe_text := ""

	view := new(RecipeProductView)
	view.Product = blender.RecipeProduct_from_new()

	view.AutoPanel = windigo.NewAutoPanel(parent)

	product_panel := windigo.NewAutoPanel(view.AutoPanel)
	recipe_panel := windigo.NewAutoPanel(view.AutoPanel)
	view.panels = append(view.panels, product_panel)
	view.panels = append(view.panels, recipe_panel)

	view.Product_Field = GUI.ComboBox_from_new(product_panel, product_text)
	product_add_button := windigo.NewPushButton(product_panel)
	product_add_button.SetText("+")
	product_add_button.SetSize(ADD_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.Recipe_sel_field = GUI.ComboBox_from_new(recipe_panel, recipe_text)
	recipe_add_button := windigo.NewPushButton(recipe_panel)
	recipe_add_button.SetText("+")
	recipe_add_button.SetSize(ADD_BUTTON_WIDTH, GUI.OFF_AXIS)

	recipe_accept_button := windigo.NewPushButton(recipe_panel)
	recipe_accept_button.SetText("OK")
	recipe_accept_button.SetSize(GUI.ACCEPT_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.Recipe = NewRecipeView(view)

	view.controls = append(view.controls, product_add_button)
	view.controls = append(view.controls, recipe_add_button)
	view.controls = append(view.controls, recipe_accept_button)

	// Dock

	view.AutoPanel.Dock(product_panel, windigo.Top)
	view.AutoPanel.Dock(recipe_panel, windigo.Top)
	view.AutoPanel.Dock(view.Recipe, windigo.Top)
	product_panel.Dock(view.Product_Field, windigo.Left)
	product_panel.Dock(product_add_button, windigo.Left)
	recipe_panel.Dock(view.Recipe_sel_field, windigo.Left)
	recipe_panel.Dock(recipe_add_button, windigo.Left)
	recipe_panel.Dock(recipe_accept_button, windigo.Left)

	//event handling
	view.Product_Field.OnSelectedChange().Bind(func(e *windigo.Event) {
		view.Product = blender.RecipeProduct_from_new()
		name := view.Product_Field.GetSelectedItem()
		view.Product.Set(view.Product_data[name])
		// view.Product.GetRecipes()
		view.Product.GetRecipes(view.Recipe_sel_field)
		view.Product.GetComponents()
		view.Recipe_sel_field.OnSelectedChange().Fire(nil)
	})

	product_add_button.OnClick().Bind(func(e *windigo.Event) {
		log.Println("TODO product_add_button")

		// clear()
		// clear_product()
	})
	recipe_add_button.OnClick().Bind(func(e *windigo.Event) {
		if view.Product.Product_id != DB.INVALID_ID {
			recipe_procedure_id := 1 // TODO

			numRecipes := len(view.Product.Recipes)
			recipe_name := view.Product_Field.Text() + ":" + strconv.Itoa(numRecipes)

			view.Recipe_sel_field.AddItem(recipe_name)
			view.Recipe_sel_field.SetSelectedItem(numRecipes)
			view.Recipe.Update(view.Product.NewRecipe(recipe_procedure_id, recipe_name))

			// log.Println("product_field", product_field.SelectedItem())

		}
	})
	recipe_accept_button.OnClick().Bind(func(e *windigo.Event) {

		log.Println("ERR: DEBUG: view.Recipe Get:", view.Recipe.Get())

	})

	view.Recipe_sel_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		view.Recipe.Update(view.Product.Recipes[view.Recipe_sel_field.GetSelectedItem()])
	})

	//TODO normalize amounts

	return view
}

func (view *RecipeProductView) Update_products(product_list []string, product_data map[string]int64) {
	view.Product_data = product_data
	view.Recipe.Update_products(product_list, product_data)

}

func (view *RecipeProductView) SetFont(font *windigo.Font) {
	view.Recipe.SetFont(font)
	view.Product_Field.SetFont(font)
	view.Recipe_sel_field.SetFont(font)

	for _, control := range view.controls {
		control.SetFont(font)
	}

}

// //TODO grow
// func (view *RecipeProductView) SetLabeledSize(label_width, control_width, height int) {
// total_height := height
// delta_height := height
// 	view.SetSize(label_width+control_width, height)
// 	for _, component := range view.Components {
// 		component.SetSize(label_width+control_width, height)
// 	}
// }

func (view *RecipeProductView) RefreshSize() {

	for _, panel := range view.panels {
		panel.SetSize(GUI.TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	}

	view.Product_Field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Recipe_sel_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Recipe.RefreshSize()
}
