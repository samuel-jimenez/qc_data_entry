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

			if err := rows.Scan(&Recipe_component.Component_id, &Recipe_component.Component_name, &Recipe_component.Component_type_id, &Recipe_component.Component_amount, &Recipe_component.Add_order); err != nil {
				log.Fatal("Crit: ProductRecipe GetComponents ", err)
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

	lookup_map := make(map[int64]*RecipeComponent)
	add_map := make(map[int64]*RecipeComponent)
	// Component_id
	// map[T]struct{}
	for _, db := range object.db_components {
		lookup_map[db.Component_id] = db
		log.Println("DEBUG: ProductRecipe db_components", db)

	}
	for _, val := range object.Components {
		log.Println("DEBUG: ProductRecipe Components", val)

		oldVal := lookup_map[val.Component_id]
		// new Component_id
		if oldVal == nil {
			// add_set = append(add_set, val)
			add_map[val.Component_type_id] = val
			continue
		}
		// no change (intersection)
		if oldVal == val {
			delete(lookup_map, val.Component_id)
			continue
		}
		// value change
		up_set = append(up_set, val)
	}
	// check if any product was deleted and re-added
	for _, val := range lookup_map {
		newVal := add_map[val.Component_type_id]
		if newVal == nil {
			del_set = append(del_set, val)
			continue
		}
		log.Println("DEBUG: ProductRecipe add_set match", newVal, val)
		newVal.Component_id = val.Component_id
		up_set = append(up_set, newVal)

		delete(add_map, val.Component_type_id)

		log.Println("DEBUG: ProductRecipe add_set matched", newVal, val)
	}

	add_set = slices.Collect(maps.Values(add_map))
	log.Println("DEBUG: ProductRecipe del_set, add_set, up_set", del_set, add_set, up_set)

	for _, val := range up_set {
		val.Update()
	}
	for _, val := range add_set {
		val.Insert(object.Recipe_id)
	}
	for _, val := range del_set {
		val.Delete()
	}

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
