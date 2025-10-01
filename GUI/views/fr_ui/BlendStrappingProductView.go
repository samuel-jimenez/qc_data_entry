package fr_ui

import (
	"github.com/samuel-jimenez/windigo"
)

/*
 * BlendStrappingProductViewer
 *
 */
type BlendStrappingProductViewer interface {
	SetFont(font *windigo.Font)
	RefreshSize()
	SetHeelVolume(float64)
}

/*
 * BlendStrappingProductView
 *
 */
type BlendStrappingProductView struct {
	*windigo.AutoPanel

	Blend_Vessel       *BlendVessel
	Blend_product_view *BlendProductView
}

func NewBlendStrappingProductView(parent windigo.Controller) *BlendStrappingProductView {

	view := new(BlendStrappingProductView)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.Blend_Vessel = NewBlendVessel(view)
	view.Blend_product_view = NewBlendProductView(view)

	// Dock
	view.Dock(view.Blend_Vessel, windigo.Top)
	view.Dock(view.Blend_product_view, windigo.Fill)

	return view
}

// functionality

func (view *BlendStrappingProductView) SetFont(font *windigo.Font) {
	view.Blend_Vessel.SetFont(windigo.DefaultFont)
	view.Blend_product_view.SetFont(windigo.DefaultFont)

}

func (view *BlendStrappingProductView) RefreshSize() {
	view.Blend_Vessel.RefreshSize()
	view.Blend_product_view.RefreshSize()
}

func (view *BlendStrappingProductView) SetHeelVolume(heel float64) {
	view.Blend_product_view.SetHeelVolume(heel)
}

func (view *BlendStrappingProductView) SetStrap(volume float64) {
	view.Blend_Vessel.SetStrap(volume)
}
