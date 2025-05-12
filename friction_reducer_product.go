package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/samuel-jimenez/windigo"
	"github.com/samuel-jimenez/windigo/w32"
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

func newFrictionReducerProduct(base_product BaseProduct, sample_point string, viscosity_field *windigo.Edit, mass_field *windigo.Edit, string_field *windigo.Edit) Product {

	base_product.Sample_point = sample_point

	viscosity, _ := strconv.ParseFloat(strings.TrimSpace(viscosity_field.Text()), 64)
	string_test, _ := strconv.ParseFloat(strings.TrimSpace(string_field.Text()), 64)
	sg := sg_from_mass(mass_field)

	return FrictionReducerProduct{base_product, sg, string_test, viscosity}.toProduct()

}

func (product FrictionReducerProduct) check_data() bool {
	return true
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
func show_fr(parent windigo.Controller, create_new_product_cb func() BaseProduct) {

	top_col := 10
	bottom_col := 320

	group_row := 25

	submit_col := SUBMIT_COL
	submit_button_width := 100
	submit_button_height := 40

	group_width := 300
	group_height := 170

	top_text := "Top"
	// bottom_text := "Bottom"
	bottom_text := "Btm"

	top_group_cb, top_group_clear := show_fr_sample_group(parent, top_text, top_col, group_row, group_width, group_height)

	bottom_group_cb, bottom_group_clear := show_fr_sample_group(parent, bottom_text, bottom_col, group_row, group_width, group_height)

	submit_button := windigo.NewPushButton(parent)

	submit_button.SetText("Submit")
	submit_button.SetPos(submit_col, SUBMIT_ROW) // (x, y)
	// submit_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	submit_button.SetSize(submit_button_width, submit_button_height) // (width, height)
	submit_button.OnClick().Bind(func(e *windigo.Event) {
		base_product := create_new_product_cb()
		top_product := top_group_cb(base_product)
		bottom_product := bottom_group_cb(base_product)
		log.Println("top", top_product)
		log.Println("btm", bottom_product)
		check_dual_data(top_product, bottom_product)

	})

	clear_button_col := CLEAR_COL
	clear_button_row := SUBMIT_ROW
	clear_button_width := 100
	clear_button_height := 40
	clear_button := windigo.NewPushButton(parent)

	clear_button.SetText("Clear")
	clear_button.SetPos(clear_button_col, clear_button_row) // (x, y)
	// clear_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	clear_button.SetSize(clear_button_width, clear_button_height) // (width, height)
	clear_button.OnClick().Bind(func(e *windigo.Event) {
		top_group_clear()
		bottom_group_clear()
	})

	top_button_col := 250
	top_button_row := SUBMIT_ROW
	top_button_width := 100
	top_button_height := 40
	top_button := windigo.NewPushButton(parent)

	top_button.SetText("Accept Top")
	top_button.SetPos(top_button_col, top_button_row) // (x, y)
	// top_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	top_button.SetSize(top_button_width, top_button_height) // (width, height)
	top_button.OnClick().Bind(func(e *windigo.Event) {
		base_product := create_new_product_cb()

		top_product := top_group_cb(base_product)
		if top_product.check_data() {
			log.Println("data", top_product)
			top_product.output_sample()
		}

	})

	btm_button_col := 350
	btm_button_row := SUBMIT_ROW
	btm_button_width := 100
	btm_button_height := 40
	btm_button := windigo.NewPushButton(parent)

	btm_button.SetText("Accept Btm")
	btm_button.SetPos(btm_button_col, btm_button_row) // (x, y)
	// btm_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	btm_button.SetSize(btm_button_width, btm_button_height) // (width, height)
	btm_button.OnClick().Bind(func(e *windigo.Event) {
		base_product := create_new_product_cb()

		bottom_product := bottom_group_cb(base_product)
		if bottom_product.check_data() {
			log.Println("data", bottom_product)
			bottom_product.output_sample()
		}
	})

}

func show_fr_sample_group(parent windigo.Controller, sample_point string, x_pos, y_pos, group_width, group_height int) (func(base_product BaseProduct) Product, func()) {

	sample_group := windigo.NewPanel(parent)
	sample_group.SetAndClearStyleBits(w32.WS_TABSTOP, 0)
	sample_group.SetPos(x_pos, y_pos)
	sample_group.SetSize(group_width, group_height)
	sample_group.SetText(sample_point)

	group_box := windigo.NewGroupBox(parent)
	group_box.SetPos(x_pos-5, y_pos-5)
	group_box.SetSize(group_width+10, group_height+10)
	group_box.SetText(sample_point)

	label_col := 10
	field_col := 120

	visual_row := 25
	viscosity_row := 50
	mass_row := 75
	string_row := 100

	visual_text := "Visual Inspection"
	viscosity_text := "Viscosity"
	mass_text := "Mass"
	string_text := "String"

	visual_field := show_checkbox(sample_group, label_col, field_col, visual_row, visual_text)

	viscosity_field := show_edit(sample_group, label_col, field_col, viscosity_row, viscosity_text)
	mass_field := show_mass_sg(sample_group, label_col, field_col, mass_row, mass_text)
	string_field := show_edit(sample_group, label_col, field_col, string_row, string_text)

	// parent.Bind(w32.WM_COPYDATA, func(arg *EventArg) {
	// 	sender := arg.Sender()
	// 	if data, ok := arg.Data().(*gform.RawMsg); ok {
	// 		println(data.Hwnd, data.Msg, data.WParam, data.LParam)
	// 	}
	// }

	visual_field.SetFocus()

	return func(base_product BaseProduct) Product {
			base_product.Visual = visual_field.Checked()
			return newFrictionReducerProduct(base_product, sample_point, viscosity_field, mass_field, string_field)

		}, func() {
			visual_field.SetChecked(false)
			viscosity_field.SetText("")
			mass_field.SetText("")
			mass_field.OnKillFocus().Fire(nil)
			string_field.SetText("")

		}

}
