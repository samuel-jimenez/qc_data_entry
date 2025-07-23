package blender

import (
	"database/sql"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
)

type ProductRecipe struct {
	Components []*RecipeComponent
	// Product_name string `json:"product_name"`
	// Lot_number               string `json:"lot_number"`
	// Sample_point             string
	// Visual                   bool
	Product_id int64
	Recipe_id  int64
	// Product_name_customer_id nullable.NullInt64
	// Product_name_customer    string `json:"customer_product_name"`
}

func (object *ProductRecipe) GetComponents() {
	DB.Forall("GetComponents",
		func() {
			object.Components = nil
		},
		func(rows *sql.Rows) {
			Recipe_component := new(RecipeComponent)

			if err := rows.Scan(&Recipe_component.Component_name, &Recipe_component.Component_id, &Recipe_component.Component_amount, &Recipe_component.Add_order); err != nil {
				log.Fatal(err)
			}
			log.Println("DEBUG: GetComponents qc_data", Recipe_component)
			object.Components = append(object.Components, Recipe_component)
		},
		DB.DB_Select_recipe_components, object.Recipe_id)
}

// TODO give our id to children when writing
// recipe_list_id integer not null,
// component_type_id not null,
// component_type_amount real,
// component_add_order not null,
// TODO
// func (object *ProductRecipe) AddComponent() *RecipeComponent {
// 	proc_name := "ProductRecipe.AddComponent"
// 	var (
// 		component_data *RecipeComponent
// 	)
// 	result, err := DB.DB_Insert_recipe_component.Exec(object.Product_id)
// 	if err != nil {
// 		log.Printf("%q: %s\n", err, proc_name)
// 		return component_data
// 	}
// 	insert_id, err := result.LastInsertId()
// 	if err != nil {
// 		log.Printf("%q: %s\n", err, proc_name)
// 		return component_data
// 	}
// 	component_data = new(RecipeComponent)
//
// 	component_data.Component_id = insert_id
// 	component_data.Product_id = object.Product_id
// 	object.Components = append(object.Components, component_data)
// 	return component_data
// }

//
//
