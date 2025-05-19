package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/samuel-jimenez/windigo"
)

type FrictionReducerProduct struct {
	BaseProduct
	sg          float64
	string_test float64
	viscosity   float64
}

func (product FrictionReducerProduct) toProduct() Product {
	return Product{product.toBaseProduct(), NewNullFloat64(product.sg, true), NewNullFloat64(0, false), NewNullFloat64(product.sg*LB_PER_GAL, true), NewNullFloat64(product.string_test, true), NewNullFloat64(product.viscosity, true)}
}

func newFrictionReducerProduct(base_product BaseProduct, sample_point string, viscosity_field windigo.LabeledEdit, mass_field windigo.LabeledEdit, string_field windigo.LabeledEdit) Product {

	base_product.Sample_point = sample_point

	viscosity, _ := strconv.ParseFloat(strings.TrimSpace(viscosity_field.Text()), 64)
	string_test, _ := strconv.ParseFloat(strings.TrimSpace(string_field.Text()), 64)
	sg := sg_from_mass(mass_field)

	return FrictionReducerProduct{base_product, sg, string_test, viscosity}.toProduct()

}

func (product FrictionReducerProduct) check_data() bool {
	return true
}

type FrictionReducerProductView struct {
	windigo.AutoPanel
	Get   func(base_product BaseProduct) Product
	Clear func()
}

func BuildNewFrictionReducerProductView(parent windigo.AutoPanel, sample_point string, group_width, group_height int) FrictionReducerProductView {

	label_width := 110
	field_width := 200
	field_height := 22

	top_spacer_height := 25
	top_spacer_width := 10
	inter_spacer_height := 5

	visual_text := "Visual Inspection"
	viscosity_text := "Viscosity"
	mass_text := "Mass"
	string_text := "String"

	group_panel := windigo.NewGroupAutoPanel(parent)
	group_panel.SetSize(group_width, group_height)
	group_panel.SetText(sample_point)
	group_panel.SetPaddings(top_spacer_height, inter_spacer_height, inter_spacer_height, top_spacer_width)

	// visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)

	visual_field := windigo.NewCheckBox(group_panel)
	visual_field.SetText(visual_text)

	viscosity_field := show_edit(group_panel, label_width, field_width, field_height, viscosity_text)
	viscosity_field.SetMarginTop(inter_spacer_height)

	mass_field := show_mass_sg(group_panel, label_width, field_width, field_height, mass_text)
	mass_field.SetMarginTop(inter_spacer_height)

	string_field := show_edit(group_panel, label_width, field_width, field_height, string_text)
	string_field.SetMarginTop(inter_spacer_height)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(viscosity_field, windigo.Top)
	group_panel.Dock(mass_field, windigo.Top)
	group_panel.Dock(string_field, windigo.Top)

	get := func(base_product BaseProduct) Product {
		base_product.Visual = visual_field.Checked()
		return newFrictionReducerProduct(base_product, sample_point, viscosity_field, mass_field, string_field)

	}
	clear := func() {
		visual_field.SetChecked(false)
		viscosity_field.SetText("")
		mass_field.SetText("")
		mass_field.OnChange().Fire(nil)
		string_field.SetText("")

	}

	return FrictionReducerProductView{group_panel, get, clear}

}

// TODO
func check_dual_data(top_product, bottom_product Product) {
	if top_product.check_data() {
		log.Println("data", top_product)
		top_product.save()
		top_product.output()

	}
	if bottom_product.check_data() {
		log.Println("data", bottom_product)
		bottom_product.save()
		bottom_product.output()
	}
}

// create table product_line (product_id integer not null primary key, product_name text);
func show_fr(parent windigo.AutoPanel, create_new_product_cb func() BaseProduct) {

	bottom_spacer_height := BUTTON_SPACER_HEIGHT
	group_width := 300
	group_height := 170
	group_margin := 5

	button_width := 100
	button_height := 40
	// 	200
	// 50

	top_text := "Top"
	// bottom_text := "Bottom"
	bottom_text := "Btm"

	top_group := BuildNewFrictionReducerProductView(parent, top_text, group_width, group_height)
	bottom_group := BuildNewFrictionReducerProductView(parent, bottom_text, group_width, group_height)
	top_group.SetMargins(group_margin, 0, 0, group_margin)
	bottom_group.SetMarginsVH(group_margin, 0)

	submit_cb := func() {
		base_product := create_new_product_cb()
		top_product := top_group.Get(base_product)
		bottom_product := bottom_group.Get(base_product)
		log.Println("top", top_product)
		log.Println("btm", bottom_product)
		check_dual_data(top_product, bottom_product)
	}

	clear_cb := func() {
		top_group.Clear()
		bottom_group.Clear()
	}

	top_cb := func() {
		base_product := create_new_product_cb()

		top_product := top_group.Get(base_product)
		if top_product.check_data() {
			log.Println("data", top_product)
			top_product.output_sample()
		}
	}

	btm_cb := func() {

		base_product := create_new_product_cb()

		bottom_product := bottom_group.Get(base_product)
		if bottom_product.check_data() {
			log.Println("data", bottom_product)
			bottom_product.output_sample()
		}
	}

	// parent.Dock(ph_field, windigo.Top)

	button_dock := build_marginal_button_dock(parent, button_width, button_height, []string{"Submit", "Clear", "Accept Top", "Accept Btm"}, []int{40, 0, 10, 0}, []func(){submit_cb, clear_cb, top_cb, btm_cb})

	// button_dock.SetPaddingsAll(10)
	//
	// // button_dock.SetPaddingLeft(10)
	// button_dock.SetPaddingBottom(100)

	// button_dock.SetMargins(10, 0, 185, 0)
	// button_dock.SetMargins(10, 0, 100, 40)
	// button_dock.SetMarginsAll(10)

	// button_dock.SetMarginLeft(10)
	button_dock.SetMarginBtm(bottom_spacer_height)

	parent.Dock(button_dock, windigo.Bottom)
	parent.Dock(top_group, windigo.Left)
	parent.Dock(bottom_group, windigo.Left)

}
