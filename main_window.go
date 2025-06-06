package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"

	QR "github.com/samuel-jimenez/qc_data_entry/qr"
	"github.com/samuel-jimenez/windigo"
)

var (
	status_bar *windigo.StatusBar
)

func show_window() {

	log.Println("Info: Process started")
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

	container_field_width := 150

	reprint_button_width := 80
	reprint_button_margin_l := 150
	reprint_button_margins := button_margin

	product_text := "Product"
	lot_text := "Lot Number"
	sample_text := "Sample Point"
	customer_text := "Customer Name"

	var qc_product QCProduct
	qc_product.product_id = INVALID_ID
	qc_product.lot_id = DEFAULT_LOT_ID

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
	sample_field := show_combobox(lot_panel, label_width, field_width, field_height, sample_text)

	lot_panel.SetMarginTop(inter_spacer_height)
	lot_panel.SetMarginLeft(hpanel_margin)
	sample_field.SetMarginLeft(inter_spacer_width)
	lot_panel.Dock(lot_field, windigo.Left)
	lot_panel.Dock(sample_field, windigo.Left)

	ranges_button := windigo.NewPushButton(product_panel)
	ranges_button.SetText("Ranges")
	ranges_button.SetMarginsAll(button_margin)
	ranges_button.SetSize(ranges_button_width, OFF_AXIS)

	container_field := BuildNewDiscreteView(product_panel, 50, 50, "Container Type", qc_product.container_type, []string{"Tote", "Railcar"})
	container_field.SetSize(container_field_width, OFF_AXIS)

	reprint_button := windigo.NewPushButton(product_panel)
	reprint_button.SetText("Reprint")
	reprint_button.SetMarginsAll(reprint_button_margins)
	reprint_button.SetMarginLeft(reprint_button_margin_l)
	reprint_button.SetSize(reprint_button_width, OFF_AXIS)

	reprint_sample_button := windigo.NewPushButton(product_panel)
	reprint_sample_button.SetText("Reprint Sample")
	reprint_sample_button.SetMarginsAll(reprint_button_margins)
	reprint_sample_button.SetSize(reprint_button_width, OFF_AXIS)

	product_panel.Dock(prod_panel, windigo.Top)
	product_panel.Dock(lot_panel, windigo.Top)
	product_panel.Dock(ranges_button, windigo.Left)
	product_panel.Dock(container_field, windigo.Left)
	product_panel.Dock(reprint_button, windigo.Left)
	product_panel.Dock(reprint_sample_button, windigo.Left)

	tabs := windigo.NewTabView(mainWindow)
	tab_wb := tabs.AddAutoPanel("Water Based")
	tab_oil := tabs.AddAutoPanel("Oil Based")
	tab_fr := tabs.AddAutoPanel("Friction Reducer")

	dock.Dock(product_panel, windigo.Top)
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
		log.Println("Debug: product_field new_product_cb", qc_product.toBaseProduct())
		return qc_product.toBaseProduct()
	}

	update_water_based := show_water_based(tab_wb, qc_product, new_product_cb)
	update_oil_based := show_oil_based(tab_oil, qc_product, new_product_cb)
	fr_panel := show_fr(tab_fr, qc_product, new_product_cb)

	qc_product.Update = func() {
		container_field.Update(qc_product.container_type)
		log.Println("Debug: update new_product_cb", qc_product)
		update_water_based(qc_product)
		update_oil_based(qc_product)
		fr_panel.Update(qc_product)
		fr_panel.ChangeContainer(qc_product)
	}

	product_field_pop_data := func(str string) {
		log.Println("Warn: Debug: product_field_pop_data product_id", qc_product.product_id)
		log.Println("Warn: Debug: product_field_pop_data lot_id", qc_product.lot_id)

		// if product_lot.product_id != product_lot.insel_product_id(str) {
		old_product_id := qc_product.product_id
		qc_product.Product_type = str
		qc_product.insel_product_self()
		if qc_product.product_id != old_product_id {
			mainWindow.SetText(qc_product.Product_type)
			qc_product.reset()
			qc_product.lot_id = DEFAULT_LOT_ID

			qc_product.select_product_details()
			qc_product.Update()

			if qc_product.product_type.Valid {
				tabs.SetCurrent(qc_product.product_type.toIndex())
			}

			fill_combobox_from_query(lot_field, db_select_lot_info, qc_product.product_id)
			fill_combobox_from_query(customer_field, db_select_product_customer_info, qc_product.product_id)
			fill_combobox_from_query(sample_field, db_select_sample_points, qc_product.product_id)

			lot_field.OnKillFocus().Fire(nil)
			qc_product.Product_name_customer = customer_field.Text()
			qc_product.Sample_point = sample_field.Text()

		}
	}
	product_field_text_pop_data := func(str string) {
		formatted_text := strings.ToUpper(strings.TrimSpace(str))
		product_field.SetText(formatted_text)
		if product_field.Text() != "" {
			product_field_pop_data(product_field.Text())
			log.Println("Debug: product_field_text_pop_data", qc_product)
		} else {
			qc_product.product_id = INVALID_ID
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
		if formatted_text != "" && qc_product.product_id != INVALID_ID {
			lot_field_pop_data(lot_field.Text())
		} else {
			qc_product.Lot_number = ""
			qc_product.lot_id = DEFAULT_LOT_ID
		}
	}

	customer_field_pop_data := func(str string) {
		qc_product.Product_name_customer = str
	}
	customer_field_text_pop_data := func(str string) {
		formatted_text := strings.ToUpper(strings.TrimSpace(str))
		customer_field.SetText(formatted_text)

		old_product_name := qc_product.Product_name_customer
		customer_field_pop_data(formatted_text)
		if qc_product.Product_name_customer != old_product_name && formatted_text != "" && qc_product.product_id != INVALID_ID {
			insert_product_name_customer(qc_product.Product_name_customer, qc_product.product_id)
		}
	}

	sample_field_pop_data := func(str string) {
		qc_product.Sample_point = str
	}
	sample_field_text_pop_data := func(str string) {
		formatted_text := strings.ToUpper(strings.TrimSpace(str))
		sample_field.SetText(formatted_text)
		sample_field_pop_data(formatted_text)
	}

	qr_pop_data := func(product QR.QRJson) {
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

	customer_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		customer_field_pop_data(customer_field.GetSelectedItem())
	})
	customer_field.OnKillFocus().Bind(func(e *windigo.Event) {
		customer_field_text_pop_data(customer_field.Text())
	})

	sample_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		sample_field_pop_data(sample_field.GetSelectedItem())
	})
	sample_field.OnKillFocus().Bind(func(e *windigo.Event) {
		sample_field_text_pop_data(sample_field.Text())
	})

	ranges_button.OnClick().Bind(func(e *windigo.Event) {
		if qc_product.Product_type != "" {
			qc_product.show_ranges_window()
			log.Println("debug: ranges_button product_lot", qc_product)
		}
	})

	reprint_button.OnClick().Bind(func(e *windigo.Event) {
		if qc_product.Lot_number != "" {
			log.Println("debug: reprint_button")
			qc_product.reprint()
		}
	})

	reprint_sample_button.OnClick().Bind(func(e *windigo.Event) {
		if qc_product.Lot_number != "" {
			log.Println("debug: Notice: reprint_sample_button")
			qc_product.reprint_sample()
			log.Println("debug: Notice: qc_product", qc_product)

			// qc_product.toBaseProduct()
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
			var product QR.QRJson

			qr_json := keygrab.Text() + "}"
			log.Println("debug: ReadFromScanner: ", qr_json)
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
