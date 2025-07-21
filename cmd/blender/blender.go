package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/windigo"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var (
	qc_db *sql.DB
)

func main() {
	//load config
	config.Main_config = config.Load_config_viewer("qc_data_blender")

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
	qc_db, err = sql.Open("sqlite3", ":memory:")
	qc_db.Exec("attach ? as 'bs'", config.DB_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer qc_db.Close()

	dbinit(qc_db)
	/*
		//setup print goroutine
		threads.PRINT_QUEUE = make(chan string, 4)
		defer close(threads.PRINT_QUEUE)
		go threads.Do_print_queue(threads.PRINT_QUEUE)*/

	//setup status_bar goroutine
	threads.STATUS_QUEUE = make(chan string, 16)
	defer close(threads.STATUS_QUEUE)
	go threads.Do_status_queue(threads.STATUS_QUEUE)

	//show main window
	show_window()

}

func dbinit(db *sql.DB) {

	DB.Check_db(db)
	DB.DBinit(db)
}

func show_window() {

	log.Println("Info: Process started")
	window_title := "QC Data Blender"

	// build window
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetText(window_title)

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(func(*windigo.Event) {
		windigo.Exit()
	})
	mainWindow.RunMainLoop() // Must call to start event loop.
}
