package blender

import (
	"database/sql"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
)

type RecipeComponent struct {
	Component_id      int64
	Component_type_id int64
	Component_name    string
	Component_amount  float64
	Add_order         int
	// Lot_id     int64
	// Product_name_customer_id nullable.NullInt64
	// Product_name_customer    string `json:"customer_product_name"`
}

func NewRecipeComponent() *RecipeComponent {
	return new(RecipeComponent)
}
func NewRecipeComponentfromSQL(row *sql.Rows) *RecipeComponent {

	Recipe_component := NewRecipeComponent()

	if err := row.Scan(
		&Recipe_component.Component_id,
		&Recipe_component.Component_type_id,
		&Recipe_component.Component_name,
		&Recipe_component.Component_amount,
		&Recipe_component.Add_order); err != nil {
		log.Fatal("Crit: [RecipeComponent NewRecipeComponentfromSQL]: ", err)
	}
	return Recipe_component
}

func (object *RecipeComponent) Insert(Recipe_id int64) {
	proc_name := "RecipeComponent.Insert"
	result, err := DB.DB_Insert_recipe_component.Exec(Recipe_id, object.Component_type_id, object.Component_amount, object.Add_order)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return
	}
	object.Component_id, err = result.LastInsertId()
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return
	}
}

func (object *RecipeComponent) Update() {
	proc_name := "RecipeComponent.Update"
	_, err := DB.DB_Update_recipe_component.Exec(object.Component_type_id, object.Component_amount, object.Add_order, object.Component_id)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return
	}
}

func (object *RecipeComponent) Delete() {
	proc_name := "RecipeComponent.Delete"
	_, err := DB.DB_Delete_recipe_component.Exec(object.Component_id)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
	}
}
