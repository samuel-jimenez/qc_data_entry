package main

import (
	"log"

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

func newFrictionReducerProduct(base_product BaseProduct, viscosity_field windigo.LabeledEdit, mass_field MassDataView, string_field windigo.LabeledEdit) Product {

	viscosity := parse_field(viscosity_field)
	string_test := parse_field(string_field)
	sg := sg_from_mass(parse_field(mass_field))

	return FrictionReducerProduct{base_product, sg, string_test, viscosity}.toProduct()

}

func (product FrictionReducerProduct) check_data() bool {
	return true
}

type FrictionReducerProductView struct {
	windigo.AutoPanel
	Get   func(base_product BaseProduct, replace_sample_point bool) Product
	Clear func()
}

func BuildNewFrictionReducerProductView(parent windigo.AutoPanel, sample_point string, group_width, group_height int, ranges_panel FrictionReducerProductRangesView) FrictionReducerProductView {

	label_width := LABEL_WIDTH
	field_width := DATA_FIELD_WIDTH
	field_height := FIELD_HEIGHT

	visual_text := "Visual Inspection"
	viscosity_text := "Viscosity"
	mass_text := "Mass"
	string_text := "String"

	group_panel := windigo.NewGroupAutoPanel(parent)
	group_panel.SetSize(group_width, group_height)
	group_panel.SetText(sample_point)
	group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, BTM_SPACER_WIDTH, BTM_SPACER_HEIGHT)

	// visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)

	visual_field := windigo.NewCheckBox(group_panel)
	visual_field.SetText(visual_text)

	viscosity_field := show_number_edit(group_panel, label_width, field_width, field_height, viscosity_text, ranges_panel.viscosity_field)

	mass_field := show_mass_sg(group_panel, label_width, field_width, field_height, mass_text, ranges_panel)

	string_field := show_number_edit(group_panel, label_width, field_width, field_height, string_text, ranges_panel.string_field)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(viscosity_field, windigo.Top)
	group_panel.Dock(mass_field, windigo.Top)
	group_panel.Dock(string_field, windigo.Top)

	get := func(base_product BaseProduct, replace_sample_point bool) Product {
		base_product.Visual = visual_field.Checked()
		if replace_sample_point {
			base_product.Sample_point = sample_point
		}
		return newFrictionReducerProduct(base_product, viscosity_field, mass_field, string_field)

	}
	clear := func() {
		visual_field.SetChecked(false)
		clear_field(viscosity_field)

		mass_field.Clear()

		clear_field(string_field)

		ranges_panel.Clear()

	}

	return FrictionReducerProductView{group_panel, get, clear}

}

type FrictionReducerProductRangesView struct {
	windigo.AutoPanel
	DerivedMassRangesView

	viscosity_field,
	// mass_field,
	// sg_field,
	// density_field,
	string_field *RangeROView

	Update func(qc_product QCProduct)
}

func (data_view FrictionReducerProductRangesView) Clear() {
	data_view.DerivedMassRangesView.Clear()
	data_view.viscosity_field.Clear()
	data_view.string_field.Clear()
}

func BuildNewFrictionReducerProductRangesView(parent windigo.AutoPanel, qc_product QCProduct, group_width, group_height int) FrictionReducerProductRangesView {

	visual_text := "Visual Inspection"
	viscosity_text := "Viscosity"
	mass_text := "Mass"
	string_text := "String"
	sg_text := "Specific Gravity"
	density_text := "Density"

	group_panel := windigo.NewAutoPanel(parent)
	group_panel.SetSize(group_width, group_height)
	group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, BTM_SPACER_WIDTH, BTM_SPACER_HEIGHT)

	visual_field := BuildNewProductAppearanceROView(group_panel, visual_text, qc_product.Appearance)

	viscosity_field := BuildNewRangeROView(group_panel, viscosity_text, qc_product.Viscosity, format_ranges_viscosity)

	string_field := BuildNewRangeROView(group_panel, string_text, qc_product.SG, format_ranges_string_test)

	mass_field := BuildNewRangeROViewMap(group_panel, mass_text, qc_product.SG, format_mass, mass_from_sg)

	sg_field := BuildNewRangeROView(group_panel, sg_text, qc_product.SG, format_ranges_sg)
	density_field := BuildNewRangeROView(group_panel, density_text, qc_product.Density, format_ranges_density)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(viscosity_field, windigo.Top)
	group_panel.Dock(mass_field, windigo.Top)
	group_panel.Dock(string_field, windigo.Top)
	group_panel.Dock(density_field, windigo.Bottom)
	group_panel.Dock(sg_field, windigo.Bottom)

	update := func(qc_product QCProduct) {
		visual_field.Update(qc_product.Appearance)
		viscosity_field.Update(qc_product.Viscosity)
		string_field.Update(qc_product.String_test)

		mass_field.Update(qc_product.SG)
		sg_field.Update(qc_product.SG)
		density_field.Update(qc_product.Density)
	}

	return FrictionReducerProductRangesView{group_panel,
		DerivedMassRangesView{&mass_field,
			&sg_field,
			&density_field},
		&viscosity_field,
		&string_field,
		update}

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

type FrictionReducerPanelView struct {
	Update          func(qc_product QCProduct)
	ChangeContainer func(qc_product QCProduct)
}

// create table product_line (product_id integer not null primary key, product_name text);
func show_fr(parent windigo.AutoPanel, qc_product QCProduct, create_new_product_cb func() BaseProduct) FrictionReducerPanelView {

	bottom_spacer_height := BUTTON_SPACER_HEIGHT

	group_width := GROUP_WIDTH
	group_height := GROUP_HEIGHT
	group_margin := GROUP_MARGIN

	button_width := BUTTON_WIDTH
	button_height := BUTTON_HEIGHT

	top_text := "Top"
	// bottom_text := "Bottom"
	bottom_text := "Btm"

	ranges_panel := BuildNewFrictionReducerProductRangesView(parent, qc_product, RANGE_WIDTH, group_height)
	ranges_panel.SetMarginTop(group_margin)

	top_group := BuildNewFrictionReducerProductView(parent, top_text, group_width, group_height, ranges_panel)
	bottom_group := BuildNewFrictionReducerProductView(parent, bottom_text, group_width, group_height, ranges_panel)
	top_group.SetMargins(group_margin, group_margin, 0, 0)
	bottom_group.SetMarginTop(group_margin)

	submit_cb := func() {
		base_product := create_new_product_cb()
		top_product := top_group.Get(base_product, true)
		bottom_product := bottom_group.Get(base_product, true)
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

		top_product := top_group.Get(base_product, true)
		if top_product.check_data() {
			log.Println("data", top_product)
			top_product.output_sample()
		}
	}

	btm_cb := func() {

		base_product := create_new_product_cb()

		bottom_product := bottom_group.Get(base_product, true)
		if bottom_product.check_data() {
			log.Println("data", bottom_product)
			bottom_product.output_sample()
		}
	}

	tote_cb := func() {
		base_product := create_new_product_cb()

		top_product := top_group.Get(base_product, false)
		if top_product.check_data() {
			log.Println("tote", top_product)
			top_product.save()
			top_product.output()

		}
	}

	button_dock_totes := build_marginal_button_dock(parent, button_width, button_height, []string{"Submit", "Clear"}, []int{40, 0}, []func(){tote_cb, clear_cb})
	button_dock_cars := build_marginal_button_dock(parent, button_width, button_height, []string{"Submit", "Clear", "Accept Top", "Accept Btm"}, []int{40, 0, 10, 0}, []func(){submit_cb, clear_cb, top_cb, btm_cb})
	button_dock_totes.Hide()

	button_dock_totes.SetMarginBtm(bottom_spacer_height)
	button_dock_cars.SetMarginBtm(bottom_spacer_height)

	parent.Dock(button_dock_totes, windigo.Bottom)
	parent.Dock(button_dock_cars, windigo.Bottom)
	parent.Dock(top_group, windigo.Left)
	parent.Dock(bottom_group, windigo.Left)
	parent.Dock(ranges_panel, windigo.Right)

	changeContainer := func(qc_product QCProduct) {
		if qc_product.container_type.Int32 == int32(CONTAINER_TOTE) {
			bottom_group.Hide()
			button_dock_cars.Hide()
			button_dock_totes.Show()
		} else {
			bottom_group.Show()
			button_dock_cars.Show()
			button_dock_totes.Hide()
		}
	}

	return FrictionReducerPanelView{ranges_panel.Update, changeContainer}

}
