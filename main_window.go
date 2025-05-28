package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/samuel-jimenez/windigo"
)

var (
	status_bar *windigo.StatusBar
)

func show_window() {

	log.Println("Process started")
	// DEBUG
	// log.Println(time.Now().UTC().UnixNano())

	window_width := 800
	window_height := 600

	top_panel_width := 750
	top_panel_height := 110

	hpanel_width := 650

	label_width := LABEL_WIDTH

	hpanel_margin := 10

	field_width := FIELD_WIDTH
	field_height := 20
	top_spacer_height := 20
	inter_spacer_width := 30
	inter_spacer_height := INTER_SPACER_HEIGHT

	button_margin := 5
	cam_button_width := 100
	cam_button_height := 100

	ranges_button_width := 80
	ranges_button_height := 20

	product_text := "Product"
	lot_text := "Lot Number"
	sample_text := "Sample Point"
	customer_text := "Customer Name"

	var qc_product QCProduct

	mainWindow := windigo.NewForm(nil)
	mainWindow.SetSize(window_width, window_height)
	mainWindow.SetText("QC Data Entry")

	keygrab := windigo.NewEdit(mainWindow)
	keygrab.Hide()

	keygrab.OnSetFocus().Bind(func(e *windigo.Event) {
		keygrab.SetText("{")
		keygrab.SelectText(1, 1)
	})

	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := windigo.NewAutoPanel(mainWindow)
	product_panel.SetSize(top_panel_width, top_panel_height)

	cam_button := windigo.NewPushButton(product_panel)
	//TODO array db_select_all_product

	prod_panel := windigo.NewAutoPanel(product_panel)
	prod_panel.SetSize(hpanel_width, field_height)

	product_field := show_combobox(prod_panel, label_width, field_width, field_height, product_text)
	customer_field := show_edit(prod_panel, label_width, field_width, field_height, customer_text)

	prod_panel.SetMarginLeft(hpanel_margin)
	prod_panel.SetMarginTop(top_spacer_height)
	prod_panel.Dock(product_field, windigo.Left)
	prod_panel.Dock(customer_field, windigo.Left)

	customer_field.SetMarginLeft(inter_spacer_width)
	//TODO insert like lot
	customer_field.OnKillFocus().Bind(func(e *windigo.Event) {
		customer_field.SetText(strings.ToUpper(strings.TrimSpace(customer_field.Text())))
		if customer_field.Text() != "" && product_field.Text() != "" {
			product_name_customer := customer_field.Text()
			upsert_product_name_customer(product_name_customer, qc_product.product_id)
			qc_product.Product_name_customer = select_product_name_customer(qc_product.product_id)
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

	lot_panel := windigo.NewAutoPanel(product_panel)
	lot_panel.SetSize(hpanel_width, field_height)

	lot_field := show_combobox(lot_panel, label_width, field_width, field_height, lot_text)
	sample_field := show_edit_with_lose_focus(lot_panel, label_width, field_width, field_height, sample_text, strings.ToUpper)

	lot_panel.SetMarginTop(inter_spacer_height)
	lot_panel.SetMarginLeft(hpanel_margin)
	sample_field.SetMarginLeft(inter_spacer_width)
	lot_panel.Dock(lot_field, windigo.Left)
	lot_panel.Dock(sample_field, windigo.Left)

	ranges_button := windigo.NewPushButton(product_panel)

	product_panel.Dock(cam_button, windigo.Right)
	product_panel.Dock(prod_panel, windigo.Top)
	product_panel.Dock(lot_panel, windigo.Top)
	product_panel.Dock(ranges_button, windigo.Left)

	ranges_button.SetText("Ranges")
	ranges_button.SetMarginsAll(button_margin)
	ranges_button.SetSize(ranges_button_width, ranges_button_height) // (width, height)
	ranges_button.OnClick().Bind(func(e *windigo.Event) {
		if qc_product.Product_type != "" {
			qc_product.show_ranges_window()
			log.Println("product_lot", qc_product)
		}
	})

	tabs := windigo.NewTabView(mainWindow)
	tab_wb := tabs.AddAutoPanel("Water Based")
	tab_oil := tabs.AddAutoPanel("Oil Based")
	tab_fr := tabs.AddAutoPanel("Friction Reducer")

	new_product_cb := func() BaseProduct {
		qc_product.Sample_point = sample_field.Text()
		log.Println("product_field new_product_cb", qc_product.toBaseProduct())
		return qc_product.toBaseProduct()
	}

	update_water_based := show_water_based(tab_wb, qc_product, new_product_cb)
	update_oil_based := show_oil_based(tab_oil, qc_product, new_product_cb)
	update_fr := show_fr(tab_fr, qc_product, new_product_cb)

	qc_product.Update = func() {
		log.Println("update new_product_cb", qc_product)
		update_water_based(qc_product)
		update_oil_based(qc_product)
		update_fr(qc_product)
	}

	lot_field.OnKillFocus().Bind(func(e *windigo.Event) {
		lot_field.SetText(strings.ToUpper(strings.TrimSpace(lot_field.Text())))
		if lot_field.Text() != "" && product_field.Text() != "" {
			qc_product.Lot_number = lot_field.Text()
			qc_product.insel_lot_self()
			mainWindow.SetText(lot_field.Text())

		}

	})

	product_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		mainWindow.SetText(lot_field.GetSelectedItem())
	})

	product_field_pop_data := func(str string) {
		log.Println("product_field_pop_data product_id", qc_product.product_id)

		// if product_lot.product_id != product_lot.insel_product_id(str) {
		old_product_id := qc_product.product_id
		qc_product.Product_type = str
		qc_product.insel_product_self()
		if qc_product.product_id != old_product_id {
			log.Println("product_field_pop_data product_id changes", qc_product.product_id)

			qc_product.reset()
			qc_product.select_product_details()
			log.Println("product_field_pop_data select_product_details", qc_product)
			qc_product.Update()

			customer_field.SetText(qc_product.Product_name_customer)
			if qc_product.product_type.Valid {

				tabs.SetCurrent(qc_product.product_type.toIndex())
			}

			// TODO lot
			rows, err := db_select_lot_info.Query(qc_product.product_id)
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
			log.Println("product_field started", qc_product)

		}

	}

	qr_pop_data := func(product QRJson) {
		product_field_text_pop_data(product.Product_type)
		lot_field.SetText(strings.ToUpper(strings.TrimSpace(product.Lot_number)))
		lot_field.OnKillFocus().Fire(nil)
	}

	mainWindow.AddShortcut(windigo.Shortcut{windigo.ModShift, windigo.KeyOEM4}, // {
		func() bool {
			keygrab.SetText("")
			keygrab.SetFocus()
			return true
		})

	mainWindow.AddShortcut(windigo.Shortcut{windigo.ModShift, windigo.KeyOEM6}, // }
		func() bool {
			var product QRJson

			qr_json := keygrab.Text() + "}"
			log.Println("ReadFromScanner: ", qr_json)
			err := json.Unmarshal([]byte(qr_json), &product)
			if err == nil {
				qr_pop_data(product)
			}
			return true
		})

	cam_button.SetText("Scan")
	cam_button.SetMarginsAll(button_margin)
	cam_button.SetSize(cam_button_width, cam_button_height)
	cam_button.OnClick().Bind(func(e *windigo.Event) {
		close(qr_done)
		qr_sync_waitgroup.Wait()

		qr_done = make(chan bool)
		qr_sync_waitgroup.Add(1)
		go do_read_qr(qr_pop_data, qr_done)

	})

	product_field.OnSelectedChange().Bind(func(e *windigo.Event) {

		log.Println("product_field OnSelectedChange GetSelectedItem", product_field.GetSelectedItem())
		product_field_pop_data(product_field.GetSelectedItem())
	})

	product_field.OnKillFocus().Bind(func(e *windigo.Event) {
		product_field_text_pop_data(product_field.Text())
	})

	dock.Dock(product_panel, windigo.Top)  // tabs should prefer docking at the top
	dock.Dock(tabs, windigo.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), windigo.Fill) // tab panels dock just below tabs and fill area

	status_bar = windigo.NewStatusBar(mainWindow)
	mainWindow.SetStatusBar(status_bar)

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)
	mainWindow.RunMainLoop() // Must call to start event loop.
}

func wndOnClose(arg *windigo.Event) {
	okPen.Dispose()
	erroredPen.Dispose()
	windigo.Exit()
}
