package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/QR"
	"github.com/samuel-jimenez/qc_data_entry/blender/blendbound"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/windigo"
)

func show_window() {
	proc_name := "main.show_window"

	log.Println("Info: Process started")
	// DEBUG
	// log.Println(time.Now().UTC().UnixNano())

	//
	//
	// definitions
	//
	//
	window_title := "QC Data Entry"

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

		container_item_width,

		reprint_button_margin_l int

		Inbound_Lot *blendbound.InboundLot
	)

	product_text := "Product"
	lot_text := "Lot Number"
	sample_text := "Sample Point"
	customer_text := "Customer Name"
	tester_text := "Tester"

	ranges_text := "Ranges"
	reprint_text := "Reprint"
	inbound_text := "Inbound"
	sample_button_text := "Sample"

	inbound_lot_text := "Inbound Lot"
	internal_text := "Internal"
	container_text := "Container"

	product_data := make(map[string]int)
	inbound_lot_data := make(map[string]*blendbound.InboundLot)
	inbound_container_data := make(map[string]*blendbound.InboundLot)

	qc_product := product.NewQCProduct()

	windigo.DefaultFont = windigo.NewFont("MS Shell Dlg 2", GUI.BASE_FONT_SIZE, windigo.FontNormal)

	inbound_test := "BSQL%"

	//
	//
	//
	// build window
	//
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetText(window_title)

	keygrab := windigo.NewEdit(mainWindow)
	keygrab.Hide()

	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := windigo.NewAutoPanel(mainWindow)

	product_panel_0 := windigo.NewAutoPanel(product_panel)

	//TODO array db_select_all_product

	product_panel_0_0 := windigo.NewAutoPanel(product_panel_0)

	internal_product_field := GUI.NewComboBox(product_panel_0_0, product_text)
	customer_field := GUI.NewComboBox(product_panel_0_0, customer_text)

	product_panel_1_0 := windigo.NewAutoPanel(product_panel_0)
	product_panel_1_0.Hide()

	testing_lot_field := GUI.NewListComboBox(product_panel_1_0, lot_text)
	// inbound_product_field := GUI.NewComboBox(product_panel_1_0, product_text)
	// inbound_product_field := windigo.NewLabeledComboBox(product_panel_1_0, product_text)
	inbound_product_field := windigo.NewLabeledEdit(product_panel_1_0, product_text)
	inbound_product_field.SetReadOnly(true)

	product_panel_0_1 := windigo.NewAutoPanel(product_panel_0)

	lot_field := GUI.NewComboBox(product_panel_0_1, lot_text)
	sample_field := GUI.NewComboBox(product_panel_0_1, sample_text)

	product_panel_1_1 := windigo.NewAutoPanel(product_panel_0)
	product_panel_1_1.Hide()

	// inbound_lot_field := GUI.NewComboBox(product_panel_1_1, product_text)
	inbound_lot_field := GUI.NewComboBox(product_panel_1_1, inbound_lot_text)
	inbound_container_field := GUI.NewComboBox(product_panel_1_1, container_text)
	// inbound_lot_field container_text

	product_panel_0_2 := windigo.NewAutoPanel(product_panel_0)

	tester_field := GUI.NewLabeledSearchBoxFromQuery(product_panel_0_2, tester_text, DB.DB_Select_all_qc_tester)

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
	ranges_button.SetText(ranges_text)

	sample_button := windigo.NewPushButton(product_panel)
	sample_button.SetText(sample_button_text)
	sample_button.Hide()

	container_field := product.BuildNewDiscreteView(product_panel, "Container Type", qc_product.Container_type, []string{"Sample", "Tote", "Railcar"}) // bs.container_types

	reprint_button := windigo.NewPushButton(product_panel)
	reprint_button.SetText(reprint_text)

	inbound_button := windigo.NewPushButton(product_panel)
	inbound_button.SetText(inbound_text)

	internal_button := windigo.NewPushButton(product_panel)
	internal_button.SetText(internal_text)
	internal_button.Hide()

	tabs := windigo.NewTabView(mainWindow)
	tab_wb := tabs.AddAutoPanel("Water Based")
	tab_oil := tabs.AddAutoPanel("Oil Based")
	tab_fr := tabs.AddAutoPanel("Friction Reducer")

	//
	//
	// Dock
	//
	//
	product_panel_0_0.Dock(internal_product_field, windigo.Left)
	product_panel_0_0.Dock(customer_field, windigo.Left)

	product_panel_0_1.Dock(lot_field, windigo.Left)
	product_panel_0_1.Dock(sample_field, windigo.Left)

	product_panel_0_2.Dock(tester_field, windigo.Left)

	product_panel_1_0.Dock(testing_lot_field, windigo.Left)
	product_panel_1_0.Dock(inbound_product_field, windigo.Left)

	product_panel_1_1.Dock(inbound_lot_field, windigo.Left)
	product_panel_1_1.Dock(inbound_container_field, windigo.Left)

	clock_panel.Dock(clock_display_now, windigo.Left)
	clock_panel.Dock(clock_display_future, windigo.Left)

	product_panel_0.Dock(clock_panel, windigo.Right)
	product_panel_0.Dock(product_panel_0_0, windigo.Top)
	product_panel_0.Dock(product_panel_1_0, windigo.Top)
	product_panel_0.Dock(product_panel_0_1, windigo.Top)
	product_panel_0.Dock(product_panel_1_1, windigo.Top)
	product_panel_0.Dock(product_panel_0_2, windigo.Top)

	product_panel.Dock(product_panel_0, windigo.Top)
	product_panel.Dock(ranges_button, windigo.Left)
	product_panel.Dock(sample_button, windigo.Left)
	product_panel.Dock(container_field, windigo.Left)
	product_panel.Dock(reprint_button, windigo.Left)
	product_panel.Dock(inbound_button, windigo.Left)
	product_panel.Dock(internal_button, windigo.Left)

	dock.Dock(product_panel, windigo.Top)
	dock.Dock(tabs, windigo.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), windigo.Fill) // tab panels dock just below tabs and fill area

	threads.Status_bar = windigo.NewStatusBar(mainWindow)
	mainWindow.SetStatusBar(threads.Status_bar)

	//
	//
	//
	//
	//
	// combobox
	GUI.Fill_combobox_from_query_rows(internal_product_field, func(row *sql.Rows) error {
		var (
			id                   int
			internal_name        string
			product_moniker_name string
		)
		if err := row.Scan(&id, &internal_name, &product_moniker_name); err != nil {
			return err
		}
		name := product_moniker_name + " " + internal_name
		product_data[name] = id

		internal_product_field.AddItem(name)
		return nil
	}, DB.DB_Select_product_info)

	GUI.Fill_combobox_from_query(testing_lot_field,
		DB.DB_Select_lot_list_name, inbound_test)

	proc_name = "main.FillInbound"
	DB.Forall_err(proc_name,
		func() {
			inbound_container_field.DeleteAllItems()
			inbound_lot_field.DeleteAllItems()
		},
		func(row *sql.Rows) error {
			Inbound, err := blendbound.NewInboundLotFromRow(row)
			if err != nil {
				return err
			}
			inbound_lot_field.AddItem(Inbound.Lot_number)
			inbound_lot_data[Inbound.Lot_number] = Inbound
			inbound_container_data[Inbound.Container_name] = Inbound
			return nil
		},
		DB.DB_Select_inbound_lot_status, blendbound.Status_AVAILABLE)

	for _, Container_name := range slices.Sorted(maps.Keys(inbound_container_data)) {
		inbound_container_field.AddItem(Container_name)
	}

	// new_product_cb := func() *product.BaseProduct {
	new_product_cb := func() product.BaseProduct {
		log.Println("Debug: product_field new_product_cb", qc_product.Base())

		// there has to be better way
		// if !qc_product.Tester.Valid{
		// 	tester_field.ShowDropdown(true)
		// 	return nil
		// }
		qc_product.Valid = qc_product.Tester.Valid

		// only last test holds dropdown
		if !qc_product.Tester.Valid {
			tester_field.ShowDropdown(true)
			tester_field.SetFocus()
		}
		// return &(qc_product.Base())
		return qc_product.Base()
	}

	panel_water_based := show_water_based(tab_wb, qc_product, new_product_cb)
	panel_oil_based := show_oil_based(tab_oil, qc_product, new_product_cb)
	panel_fr := show_fr(tab_fr, qc_product, new_product_cb)

	//
	//
	//
	//
	// sizing
	refresh_vars := func(font_size int) {

		refresh_globals(font_size)

		num_rows := 3

		top_spacer_height = 20
		top_subpanel_height = top_spacer_height + num_rows*(GUI.PRODUCT_FIELD_HEIGHT+INTER_SPACER_HEIGHT) + BTM_SPACER_HEIGHT

		top_panel_height = top_subpanel_height + num_rows*GUI.GROUPBOX_CUSHION + GUI.PRODUCT_FIELD_HEIGHT

		hpanel_margin = 10
		hpanel_width = TOP_PANEL_WIDTH - clock_width

		container_item_width = 6 * font_size

		reprint_button_margin_l = 2*GUI.LABEL_WIDTH + GUI.PRODUCT_FIELD_WIDTH + TOP_PANEL_INTER_SPACER_WIDTH - GUI.SMOL_BUTTON_WIDTH - GUI.DISCRETE_FIELD_WIDTH
	}

	refresh := func(font_size int) {
		refresh_vars(font_size)

		mainWindow.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)

		product_panel.SetSize(TOP_PANEL_WIDTH, top_panel_height)

		product_panel_0.SetSize(TOP_PANEL_WIDTH, top_subpanel_height)

		product_panel_0_0.SetSize(hpanel_width, GUI.PRODUCT_FIELD_HEIGHT)
		product_panel_0_0.SetMarginTop(top_spacer_height)
		product_panel_0_0.SetMarginLeft(hpanel_margin)

		customer_field.SetMarginLeft(TOP_PANEL_INTER_SPACER_WIDTH)
		sample_field.SetMarginLeft(TOP_PANEL_INTER_SPACER_WIDTH)

		product_panel_0_1.SetSize(hpanel_width, GUI.PRODUCT_FIELD_HEIGHT)
		product_panel_0_1.SetMarginTop(INTER_SPACER_HEIGHT)
		product_panel_0_1.SetMarginLeft(hpanel_margin)

		internal_product_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		customer_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		sample_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

		product_panel_1_0.SetSize(hpanel_width, GUI.PRODUCT_FIELD_HEIGHT)
		product_panel_1_0.SetMarginTop(top_spacer_height)
		product_panel_1_0.SetMarginLeft(hpanel_margin)

		testing_lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

		inbound_product_field.SetMarginLeft(TOP_PANEL_INTER_SPACER_WIDTH)
		inbound_product_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

		product_panel_1_1.SetSize(hpanel_width, GUI.PRODUCT_FIELD_HEIGHT)
		product_panel_1_1.SetMarginTop(INTER_SPACER_HEIGHT)
		product_panel_1_1.SetMarginLeft(hpanel_margin)

		inbound_lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		inbound_container_field.SetMarginLeft(TOP_PANEL_INTER_SPACER_WIDTH)
		inbound_container_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

		product_panel_0_2.SetSize(hpanel_width, GUI.PRODUCT_FIELD_HEIGHT)
		product_panel_0_2.SetMarginTop(INTER_SPACER_HEIGHT)
		product_panel_0_2.SetMarginLeft(hpanel_margin)

		tester_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

		ranges_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
		ranges_button.SetMarginsAll(BUTTON_MARGIN)

		sample_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
		sample_button.SetMarginsAll(BUTTON_MARGIN)

		container_field.SetSize(GUI.DISCRETE_FIELD_WIDTH, GUI.OFF_AXIS)
		container_field.SetItemSize(container_item_width)
		container_field.SetPaddingsAll(GUI.GROUPBOX_CUSHION)
		container_field.SetPaddingLeft(0)

		reprint_button.SetMarginsAll(BUTTON_MARGIN)
		reprint_button.SetMarginLeft(reprint_button_margin_l)
		reprint_button.SetSize(REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

		inbound_button.SetMarginsAll(BUTTON_MARGIN)
		inbound_button.SetSize(REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

		internal_button.SetMarginsAll(BUTTON_MARGIN)
		internal_button.SetSize(REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

		tabs.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

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
		internal_product_field.SetFont(windigo.DefaultFont)
		customer_field.SetFont(windigo.DefaultFont)
		lot_field.SetFont(windigo.DefaultFont)
		sample_field.SetFont(windigo.DefaultFont)
		tester_field.SetFont(windigo.DefaultFont)
		inbound_lot_field.SetFont(windigo.DefaultFont)
		testing_lot_field.SetFont(windigo.DefaultFont)
		inbound_product_field.SetFont(windigo.DefaultFont)
		inbound_container_field.SetFont(windigo.DefaultFont)

		ranges_button.SetFont(windigo.DefaultFont)
		sample_button.SetFont(windigo.DefaultFont)
		container_field.SetFont(windigo.DefaultFont)
		reprint_button.SetFont(windigo.DefaultFont)
		inbound_button.SetFont(windigo.DefaultFont)
		internal_button.SetFont(windigo.DefaultFont)
		tabs.SetFont(windigo.DefaultFont)

		panel_water_based.SetFont(windigo.DefaultFont)
		panel_oil_based.SetFont(windigo.DefaultFont)
		panel_fr.SetFont(windigo.DefaultFont)
		threads.Status_bar.SetFont(windigo.DefaultFont)
		refresh(font_size)

	}
	set_font(GUI.BASE_FONT_SIZE)

	//
	//
	//
	//
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
		log.Println("Warn: Debug: product_field_pop_data product_lot_id", qc_product.Product_Lot_id)

		// if product_lot.product_id != product_lot.insel_product_id(str) {
		old_product_id := qc_product.Product_id
		qc_product.Product_name = str
		qc_product.Insel_product_self()
		if qc_product.Product_id != old_product_id {
			mainWindow.SetText(qc_product.Product_name)
			qc_product.ResetQC()

			qc_product.Select_product_details()
			qc_product.Update()

			if qc_product.Product_type.Valid {
				tabs.SetCurrent(qc_product.Product_type.Index())
			}

			GUI.Fill_combobox_from_query(lot_field, DB.DB_Select_product_lot_product, qc_product.Product_id)
			GUI.Fill_combobox_from_query(customer_field, DB.DB_Select_product_customer_info, qc_product.Product_id)
			GUI.Fill_combobox_from_query(sample_field, DB.DB_Select_product_sample_points, qc_product.Product_id)

			qc_product.Update_lot(lot_field.Text(), customer_field.Text())

			qc_product.Sample_point = sample_field.Text()

		}
	}
	product_field_text_pop_data := func(str string) {
		formatted_text := strings.ToUpper(strings.TrimSpace(str))
		internal_product_field.SetText(formatted_text)
		if internal_product_field.Text() != "" {
			product_field_pop_data(internal_product_field.Text())
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

	tester_field_pop_data := func(str string) {
		qc_product.SetTester(str)
	}
	tester_field_text_pop_data := func(str string) {
		formatted_text := strings.ToUpper(strings.TrimSpace(str))
		tester_field.SetText(formatted_text)

		tester_field_pop_data(formatted_text)
	}

	qr_pop_data := func(product QR.QRJson) {
		product_field_text_pop_data(product.Product_type)
		lot_field_text_pop_data(product.Lot_number)
	}

	internal_product_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		product_field_pop_data(internal_product_field.GetSelectedItem())

	})
	internal_product_field.OnKillFocus().Bind(func(e *windigo.Event) {
		product_field_text_pop_data(internal_product_field.Text())
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

	testing_lot_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		Inbound_Lot = nil
		qc_product.Update_testing_lot(testing_lot_field.GetSelectedItem())
		inbound_product_field.SetText(qc_product.Product_name)

		blend := qc_product.Blend
		if blend == nil || len(blend.Components) < 1 {
			log.Println("Err: testing_lot_field.OnSelectedChange", qc_product, blend)
			return
		}

		component := blend.Components[0]
		inbound_lot_field.SetText(component.Lot_name)
		inbound_container_field.SetText(component.Container_name)
		// inbound_product_field.SetText(component.Component_name)

		// TODO update component list

	})

	inbound_lot_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		Inbound_Lot = inbound_lot_data[inbound_lot_field.GetSelectedItem()]
		if Inbound_Lot == nil {
			return
		}

		inbound_product_field.SetText(Inbound_Lot.Product_name)
		inbound_container_field.SetText(Inbound_Lot.Container_name)
	})
	inbound_container_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		Inbound_Lot = inbound_container_data[inbound_container_field.GetSelectedItem()]
		if Inbound_Lot == nil {
			return
		}

		inbound_product_field.SetText(Inbound_Lot.Product_name)
		inbound_lot_field.SetText(Inbound_Lot.Lot_number)
	})

	tester_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		tester_field_pop_data(tester_field.GetSelectedItem())
	})
	tester_field.OnKillFocus().Bind(func(e *windigo.Event) {
		tester_field_text_pop_data(tester_field.Text())
	})

	ranges_button.OnClick().Bind(func(e *windigo.Event) {
		if qc_product.Product_name != "" {
			views.ShowNewQCProductRangesView(qc_product)
			log.Println("debug: ranges_button-product_lot", qc_product)
		}
	})

	sample_button.OnClick().Bind(func(e *windigo.Event) {
		if Inbound_Lot == nil {
			return
		}

		Inbound_Lot_ := Inbound_Lot
		Inbound_Lot = nil

		Inbound_Lot_.Update_status(blendbound.Status_SAMPLED)
		Inbound_Lot_.Quality_test()
		GUI.Fill_combobox_from_query(testing_lot_field,
			DB.DB_Select_lot_list_name, inbound_test)
		delete(inbound_container_data, Inbound_Lot_.Container_name)
		delete(inbound_lot_data, Inbound_Lot_.Lot_number)

	})

	reprint_button.OnClick().Bind(func(e *windigo.Event) {
		if qc_product.Lot_number != "" {
			log.Println("debug: reprint_button")
			qc_product.Reprint()
		}
	})

	inbound_button.OnClick().Bind(func(e *windigo.Event) {
		product_panel_1_0.Show()
		product_panel_1_1.Show()
		sample_button.Show()
		internal_button.Show()

		product_panel_0_0.Hide()
		product_panel_0_1.Hide()
		ranges_button.Hide()
		inbound_button.Hide()

		qc_product.Container_type = product.DiscreteFromInt(CONTAINER_SAMPLE)
		panel_fr.ChangeContainer(qc_product)

	})

	internal_button.OnClick().Bind(func(e *windigo.Event) {
		product_panel_1_0.Hide()
		product_panel_1_1.Hide()
		sample_button.Hide()
		internal_button.Hide()

		product_panel_0_0.Show()
		product_panel_0_1.Show()
		ranges_button.Show()
		inbound_button.Show()
	})

	container_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		qc_product.Container_type = container_field.Get()
		panel_fr.ChangeContainer(qc_product)
	})

	//
	//
	//
	//
	// Clock

	go func() {
		mix_time := 20 * time.Second
		for tock := range clock_ticker {
			clock_display_now.SetText(tock.Format("05"))
			clock_display_future.SetText(tock.Add(mix_time).Format("05"))
		}
	}()

	AddShortcuts(mainWindow, keygrab, set_font, qr_pop_data)

	dbError := windigo.NewDialog(nil)
	dbError.Center()
	dbError.Show()
	dbError.OnClose().Bind(func(e *windigo.Event) { windigo.Exit() })
	dbError.RunMainLoop() // Must call to start event loop.

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)
	mainWindow.RunMainLoop() // Must call to start event loop.
}

func AddShortcuts(mainWindow *windigo.Form, keygrab *windigo.Edit, set_font func(int), qr_pop_data func(QR.QRJson)) {
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
				log.Printf("error: [%s]: %q\n", "qr_json_keygrab", err)
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
}

func wndOnClose(arg *windigo.Event) {
	GUI.OKPen.Dispose()
	GUI.ErroredPen.Dispose()
	windigo.Exit()
}
