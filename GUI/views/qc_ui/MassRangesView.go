package qc_ui

import (
	"github.com/samuel-jimenez/qc_data_entry/product"
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

func MassRangesView_from_new(mass_field,
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
