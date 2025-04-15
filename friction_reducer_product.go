package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/samuel-jimenez/winc"
	"github.com/samuel-jimenez/winc/w32"
)

type FrictionReducerProduct struct {
	BaseProduct
	sg          float64
	string_test float64
	viscosity   float64
}

func (product FrictionReducerProduct) toProduct() Product {
	return Product{product.toBaseProduct(), sql.NullFloat64{product.sg, true}, sql.NullFloat64{0, false}, sql.NullFloat64{product.sg * LB_PER_GAL, true}, sql.NullFloat64{product.string_test, true}, sql.NullFloat64{product.viscosity, true}}
}

func newFrictionReducerProduct(base_product BaseProduct, sample_point string, viscosity_field *winc.Edit, mass_field *winc.Edit, string_field *winc.Edit) Product {

	base_product.sample_point = sample_point

	viscosity, _ := strconv.ParseFloat(strings.TrimSpace(viscosity_field.Text()), 64)
	mass, _ := strconv.ParseFloat(strings.TrimSpace(mass_field.Text()), 64)
	// if !err.Error(){fmt.Println("error",err)}
	string_test, _ := strconv.ParseFloat(strings.TrimSpace(string_field.Text()), 64)
	sg := mass / SAMPLE_VOLUME

	return FrictionReducerProduct{base_product, sg, string_test, viscosity}.toProduct()

}

func (product FrictionReducerProduct) check_data() bool {
	return true
}

// create table product_line (product_id integer not null primary key, product_name text);
// func show_fr(parent winc.Controller) {
func show_fr(parent winc.Controller, create_new_product_cb func() BaseProduct) {

	top_col := 10
	bottom_col := 320

	group_row := 25
	// 		lot_row := 45
	// 		sample_row := 70
	//
	// 		visual_row := 125
	// 		viscosity_row := 150
	// 		mass_row := 175
	// 		string_row := 200
	// group_row := 120

	submit_col := 40
	submit_button_width := 100
	submit_button_height := 40

	// visual_text := "Visual Inspection"
	// viscosity_text := "Viscosity"
	// mass_text := "Mass"
	// string_text := "String"
	// sample_text := "Sample Point"
	group_width := 300
	group_height := 170

	top_text := "Top"
	// bottom_text := "Bottom"
	bottom_text := "Btm"

	top_group_cb := show_fr_sample_group(parent, top_text, top_col, group_row, group_width, group_height)
	// show_fr_sample_group(parent, top_text, top_col, group_row, group_width, group_height)

	// bottom_group := show_fr_sample_group(parent, bottom_text, bottom_col, group_row, group_width, group_height, top_group)
	// show_fr_sample_group(parent, bottom_text, bottom_col, group_row, group_width, group_height, top_group)
	bottom_group_cb := show_fr_sample_group(parent, bottom_text, bottom_col, group_row, group_width, group_height)

	// string_field := show_edit(mainWindow, label_col, field_col, string_row, string_text)
	// mass_field := show_edit(mainWindow, label_col, field_col, mass_row, mass_text)
	// viscosity_field := show_edit(mainWindow, label_col, field_col, viscosity_row, viscosity_text)
	// visual_field := show_checkbox(mainWindow, label_col, field_col, visual_row, visual_text)
	// sample_field := show_edit(mainWindow, label_col, field_col, sample_row, sample_text)

	// lot_field :=
	// show_edit(mainWindow, label_col, field_col, lot_row, lot_text)
	// // product_field :=
	// show_edit(mainWindow, label_col, field_col, product_row, product_text)

	//

	// 	product_row := 20
	// product_text := "Product"
	// product_field := show_edit(mainWindow, label_col, field_col, product_row, product_text)

	submit_button := winc.NewPushButton(parent)

	submit_button.SetText("Submit")
	submit_button.SetPos(submit_col, SUBMIT_ROW) // (x, y)
	// submit_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	submit_button.SetSize(submit_button_width, submit_button_height) // (width, height)
	submit_button.OnClick().Bind(func(e *winc.Event) {
		base_product := create_new_product_cb()
		top_product := top_group_cb(base_product)
		bottom_product := bottom_group_cb(base_product)
		fmt.Println("top", top_product)
		fmt.Println("btm", bottom_product)
		if top_product.check_data() {
			fmt.Println("data", top_product)
			top_product.save()
			top_product.print()

		}
		if bottom_product.check_data() {
			fmt.Println("data", bottom_product)
			bottom_product.save()
			bottom_product.print()

		}
	})

	top_button_col := 150
	top_button_row := SUBMIT_ROW
	top_button_width := 100
	top_button_height := 40
	top_button := winc.NewPushButton(parent)

	top_button.SetText("Accept Top")
	top_button.SetPos(top_button_col, top_button_row) // (x, y)
	// top_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	top_button.SetSize(top_button_width, top_button_height) // (width, height)
	top_button.OnClick().Bind(func(e *winc.Event) {
		base_product := create_new_product_cb()
		top_product := top_group_cb(base_product)
		top_product.sample_point = ""
		if top_product.check_data() {
			fmt.Println("data", top_product)
			top_product.print()
		}

	})

	btm_button_col := 250
	btm_button_row := SUBMIT_ROW
	btm_button_width := 100
	btm_button_height := 40
	btm_button := winc.NewPushButton(parent)

	btm_button.SetText("Accept Btm")
	btm_button.SetPos(btm_button_col, btm_button_row) // (x, y)
	// btm_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	btm_button.SetSize(btm_button_width, btm_button_height) // (width, height)
	btm_button.OnClick().Bind(func(e *winc.Event) {
		base_product := create_new_product_cb()
		bottom_product := bottom_group_cb(base_product)
		bottom_product.sample_point = ""
		if bottom_product.check_data() {
			fmt.Println("data", bottom_product)
			bottom_product.print()
		}
	})

}

// func show_fr_sample_group(parent winc.Controller, sample_point string, x_pos, y_pos, group_width, group_height int) winc.Controller {

func show_fr_sample_group(parent winc.Controller, sample_point string, x_pos, y_pos, group_width, group_height int) func(base_product BaseProduct) Product {

	// func show_fr_sample_group(parent winc.Controller, sample_point string, x_pos, y_pos, group_width, group_height int, after winc.Controller) winc.Controller {

	sample_group := winc.NewPanel(parent)
	sample_group.SetAndClearStyleBits(w32.WS_TABSTOP, 0)
	sample_group.SetPos(x_pos, y_pos)
	sample_group.SetSize(group_width, group_height)
	sample_group.SetText(sample_point)

	bottom_group := winc.NewGroupBox(parent)
	bottom_group.SetPos(x_pos-5, y_pos-5)
	bottom_group.SetSize(group_width+10, group_height+10)
	bottom_group.SetText(sample_point)

	return show_fr_sample(sample_group, sample_point)
	// return sample_group

}

func show_fr_sample(parent winc.Controller, sample_point string) func(base_product BaseProduct) Product {
	label_col := 10
	field_col := 120

	visual_row := 25
	viscosity_row := 50
	mass_row := 75
	string_row := 100

	// group_row := 120

	visual_text := "Visual Inspection"
	viscosity_text := "Viscosity"
	mass_text := "Mass"
	string_text := "String"

	visual_field := show_checkbox(parent, label_col, field_col, visual_row, visual_text)

	viscosity_field := show_edit(parent, label_col, field_col, viscosity_row, viscosity_text)
	mass_field := show_mass_sg(parent, label_col, field_col, mass_row, mass_text)
	string_field := show_edit(parent, label_col, field_col, string_row, string_text)

	// parent.Bind(w32.WM_COPYDATA, func(arg *EventArg) {
	// 	sender := arg.Sender()
	// 	if data, ok := arg.Data().(*gform.RawMsg); ok {
	// 		println(data.Hwnd, data.Msg, data.WParam, data.LParam)
	// 	}
	// }

	visual_field.SetFocus()

	return func(base_product BaseProduct) Product {
		base_product.visual = visual_field.Checked()
		return newFrictionReducerProduct(base_product, sample_point, viscosity_field, mass_field, string_field)

	}

}
