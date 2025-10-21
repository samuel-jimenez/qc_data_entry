package fr_ui

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/windigo"
)

// RecipeHeader
type BlendHeader struct {
	*GUI.TextDock
}

func Header_units(title, units string) string {
	return title + "\n(" + units + ")"
}

func BlendHeader_from_new(parent windigo.Controller) *BlendHeader {
	view := new(BlendHeader)

	// 	TODO!!! RecipeHeader
	// SG_UNITS        = "g/mL"

	view.TextDock = GUI.NewTextDock(
		parent,
		"Component",
		Header_units("Weight", "lb"),
		"Component Lot",
		Header_units("Density", formats.DENSITY_UNITS),
		Header_units("Volume", "gal"),
		Header_units("Strap", "in"),
	)

	return view
}

func (view *BlendHeader) RefreshSize() {
	view.SetSize(GUI.OFF_AXIS, 2*GUI.PRODUCT_FIELD_HEIGHT)
	for i, label := range view.Labels {
		switch i {
		case 0:
			label.SetSize(GUI.SOURCES_LABEL_WIDTH, GUI.OFF_AXIS)
		case 2:
			label.SetSize(GUI.SOURCES_FIELD_WIDTH, GUI.OFF_AXIS)
		default:
			label.SetSize(GUI.DATA_SUBFIELD_WIDTH, GUI.OFF_AXIS)
		}
	}
}
