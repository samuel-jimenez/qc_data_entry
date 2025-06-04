package main

/* DerivedMassRangesView
 *
 */
type DerivedMassRangesView struct {
	mass_field,
	sg_field,
	density_field *RangeROView
}

func (view DerivedMassRangesView) CheckMass(data float64) bool {
	return view.mass_field.Check(data)
}

func (view DerivedMassRangesView) CheckSG(data float64) bool {
	return view.sg_field.Check(data)
}

func (view DerivedMassRangesView) CheckDensity(data float64) bool {
	return view.density_field.Check(data)
}

func (data_view DerivedMassRangesView) Clear() {
	data_view.mass_field.Clear()
	data_view.sg_field.Clear()
	data_view.density_field.Clear()
}

/* MassRangesView
 *
 */
type MassRangesView interface {
	CheckMass(data float64) bool
	CheckSG(data float64) bool
	CheckDensity(data float64) bool
}

/* MassDataView
 *
 */
type MassDataView struct {
	NumberEditView
	Clear func()
}

// type MassDataView interface {
// 	windigo.LabeledEdit
// 	Clear()
// }

// func (data_view MassDataView) Clear() {
// 	data_view.mass_field.Clear()
// 	data_view.sg_field.Clear()
// 	data_view.density_field.Clear()
// }

// func (data_view windigo.LabeledEdit) Clear() {
// 	data_view.mass_field.Clear()
// 	data_view.sg_field.Clear()
// 	data_view.density_field.Clear()
// }
