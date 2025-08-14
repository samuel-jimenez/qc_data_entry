package formats

import (
	"strconv"
)

var SAMPLE_VOLUME = 83.2
var LB_PER_GAL = 8.34 // g/mL

var (
	SG_UNITS        = "g/mL"
	DENSITY_UNITS   = "lb/gal"
	STRING_UNITS    = "s"
	VISCOSITY_UNITS = "cP"
)

func SG_from_density(density float64) float64 {
	sg := density / LB_PER_GAL
	return sg
}

func SG_from_mass(mass float64) float64 {
	sg, _ := strconv.ParseFloat(strconv.FormatFloat(mass/SAMPLE_VOLUME, 'G', 4, 64), 64)
	return sg
}

func Density_from_sg(sg float64) float64 {
	density := sg * LB_PER_GAL
	return density
}

func Mass_from_sg(sg float64) float64 {
	mass := sg * SAMPLE_VOLUME
	return mass
}

func Format_mass(mass float64) string {
	return strconv.FormatFloat(mass, 'f', 2, 64)
}

func Format_sg(sg float64, fixed_precision bool) string {
	var precision int
	if fixed_precision || sg < 1 {
		precision = 4
	} else {
		precision = 3
	}
	return strconv.FormatFloat(sg, 'f', precision, 64)
}

func Format_ph(ph float64) string {
	return strconv.FormatFloat(ph, 'f', 2, 64)
}

func Format_density(density float64) string {
	return strconv.FormatFloat(density, 'f', 3, 64)
}

func Format_string_test(string_test float64) string {
	return strconv.FormatFloat(string_test, 'f', 0, 64)
}

func Format_viscosity(viscosity float64) string {
	return strconv.FormatFloat(viscosity, 'f', 0, 64)
}

func FormatInt(_int int64) string {
	return strconv.FormatInt(_int, 10)
}

func Format_ranges_sg(sg float64) string {
	// return strconv.FormatFloat(sg, 'f', 2, 64)
	//when not fr
	return strconv.FormatFloat(sg, 'f', 3, 64)
}

func Format_ranges_ph(ph float64) string {
	return strconv.FormatFloat(ph, 'f', 2, 64)
}

func Format_ranges_density(density float64) string {
	return strconv.FormatFloat(density, 'f', 3, 64)
}

func Format_ranges_string_test(string_test float64) string {
	return strconv.FormatFloat(string_test, 'f', 0, 64)
}

func Format_ranges_viscosity(viscosity float64) string {
	return strconv.FormatFloat(viscosity, 'f', 0, 64)
}
