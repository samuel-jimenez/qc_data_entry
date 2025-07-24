package blender

type BlendProduct struct {
	Product_name string `json:"product_name"`
	// Lot_number               string `json:"lot_number"`
	// Sample_point             string
	// Visual                   bool
	Product_id int64
	Lot_id     int64
	Recipe     ProductBlend
	// Product_name_customer_id nullable.NullInt64
	// Product_name_customer    string `json:"customer_product_name"`
}
