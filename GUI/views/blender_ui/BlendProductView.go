package blender_ui

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/windigo"
)

type BlendProductViewer interface {
	windigo.Pane
	Get() *blender.BlendProduct
	// Update(blend *blender.BlendProduct)
	Update_products(product_data map[string]int64)
	AddComponent()
	SetFont(font *windigo.Font)
	RefreshSize()
}

type BlendProductView struct {
	*windigo.AutoPanel

	RecipeProduct *blender.RecipeProduct
	// BlendProduct                   *blender.BlendProduct
	Blend           *BlendView
	Product_data    map[string]int64
	Product_Field   *GUI.SearchBox
	Blend_sel_field *GUI.ComboBox
	panels          []*windigo.AutoPanel
	controls        []windigo.Controller
}

func NewBlendProductView(parent windigo.Controller) *BlendProductView {

	product_text := "Product"
	recipe_text := ""

	view := new(BlendProductView)
	view.RecipeProduct = blender.NewRecipeProduct()
	// view.BlendProduct = NewBlendProduct()

	view.Product_data = make(map[string]int64)

	view.AutoPanel = windigo.NewAutoPanel(parent)

	product_panel := windigo.NewAutoPanel(view.AutoPanel)
	recipe_panel := windigo.NewAutoPanel(view.AutoPanel)
	view.panels = append(view.panels, product_panel)
	view.panels = append(view.panels, recipe_panel)

	view.Product_Field = GUI.NewLabeledListSearchBox(product_panel, product_text)

	view.Blend_sel_field = GUI.NewListComboBox(recipe_panel, recipe_text)

	recipe_accept_button := windigo.NewPushButton(recipe_panel)
	recipe_accept_button.SetText("OK")
	recipe_accept_button.SetSize(GUI.ACCEPT_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.Blend = NewBlendView(view.AutoPanel)

	view.controls = append(view.controls, recipe_accept_button)

	// Dock

	view.AutoPanel.Dock(product_panel, windigo.Top)
	view.AutoPanel.Dock(recipe_panel, windigo.Top)
	view.AutoPanel.Dock(view.Blend, windigo.Top)
	product_panel.Dock(view.Product_Field, windigo.Left)
	recipe_panel.Dock(view.Blend_sel_field, windigo.Left)
	recipe_panel.Dock(recipe_accept_button, windigo.Left)

	// combobox
	GUI.Fill_combobox_from_query_rows(view.Product_Field, func(row *sql.Rows) error {
		var (
			id                   int64
			internal_name        string
			product_moniker_name string
		)
		if err := row.Scan(&id, &internal_name, &product_moniker_name); err != nil {
			return err
		}
		name := product_moniker_name + " " + internal_name
		view.Product_data[name] = id
		view.Product_Field.AddItem(name)
		return nil
	},
		DB.DB_Select_product_info_all)

	//event handling
	view.Product_Field.OnSelectedChange().Bind(func(e *windigo.Event) {
		view.RecipeProduct = blender.NewRecipeProduct()
		// view.BlendProduct = NewBlendProduct()
		name := view.Product_Field.GetSelectedItem()
		view.RecipeProduct.Set(view.Product_data[name])
		view.RecipeProduct.GetRecipes(view.Blend_sel_field)
		view.RecipeProduct.GetComponents()
		view.Blend_sel_field.OnSelectedChange().Fire(nil)
	})

	recipe_accept_button.OnClick().Bind(func(e *windigo.Event) {
		log.Println("ERR: DEBUG: view Get:", view.Get())

	})

	view.Blend_sel_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		i, err := strconv.Atoi(view.Blend_sel_field.GetSelectedItem())
		if err != nil {
			log.Println("ERR: recipe_field strconv", err)
			view.Blend.UpdateRecipe(nil)
			return
		}
		view.Blend.UpdateRecipe(view.RecipeProduct.Recipes[i])

		// view.Product = NewBlendProduct()
		// name := product_field.GetSelectedItem()
		// view.Product.Set(name, product_data[name])
		// view.Product.LoadBlendCombo(recipe_field)
	})

	return view
}

func (view *BlendProductView) Get() *blender.BlendProduct {

	ProductBlend := view.Blend.Get()
	log.Println("ERR: DEBUG: view.Blend Get:", ProductBlend)
	if ProductBlend == nil {
		return nil
	}
	BlendProduct := blender.NewBlendProductFromRecipe(view.RecipeProduct)

	if BlendProduct.Product_id == DB.INVALID_ID {
		//TODO make error
		return nil
	}

	BlendProduct.Recipe = view.Blend.Recipe
	BlendProduct.Blend = ProductBlend
	BlendProduct.Save()
	return BlendProduct

}

/*
func (view *BlendProductView) Update(recipe *BlendProduct) {
	if view.Product == recipe {
		return
	}

	// height := FIELD_HEIGHT
	// delta_height := FIELD_HEIGHT
	// 	total_height := height
	// delta_height := height
	// 	height := FIELD_HEIGHT
	// delta_height := FIELD_HEIGHT

	for _, component := range view.Blend {
		component.Close()
	}
	view.Blend = nil
	view.Product = recipe
	log.Println("BlendProductView Update", view.Product)
	if view.Product == nil {
		return
	}

	for _, component := range view.Product.Components {
		log.Println("DEBUG: BlendProductView update_components", component)
		// view.AddComponent()
		//TODO
		component_view := view.AddComponent()
		component_view.Update(&component)
	}

}*/

func (view *BlendProductView) Update_products(product_data map[string]int64) {
	view.Product_data = product_data
}

func (view *BlendProductView) SetFont(font *windigo.Font) {
	view.Blend.SetFont(font)
	view.Product_Field.SetFont(font)
	view.Blend_sel_field.SetFont(font)

	for _, control := range view.controls {
		control.SetFont(font)
	}

}

func (view *BlendProductView) RefreshSize() {

	for _, panel := range view.panels {
		panel.SetSize(GUI.TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	}

	view.Product_Field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Blend_sel_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Blend.RefreshSize()

}
