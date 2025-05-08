package main

import (
	"log"
	"strings"

	"github.com/samuel-jimenez/windigo"
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

	mainWindow := windigo.NewForm(nil)
	mainWindow.SetSize(800, 600) // (width, height)
	mainWindow.SetText("QC Data Entry")

	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := windigo.NewPanel(mainWindow)
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

	var product_lot QCProduct

	cam_button_col := 675
	cam_button_row := 5
	cam_button_width := 100
	cam_button_height := 100

	ranges_button_col := 0
	ranges_button_row := 80
	ranges_button_width := 80
	ranges_button_height := 20

	cam_button := windigo.NewPushButton(product_panel)
	//TODO array db_select_all_product
	product_field := show_combobox(product_panel, label_col_0, field_col_0, product_row, product_text)
	customer_field := show_edit(product_panel, label_col_1, field_col_1, customer_row, customer_text)
	customer_field.OnKillFocus().Bind(func(e *windigo.Event) {
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
	sample_field := show_edit_with_lose_focus(product_panel, label_col_1, field_col_1, sample_row, sample_text, strings.ToUpper)

	ranges_button := windigo.NewPushButton(product_panel)
	ranges_button.SetText("Ranges")
	ranges_button.SetPos(ranges_button_col, ranges_button_row) // (x, y)
	// ranges_button.SetPosAfter(submit_col, submit_row, bottom_group)  // (x, y)
	ranges_button.SetSize(ranges_button_width, ranges_button_height) // (width, height)
	ranges_button.OnClick().Bind(func(e *windigo.Event) {
		if product_lot.Product_type != "" {
			product_lot.show_ranges_window()
			log.Println("product_lot", product_lot)
		}
	})

	tabs := windigo.NewTabView(mainWindow)
	// tabs.SetPos(20, 20)
	// tabs.SetSize(750, 500)
	tab_wb := tabs.AddPanel("Water Based")
	tab_oil := tabs.AddPanel("Oil Based")
	tab_fr := tabs.AddPanel("Friction Reducer")

	lot_field.OnKillFocus().Bind(func(e *windigo.Event) {
		lot_field.SetText(strings.ToUpper(strings.TrimSpace(lot_field.Text())))
		if lot_field.Text() != "" && product_field.Text() != "" {
			product_lot.Lot_number = lot_field.Text()
			product_lot.insel_lot_self()
		}
	})

	product_field_pop_data := func(str string) {
		log.Println("product_field_pop_data product_id", product_lot.product_id)

		// if product_lot.product_id != product_lot.insel_product_id(str) {
		old_product_id := product_lot.product_id
		product_lot.Product_type = str
		product_lot.insel_product_self()
		if product_lot.product_id != old_product_id {
			log.Println("product_field_pop_data product_id changes", product_lot.product_id)

			product_lot.reset()
			product_lot.select_product_details()
			log.Println("product_field_pop_data select_product_details", product_lot)

			customer_field.SetText(product_lot.Product_name_customer)
			if product_lot.product_type.Valid {

				tabs.SetCurrent(product_lot.product_type.toIndex())
			}

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
	cam_button.OnClick().Bind(func(e *windigo.Event) {
		close(qr_done)
		qr_done = make(chan bool)
		qr_sync_waitgroup.Add(1)
		go do_read_qr(qr_pop_data, qr_done)

	})

	product_field.OnSelectedChange().Bind(func(e *windigo.Event) {

		log.Println("product_field OnSelectedChange GetSelectedItem", product_field.GetSelectedItem())

		log.Println("product_field OnSelectedChange Text", product_field.Text())
		product_field_pop_data(product_field.GetSelectedItem())
	})

	product_field.OnKillFocus().Bind(func(e *windigo.Event) {
		product_field_text_pop_data(product_field.Text())
	})

	new_product_cb := func() BaseProduct {
		// return newProduct_0(product_field, lot_field).copy_ids(product_lot)
		log.Println("product_field new_product_cb", product_lot)

		base_product := NewBaseProduct(product_field, lot_field, sample_field)

		log.Println("product_field new_product_cb copy_ids", product_lot.toBaseProduct())
		log.Println("product_field base_product copy_ids", base_product)

		base_product.copy_ids(product_lot.toBaseProduct())
		log.Println("base_product new_product_cb", base_product)

		return base_product
	}

	show_water_based(tab_wb, new_product_cb)
	show_oil_based(tab_oil, new_product_cb)
	show_fr(tab_fr, new_product_cb)

	// dock.Dock(quux, windigo.Top)        // toolbars always dock to the top
	dock.Dock(product_panel, windigo.Top)  // tabs should prefer docking at the top
	dock.Dock(tabs, windigo.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), windigo.Fill) // tab panels dock just below tabs and fill area

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)
	mainWindow.RunMainLoop() // Must call to start event loop.
}

func wndOnClose(arg *windigo.Event) {
	windigo.Exit()
}
