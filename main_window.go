package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/view"
	"github.com/samuel-jimenez/qc_data_entry/QR"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/windigo"
)

func show_window() {

	log.Println("Info: Process started")
	// DEBUG
	// log.Println(time.Now().UTC().UnixNano())

	clock_width := 90
	clock_timer_width := 35
	clock_timer_offset_h := 10
	clock_timer_offset_v := 30
	clock_font_size := 25
	clock_font_style_flags := byte(0x0) //FontNormal

	var (
		top_panel_height,

		hpanel_width,

		hpanel_margin,

		top_spacer_height,
		top_subpanel_height,

		ranges_button_width,

		container_item_width,

		reprint_button_margin_l int
	)

	product_text := "Product"
	lot_text := "Lot Number"
	sample_text := "Sample Point"
	customer_text := "Customer Name"

	qc_product := product.NewQCProduct()
	qc_product.Product_id = DB.INVALID_ID
	qc_product.Lot_id = DB.DEFAULT_LOT_ID

	windigo.DefaultFont = windigo.NewFont("MS Shell Dlg 2", GUI.BASE_FONT_SIZE, windigo.FontNormal)

	// build window
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetText("QC Data Entry")

	keygrab := windigo.NewEdit(mainWindow)
	keygrab.Hide()

	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := windigo.NewAutoPanel(mainWindow)

	product_panel_0 := windigo.NewAutoPanel(product_panel)

	//TODO array db_select_all_product

	product_panel_0_0 := windigo.NewAutoPanel(product_panel_0)

	product_field := GUI.NewComboBox(product_panel_0_0, product_text)
	customer_field := GUI.NewComboBox(product_panel_0_0, customer_text)

	product_panel_0_1 := windigo.NewAutoPanel(product_panel_0)

	lot_field := GUI.NewComboBox(product_panel_0_1, lot_text)
	sample_field := GUI.NewComboBox(product_panel_0_1, sample_text)

	clock_panel := windigo.NewAutoPanel(product_panel_0)
	clock_panel.SetSize(clock_width, GUI.OFF_AXIS)
	clock_display_now := windigo.NewSizedLabeledLabel(clock_panel, clock_timer_width, GUI.OFF_AXIS, "")
	clock_display_future := windigo.NewSizedLabeledLabel(clock_panel, clock_timer_width, GUI.OFF_AXIS, "")
	clock_display_now.SetMarginTop(clock_timer_offset_v)
	clock_display_future.SetMarginLeft(clock_timer_offset_h)
	clock_display_future.SetMarginRight(clock_timer_offset_h)
	clock_font := windigo.NewFont(clock_display_now.Font().Family(), clock_font_size, clock_font_style_flags)

	clock_display_now.SetFont(clock_font)
	clock_display_future.SetFont(clock_font)

	clock_ticker := time.Tick(time.Second)

	ranges_button := windigo.NewPushButton(product_panel)
	ranges_button.SetText("Ranges")

	container_field := product.BuildNewDiscreteView(product_panel, "Container Type", qc_product.Container_type, []string{"Tote", "Railcar"})

	reprint_button := windigo.NewPushButton(product_panel)
	reprint_button.SetText("Reprint")

	reprint_sample_button := windigo.NewPushButton(product_panel)
	reprint_sample_button.SetText("Reprint Sample")

	tabs := windigo.NewTabView(mainWindow)
	tab_wb := tabs.AddAutoPanel("Water Based")
	tab_oil := tabs.AddAutoPanel("Oil Based")
	tab_fr := tabs.AddAutoPanel("Friction Reducer")

	// Dock

	product_panel_0_0.Dock(product_field, windigo.Left)
	product_panel_0_0.Dock(customer_field, windigo.Left)

	product_panel_0_1.Dock(lot_field, windigo.Left)
	product_panel_0_1.Dock(sample_field, windigo.Left)

	clock_panel.Dock(clock_display_now, windigo.Left)
	clock_panel.Dock(clock_display_future, windigo.Left)

	product_panel_0.Dock(clock_panel, windigo.Right)
	product_panel_0.Dock(product_panel_0_0, windigo.Top)
	product_panel_0.Dock(product_panel_0_1, windigo.Top)

	product_panel.Dock(product_panel_0, windigo.Top)
	product_panel.Dock(ranges_button, windigo.Left)
	product_panel.Dock(container_field, windigo.Left)
	product_panel.Dock(reprint_button, windigo.Left)
	product_panel.Dock(reprint_sample_button, windigo.Left)

	dock.Dock(product_panel, windigo.Top)
	dock.Dock(tabs, windigo.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), windigo.Fill) // tab panels dock just below tabs and fill area

	threads.Status_bar = windigo.NewStatusBar(mainWindow)
	mainWindow.SetStatusBar(threads.Status_bar)

	rows, err := DB.DB_Select_product_info.Query()
	GUI.Fill_combobox_from_query_rows(product_field, rows, err, func(rows *sql.Rows) {
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

	new_product_cb := func() product.BaseProduct {
		log.Println("Debug: product_field new_product_cb", qc_product.Base())
		return qc_product.Base()
	}

	panel_water_based := show_water_based(tab_wb, qc_product, new_product_cb)
	panel_oil_based := show_oil_based(tab_oil, qc_product, new_product_cb)
	panel_fr := show_fr(tab_fr, qc_product, new_product_cb)

	// sizing
	refresh_vars := func(font_size int) {

		refresh_globals(font_size)

		top_spacer_height = 20
		top_subpanel_height = top_spacer_height + 2*PRODUCT_FIELD_HEIGHT + 2*INTER_SPACER_HEIGHT + BTM_SPACER_HEIGHT

		top_panel_height = top_subpanel_height + 2*GUI.GROUPBOX_CUSHION + PRODUCT_FIELD_HEIGHT

		hpanel_margin = 10
		hpanel_width = TOP_PANEL_WIDTH - clock_width

		ranges_button_width = 8 * font_size

		container_item_width = 6 * font_size

		reprint_button_margin_l = 2*GUI.LABEL_WIDTH + PRODUCT_FIELD_WIDTH + TOP_PANEL_INTER_SPACER_WIDTH - ranges_button_width - GUI.DISCRETE_FIELD_WIDTH
	}

	refresh := func(font_size int) {
		refresh_vars(font_size)

		mainWindow.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)

		product_panel.SetSize(TOP_PANEL_WIDTH, top_panel_height)

		product_panel_0.SetSize(TOP_PANEL_WIDTH, top_subpanel_height)

		product_panel_0_0.SetSize(hpanel_width, PRODUCT_FIELD_HEIGHT)
		product_panel_0_0.SetMarginLeft(hpanel_margin)
		product_panel_0_0.SetMarginTop(top_spacer_height)

		customer_field.SetMarginLeft(TOP_PANEL_INTER_SPACER_WIDTH)
		sample_field.SetMarginLeft(TOP_PANEL_INTER_SPACER_WIDTH)

		product_panel_0_1.SetSize(hpanel_width, PRODUCT_FIELD_HEIGHT)
		product_panel_0_1.SetMarginTop(INTER_SPACER_HEIGHT)
		product_panel_0_1.SetMarginLeft(hpanel_margin)

		product_field.SetLabeledSize(GUI.LABEL_WIDTH, PRODUCT_FIELD_WIDTH, PRODUCT_FIELD_HEIGHT)
		customer_field.SetLabeledSize(GUI.LABEL_WIDTH, PRODUCT_FIELD_WIDTH, PRODUCT_FIELD_HEIGHT)
		lot_field.SetLabeledSize(GUI.LABEL_WIDTH, PRODUCT_FIELD_WIDTH, PRODUCT_FIELD_HEIGHT)
		sample_field.SetLabeledSize(GUI.LABEL_WIDTH, PRODUCT_FIELD_WIDTH, PRODUCT_FIELD_HEIGHT)

		ranges_button.SetSize(ranges_button_width, GUI.OFF_AXIS)
		ranges_button.SetMarginsAll(BUTTON_MARGIN)

		container_field.SetSize(GUI.DISCRETE_FIELD_WIDTH, GUI.OFF_AXIS)
		container_field.SetItemSize(container_item_width)
		container_field.SetPaddingsAll(GUI.GROUPBOX_CUSHION)

		reprint_button.SetMarginsAll(BUTTON_MARGIN)
		reprint_button.SetMarginLeft(reprint_button_margin_l)
		reprint_button.SetSize(REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

		reprint_sample_button.SetMarginsAll(BUTTON_MARGIN)
		reprint_sample_button.SetSize(REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

		tabs.SetSize(PRODUCT_FIELD_WIDTH, PRODUCT_FIELD_HEIGHT)

		panel_water_based.Refresh()
		panel_oil_based.Refresh()
		panel_fr.Refresh()

	}

	set_font := func(font_size int) {
		GUI.BASE_FONT_SIZE = font_size

		config.Main_config.Set("font_size", GUI.BASE_FONT_SIZE)
		config.Write_config(config.Main_config)

		old_font := windigo.DefaultFont
		windigo.DefaultFont = windigo.NewFont(old_font.Family(), GUI.BASE_FONT_SIZE, 0)
		old_font.Dispose()

		mainWindow.SetFont(windigo.DefaultFont)
		product_field.SetFont(windigo.DefaultFont)
		customer_field.SetFont(windigo.DefaultFont)
		lot_field.SetFont(windigo.DefaultFont)
		sample_field.SetFont(windigo.DefaultFont)

		ranges_button.SetFont(windigo.DefaultFont)
		container_field.SetFont(windigo.DefaultFont)
		reprint_button.SetFont(windigo.DefaultFont)
		reprint_sample_button.SetFont(windigo.DefaultFont)
		tabs.SetFont(windigo.DefaultFont)

		panel_water_based.SetFont(windigo.DefaultFont)
		panel_oil_based.SetFont(windigo.DefaultFont)
		panel_fr.SetFont(windigo.DefaultFont)
		threads.Status_bar.SetFont(windigo.DefaultFont)
		refresh(font_size)

	}
	set_font(GUI.BASE_FONT_SIZE)

	// functionality

	qc_product.Update = func() {
		container_field.Update(qc_product.Container_type)
		log.Println("Debug: update new_product_cb", qc_product)
		panel_water_based.Update(qc_product)
		panel_oil_based.Update(qc_product)
		panel_fr.Update(qc_product)
		panel_fr.ChangeContainer(qc_product)
	}

	product_field_pop_data := func(str string) {
		log.Println("Warn: Debug: product_field_pop_data product_id", qc_product.Product_id)
		log.Println("Warn: Debug: product_field_pop_data lot_id", qc_product.Lot_id)

		// if product_lot.product_id != product_lot.insel_product_id(str) {
		old_product_id := qc_product.Product_id
		qc_product.Product_name = str
		qc_product.Insel_product_self()
		if qc_product.Product_id != old_product_id {
			mainWindow.SetText(qc_product.Product_name)
			qc_product.Reset()
			qc_product.Lot_id = DB.DEFAULT_LOT_ID
			qc_product.Select_product_details()
			qc_product.Update()

			if qc_product.Product_type.Valid {
				tabs.SetCurrent(qc_product.Product_type.Index())
			}

			GUI.Fill_combobox_from_query(lot_field, DB.DB_Select_lot_info, qc_product.Product_id)
			GUI.Fill_combobox_from_query(customer_field, db_select_product_customer_info, qc_product.Product_id)
			GUI.Fill_combobox_from_query(sample_field, db_select_sample_points, qc_product.Product_id)

			qc_product.Update_lot(lot_field.Text(), customer_field.Text())

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
			qc_product.Product_id = DB.INVALID_ID
		}
	}

	lot_field_pop_data := func(str string) {
		qc_product.Update_lot(str, customer_field.Text())
		mainWindow.SetText(str)
	}
	lot_field_text_pop_data := func(str string) {
		formatted_text := strings.ToUpper(strings.TrimSpace(str))
		lot_field.SetText(formatted_text)

		lot_field_pop_data(formatted_text)
	}

	customer_field_pop_data := func(str string) {
		qc_product.Update_lot(lot_field.Text(), str)
	}
	customer_field_text_pop_data := func(str string) {
		formatted_text := strings.ToUpper(strings.TrimSpace(str))
		customer_field.SetText(formatted_text)

		customer_field_pop_data(formatted_text)
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
		if qc_product.Product_name != "" {
			view.ShowNewQCProductRangesView(qc_product)
			log.Println("debug: ranges_button product_lot", qc_product)
		}
	})

	reprint_button.OnClick().Bind(func(e *windigo.Event) {
		if qc_product.Lot_number != "" {
			log.Println("debug: reprint_button")
			qc_product.Reprint()
		}
	})

	reprint_sample_button.OnClick().Bind(func(e *windigo.Event) {
		if qc_product.Lot_number != "" {
			log.Println("debug: Notice: reprint_sample_button")
			qc_product.Reprint_sample()
			log.Println("debug: Notice: qc_product", qc_product)

			// qc_product.Base()
		}
	})

	container_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		qc_product.Container_type = container_field.Get()
		panel_fr.ChangeContainer(qc_product)
	})

	// Clock

	go func() {
		mix_time := 20 * time.Second
		for tock := range clock_ticker {
			clock_display_now.SetText(tock.Format("05"))
			clock_display_future.SetText(tock.Add(mix_time).Format("05"))
		}
	}()

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
			} else {
				log.Printf("error: %q: %s\n", err, "qr_json_keygrab")
			}
			keygrab.SetText("")
			mainWindow.SetFocus()
			return true
		})

	// Resize handling
	mainWindow.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModControl, Key: windigo.KeyOEMPlus}, // {
		func() bool {
			set_font(GUI.BASE_FONT_SIZE + 1)
			return true
		})

	mainWindow.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModControl, Key: windigo.KeyOEMMinus}, // {
		func() bool {
			set_font(GUI.BASE_FONT_SIZE - 1)
			return true
		})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)
	mainWindow.RunMainLoop() // Must call to start event loop.
}

func wndOnClose(arg *windigo.Event) {
	GUI.OKPen.Dispose()
	GUI.ErroredPen.Dispose()
	windigo.Exit()
}
