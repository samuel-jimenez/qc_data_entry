package blender_ui

import (
	"database/sql"
	"math"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
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
	top, btm                     *windigo.AutoPanel
	Vessel_field, Capacity_field *GUI.SearchBox
	Strap_field                  *GUI.NumbSearchView
	heel_field                   *GUI.NumbestEditView
	volume_data_strap            map[float64]float64
	strap_data_volume            map[float64]float64
	capacity_data                map[string]int
	parent                       *BlendStrappingProductView
	MinStrap, MaxStrap           float64
	Strap, HeelVolume, HeelMass  float64
	Capacity_amount              float64
	Capacity_description         string
}

func BlendVessel_from_new(parent *BlendStrappingProductView) *BlendVessel {

	Vessel_text := "Vessel"
	Capacity_text := "Capacity"
	Strap_text := "Heel Strap"
	heel_text := "Heel Weight"

	view := new(BlendVessel)
	view.parent = parent

	view.volume_data_strap = make(map[float64]float64)
	view.strap_data_volume = make(map[float64]float64)
	view.capacity_data = make(map[string]int)

	view.AutoPanel = windigo.NewAutoPanel(parent)
	view.top = windigo.NewAutoPanel(view)
	view.btm = windigo.NewAutoPanel(view)

	view.Vessel_field = GUI.NewLabeledListSearchBox(view.top, Vessel_text)
	view.Capacity_field = GUI.NewLabeledListSearchBox(view.top, Capacity_text)
	view.Strap_field = GUI.NumbSearchView_From_SearchBox(GUI.NewLabeledListSearchBox(view.btm, Strap_text))
	view.heel_field = GUI.NumbestEditView_from_new(view.btm, heel_text)

	// bad things happen if this is not true
	view.MinStrap = 0

	// combobox
	GUI.Fill_combobox_from_query_fn(
		view.Capacity_field,
		func(id int, name string) {
			view.capacity_data[name] = id
			view.Capacity_field.AddItem(name)
		},
		DB.DB_Select_container_capacity_all,
	)
	view.Capacity_field.SetSelectedItem(0)
	view.OnChange_Capacity_field(nil)

	view.Dock(view.top, windigo.Top)
	view.Dock(view.btm, windigo.Top)

	view.top.Dock(view.Vessel_field, windigo.Left)
	view.top.Dock(view.Capacity_field, windigo.Left)

	view.btm.Dock(view.Strap_field, windigo.Left)
	view.btm.Dock(view.heel_field, windigo.Left)

	view.Strap_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		view.Strap = view.Strap_field.Get()
		view.HeelVolume = view.volume_data_strap[view.Strap]
		view.SetHeelVolume(view.HeelVolume)
		// TODO broken due to... reasons
		// view.Strap_field.Set(strap)
	})
	view.heel_field.OnChange().Bind(func(e *windigo.Event) {
		view.SetHeel(view.heel_field.Get())
	})
	view.Capacity_field.OnSelectedChange().Bind(view.OnChange_Capacity_field)
	// view.Vessel_field.OnChange().Bind(func(e *windigo.Event) {
	// 	view.SetHeel(view.Vessel_field.Get())
	// })

	return view
}

func (view *BlendVessel) OnChange_Capacity_field(*windigo.Event) {
	clear(view.volume_data_strap)
	clear(view.strap_data_volume)
	cap_id := view.capacity_data[view.Capacity_field.GetSelectedItem()]
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
		DB.DB_Select_container_strap_container_capacity, cap_id,
	)
	proc_name := "BlendVessel-OnChange_Capacity_field"
	DB.Select_Panic(proc_name,
		DB.DB_Select_container_capacity_info.QueryRow(cap_id),
		&view.Capacity_description,
		&view.Capacity_amount,
	)

}

func (view *BlendVessel) SetFont(font *windigo.Font) {
	view.Vessel_field.SetFont(font)
	view.Capacity_field.SetFont(font)
	view.Strap_field.SetFont(font)
	view.heel_field.SetFont(font)
}

func (view *BlendVessel) RefreshSize() {
	view.SetSize(GUI.TOP_PANEL_WIDTH, 2*GUI.PRODUCT_FIELD_HEIGHT)
	view.top.SetSize(GUI.TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.btm.SetSize(GUI.TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.Vessel_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Capacity_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Capacity_field.SetPaddingLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)

	view.Strap_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.heel_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.heel_field.SetPaddingLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)

}

func (view *BlendVessel) SetHeelVolume(heel float64) {
	view.parent.SetHeelVolume(heel)
}

func (view *BlendVessel) SetHeel(heel float64) {
	view.HeelMass = heel
	view.parent.SetHeel(heel)
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
	view.HeelVolume = volume
	view.Strap = view.GetStrap(volume)

	view.Strap_field.Set(view.Strap)
}

func (view *BlendVessel) Get() blender.BlendWessel {
	return blender.BlendWessel{
		Vessel: view.Vessel_field.Text(),
		// view.Capacity_amount,
		Capacity:   view.Capacity_description,
		Strap:      view.Strap,
		HeelVolume: view.HeelVolume,
		HeelMass:   view.HeelMass}
}
