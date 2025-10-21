package blender

/*
 * BlendWessel
 *
 */
type BlendWessel struct {
	Vessel   string
	Capacity string
	// Capacity float64

	// Capacity,
	Strap, HeelVolume, HeelMass float64
}

func (view BlendWessel) Get() (string, string, float64, float64, float64) {
	return view.Vessel,
		// view.Capacity_amount,
		view.Capacity,
		view.Strap,
		view.HeelVolume,
		view.HeelMass
}
