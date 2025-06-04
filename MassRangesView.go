package main

import "github.com/samuel-jimenez/windigo"

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
	density_field *NumberEditView
}

func (data_view MassDataView) Clear() {
	data_view.NumberEditView.Clear()
	data_view.sg_field.Clear()
	data_view.density_field.Clear()
}

func NewMassDataView(parent windigo.AutoPanel, label_width, control_width, height int, field_text string, ranges_panel MassRangesViewable) *MassDataView {

	field_width := DATA_FIELD_WIDTH

	sg_text := "Specific Gravity"
	density_text := "Density"

	sg_units := "g/mL"
	density_units := "lb/gal"

	mass_field := NewNumberEditView(parent, label_width, control_width, height, field_text)

	//PUSH TO BOTTOM
	density_field := NewNumberEditViewWithUnits(parent, label_width, field_width, height, density_text, density_units)
	sg_field := NewNumberEditViewWithUnits(parent, label_width, field_width, height, sg_text, sg_units)

	check_or_error_mass := func(mass, sg, density float64) {
		mass_field.Check(ranges_panel.CheckMass(mass))
		sg_field.Check(ranges_panel.CheckSG(sg))
		density_field.Check(ranges_panel.CheckDensity(density))
	}

	mass_field.OnChange().Bind(func(e *windigo.Event) {
		mass := mass_field.GetFixed()
		sg := sg_from_mass(mass)
		density := density_from_sg(sg)

		sg_field.SetText(format_sg(sg, false))
		density_field.SetText(format_density(density))

		check_or_error_mass(mass, sg, density)

	})

	sg_field.OnChange().Bind(func(e *windigo.Event) {
		sg := sg_field.GetFixed()
		mass := mass_from_sg(sg)
		density := density_from_sg(sg)

		mass_field.SetText(format_mass(mass))
		density_field.SetText(format_density(density))

		check_or_error_mass(mass, sg, density)
	})

	density_field.OnChange().Bind(func(e *windigo.Event) {

		density := density_field.GetFixed()
		sg := sg_from_density(density)
		mass := mass_from_sg(sg)

		sg_field.SetText(format_sg(sg, false))
		mass_field.SetText(format_mass(mass))

		check_or_error_mass(mass, sg, density)

	})

	return &MassDataView{mass_field, sg_field, density_field}
}
