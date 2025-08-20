package blender

import (
	"database/sql"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
)

type RecipeComponent struct {
	//recipe_components_id
	Component_id      int64
	Component_type_id int64
	Component_name    string
	Component_amount  float64
	Add_order         int
}

func NewRecipeComponent() *RecipeComponent {
	return new(RecipeComponent)
}
func NewRecipeComponentfromSQL(row *sql.Rows) RecipeComponent {

	Recipe_component := RecipeComponent{}

	if err := row.Scan(
		&Recipe_component.Component_id,
		&Recipe_component.Component_type_id,
		&Recipe_component.Component_name,
		&Recipe_component.Component_amount,
		&Recipe_component.Add_order,
	); err != nil {
		log.Fatal("Crit: [RecipeComponent NewRecipeComponentfromSQL]: ", err)
	}
	return Recipe_component
}

func (object *RecipeComponent) Insert(Recipe_id int64) {
	proc_name := "RecipeComponent.Insert"
	object.Component_id = DB.Insert(proc_name, DB.DB_Insert_recipe_component, Recipe_id, object.Component_type_id, object.Component_amount, object.Add_order)
}

func (object *RecipeComponent) Update() {
	proc_name := "RecipeComponent.Update"
	DB.Update(proc_name, DB.DB_Update_recipe_component, object.Component_type_id, object.Component_amount, object.Add_order, object.Component_id)
}

func (object *RecipeComponent) Delete() {
	proc_name := "RecipeComponent.Delete"
	DB.Delete(proc_name, DB.DB_Delete_recipe_component, object.Component_id)
}
