package views

import (
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

/* MassRangesViewable
 *
 */
type MassRangesViewable interface {
	CheckMass(data float64) bool
	CheckSG(data float64) bool
	CheckDensity(data float64) bool
	Clear()
	// Update(qc_product *product.QCProduct)
}

/* MassRangesView
 *
 */
type MassRangesView struct {
	Mass_field,
	SG_field,
	Density_field *RangeROView
}

func NewMassRangesView(mass_field,
	sg_field,
	density_field *RangeROView) *MassRangesView {
	return &MassRangesView{Mass_field: mass_field,
		SG_field:      sg_field,
		Density_field: density_field}

	// view := new(MassRangesView)
	// view.Mass_field = mass_field
	// view.SG_field = sg_field
	// view.Density_field = density_field
	// return view
}

func (view *MassRangesView) CheckMass(data float64) bool {
	return view.Mass_field.Check(data)
}

func (view *MassRangesView) CheckSG(data float64) bool {
	return view.SG_field.Check(data)
}

func (view *MassRangesView) CheckDensity(data float64) bool {
	return view.Density_field.Check(data)
}

func (view *MassRangesView) CheckMassAll(data []float64) []bool {
	return view.Mass_field.CheckAll(data...)
}

func (view *MassRangesView) CheckSGAll(data []float64) []bool {
	return view.SG_field.CheckAll(data...)
}

func (view *MassRangesView) CheckDensityAll(data []float64) []bool {
	return view.Density_field.CheckAll(data...)
}

func (view *MassRangesView) Clear() {
	view.Mass_field.Clear()
	view.SG_field.Clear()
	view.Density_field.Clear()
}

func (view *MassRangesView) Update(qc_product *product.QCProduct) {
	view.Mass_field.Update(qc_product.SG)
	view.SG_field.Update(qc_product.SG)
	view.Density_field.Update(qc_product.Density)
}

/* MassDataViewable
 *
 */
type MassDataViewable interface {
	NumberEditViewable
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
	*NumberEditView
	sg_field,
	density_field *NumberUnitsEditView
	ranges_panel MassRangesViewable
}

// func NewMassDataView(parent *windigo.AutoPanel, ranges_panel MassRangesViewable, mass_field *NumberEditView) *MassDataView {

func NewMassDataView(parent *windigo.AutoPanel, ranges_panel MassRangesViewable) *MassDataView {

	view := new(MassDataView)

	//TAB ORDER
	mass_field := NewNumberEditView(parent, MASS_TEXT)

	sg_field := NewNumberEditViewWithUnits(parent, SG_TEXT, formats.SG_UNITS)
	density_field := NewNumberEditViewWithUnits(parent, DENSITY_TEXT, formats.DENSITY_UNITS)

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
