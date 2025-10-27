package qc_ui

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/windigo"
)

/* MassDataViewable
 *
 */
type MassDataViewable interface {
	GUI.NumberEditViewable
	Clear()
	SetFont(font *windigo.Font)
	// SetLabeledSize(label_width, field_width, subfield_width, unit_width, height int)
	check_or_error_mass(mass, sg, density float64)
	OnChangeMass(e *windigo.Event)
	OnChangeSG(e *windigo.Event)
	OnChangeDensity(e *windigo.Event)
}

/* MassDataView
 *
 */
type MassDataView struct {
	*GUI.NumberEditView
	sg_field,
	density_field *GUI.NumberUnitsEditView
	ranges_panel MassRangesViewable
}

// func MassDataView_from_new(parent *windigo.AutoPanel, ranges_panel MassRangesViewable, mass_field *GUI.NumberEditView) *MassDataView {

func MassDataView_from_new(parent *windigo.AutoPanel, ranges_panel MassRangesViewable) *MassDataView {

	view := new(MassDataView)

	//TAB ORDER
	mass_field := GUI.NumberEditView_from_new(parent, formats.MASS_TEXT)

	sg_field := GUI.NumberEditView_with_Units_from_new(parent, formats.SG_TEXT, formats.SG_UNITS)
	density_field := GUI.NumberEditView_with_Units_from_new(parent, formats.DENSITY_TEXT, formats.DENSITY_UNITS)

	//PUSH TO BOTTOM
	parent.Dock(density_field, windigo.Bottom)
	parent.Dock(sg_field, windigo.Bottom)

	mass_field.OnChange().Bind(view.OnChangeMass)
	sg_field.OnChange().Bind(view.OnChangeSG)
	density_field.OnChange().Bind(view.OnChangeDensity)

	view.NumberEditView = mass_field
	view.sg_field = sg_field
	view.density_field = density_field
	view.ranges_panel = ranges_panel

	return view
}

func (view *MassDataView) Clear() {
	view.NumberEditView.Clear()
	view.sg_field.Clear()
	view.density_field.Clear()
}

func (view *MassDataView) SetFont(font *windigo.Font) {
	view.NumberEditView.SetFont(font)
	view.sg_field.SetFont(font)
	view.density_field.SetFont(font)
}

func (view *MassDataView) SetLabeledSize(label_width, field_width, subfield_width, unit_width, height int) {
	view.NumberEditView.SetLabeledSize(label_width, field_width, height)
	view.sg_field.SetLabeledSize(label_width, subfield_width, unit_width, height)
	view.density_field.SetLabeledSize(label_width, subfield_width, unit_width, height)
}

func (view *MassDataView) check_or_error_mass(mass, sg, density float64) {
	view.NumberEditView.Check(view.ranges_panel.CheckMass(mass))
	view.sg_field.Check(view.ranges_panel.CheckSG(sg))
	view.density_field.Check(view.ranges_panel.CheckDensity(density))
}

func (view *MassDataView) FixMass() (mass, sg, density float64) {
	mass = view.NumberEditView.GetFixed()
	sg = formats.SG_from_mass(mass)
	density = formats.Density_from_sg(sg)

	view.sg_field.SetText(formats.Format_sg(sg, false))
	view.density_field.SetText(formats.Format_density(density))

	return mass, sg, density

}

func (view *MassDataView) OnChangeMass(e *windigo.Event) {
	mass, sg, density := view.FixMass()

	view.check_or_error_mass(mass, sg, density)

}

//TODO FixMass FOR EACH
// CHECKALL
// checks := range_field.CheckAll(this_val, other_val)
// 	this_check, other_check := checks[0], checks[1]
// 	this_field.Check(this_check)
// 	other_field.Check(other_check && diff_val)
// check_or_error_mass_all
// CheckMassAll

func (view *MassDataView) OnChangeSG(e *windigo.Event) {
	sg := view.sg_field.GetFixed()
	mass := formats.Mass_from_sg(sg)
	density := formats.Density_from_sg(sg)

	view.NumberEditView.SetText(formats.Format_mass(mass))
	view.density_field.SetText(formats.Format_density(density))

	view.check_or_error_mass(mass, sg, density)

}

func (view *MassDataView) OnChangeDensity(e *windigo.Event) {
	density := view.density_field.GetFixed()
	sg := formats.SG_from_density(density)
	mass := formats.Mass_from_sg(sg)

	view.sg_field.SetText(formats.Format_sg(sg, false))
	view.NumberEditView.SetText(formats.Format_mass(mass))

	view.check_or_error_mass(mass, sg, density)

}
