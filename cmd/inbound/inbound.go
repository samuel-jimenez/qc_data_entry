package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/xuri/excelize/v2"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/blender/blendbound"
	"github.com/samuel-jimenez/qc_data_entry/config"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var (
	qc_db *sql.DB
)

func dbinit(db *sql.DB) {

	DB.Check_db(db)
	DB.DBinit(db)

}

func main() {

	//load config
	config.Main_config = config.Load_config_inbound("qc_data_inbound")

	// log to file
	log_file, err := os.OpenFile(config.LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Crit: error opening file: %v", err)
	}
	defer log_file.Close()

	log.SetOutput(log_file)
	log.Println("Info: Logging Started")
	log.Println("Info: config.DB_FILE", config.DB_FILE)

	//open_db
	// qc_db, err := sql.Open("sqlite3", DB_FILE)
	qc_db, err = sql.Open("sqlite3", ":memory:")
	qc_db.Exec("attach ? as 'bs'", config.DB_FILE)
	if err != nil {
		log.Fatal("Crit: error opening database: ", err)
	}
	defer qc_db.Close()
	dbinit(qc_db)

	// get_sheets(config.PRODUCTION_SCHEDULE_FILE_NAME)

	get_sched(config.PRODUCTION_SCHEDULE_FILE_NAME, config.PRODUCTION_SCHEDULE_WORKSHEET_NAME)

}

func withOpenFile(file_name string, FN func(*excelize.File) error) error {
	xl_file, err := excelize.OpenFile(file_name)
	if err != nil {
		log.Printf("Error: [%s]: %q\n", "withOpenFile", err)
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := xl_file.Close(); err != nil {
			log.Printf("Error: [%s]: %q\n", "withOpenFile", err)
		}
	}()
	return FN(xl_file)
}

func get_sheets(file_name string) {
	withOpenFile(file_name, func(xl_file *excelize.File) error {
		for index, name := range xl_file.GetSheetMap() {
			log.Println(index, name)
			visible, _ := xl_file.GetSheetVisible(name)
			log.Println(visible)

		}
		return nil
	})
}

func get_sched(file_name, worksheet_name string) {
	InboundLotMap0 := blendbound.NewInboundLotMapFromQuery()
	InboundLotMap1 := make(map[string]*blendbound.InboundLot)

	withOpenFile(file_name, func(xl_file *excelize.File) error {
		status_here := "ARRIVED"

		// Get all the rows in the worksheet.
		rows, err := xl_file.GetRows(worksheet_name)
		if err != nil {
			log.Println(err)
			return err
		}
		for _, row := range rows {
			// // info
			// 		for i, col := range row {
			// 	log.Printf("%v: %-20s,\t", i, col)
			// }
			// log.Printf("\n")

			// release_date := row[16]
			released_row := 16

			container := row[2]
			product := row[5]
			lot := row[7]
			status := row[11]
			provider := "SNF" // TODO inbound_provider_list

			// sync db to  schedule
			//schedule is authorative source, so if an item does not exist in it, we should remove it
			if status == status_here {
				if len(row) > released_row && row[released_row] != "" && InboundLotMap0[lot] != nil {
					continue
				}

				if InboundLotMap0[lot] == nil {
					inby := blendbound.NewInboundLotFromValues(lot, product, provider, container, blendbound.Status_AVAILABLE)
					if inby == nil {
						log.Printf("error: [%s invalid product]: %q\n", "NewInboundLotFromValues", product)
						// invalid product
						continue
					}
					inby.Insert()
					InboundLotMap1[lot] = inby
				} else {
					InboundLotMap1[lot] = InboundLotMap0[lot]
					delete(InboundLotMap0, lot)
				}
			}
		}

		// items not found as "available"
		for key, val := range InboundLotMap0 {
			val.Update_status(blendbound.Status_UNAVAILABLE)
			delete(InboundLotMap0, key)
		}

		// get all available
		proc_name := "InboundSync DB_Select_inbound_lot_status"
		DB.Forall(proc_name,
			func() {},
			func(row *sql.Rows) {
				var lot_name string
				if err := row.Scan(&lot_name); err != nil {
					log.Printf("error: [%s]: %q\n", proc_name, err)
					return
				} else {
					// queries cannot be nested, so just dump the results into the map
					InboundLotMap0[lot_name] = InboundLotMap1[lot_name]
				}
			},
			DB.DB_Select_inbound_lot_status, blendbound.Status_AVAILABLE)
		// process new entries
		for _, val := range InboundLotMap0 {
			val.Quality_test()
		}
		return nil
	})
}
