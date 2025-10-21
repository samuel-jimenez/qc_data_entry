package blender

import (
	"database/sql"
	"log"
	"maps"
	"math"
	"slices"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
)

// TODO []*
type ProductRecipe struct {
	Components, db_components []RecipeComponent
	Recipe_id                 int64
	Recipe_name               string
	Total_Default,            //  nullable.NullInt64
	Specific_gravity_Default nullable.NullFloat64
}

func ProductRecipe_from_Recipe_id_name(Recipe_id int64, Recipe_name string) *ProductRecipe {
	object := new(ProductRecipe)
	object.Recipe_id = Recipe_id
	object.Recipe_name = Recipe_name
	return object
}

func ProductRecipe_from_SQL(row *sql.Rows) (*ProductRecipe, error) {
	object := new(ProductRecipe)
	err := row.Scan(
		&object.Recipe_id,
		&object.Recipe_name,
		&object.Total_Default,
		&object.Specific_gravity_Default,
	)
	if err != nil {
		log.Println("Crit: [ProductRecipe-ProductRecipe_from_SQL]: ", err)
	}
	return object, err
}

func (object *ProductRecipe) GetComponents() {
	DB.Forall("GetComponents",
		func() {
			object.Components = nil
		},
		func(row *sql.Rows) {
			Recipe_component := NewRecipeComponentfromSQL(row)
			object.Components = append(object.Components, Recipe_component)
		},
		DB.DB_Select_recipe_components_id, object.Recipe_id)
	object.db_components = slices.Clone(object.Components)
}
func (object *ProductRecipe) AddComponent(component_data RecipeComponent) {
	object.Components = append(object.Components, component_data)
}

func (object *ProductRecipe) NormalizeComponents(total_Component_amount float64) {
	if math.Abs(total_Component_amount-1) > 1e-3 {
		unnormalizedComponents := object.Components
		object.Components = nil
		// normalize amounts
		for _, Recipe_component := range unnormalizedComponents {
			Recipe_component.Component_amount /= total_Component_amount
			object.Components = append(object.Components, Recipe_component)
		}
	}
}

func (object *ProductRecipe) SaveComponents() {
	log.Println("DEBUG: [ProductRecipe SaveComponents]: ", object.db_components, object.Components)
	//TODO set diff
	// map[T]struct{}
	var del_set, add_set, up_set []*RecipeComponent

	lookup_map := make(map[int64]*RecipeComponent)
	add_map := make(map[int64]*RecipeComponent)
	// Component_id
	// map[T]struct{}
	for _, db := range object.db_components {
		lookup_map[db.Component_id] = &db
		log.Println("DEBUG: ProductRecipe db_components", db)

	}
	for _, val := range object.Components {
		log.Println("DEBUG: ProductRecipe Components", val)

		oldVal := lookup_map[val.Component_id]
		log.Println("DEBUG: ProductRecipe oldVal", oldVal)

		// new Component_id
		if oldVal == nil {
			log.Println("DEBUG: ProductRecipe tryadd")
			// add_set = append(add_set, val)
			add_map[val.Component_type_id] = &val
			continue
		}
		if *oldVal != val {
			// value change
			log.Println("DEBUG: ProductRecipe up_set")
			up_set = append(up_set, &val)
		}
		// no change (intersection)
		log.Println("DEBUG: ProductRecipe intersection")
		delete(lookup_map, val.Component_id)
	}
	// check if any product was deleted and re-added
	//TODO test this
	for _, val := range lookup_map {
		newVal := add_map[val.Component_type_id]
		log.Println("DEBUG: ProductRecipe val,newVal", val, newVal)
		if newVal == nil {
			log.Println("DEBUG: ProductRecipe del")
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

	for _, val := range del_set {
		log.Println("DEBUG: ProductRecipe Delete", *val)
		val.Delete()
	}
	for _, val := range up_set {
		log.Println("DEBUG: ProductRecipe Update", *val)
		val.Update()
	}
	for _, val := range add_set {
		log.Println("DEBUG: ProductRecipe Insert", *val)
		val.Insert(object.Recipe_id)
	}
	object.db_components = slices.Clone(object.Components)

}
