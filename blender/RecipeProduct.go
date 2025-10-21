package blender

import (
	"database/sql"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
)

type RecipeProduct struct {
	Product_id int64
	Recipes    map[string]*ProductRecipe
}

func RecipeProduct_from_new() *RecipeProduct {
	object := new(RecipeProduct)
	object.Recipes = make(map[string]*ProductRecipe)
	return object
}

func RecipeProduct_from_Product_id(Product_id int64) *RecipeProduct {
	object := new(RecipeProduct)
	object.Product_id = Product_id
	object.Recipes = make(map[string]*ProductRecipe)
	return object
}

func (object *RecipeProduct) Set(Product_id int64) {
	object.Product_id = Product_id
}

func (object *RecipeProduct) GetRecipes(combo_field *GUI.ComboBox) {

	proc_name := "RecipeProduct.GetRecipes"
	i := 0
	DB.Forall_exit(proc_name,
		func() {
			clear(object.Recipes)
			combo_field.DeleteAllItems()
		},
		func(row *sql.Rows) error {
			recipe_data, err := ProductRecipe_from_SQL(row)
			if err != nil {
				return err
			}
			object.Recipes[recipe_data.Recipe_name] = recipe_data
			combo_field.AddItem(recipe_data.Recipe_name)
			i++
			return nil
		},
		DB.DB_Select_product_recipe_defaults, object.Product_id)
	if i != 0 {
		combo_field.SetSelectedItem(0)
	}
}

func (object *RecipeProduct) NewRecipe(recipe_procedure_id int, recipe_name string) *ProductRecipe {
	proc_name := "RecipeProduct.NewRecipe"

	Recipe_id := DB.Insert(proc_name, DB.DB_Insert_product_recipe, object.Product_id, recipe_procedure_id, recipe_name)
	if Recipe_id == DB.INVALID_ID {
		return nil
	}

	recipe_data := ProductRecipe_from_Recipe_id_name(Recipe_id, recipe_name)
	object.Recipes[recipe_data.Recipe_name] = recipe_data
	return recipe_data
}

func (object *RecipeProduct) GetComponents() {
	for _, recipe_data := range object.Recipes {
		recipe_data.GetComponents()
	}
}

// /?TODO
func (object *RecipeProduct) GetDefaults() {
	proc_name := "RecipeProduct.GetDefaults"
	DB.Forall_exit(proc_name,
		func() {
			object.Recipes = nil
		},
		func(row *sql.Rows) error {
			var (
				amount_total_default, //  nullable.NullInt64
				specific_gravity_default nullable.NullFloat64
			)
			if err := row.Scan(
				&amount_total_default,
				&specific_gravity_default,
				// &recipe_data.Recipe_id,
			); err != nil {
				return err
			}
			// object.Recipes = append(object.Recipes, &recipe_data)
			// combo_field.AddItem(strconv.Itoa(i))
			return nil
		},
		DB.DB_Select_product_defaults, object.Product_id)

}
