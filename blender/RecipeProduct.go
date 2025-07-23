package blender

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
)

type RecipeProduct struct {
	Product_name string `json:"product_name"`
	// Lot_number               string `json:"lot_number"`
	// Sample_point             string
	// Visual                   bool
	Product_id int64
	// Lot_id     int64
	Recipes []*ProductRecipe
	// Product_name_customer_id nullable.NullInt64
	// Product_name_customer    string `json:"customer_product_name"`
}

func (object *RecipeProduct) Set(name string, i int64) {
	object.Product_name = name
	object.Product_id = i
}

// func (r *RecipeProduct) GetRecipes() (rows *sql.Rows,error){
// 	return DB.DB_Select_product_recipe.Query(r.Product_id)
// 	// r.Recipes =
//
// 	// rows, err :=
// 	//TODO
// }

func (object *RecipeProduct) GetRecipes() {
	DB.Forall("GetRecipes",
		func() { object.Recipes = nil },
		func(rows *sql.Rows) {
			var (
				recipe_data ProductRecipe
			)

			if err := rows.Scan(&recipe_data.Recipe_id); err != nil {
				log.Fatal(err)
			}
			recipe_data.Product_id = object.Product_id
			log.Println("DEBUG: GetRecipes qc_data", recipe_data)
			object.Recipes = append(object.Recipes, &recipe_data)

		},
		DB.DB_Select_product_recipe, object.Product_id)
}

func (object *RecipeProduct) LoadRecipeCombo(combo_field *GUI.ComboBox) {
	object.Recipes = nil
	rows, err := DB.DB_Select_product_recipe.Query(object.Product_id)
	i := 0
	GUI.Fill_combobox_from_query_rows(combo_field, rows, err, func(rows *sql.Rows) {

		var (
			recipe_data ProductRecipe
		)

		if err := rows.Scan(&recipe_data.Recipe_id); err != nil {
			log.Fatal(err)
		}
		recipe_data.Product_id = object.Product_id
		log.Println("DEBUG: GetRecipes qc_data", recipe_data)
		object.Recipes = append(object.Recipes, &recipe_data)
		// combo_field.AddItem(strconv.FormatInt(i, 10))
		combo_field.AddItem(strconv.Itoa(i))
		i++
	})
	if i != 0 {
		combo_field.SetSelectedItem(0)
	}
}

func (object *RecipeProduct) NewRecipe() *ProductRecipe {
	proc_name := "RecipeProduct.NewRecipe"
	var (
		recipe_data *ProductRecipe
	)
	result, err := DB.DB_Insert_product_recipe.Exec(object.Product_id)
	if err != nil {
		log.Printf("%q: %s\n", err, proc_name)
		return recipe_data
	}
	insert_id, err := result.LastInsertId()
	if err != nil {
		log.Printf("%q: %s\n", err, proc_name)
		return recipe_data
	}
	recipe_data = new(ProductRecipe)

	recipe_data.Recipe_id = insert_id
	recipe_data.Product_id = object.Product_id
	object.Recipes = append(object.Recipes, recipe_data)
	return recipe_data
}

func (object *RecipeProduct) GetComponents() {
	for _, recipe_data := range object.Recipes {
		recipe_data.GetComponents()
	}
}

/*
func (object *RecipeProduct) AddRecipe() {
}*/

func NewRecipeProduct() *RecipeProduct {
	return new(RecipeProduct)
}
