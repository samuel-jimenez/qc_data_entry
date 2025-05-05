package main

import (
	"log"
	"strings"

	"github.com/samuel-jimenez/winc"
)

var (
	SUBMIT_ROW = 200
	SUBMIT_COL = 40
	CLEAR_COL  = 140
)

func show_window() {

	log.Println("Process started")
	// DEBUG
	// log.Println(time.Now().UTC().UnixNano())
	label_size := 110

	mainWindow := winc.NewForm(nil)
	mainWindow.SetSize(800, 600) // (width, height)
	mainWindow.SetText("QC Data Entry")

	dock := winc.NewSimpleDock(mainWindow)

	product_panel := winc.NewPanel(mainWindow)
	product_panel.SetSize(750, 105)

	label_col_0 := 10
	field_col_0 := label_col_0 + label_size

	label_col_1 := 350
	field_col_1 := label_col_1 + label_size

	product_row := 20
	lot_row := 45
	customer_row := product_row
	sample_row := lot_row

	product_text := "Product"
	lot_text := "Lot Number"
	sample_text := "Sample Point"
	customer_text := "Customer Name"

	var product_lot BaseProduct

	cam_button_col := 675
	cam_button_row := 5
	cam_button_width := 100
	cam_button_height := 100

	cam_button := winc.NewPushButton(product_panel)
	//TODO array db_select_all_product
	product_field := show_combobox(product_panel, label_col_0, field_col_0, product_row, product_text)
	customer_field := show_edit(product_panel, label_col_1, field_col_1, customer_row, customer_text)
	customer_field.OnKillFocus().Bind(func(e *winc.Event) {
		customer_field.SetText(strings.ToUpper(strings.TrimSpace(customer_field.Text())))
		if customer_field.Text() != "" && product_field.Text() != "" {
			product_name_customer := customer_field.Text()
			upsert_product_name_customer(product_name_customer, product_lot.product_id)
			product_lot.Product_name_customer = select_product_name_customer(product_lot.product_id)
			// product_lot.product_name_customer = product_name_customer

		}
	})

	//TODO extract
	rows, err := db_select_product_info.Query()
	if err != nil {
		log.Printf("%q: %s\n", err, "insel_lot_id")
		// return -1
	}
	for rows.Next() {
		var (
			id                   uint8
			internal_name        string
			product_moniker_name string
		)

		if error := rows.Scan(&id, &internal_name, &product_moniker_name); error == nil {
			// data[id] = internal_name
			product_field.AddItem(product_moniker_name + " " + internal_name)
		}
	}

	lot_field := show_combobox(product_panel, label_col_0, field_col_0, lot_row, lot_text)

	lot_field.OnKillFocus().Bind(func(e *winc.Event) {
		lot_field.SetText(strings.ToUpper(strings.TrimSpace(lot_field.Text())))
		if lot_field.Text() != "" && product_field.Text() != "" {
			product_lot.insel_lot_id(lot_field.Text())
		}
	})

	product_field_pop_data := func(str string) {
		log.Println("product_field_pop_data product_id", product_lot.product_id)

		// if product_lot.product_id != product_lot.insel_product_id(str) {
		old_product_id := product_lot.product_id
		product_lot.insel_product_id(str)
		if product_lot.product_id != old_product_id {
			log.Println("product_field_pop_data product_id changes", product_lot.product_id)
			product_lot.Product_name_customer = select_product_name_customer(product_lot.product_id)
			customer_field.SetText(product_lot.Product_name_customer)

			// TODO lot
			rows, err := db_select_lot_info.Query(product_lot.product_id)
			if err != nil {
				log.Printf("%q: %s\n", err, "insel_lot_id")
				// return -1
			}
			lot_field.DeleteAllItems()
			for rows.Next() {
				var (
					id       uint8
					lot_name string
				)

				if error := rows.Scan(&id, &lot_name); error == nil {
					// data[id] = value
					lot_field.AddItem(lot_name)
				}
			}

		}

	}

	product_field_text_pop_data := func(str string) {
		product_field.SetText(strings.ToUpper(strings.TrimSpace(str)))
		if product_field.Text() != "" {

			log.Println("product_field OnKillFocus Text", product_field.Text())
			product_field_pop_data(product_field.Text())
			log.Println("product_field started", product_lot)

		}

	}

	qr_pop_data := func(product QRJson) {

		lot_field.SetText(strings.ToUpper(strings.TrimSpace(product.Lot_number)))

		product_field_text_pop_data(product.Product_type)

	}

	cam_button.SetText("Scan")
	cam_button.SetPos(cam_button_col, cam_button_row) // (x, y)
	// cam_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	cam_button.SetSize(cam_button_width, cam_button_height) // (width, height)
	cam_button.OnClick().Bind(func(e *winc.Event) {
		close(qr_done)
		qr_done = make(chan bool)
		qr_sync_waitgroup.Add(1)
		go do_read_qr(qr_pop_data, qr_done)
	})

	product_field.OnSelectedChange().Bind(func(e *winc.Event) {

		log.Println("product_field OnSelectedChange GetSelectedItem", product_field.GetSelectedItem())

		log.Println("product_field OnSelectedChange Text", product_field.Text())
		product_field_pop_data(product_field.GetSelectedItem())
	})

	product_field.OnKillFocus().Bind(func(e *winc.Event) {
		product_field_text_pop_data(product_field.Text())
	})

	sample_field := show_edit_with_lose_focus(product_panel, label_col_1, field_col_1, sample_row, sample_text, strings.ToUpper)

	new_product_cb := func() BaseProduct {
		// return newProduct_0(product_field, lot_field).copy_ids(product_lot)
		log.Println("product_field new_product_cb", product_lot)

		// base_product := newProduct_0(product_field, lot_field)
		base_product := NewBaseProduct(product_field, lot_field, sample_field)

		base_product.copy_ids(product_lot)
		log.Println("base_product new_product_cb", base_product)

		return base_product
	}

	tabs := winc.NewTabView(mainWindow)
	// tabs.SetPos(20, 20)
	// tabs.SetSize(750, 500)
	tab_wb := tabs.AddPanel("Water Based")
	tab_oil := tabs.AddPanel("Oil Based")
	tab_fr := tabs.AddPanel("Friction Reducer")

	show_water_based(tab_wb, new_product_cb)
	show_oil_based(tab_oil, new_product_cb)
	show_fr(tab_fr, new_product_cb)

	// dock.Dock(quux, winc.Top)        // toolbars always dock to the top
	dock.Dock(product_panel, winc.Top)  // tabs should prefer docking at the top
	dock.Dock(tabs, winc.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), winc.Fill) // tab panels dock just below tabs and fill area

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)
	mainWindow.RunMainLoop() // Must call to start event loop.
}

func show_checkbox(parent winc.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *winc.CheckBox {
	checkbox_label := winc.NewLabel(parent)
	checkbox_label.SetPos(x_label_pos, y_pos)

	checkbox_label.SetText(field_text)

	checkbox_field := winc.NewCheckBox(parent)
	checkbox_field.SetText("")

	checkbox_field.SetPos(x_field_pos, y_pos)
	// visual_label.OnClick().Bind(func(e *winc.Event) {
	// 		visual_field.SetFocus()
	// })
	return checkbox_field
}

func show_combobox(parent winc.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *winc.ComboBox {
	// func show_combobox(parent winc.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *winc.Edit {
	combobox_label := winc.NewLabel(parent)
	combobox_label.SetPos(x_label_pos, y_pos)
	combobox_label.SetText(field_text)

	// combobox_field := combobox_label.NewEdit(mainWindow)
	// combobox_field := winc.NewEdit(parent)

	combobox_field := winc.NewComboBox(parent)

	combobox_field.SetPos(x_field_pos, y_pos)
	// Most Controls have default size unless SetSize is called.
	combobox_field.SetText("")
	// combobox_field.SetParent(combobox_label)
	// combobox_label.SetParent(combobox_field)

	// combobox_label.OnClick().Bind(func(e *winc.Event) {
	// 		combobox_field.SetFocus()
	// })
	combobox_field.OnKillFocus().Bind(func(e *winc.Event) {
		combobox_field.SetText(strings.TrimSpace(combobox_field.Text()))
	})

	return combobox_field
}

func show_edit(parent winc.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *winc.Edit {
	edit_label := winc.NewLabel(parent)
	edit_label.SetPos(x_label_pos, y_pos)
	edit_label.SetText(field_text)

	// edit_field := edit_label.NewEdit(mainWindow)
	edit_field := winc.NewEdit(parent)
	edit_field.SetPos(x_field_pos, y_pos)
	// Most Controls have default size unless SetSize is called.
	edit_field.SetText("")
	// edit_field.SetParent(edit_label)
	// edit_label.SetParent(edit_field)

	// edit_label.OnClick().Bind(func(e *winc.Event) {
	// 		edit_field.SetFocus()
	// })
	edit_field.OnKillFocus().Bind(func(e *winc.Event) {
		edit_field.SetText(strings.TrimSpace(edit_field.Text()))
	})

	return edit_field
}

func show_edit_with_lose_focus(parent winc.Controller, x_label_pos, x_field_pos, y_pos int, field_text string, focus_cb func(string) string) *winc.Edit {
	edit_label := winc.NewLabel(parent)
	edit_label.SetPos(x_label_pos, y_pos)
	edit_label.SetText(field_text)

	// edit_field := edit_label.NewEdit(mainWindow)
	edit_field := winc.NewEdit(parent)
	edit_field.SetPos(x_field_pos, y_pos)
	// Most Controls have default size unless SetSize is called.
	edit_field.SetText("")
	// edit_field.SetParent(edit_label)
	// edit_label.SetParent(edit_field)

	// edit_label.OnClick().Bind(func(e *winc.Event) {
	// 		edit_field.SetFocus()
	// })
	edit_field.OnKillFocus().Bind(func(e *winc.Event) {
		edit_field.SetText(focus_cb(strings.TrimSpace(edit_field.Text())))
	})

	return edit_field
}

func show_text(parent winc.Controller, x_label_pos, x_field_pos, y_pos, field_width int, field_text string, field_units string) *winc.Label {
	field_height := 20
	text_label := winc.NewLabel(parent)
	text_label.SetPos(x_label_pos, y_pos)
	text_label.SetText(field_text)

	// text_field := text_label.NewEdit(mainWindow)
	text_field := winc.NewLabel(parent)
	text_field.SetPos(x_field_pos, y_pos)
	text_field.SetSize(field_width, field_height)
	// Most Controls have default size unless  is called.
	text_field.SetText("0.000")

	text_units := winc.NewLabel(parent)
	text_units.SetPos(x_field_pos+field_width, y_pos)
	text_units.SetText(field_units)

	return text_field
}

func show_mass_sg(parent winc.Controller, x_label_pos, x_field_pos, y_pos int, field_text string) *winc.Edit {

	field_width := 50

	density_row := 125
	sg_row := 150

	sg_text := "Specific Gravity"
	density_text := "Density"

	sg_units := "g/mL"
	density_units := "lb/gal"

	mass_field := show_edit(parent, x_label_pos, x_field_pos, y_pos, field_text)

	sg_field := show_text(parent, x_label_pos, x_field_pos, sg_row, field_width, sg_text, sg_units)
	density_field := show_text(parent, x_label_pos, x_field_pos, density_row, field_width, density_text, density_units)

	mass_field.OnKillFocus().Bind(func(e *winc.Event) {
		mass_field.SetText(strings.TrimSpace(mass_field.Text()))
		sg := sg_from_mass(mass_field)
		density := density_from_sg(sg)
		sg_field.SetText(format_sg(sg))
		density_field.SetText(format_density(density))
	})

	return mass_field
}

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}
