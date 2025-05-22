package main

import (
	"strconv"
)

var SAMPLE_VOLUME = 83.2
var LB_PER_GAL = 8.345 // g/mL

func sg_from_density(density float64) float64 {
	sg := density / LB_PER_GAL
	return sg
}

func sg_from_mass(mass float64) float64 {
	sg, _ := strconv.ParseFloat(strconv.FormatFloat(mass/SAMPLE_VOLUME, 'G', 4, 64), 64)
	return sg
}

func density_from_sg(sg float64) float64 {
	density := sg * LB_PER_GAL
	return density
}

func mass_from_sg(sg float64) float64 {
	mass := sg * SAMPLE_VOLUME
	return mass
}

func format_mass(mass float64) string {
	return strconv.FormatFloat(mass, 'f', 2, 64)
}

func format_sg(sg float64, fixed_precision bool) string {
	var precision int
	if fixed_precision || sg < 1 {
		precision = 4
	} else {
		precision = 3
	}
	return strconv.FormatFloat(sg, 'f', precision, 64)
}

func format_ph(ph float64) string {
	return strconv.FormatFloat(ph, 'f', 2, 64)
}

func format_density(density float64) string {
	return strconv.FormatFloat(density, 'f', 3, 64)
}

func format_string_test(string_test float64) string {
	return strconv.FormatFloat(string_test, 'f', 0, 64)
}

func format_viscosity(viscosity float64) string {
	return strconv.FormatFloat(viscosity, 'f', 0, 64)
}

func format_ranges_sg(sg float64) string {
	// return strconv.FormatFloat(sg, 'f', 2, 64)
	//when not fr
	return strconv.FormatFloat(sg, 'f', 3, 64)
}

func format_ranges_ph(ph float64) string {
	return strconv.FormatFloat(ph, 'f', 2, 64)
}

func format_ranges_density(density float64) string {
	return strconv.FormatFloat(density, 'f', 3, 64)
}

func format_ranges_string_test(string_test float64) string {
	return strconv.FormatFloat(string_test, 'f', 0, 64)
}

func format_ranges_viscosity(viscosity float64) string {
	return strconv.FormatFloat(viscosity, 'f', 0, 64)
}
