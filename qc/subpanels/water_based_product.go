package subpanels

import (
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/product"
)

type WaterBasedProduct struct {
	product.BaseProduct
	sg float64
	ph float64
}

func (wb_product WaterBasedProduct) toProduct() product.Product {
	return product.Product{
		BaseProduct: wb_product.Base(),
		PH:          nullable.NewNullFloat64(wb_product.ph, true),
		SG:          nullable.NewNullFloat64(wb_product.sg, true),
		Density:     nullable.NewNullFloat64(0, false),
		String_test: nullable.NullInt64Default(),
		Viscosity:   nullable.NullInt64Default(),
	}

	//TODO Option?
}

func newWaterBasedProduct(base_product product.BaseProduct, have_visual bool, sg, ph float64) product.Product {

	base_product.Visual = have_visual

	return WaterBasedProduct{base_product, sg, ph}.toProduct()

}

func (product WaterBasedProduct) Check_data() bool {
	return true
}
