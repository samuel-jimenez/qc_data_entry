package blender

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/util"
)

type ProductBlend struct {
	Components []BlendComponent
	Recipe_id  int64
	// Total, Heel, Amount float64
	Total, Heel, Amount int
}

func ProductBlend_from_new() *ProductBlend {
	return new(ProductBlend)
}

func NewProductBlend_from_Recipe(ProductRecipe *ProductRecipe) *ProductBlend {
	ProductBlend := ProductBlend_from_new()
	ProductBlend.Recipe_id = ProductRecipe.Recipe_id
	return ProductBlend
}

func (object *ProductBlend) AddComponent(component_data BlendComponent) {
	object.Components = append(object.Components, component_data)
}

func (object *ProductBlend) Save(Product_Lot_id int64) {
	// only save once
	proc_name := "ProductBlend.Save"
	count := 0
	DB.Select_Error(proc_name,
		DB.DB_Select_Product_count_blend_components.QueryRow(Product_Lot_id),
		&count,
	)
	if count != 0 {
		return
	}

	for _, val := range object.Components {
		val.Save(Product_Lot_id)
	}
}

func (object *ProductBlend) GetProcedure() (Procedure []string) {
	proc_name := "ProductBlend-GetProcedure"
	find := regexp.MustCompile("{{(.*?)}}")
	DB.Forall_err(proc_name,
		util.NOOP,
		func(row *sql.Rows) error {
			var step string
			if err := row.Scan(
				&step,
			); err != nil {
				return err
			}
			step = find.ReplaceAllStringSubmatchFunc(step, func(match []string) string {
				val := util.VALUES(strconv.Atoi(strings.Replace(match[1], "comp_", "", 1))) - 1
				if val < 0 || val > len(object.Components) {
					util.LogCry(proc_name,
						"Component mismatch for Recipe %v:\n Available: %v, found: %v", object.Recipe_id, len(object.Components), val)
				}
				return object.Components[val].Component_name
			})

			// match := find.FindStringSubmatch(step)
			// for match != nil { // while
			// 	val := util.VALUES(strconv.Atoi(strings.Replace(match[1], "comp_", "", 1))) - 1
			// 	if val < 0 || val > len(object.Components) {
			// 		util.LogCry(proc_name,
			// 			"Component mismatch for Recipe %v:\n Available: %v, found: %v", object.Recipe_id, len(object.Components), val)
			// 	}
			// 	step = strings.Replace(step, match[0], object.Components[val].Component_name, 1)
			// 	match = find.FindStringSubmatch(step)
			// }
			Procedure = append(Procedure, step)
			return nil
		},
		DB.DB_Select_recipe_procedure_steps, object.Recipe_id)
	return
}
