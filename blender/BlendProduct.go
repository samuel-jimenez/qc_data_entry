package blender

import (
	"log"
	"slices"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/DB"
)

type BlendProduct struct {
	Product_id     int64
	Product_Lot_id int64
	Blend          *ProductBlend
	Recipe         *ProductRecipe
}

func NewBlendProduct() *BlendProduct {
	return new(BlendProduct)
}

func NewBlendProductFromRecipe(RecipeProduct *RecipeProduct) *BlendProduct {
	BlendProduct := NewBlendProduct()
	BlendProduct.Product_id = RecipeProduct.Product_id
	return BlendProduct
}

func (object *BlendProduct) Save() {
	proc_name := "BlendProduct.Save"

	operations_group := "BSQL"

	Lot_Id := Next_Lot_Id(operations_group)
	object.Product_Lot_id = DB.Insert(proc_name, DB.DB_Insert_blend_lot, Lot_Id, object.Product_id, nil, object.Blend.Recipe_id) //TODO product_customer_id
	object.Blend.Save(object.Product_Lot_id)
}

// convert lot suffix string back into integer
func Beta(val string) int {
	charSize := 26
	total := 0
	for _, char := range val {
		total *= charSize
		total += int(char - 'A')
		total++
	}
	return total
}

// TODO extract, rename
// not a true base shift since
// the size of the string determines the available symbols
//
//	that is, blank characters only appear at the beginning
func Alpha(val int) string {
	charSize := 26
	var runes []rune
	for val >= charSize {
		runes = append(runes, arfa(val%charSize))
		val /= charSize
		val--
	}
	runes = append(runes, arfa(val))
	slices.Reverse(runes)
	return string(runes)
}
func arfa(val int) rune {
	return rune('A') + rune(val)
}

func _arfa(val int) string {
	return string(rune('A') + rune(val))
}

// TODO extract, rename
func BlendProductLOTS() string {
	t := time.Now()
	return t.Format("060102") //YYMMDD
}

// TODO extract, rename
func Next_Lot_Id(operations_group string) int64 {
	proc_name := "Next_Lot_Id"
	lot_date := BlendProductLOTS()
	lot_count := 0
	err := DB.DB_Select_blend_lot.QueryRow(operations_group + lot_date + "%").Scan(&lot_count)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return DB.INVALID_ID
	}
	lot_suffix := Alpha(lot_count)
	Lot_number := operations_group + lot_date + lot_suffix
	log.Println("Info: Next_Lot_Id: ", Lot_number)
	return DB.Insel_lot_id(Lot_number)
}
