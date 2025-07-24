package blender

import (
	"database/sql"
	"log"
	"maps"
	"slices"

	"github.com/samuel-jimenez/qc_data_entry/DB"
)

// TODO []*
type ProductRecipe struct {
	Components, db_components []*RecipeComponent
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
			Recipe_component := NewRecipeComponent()

			if err := rows.Scan(&Recipe_component.Component_name, &Recipe_component.Component_id, &Recipe_component.Component_amount, &Recipe_component.Add_order); err != nil {
				log.Fatal(err)
			}
			log.Println("DEBUG: GetComponents qc_data", Recipe_component)
			object.Components = append(object.Components, Recipe_component)
		},
		DB.DB_Select_recipe_components, object.Recipe_id)
	object.db_components = object.Components
}
func (object *ProductRecipe) AddComponent(component_data *RecipeComponent) {
	object.Components = append(object.Components, component_data)
}

// TODO give our id to children when writing
// recipe_list_id integer not null,
// component_type_id not null,
// component_type_amount real,
// component_add_order not null,
// TODO
// func (object *ProductRecipe) SaveComponents(component_data []*RecipeComponent) {
func (object *ProductRecipe) SaveComponents() {
	log.Println("DEBUG: ProductRecipe SaveComponents", object.db_components, object.Components)
	//TODO set diff
	// map[T]struct{}
	// old_min := 0
	// old_max := len(object.db_components)
	// new_min := 0
	// new_max := len(object.Components)
	var del_set, add_set, up_set []*RecipeComponent

	lookup := make(map[int64]*RecipeComponent)
	// Component_id
	// map[T]struct{}
	for _, db := range object.db_components {
		lookup[db.Component_id] = db
		log.Println("DEBUG: ProductRecipe db_components", db)

	}
	for _, val := range object.Components {
		log.Println("DEBUG: ProductRecipe Components", val)

		oldVal := lookup[val.Component_id]
		if oldVal == nil {
			add_set = append(add_set, val)
			continue
		}
		if oldVal == val {
			delete(lookup, val.Component_id)
			continue
		}
		up_set = append(up_set, val)
	}
	del_set = slices.Collect(maps.Values(lookup))
	log.Println("DEBUG: ProductRecipe del_set, add_set, up_set", del_set, add_set, up_set)

	// 	for i := old_min; i< old_max
	//
	// 	}
	//
	// 	proc_name := "ProductRecipe.AddComponent"
	//
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

}
