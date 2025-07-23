package blender

type RecipeComponent struct {
	Component_name   string
	Component_amount float64
	Component_id     int64
	Add_order        int
	// Lot_id     int64
	// Product_name_customer_id nullable.NullInt64
	// Product_name_customer    string `json:"customer_product_name"`
}

func NewRecipeComponent() *RecipeComponent {
	return new(RecipeComponent)
}
