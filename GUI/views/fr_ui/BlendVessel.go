package fr_ui

import (
	"database/sql"
	"math"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/*
 * BlendVesseler
 *
 */
type BlendVesseler interface {
	windigo.Pane

	SetFont(font *windigo.Font)
	RefreshSize()
}

/*
 * BlendVessel
 *
 */
type BlendVessel struct {
	*windigo.AutoPanel
	Vessel_field       *GUI.SearchBox
	Strap_field        *GUI.NumbSearchView
	volume_data_strap  map[float64]float64
	strap_data_volume  map[float64]float64
	parent             *BlendStrappingProductView
	MinStrap, MaxStrap float64
}

func NewBlendVessel(parent *BlendStrappingProductView) *BlendVessel {

	Vessel_text := "Vessel"
	Strap_text := "Strap"
	container_capacity := 2

	view := new(BlendVessel)
	view.parent = parent

	view.volume_data_strap = make(map[float64]float64)
	view.strap_data_volume = make(map[float64]float64)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.Vessel_field = GUI.NewLabeledListSearchBox(view, Vessel_text)
	view.Strap_field = GUI.NumbSearchView_From_SearchBox(GUI.NewLabeledListSearchBox(view, Strap_text))

	// bad things happen if this is not true
	view.MinStrap = 0

	// combobox
	GUI.Fill_combobox_from_query_rows(
		view.Strap_field,
		func(row *sql.Rows) error {
			var (
				strap  float64
				volume float64
			)
			if err := row.Scan(
				&strap, &volume,
			); err != nil {
				return err
			}
			view.volume_data_strap[strap] = volume
			view.strap_data_volume[volume] = strap
			view.Strap_field.Add(strap)
			view.MaxStrap = strap
			return nil
		},
		DB.DB_Select_container_strap_container_capacity, container_capacity)

	view.Dock(view.Vessel_field, windigo.Left)
	view.Dock(view.Strap_field, windigo.Left)

	view.Strap_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		strap := view.Strap_field.Get()
		view.SetHeelVolume(view.volume_data_strap[strap])
		// TODO broken due to... reasons
		// view.Strap_field.Set(strap)
	})
	// view.Vessel_field.OnChange().Bind(func(e *windigo.Event) {
	// 	view.SetHeel(view.Vessel_field.Get())
	// })

	return view
}

func (view *BlendVessel) SetFont(font *windigo.Font) {
	view.Vessel_field.SetFont(font)
	view.Strap_field.SetFont(font)
}

func (view *BlendVessel) RefreshSize() {
	view.SetSize(GUI.TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Vessel_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Strap_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
}

func (view *BlendVessel) SetHeelVolume(heel float64) {
	view.parent.SetHeelVolume(heel)
}

// TODO Round function
// this assumes a 1/4" strap
func NormalizeStrap(strap float64) float64 {
	return math.Round(strap*4) / 4
}
func (view *BlendVessel) GetStrap(volume float64) float64 {
	if volume == 0 { // degenerate case
		return 0
	}
	strap := view.strap_data_volume[volume]
	if strap != 0 { // what I wouldn't give for some Rust Options
		return strap
	}
	var (
		s0, s1, s2, Δv0, Δv1 float64
	)

	// Regula falsi, 5 iterations
	s0 = view.MinStrap
	s1 = view.MaxStrap
	Δv0 = view.volume_data_strap[s0] - volume
	Δv1 = view.volume_data_strap[s1] - volume
	for i := 0; i < 5 && s0 != s1; i++ {
		s2 = NormalizeStrap((s0*Δv1 - s1*Δv0) / (Δv1 - Δv0))
		if math.Abs(Δv0) < math.Abs(Δv1) {
			s1 = s2
			Δv1 = view.volume_data_strap[s1] - volume
		} else {
			s0 = s2
			Δv0 = view.volume_data_strap[s0] - volume
		}
	}

	if math.Abs(Δv0) < math.Abs(Δv1) {
		return s0
	} else {
		return s1
	}
}

func (view *BlendVessel) SetStrap(volume float64) {
	view.Strap_field.Set(view.GetStrap(volume))
}
