package main

import (
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/windigo"
)

/* MassRangesViewable
 *
 */
type MassRangesViewable interface {
	CheckMass(data float64) bool
	CheckSG(data float64) bool
	CheckDensity(data float64) bool
}

/* MassRangesView
 *
 */
type MassRangesView struct {
	mass_field,
	sg_field,
	density_field *RangeROView
}

func (view MassRangesView) CheckMass(data float64) bool {
	return view.mass_field.Check(data)
}

func (view MassRangesView) CheckSG(data float64) bool {
	return view.sg_field.Check(data)
}

func (view MassRangesView) CheckDensity(data float64) bool {
	return view.density_field.Check(data)
}

func (data_view MassRangesView) Clear() {
	data_view.mass_field.Clear()
	data_view.sg_field.Clear()
	data_view.density_field.Clear()
}

/* MassDataViewable
 *
 */
type MassDataViewable interface {
	NumberEditViewable
	Clear()
}

/* MassDataView
 *
 */
type MassDataView struct {
	*NumberEditView
	sg_field,
	density_field *NumberUnitsEditView
}

func (data_view MassDataView) Clear() {
	data_view.NumberEditView.Clear()
	data_view.sg_field.Clear()
	data_view.density_field.Clear()
}

func (data_view MassDataView) SetFont(font *windigo.Font) {
	data_view.NumberEditView.SetFont(font)
	data_view.sg_field.SetFont(font)
	data_view.density_field.SetFont(font)
}

func (data_view MassDataView) SetLabeledSize(label_width, field_width, subfield_width, unit_width, height int) {
	data_view.NumberEditView.SetLabeledSize(label_width, field_width, height)
	data_view.sg_field.SetLabeledSize(label_width, subfield_width, unit_width, height)
	data_view.density_field.SetLabeledSize(label_width, subfield_width, unit_width, height)
}

func NewMassDataView(parent *windigo.AutoPanel, ranges_panel MassRangesViewable) *MassDataView {

	mass_text := "Mass"
	sg_text := "Specific Gravity"
	density_text := "Density"

	mass_field := NewNumberEditView(parent, mass_text)

	//PUSH TO BOTTOM
	density_field := NewNumberEditViewWithUnits(parent, density_text, formats.DENSITY_UNITS)
	sg_field := NewNumberEditViewWithUnits(parent, sg_text, formats.SG_UNITS)

	check_or_error_mass := func(mass, sg, density float64) {
		mass_field.Check(ranges_panel.CheckMass(mass))
		sg_field.Check(ranges_panel.CheckSG(sg))
		density_field.Check(ranges_panel.CheckDensity(density))
	}

	mass_field.OnChange().Bind(func(e *windigo.Event) {
		mass := mass_field.GetFixed()
		sg := formats.SG_from_mass(mass)
		density := formats.Density_from_sg(sg)

		sg_field.SetText(formats.Format_sg(sg, false))
		density_field.SetText(formats.Format_density(density))

		check_or_error_mass(mass, sg, density)

	})

	sg_field.OnChange().Bind(func(e *windigo.Event) {
		sg := sg_field.GetFixed()
		mass := formats.Mass_from_sg(sg)
		density := formats.Density_from_sg(sg)

		mass_field.SetText(formats.Format_mass(mass))
		density_field.SetText(formats.Format_density(density))

		check_or_error_mass(mass, sg, density)
	})

	density_field.OnChange().Bind(func(e *windigo.Event) {

		density := density_field.GetFixed()
		sg := formats.SG_from_density(density)
		mass := formats.Mass_from_sg(sg)

		sg_field.SetText(formats.Format_sg(sg, false))
		mass_field.SetText(formats.Format_mass(mass))

		check_or_error_mass(mass, sg, density)

	})

	return &MassDataView{mass_field, sg_field, density_field}
}
