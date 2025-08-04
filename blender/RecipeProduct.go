package blender

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
)

type RecipeProduct struct {
	Product_name string
	Product_id   int64
	Recipes      []*ProductRecipe
}

func NewRecipeProduct() *RecipeProduct {
	return new(RecipeProduct)
}

func (object *RecipeProduct) Set(name string, Product_id int64) {
	object.Product_name = name
	object.Product_id = Product_id
}

func (object *RecipeProduct) GetRecipes(combo_field *GUI.ComboBox) {
	proc_name := "RecipeProduct.GetRecipes"
	i := 0

	DB.Forall(proc_name,
		func() {
			object.Recipes = nil
			combo_field.DeleteAllItems()
		},
		func(row *sql.Rows) {
			var (
				recipe_data ProductRecipe
			)

			if err := row.Scan(
				&recipe_data.Recipe_id,
			); err != nil {
				log.Fatal("Crit: [RecipeProduct GetRecipes]: ", proc_name, err)
			}
			log.Println("DEBUG: GetRecipes qc_data", proc_name, recipe_data)
			object.Recipes = append(object.Recipes, &recipe_data)
			combo_field.AddItem(strconv.Itoa(i))
			i++

		},
		DB.DB_Select_product_recipe, object.Product_id)
	if i != 0 {
		combo_field.SetSelectedItem(0)
	}
}

func (object *RecipeProduct) NewRecipe() *ProductRecipe {
	proc_name := "RecipeProduct.NewRecipe"

	Recipe_id := DB.Insert(proc_name, DB.DB_Insert_product_recipe, object.Product_id)
	if Recipe_id == DB.INVALID_ID {
		return nil
	}

	recipe_data := NewProductRecipe(Recipe_id)
	object.Recipes = append(object.Recipes, recipe_data)
	return recipe_data
}

func (object *RecipeProduct) GetComponents() {
	for _, recipe_data := range object.Recipes {
		recipe_data.GetComponents()
	}
}
