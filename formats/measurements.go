package formats

import (
	"strconv"

	"github.com/samuel-jimenez/qc_data_entry/util/math"
)

var (
	SAMPLE_VOLUME = 83.2
	LB_PER_GAL    = 8.34 // g/mL
)

var (
	SG_UNITS        = "g/mL"
	DENSITY_UNITS   = "lb/gal"
	STRING_UNITS    = "s"
	VISCOSITY_UNITS = "cP"

	SG_PRECISION_FIXED = 4
	SG_PRECISION_VARIABLE = 3
	PH_PRECISION          = 2
	DENSITY_PRECISION     = 3
	STRING_TEST_PRECISION = 0
	VISCOSITY_PRECISION   = 0
)

// TODO fix math.Round()
func Round(val float64, precision int) float64 {
	out, _ := strconv.ParseFloat(Format_float(val, precision), 64)
	return out
}

func Format_float(val float64, precision int) string {
	return strconv.FormatFloat(val, 'f', precision, 64)
}

func SG_from_density(density float64) float64 {
	sg := density / LB_PER_GAL
	return sg
}

func SG_from_mass(mass float64) float64 {
	sg := math.SigFig(mass/SAMPLE_VOLUME, 4)
	return sg
}

func Density_from_sg(sg float64) float64 {
	density := math.SigFig(sg*LB_PER_GAL, 4)
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
	if fixed_precision || sg < 1 {
		return Format_fixed_sg(sg)
	} else {
		return Format_float(sg, SG_PRECISION_VARIABLE)
	}
}

func Format_fixed_sg(sg float64) string {
	return Format_float(sg, SG_PRECISION_FIXED)
}



func Format_ph(ph float64) string {
	return Format_float(ph, PH_PRECISION)
}

func Format_density(density float64) string {
	return Format_float(density, DENSITY_PRECISION)
}

func Format_string_test(string_test float64) string {
	return Format_float(string_test, STRING_TEST_PRECISION)
}

func Format_viscosity(viscosity float64) string {
	return Format_float(viscosity, VISCOSITY_PRECISION)
}

func FormatInt(_int int64) string {
	return strconv.FormatInt(_int, 10)
}

func Format_ranges_sg(sg float64) string {
	// return strconv.FormatFloat(sg, 'f', 2, 64)
	// when not fr
	return strconv.FormatFloat(sg, 'f', 3, 64)
}

func Format_ranges_ph(ph float64) string {
	return Format_float(ph, PH_PRECISION)
}

func Format_ranges_density(density float64) string {
	return Format_float(density, DENSITY_PRECISION)
}

func Format_ranges_string_test(string_test float64) string {
	return Format_float(string_test, STRING_TEST_PRECISION)
}

func Format_ranges_viscosity(viscosity float64) string {
	return Format_float(viscosity, VISCOSITY_PRECISION)
}
