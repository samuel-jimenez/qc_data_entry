package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/qc"
	"github.com/samuel-jimenez/qc_data_entry/threads"
)

func main() {
	//load config
	config.Main_config = config.Load_config_entry("qc_data_entry")

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
	qc_db, err := sql.Open("sqlite3", ":memory:")
	qc_db.Exec("attach ? as 'bs'", config.DB_FILE)
	if err != nil {
		log.Fatal("Crit: error opening database: ", err)
	}
	defer qc_db.Close()
	qc.DBinit(qc_db)

	//setup print goroutine
	threads.PRINT_QUEUE = make(chan string, 4)
	defer close(threads.PRINT_QUEUE)
	go threads.Do_print_queue(threads.PRINT_QUEUE)

	//setup status_bar goroutine
	threads.STATUS_QUEUE = make(chan string, 16)
	defer close(threads.STATUS_QUEUE)
	go threads.Do_status_queue(threads.STATUS_QUEUE)

	//show main window
	qc.Show_window()

}
