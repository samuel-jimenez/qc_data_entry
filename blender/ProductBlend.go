package blender

type ProductBlend struct {
	Components []BlendComponent
	Recipe_id  int64
}

func NewProductBlend() *ProductBlend {
	return new(ProductBlend)
}

func NewProductBlendFromRecipe(ProductRecipe *ProductRecipe) *ProductBlend {
	ProductBlend := NewProductBlend()
	ProductBlend.Recipe_id = ProductRecipe.Recipe_id
	return ProductBlend
}

func (object *ProductBlend) AddComponent(component_data BlendComponent) {
	object.Components = append(object.Components, component_data)
}

func (object *ProductBlend) Save(Lot_id int64) {
	for _, val := range object.Components {
		val.Save(Lot_id, object.Recipe_id)
	}
}
