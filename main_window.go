package main

import (
	"database/sql"
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

	window_width := 650
	window_height := 600

	top_panel_width := window_width
	top_panel_height := 110

	hpanel_width := top_panel_width

	label_width := LABEL_WIDTH

	hpanel_margin := 10

	field_width := PRODUCT_FIELD_WIDTH
	field_height := 20
	top_spacer_height := 20
	inter_spacer_width := 30
	inter_spacer_height := INTER_SPACER_HEIGHT

	button_margin := 5

	ranges_button_width := 80
	ranges_button_height := 20

	product_text := "Product"
	lot_text := "Lot Number"
	sample_text := "Sample Point"
	customer_text := "Customer Name"

	var qc_product QCProduct

	// build window
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetSize(window_width, window_height)
	mainWindow.SetText("QC Data Entry")

	keygrab := windigo.NewEdit(mainWindow)
	keygrab.Hide()

	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := windigo.NewAutoPanel(mainWindow)
	product_panel.SetSize(top_panel_width, top_panel_height)

	//TODO array db_select_all_product

	prod_panel := windigo.NewAutoPanel(product_panel)
	prod_panel.SetSize(hpanel_width, field_height)

	product_field := show_combobox(prod_panel, label_width, field_width, field_height, product_text)
	customer_field := show_combobox(prod_panel, label_width, field_width, field_height, customer_text)

	customer_field.SetMarginLeft(inter_spacer_width)
	prod_panel.SetMarginLeft(hpanel_margin)
	prod_panel.SetMarginTop(top_spacer_height)
	prod_panel.Dock(product_field, windigo.Left)
	prod_panel.Dock(customer_field, windigo.Left)

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
	ranges_button.SetText("Ranges")
	ranges_button.SetMarginsAll(button_margin)
	ranges_button.SetSize(ranges_button_width, ranges_button_height) // (width, height)

	container_field := BuildNewDiscreteView(product_panel, 50, 50, "Container Type", qc_product.container_type, []string{"Tote", "Railcar"})
	container_field.SetSize(150, OFF_AXIS)

	product_panel.Dock(prod_panel, windigo.Top)
	product_panel.Dock(lot_panel, windigo.Top)
	product_panel.Dock(ranges_button, windigo.Left)
	// product_panel.Dock(radio_dock, windigo.Top)
	product_panel.Dock(container_field, windigo.Left)

	tabs := windigo.NewTabView(mainWindow)
	tab_wb := tabs.AddAutoPanel("Water Based")
	tab_oil := tabs.AddAutoPanel("Oil Based")
	tab_fr := tabs.AddAutoPanel("Friction Reducer")

	dock.Dock(product_panel, windigo.Top)  // tabs should prefer docking at the top
	dock.Dock(tabs, windigo.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), windigo.Fill) // tab panels dock just below tabs and fill area

	status_bar = windigo.NewStatusBar(mainWindow)
	mainWindow.SetStatusBar(status_bar)

	// functionality

	rows, err := db_select_product_info.Query()
	fill_combobox_from_query_rows(product_field, rows, err, func(rows *sql.Rows) {
		var (
			id                   uint8
			internal_name        string
			product_moniker_name string
		)
		if err := rows.Scan(&id, &internal_name, &product_moniker_name); err == nil {
			// data[id] = internal_name
			product_field.AddItem(product_moniker_name + " " + internal_name)
		}
	})

	new_product_cb := func() BaseProduct {
		qc_product.Sample_point = sample_field.Text()
		log.Println("product_field new_product_cb", qc_product.toBaseProduct())
		return qc_product.toBaseProduct()
	}

	update_water_based := show_water_based(tab_wb, qc_product, new_product_cb)
	update_oil_based := show_oil_based(tab_oil, qc_product, new_product_cb)
	fr_panel := show_fr(tab_fr, qc_product, new_product_cb)

	qc_product.Update = func() {
		container_field.Update(qc_product.container_type)
		log.Println("update new_product_cb", qc_product)
		update_water_based(qc_product)
		update_oil_based(qc_product)
		fr_panel.Update(qc_product)
		fr_panel.ChangeContainer(qc_product)
	}

	product_field_pop_data := func(str string) {
		log.Println("product_field_pop_data product_id", qc_product.product_id)

		// if product_lot.product_id != product_lot.insel_product_id(str) {
		old_product_id := qc_product.product_id
		qc_product.Product_type = str
		qc_product.insel_product_self()
		if qc_product.product_id != old_product_id {
			mainWindow.SetText(qc_product.Product_type)
			log.Println("product_field_pop_data product_id changes", qc_product.product_id)

			qc_product.reset()
			qc_product.select_product_details()
			log.Println("product_field_pop_data select_product_details", qc_product)
			qc_product.Update()

			if qc_product.product_type.Valid {
				tabs.SetCurrent(qc_product.product_type.toIndex())
			}

			fill_combobox_from_query(customer_field, db_select_product_customer_info, qc_product.product_id)
			fill_combobox_from_query(lot_field, db_select_lot_info, qc_product.product_id)
		}
	}
	product_field_text_pop_data := func(str string) {
		formatted_text := strings.ToUpper(strings.TrimSpace(str))
		product_field.SetText(formatted_text)
		qc_product.product_id = INVALID_ID
		if product_field.Text() != "" {
			log.Println("product_field OnKillFocus Text", product_field.Text())
			product_field_pop_data(product_field.Text())
			log.Println("product_field started", qc_product)
		}
	}

	lot_field_pop_data := func(str string) {
		qc_product.Lot_number = str
		qc_product.insel_lot_self()
		mainWindow.SetText(str)
	}
	lot_field_text_pop_data := func(str string) {
		formatted_text := strings.ToUpper(strings.TrimSpace(str))
		lot_field.SetText(formatted_text)
		qc_product.lot_id = INVALID_ID
		log.Println("lot_field_text_pop_data product_id", qc_product.product_id)
		if lot_field.Text() != "" && qc_product.product_id != INVALID_ID {
			lot_field_pop_data(lot_field.Text())
		}
		log.Println("lot_field_text_pop_data lot_field", qc_product.lot_id)
	}

	qr_pop_data := func(product QRJson) {
		product_field_text_pop_data(product.Product_type)
		lot_field_text_pop_data(product.Lot_number)
	}

	product_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		product_field_pop_data(product_field.GetSelectedItem())
	})
	product_field.OnKillFocus().Bind(func(e *windigo.Event) {
		product_field_text_pop_data(product_field.Text())
	})

	lot_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		lot_field_pop_data(lot_field.GetSelectedItem())
	})
	lot_field.OnKillFocus().Bind(func(e *windigo.Event) {
		lot_field_text_pop_data(lot_field.Text())
	})

	customer_field.OnKillFocus().Bind(func(e *windigo.Event) {
		customer_field.SetText(strings.ToUpper(strings.TrimSpace(customer_field.Text())))
		qc_product.Product_name_customer = customer_field.Text()

		if qc_product.Product_name_customer != "" && qc_product.product_id != INVALID_ID {
			insert_product_name_customer(qc_product.Product_name_customer, qc_product.product_id)
		}
	})

	ranges_button.OnClick().Bind(func(e *windigo.Event) {
		if qc_product.Product_type != "" {
			qc_product.show_ranges_window()
			log.Println("product_lot", qc_product)
		}
	})

	container_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		qc_product.container_type = container_field.Get()
		fr_panel.ChangeContainer(qc_product)
	})

	// QR keyboard handling

	keygrab.OnSetFocus().Bind(func(e *windigo.Event) {
		keygrab.SetText("{")
		keygrab.SelectText(1, 1)
	})

	mainWindow.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModShift, Key: windigo.KeyOEM4}, // {
		func() bool {
			keygrab.SetText("")
			keygrab.SetFocus()
			return true
		})

	mainWindow.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModShift, Key: windigo.KeyOEM6}, // }
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
