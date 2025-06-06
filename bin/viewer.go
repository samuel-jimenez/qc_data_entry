package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/windigo"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func main() {
	//load config
	config.Main_config = config.Load_config("qc_data_viewer")

	// log to file
	log_file, err := os.OpenFile(config.LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer log_file.Close()

	log.SetOutput(log_file)
	log.Println("Info: Logging Started")

	//open_db
	// qc_db, err := sql.Open("sqlite3", DB_FILE)
	qc_db, err := sql.Open("sqlite3", ":memory:")
	qc_db.Exec("attach ? as 'bs'", config.DB_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer qc_db.Close()

	dbinit(qc_db)
	/*
		//setup print goroutine
		print_queue = make(chan string, 4)
		defer close(print_queue)
		go do_print_queue(print_queue)
	*/

	//show main window
	show_window()

}

// package viewer
var (
	db_select_product_info,
	db_select_product_samples *sql.Stmt
)

type QCData struct {
	lot_name     string
	sample_point sql.NullString
	time_stamp   time.Time
	ph,
	specific_gravity,
	string_test,
	viscosity sql.NullFloat64
}

func select_product_samples(product_id int) []QCData {
	rows, err := db_select_product_samples.Query(product_id)
	if err != nil {
		log.Printf("error: %q: %s\n", err, "select_product_samples")
		// return -1
	}

	data := make([]QCData, 0)
	for rows.Next() {
		var (
			qc_data    QCData
			_timestamp int64
		)

		if err := rows.Scan(&qc_data.lot_name,
			&qc_data.sample_point,
			&_timestamp,
			&qc_data.ph,
			&qc_data.specific_gravity,
			&qc_data.string_test,
			&qc_data.viscosity); err != nil {
			log.Fatal(err)
		}
		qc_data.time_stamp = time.Unix(0, _timestamp)

		// time.Unix(unixTime.Unix(), 0).UTC()
		data = append(data, qc_data)
	}
	return data
}

func dbinit(db *sql.DB) {

	DB.Check_db(db)

	db_select_product_info = DB.PrepareOrElse(db, `
	select product_id, product_name_internal, product_moniker_name
		from bs.product_line
		join bs.product_moniker using (product_moniker_id)
		order by product_moniker_name,product_name_internal

	`)

	db_select_product_samples = DB.PrepareOrElse(db, `
	select
		lot_name,
		sample_point,
		time_stamp,
		ph ,
		specific_gravity ,
		string_test ,
		viscosity
	from bs.qc_samples
		join bs.product_lot using (lot_id)
		left join bs.product_sample_points using (sample_point_id)
	where product_id = ?
	`)
}

var (
	RANGE_WIDTH        = 200
	GROUP_WIDTH        = 210
	GROUP_HEIGHT       = 170
	GROUP_MARGIN       = 5
	PRODUCT_TYPE_WIDTH = 150

	LABEL_WIDTH         = 100
	PRODUCT_FIELD_WIDTH = 150
	DATA_FIELD_WIDTH    = 60
	FIELD_HEIGHT        = 28

	RANGES_PADDING         = 5
	RANGES_FIELD_HEIGHT    = 50
	OFF_AXIS               = 0
	RANGES_RO_FIELD_WIDTH  = 40
	RANGES_RO_SPACER_WIDTH = 20
	RANGES_RO_FIELD_HEIGHT = FIELD_HEIGHT

	BUTTON_WIDTH  = 100
	BUTTON_HEIGHT = 40
	// 	200
	// 50

	ERROR_MARGIN = 3

	TOP_SPACER_WIDTH     = 7
	TOP_SPACER_HEIGHT    = 17
	INTER_SPACER_HEIGHT  = 2
	BTM_SPACER_WIDTH     = 2
	BTM_SPACER_HEIGHT    = 2
	BUTTON_SPACER_HEIGHT = 195
)

func show_window() {

	log.Println("Info: Process started")

	// var product_data map [int]string
	// var product_data map[string]int
	product_data := make(map[string]int)
	// var qc_data QCData

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

	// button_margin := 5

	product_text := "Product"
	lot_text := "Lot Number"
	sample_text := "Sample Point"
	customer_text := "Customer Name"

	// build window
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetSize(window_width, window_height)
	mainWindow.SetText("QC Data Viewer")

	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := windigo.NewAutoPanel(mainWindow)
	product_panel.SetSize(top_panel_width, top_panel_height)

	//TODO array db_select_all_product

	prod_panel := windigo.NewAutoPanel(product_panel)
	prod_panel.SetSize(hpanel_width, field_height)

	product_field := GUI.Show_combobox(prod_panel, label_width, field_width, field_height, product_text)
	customer_field := GUI.Show_combobox(prod_panel, label_width, field_width, field_height, customer_text)

	customer_field.SetMarginLeft(inter_spacer_width)
	prod_panel.SetMarginLeft(hpanel_margin)
	prod_panel.SetMarginTop(top_spacer_height)
	prod_panel.Dock(product_field, windigo.Left)
	prod_panel.Dock(customer_field, windigo.Left)

	lot_panel := windigo.NewAutoPanel(product_panel)
	lot_panel.SetSize(hpanel_width, field_height)

	lot_field := GUI.Show_combobox(lot_panel, label_width, field_width, field_height, lot_text)
	sample_field := GUI.Show_combobox(lot_panel, label_width, field_width, field_height, sample_text)

	lot_panel.SetMarginTop(inter_spacer_height)
	lot_panel.SetMarginLeft(hpanel_margin)
	sample_field.SetMarginLeft(inter_spacer_width)
	lot_panel.Dock(lot_field, windigo.Left)
	lot_panel.Dock(sample_field, windigo.Left)

	// reprint_button := windigo.NewPushButton(product_panel)
	// reprint_button.SetText("Reprint")
	// reprint_button.SetMarginsAll(reprint_button_margins)
	// reprint_button.SetMarginLeft(reprint_button_margin_l)
	// reprint_button.SetSize(reprint_button_width, OFF_AXIS)
	//
	// reprint_sample_button := windigo.NewPushButton(product_panel)
	// reprint_sample_button.SetText("Reprint Sample")
	// reprint_sample_button.SetMarginsAll(reprint_button_margins)
	// reprint_sample_button.SetSize(reprint_button_width, OFF_AXIS)

	product_panel.Dock(prod_panel, windigo.Top)
	product_panel.Dock(lot_panel, windigo.Top)
	// product_panel.Dock(reprint_button, windigo.Left)
	// product_panel.Dock(reprint_sample_button, windigo.Left)

	// tabs := windigo.NewTabView(mainWindow)
	// tab_wb := tabs.AddAutoPanel("Water Based")
	// tab_oil := tabs.AddAutoPanel("Oil Based")
	// tab_fr := tabs.AddAutoPanel("Friction Reducer")

	// table := windigo.ListView(mainWindow)
	table := windigo.NewListView(mainWindow)
	table.AddColumn("heelo", 150)
	table.AddColumn("nurse", 50)
	table.DeleteAllItems()

	dock.Dock(product_panel, windigo.Top)
	dock.Dock(table, windigo.Fill)

	// status_bar = windigo.NewStatusBar(mainWindow)
	// mainWindow.SetStatusBar(status_bar)

	// functionality

	rows, err := db_select_product_info.Query()
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

	product_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		// product_field_pop_data(product_field.GetSelectedItem())

		product_id := product_data[product_field.GetSelectedItem()]
		log.Println(select_product_samples(product_id))

	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(func(*windigo.Event) {
		windigo.Exit()
	})
	mainWindow.RunMainLoop() // Must call to start event loop.
}
