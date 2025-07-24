package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/qc_data_entry/viewer"
	"github.com/samuel-jimenez/windigo"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var (
	qc_db *sql.DB
)

func main() {
	//load config
	config.Main_config = config.Load_config_viewer("qc_data_viewer")

	// log to file
	log_file, err := os.OpenFile(config.LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Crit: error opening file: %v", err)
	}
	defer log_file.Close()

	log.SetOutput(log_file)
	log.Println("Info: Logging Started")

	//open_db
	// qc_db, err := sql.Open("sqlite3", DB_FILE)
	qc_db, err = sql.Open("sqlite3", ":memory:")
	qc_db.Exec("attach ? as 'bs'", config.DB_FILE)
	if err != nil {
		log.Fatal("Crit: error opening database: ", err)
	}
	defer qc_db.Close()

	dbinit(qc_db)

	//setup print goroutine
	threads.PRINT_QUEUE = make(chan string, 4)
	defer close(threads.PRINT_QUEUE)
	go threads.Do_print_queue(threads.PRINT_QUEUE)

	//setup status_bar goroutine
	threads.STATUS_QUEUE = make(chan string, 16)
	defer close(threads.STATUS_QUEUE)
	go threads.Do_status_queue(threads.STATUS_QUEUE)

	//show main window
	show_window()

}

func _select_samples(rows *sql.Rows, err error, fn string) []viewer.QCData {
	if err != nil {
		log.Printf("error: [%s]: %q\n",  fn,  err)
		// return -1
	}

	data := make([]viewer.QCData, 0)
	for rows.Next() {
		var (
			qc_data              viewer.QCData
			_timestamp           int64
			product_moniker_name string
			internal_name        string
		)

		if err := rows.Scan(&product_moniker_name, &internal_name,
			&qc_data.Product_name_customer,
			&qc_data.Lot_name,
			&qc_data.Sample_point,
			&_timestamp,
			&qc_data.PH,
			&qc_data.Specific_gravity,
			&qc_data.String_test,
			&qc_data.Viscosity); err != nil {
			log.Fatalf("Crit: [%s]: %v", fn, err)
		}
		qc_data.Product_name = product_moniker_name + " " + internal_name

		qc_data.Time_stamp = time.Unix(0, _timestamp)
		data = append(data, qc_data)
	}
	return data
}

func select_samples() []viewer.QCData {
	rows, err := DB_Select_samples.Query()
	return _select_samples(rows, err, "select_samples")
}

func select_product_samples(product_id int) []viewer.QCData {
	rows, err := DB_Select_product_samples.Query(product_id)
	return _select_samples(rows, err, "select_product_samples")
}

var (
	SAMPLE_SELECT_STRING string
	DB_Select_product_samples,
	DB_Select_samples *sql.Stmt
)

func dbinit(db *sql.DB) {

	DB.Check_db(db)
	DB.DBinit(db)

	SAMPLE_SELECT_STRING = `
	select
		product_moniker_name,
		product_name_internal,
		product_name_customer,
		lot_name,
		sample_point,
		time_stamp,
		ph ,
		specific_gravity ,
		string_test ,
		viscosity
	from bs.qc_samples
		join bs.product_lot using (lot_id)
		join bs.product_line using (product_id)
		join bs.product_moniker using (product_moniker_id)
		left join bs.product_sample_points using (sample_point_id)
		left join bs.product_customer_line using (product_customer_id,product_id)
	`

	DB_Select_product_samples = DB.PrepareOrElse(db, SAMPLE_SELECT_STRING+`
	where product_id = ?
	`)

	DB_Select_samples = DB.PrepareOrElse(db, SAMPLE_SELECT_STRING+`
	order by product_moniker_name,product_name_internal
	`)

}

func show_window() {

	log.Println("Info: Process started")

	window_title := "QC Data Viewer"

	product_data := make(map[string]int)
	// lot_data := make(map[string]int)
	var lot_data []string

	top_panel_width := viewer.WINDOW_WIDTH
	top_panel_height := 110

	hpanel_width := top_panel_width

	label_width := GUI.LABEL_WIDTH

	hpanel_margin := 10

	field_width := viewer.PRODUCT_FIELD_WIDTH
	field_height := 20
	top_spacer_height := 20
	inter_spacer_width := 30
	inter_spacer_height := viewer.INTER_SPACER_HEIGHT

	// button_margin := 5
	clear_button_width := 20
	filter_button_width := 50
	search_button_width := 50
	reprint_sample_button_margin := hpanel_margin + label_width + field_width + clear_button_width + inter_spacer_width - filter_button_width - search_button_width
	reprint_sample_button_width := 100
	regen_sample_button_width := 90

	product_text := "Product"
	lot_text := "Lot Number"
	sample_text := "Sample Point"
	customer_text := "Customer Name"

	// build window

	windigo.DefaultFont = windigo.NewFont("MS Shell Dlg 2", GUI.BASE_FONT_SIZE, windigo.FontNormal)

	mainWindow := windigo.NewForm(nil)
	mainWindow.SetSize(viewer.WINDOW_WIDTH, viewer.WINDOW_HEIGHT)
	mainWindow.SetText(window_title)

	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := windigo.NewAutoPanel(mainWindow)
	product_panel.SetSize(top_panel_width, top_panel_height)

	//TODO array db_select_all_product

	prod_panel := windigo.NewAutoPanel(product_panel)
	prod_panel.SetSize(hpanel_width, field_height)

	product_field := GUI.NewSizedListComboBox(prod_panel, label_width, field_width, field_height, product_text)
	product_clear_button := windigo.NewPushButton(prod_panel)
	product_clear_button.SetText("-")
	product_clear_button.SetSize(clear_button_width, GUI.OFF_AXIS)
	customer_field := GUI.NewSizedListComboBox(prod_panel, label_width, field_width, field_height, customer_text)

	customer_field.SetMarginLeft(inter_spacer_width)
	prod_panel.SetMarginLeft(hpanel_margin)
	prod_panel.SetMarginTop(top_spacer_height)
	prod_panel.Dock(product_field, windigo.Left)
	prod_panel.Dock(product_clear_button, windigo.Left)
	prod_panel.Dock(customer_field, windigo.Left)

	lot_panel := windigo.NewAutoPanel(product_panel)
	lot_panel.SetSize(hpanel_width, field_height)

	//TODO fix ListComboBox sizes so it works like for size 10
	//TODO change this to search like the filter
	lot_field := GUI.NewSizedListComboBox(lot_panel, label_width, field_width, field_height, lot_text)
	lot_clear_button := windigo.NewPushButton(lot_panel)
	lot_clear_button.SetText("-")
	lot_clear_button.SetSize(clear_button_width, GUI.OFF_AXIS)
	sample_field := GUI.NewSizedListComboBox(lot_panel, label_width, field_width, field_height, sample_text)

	lot_panel.SetMarginTop(inter_spacer_height)
	lot_panel.SetMarginLeft(hpanel_margin)
	sample_field.SetMarginLeft(inter_spacer_width)
	lot_panel.Dock(lot_field, windigo.Left)
	lot_panel.Dock(lot_clear_button, windigo.Left)
	lot_panel.Dock(sample_field, windigo.Left)

	filter_button := windigo.NewPushButton(product_panel)
	filter_button.SetText("Filter")
	filter_button.SetSize(filter_button_width, GUI.OFF_AXIS)

	search_button := windigo.NewPushButton(product_panel)
	search_button.SetText("Search")
	search_button.SetSize(search_button_width, GUI.OFF_AXIS)

	reprint_sample_button := windigo.NewPushButton(product_panel)
	reprint_sample_button.SetText("Reprint Sample")
	reprint_sample_button.SetSize(reprint_sample_button_width, GUI.OFF_AXIS)
	reprint_sample_button.SetMarginLeft(reprint_sample_button_margin)

	regen_sample_button := windigo.NewPushButton(product_panel)
	regen_sample_button.SetText("Regen Sample")
	regen_sample_button.SetSize(regen_sample_button_width, GUI.OFF_AXIS)

	// FilterListView := NewSQLFilterListView(product_panel)
	FilterListView := viewer.NewSQLFilterListView(mainWindow)

	product_panel.Dock(prod_panel, windigo.Top)
	product_panel.Dock(lot_panel, windigo.Top)
	product_panel.Dock(filter_button, windigo.Left)
	product_panel.Dock(search_button, windigo.Left)
	product_panel.Dock(reprint_sample_button, windigo.Left)
	product_panel.Dock(regen_sample_button, windigo.Left)

	// product_panel.Dock(FilterListView, windigo.Left)
	// product_panel.Dock(FilterListView, windigo.Top)

	// product_panel.Dock(reprint_button, windigo.Left)

	// tabs := windigo.NewTabView(mainWindow)
	// tab_wb := tabs.AddAutoPanel("Water Based")
	// tab_oil := tabs.AddAutoPanel("Oil Based")
	// tab_fr := tabs.AddAutoPanel("Friction Reducer")

	table := viewer.NewQCDataView(mainWindow)
	table.Set(select_samples())
	table.Update()

	threads.Status_bar = windigo.NewStatusBar(mainWindow)
	mainWindow.SetStatusBar(threads.Status_bar)
	dock.Dock(threads.Status_bar, windigo.Bottom)

	dock.Dock(product_panel, windigo.Top)
	dock.Dock(FilterListView, windigo.Top)

	dock.Dock(table, windigo.Fill)

	// functionality
	update_lot := func(id int, name string) {
		lot_field.AddItem(name)
		viewer.COL_ITEMS_LOT = append(viewer.COL_ITEMS_LOT, name)
	}

	update_sample := func(id int, name string) {
		sample_field.AddItem(name)
		viewer.COL_ITEMS_SAMPLE = append(viewer.COL_ITEMS_SAMPLE, name)
	}

	clear_product := func() {
		table.Set(select_samples())
		table.Update()
		product_field.SetSelectedItem(-1)

		viewer.COL_ITEMS_LOT = nil
		lot_field.DeleteAllItems()
		// for name := range lot_data {
		for id, name := range lot_data {
			update_lot(id, name)
		}

		viewer.COL_ITEMS_SAMPLE = nil
		GUI.Fill_combobox_from_query_0_2(sample_field, DB.DB_Select_all_sample_points, update_sample)

		// FilterListView[viewer.COL_KEY_SAMPLE].Update(viewer.COL_ITEMS_SAMPLE)

		FilterListView.Update(viewer.COL_KEY_LOT, viewer.COL_ITEMS_LOT)
		FilterListView.Update(viewer.COL_KEY_SAMPLE, viewer.COL_ITEMS_SAMPLE)
	}

	clear_lot := func() {
		//TODO
		product_id := product_data[product_field.GetSelectedItem()]
		if product_id > 0 {
			table.Set(select_product_samples(product_id))
		}
		table.Update()
		lot_field.SetSelectedItem(-1)
	}

	// clear := func() {
	// 	clear_product()
	// 	clear_lot()
	// }

	// combobox
	rows, err := DB.DB_Select_product_info.Query()
	GUI.Fill_combobox_from_query_rows(product_field, rows, err, func(rows *sql.Rows) {
		var (
			id                   int
			internal_name        string
			product_moniker_name string
		)
		if err := rows.Scan(&id, &internal_name, &product_moniker_name); err == nil {
			name := product_moniker_name + " " + internal_name
			product_data[name] = id

			product_field.AddItem(name)
		}
	})

	GUI.Fill_combobox_from_query_0_2(lot_field, DB.DB_Select_lot_all, func(id int, name string) {
		// lot_data[name] = id
		lot_data = append(lot_data, name)
		update_lot(id, name)
	})

	GUI.Fill_combobox_from_query_0_2(sample_field, DB.DB_Select_all_sample_points, update_sample)

	FilterListView.AddContinuous(viewer.COL_KEY_TIME, viewer.COL_LABEL_TIME)
	FilterListView.AddDiscreteSearch(viewer.COL_KEY_LOT, viewer.COL_LABEL_LOT, viewer.COL_ITEMS_LOT)
	FilterListView.AddDiscreteMulti(viewer.COL_KEY_SAMPLE, viewer.COL_LABEL_SAMPLE, viewer.COL_ITEMS_SAMPLE)
	FilterListView.AddContinuous(viewer.COL_KEY_PH, viewer.COL_LABEL_PH)
	FilterListView.AddContinuous(viewer.COL_KEY_SG, viewer.COL_LABEL_SG)
	//TODO FilterListView.AddContinuous(viewer.COL_KEY_DENSITY, viewer.COL_LABEL_DENSITY)
	FilterListView.AddContinuous(viewer.COL_KEY_STRING, viewer.COL_LABEL_STRING)
	FilterListView.AddContinuous(viewer.COL_KEY_VISCOSITY, viewer.COL_LABEL_VISCOSITY)
	FilterListView.Hide()

	// listeners
	product_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		// product_field_pop_data(product_field.GetSelectedItem())

		product_id := product_data[product_field.GetSelectedItem()]
		samples := select_product_samples(product_id)
		table.Set(samples)
		table.Update()

		viewer.COL_ITEMS_LOT = nil
		GUI.Fill_combobox_from_query_1_2(lot_field, DB.DB_Select_lot_info, int64(product_id), update_lot)

		viewer.COL_ITEMS_SAMPLE = nil
		GUI.Fill_combobox_from_query_1_2(sample_field, DB.DB_Select_product_sample_points, int64(product_id), update_sample)

		// FilterListView[viewer.COL_KEY_LOT].Update(viewer.COL_ITEMS_LOT)
		FilterListView.Update(viewer.COL_KEY_LOT, viewer.COL_ITEMS_LOT)
		FilterListView.Update(viewer.COL_KEY_SAMPLE, viewer.COL_ITEMS_SAMPLE)

		// log.Println()

	})

	filter_button.OnClick().Bind(func(e *windigo.Event) {
		if FilterListView.Visible() {
			FilterListView.Hide()
			// log.Println(" Hide ize: ", FilterListView.Width(), FilterListView.ClientWidth(), FilterListView.Height())
			// product_panel.Update()

		} else {
			FilterListView.Show()
			// log.Println(" Show ize: ", FilterListView.Width(), FilterListView.ClientWidth(), FilterListView.Height())
			// product_panel.Update()
			// FilterListView.SetSize(viewer.WINDOW_WIDTH, FilterListView.Height())
			// log.Println(" Show ize: ", FilterListView.Width(), FilterListView.ClientWidth(), FilterListView.Height())

		}
	})

	search_button.OnClick().Bind(func(e *windigo.Event) {
		//TODO
		log.Println("FilterListView:", FilterListView.Get().Get())
		log.Println("SAMPLE_SELECT_STRING", SAMPLE_SELECT_STRING+FilterListView.Get().Get())
		rows, err := qc_db.Query(SAMPLE_SELECT_STRING + FilterListView.Get().Get())
		table.Set(_select_samples(rows, err, "search samples"))
		table.Update()

	})

	product_clear_button.OnClick().Bind(func(e *windigo.Event) {
		// log.Println("product_field", product_field.SelectedItem())

		// clear()
		clear_product()
	})

	lot_clear_button.OnClick().Bind(func(e *windigo.Event) {
		clear_lot()
	})

	reprint_sample_button.OnClick().Bind(func(e *windigo.Event) {
		for _, data := range table.SelectedItems() {
			data.(viewer.QCData).Product().Reprint_sample()
		}

	})
	regen_sample_button.OnClick().Bind(func(e *windigo.Event) {
		for _, data := range table.SelectedItems() {
			if err := data.(viewer.QCData).Product().Output_sample(); err != nil {
				log.Printf("Error: [%s]: %q\n",  "regen_sample_button",  err)
				log.Printf("Debug: %q: %v\n", err, data)
				threads.Show_status("Error Creating Label")
			}
		}

	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(func(*windigo.Event) {
		windigo.Exit()
	})
	mainWindow.RunMainLoop() // Must call to start event loop.
}
