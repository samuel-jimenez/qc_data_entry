package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/blender/inbound"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/util"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var (
	qc_db               *sql.DB
	CONTAINER_CHANGELOG *os.File
	err                 error
)

func dbinit(db *sql.DB) {
	DB.Check_db(db, false)
	DB.DBinit(db)
}

func main() {
	// load config
	config.Main_config = config.Load_config_inbound("qc_data_inbound")
	defer config.Write_config(config.Main_config)

	CONTAINER_CHANGELOG, err = os.Create(config.INBOUND_LOG)
	if err != nil {
		log.Fatalf("Crit: error opening file: %v", err)
	}
	defer CONTAINER_CHANGELOG.Close()

	// log to file
	log_file, err := os.OpenFile(config.LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("Crit: error opening file: %v", err)
	}
	defer log_file.Close()
	log.Println("Info: Logging to logfile:", config.LOG_FILE)

	log.SetOutput(log_file)
	log.Println("Info: Using config:", config.Main_config.ConfigFileUsed())

	// open_db
	// qc_db, err := sql.Open("sqlite3", DB_FILE)
	qc_db, err = sql.Open("sqlite3", ":memory:")
	qc_db.Exec("attach ? as 'bs'", config.DB_FILE)
	if err != nil {
		log.Fatal("Crit: error opening database: ", err)
	}
	defer qc_db.Close()
	log.Println("Info: Using db:", config.DB_FILE)
	dbinit(qc_db)

	// get_sheets(config.PRODUCTION_SCHEDULE_FILE_NAME)

	get_sched(config.PRODUCTION_SCHEDULE_FILE_NAME, config.PRODUCTION_SCHEDULE_WORKSHEET_NAME)
}

// format containers consistently
func format_container(container_name string) (string, product.ProductContainerType) {
	container_type := product.CONTAINER_RAILCAR
	iso_detect_re := regexp.MustCompile(`^ISO `)
	iso_replace_re := regexp.MustCompile(`\s`)
	container_re := regexp.MustCompile(`(?:.*-|[\s])`)

	if iso_detect_re.MatchString(container_name) {
		container_name = iso_detect_re.ReplaceAllString(container_name, "")
		container_name = iso_replace_re.ReplaceAllString(container_name, "")

		container_type = product.CONTAINER_ISO

	} else {
		container_name = container_re.ReplaceAllString(container_name, "")
	}

	if len(container_name) > 4 {
		container_name = fmt.Sprintf("%s %s", container_name[:4], container_name[4:])
	}
	return container_name, container_type
}

func get_sched(file_name, worksheet_name string) {
	proc_name := "InboundSync.DB_Select_inbound_lot_status"
	InboundLotMap0 := product.InboundLotMap_from_Query()
	InboundLotMap1 := make(map[string]*product.InboundLot)
	InboundContainerMap := make(map[string]*product.InboundLot)

	product.WithOpenFile(file_name, func(xl_file *excelize.File) error {
		// status_here := "ARRIVED"
		status_here := []string{"ARRIVED", "OPEN"}
		amount_re := regexp.MustCompile(`[,\s]`)

		// ISO RLTU205426-8

		// Get all the rows in the worksheet.
		rows, err := xl_file.GetRows(worksheet_name)
		if err != nil {
			log.Println(err)
			return err
		}
		// log.Printf("Info: [%s]:  New railcars:  \n", proc_name)
		for _, row := range rows {
			// // info
			// 		for i, col := range row {
			// 	log.Printf("%v: %-20s,\t", i, col)
			// }
			// log.Printf("\n")

			if len(row) <= 1 {
				continue
			}

			// release_date := row[16]
			released_row := 16
			comments_row := 16
			amount_threshold := 5000

			asn := row[1]
			// format containers consistently
			container_name, container_type := format_container(row[2])

			product_name := row[5]
			lot := row[7]
			arrival := row[10]
			status := row[11]
			amount, _ := strconv.Atoi(amount_re.ReplaceAllString(row[12], ""))

			provider := "SNF" // TODO inbound_provider_list

			// sync db to  schedule
			// schedule is authorative source, so if an item does not exist in it, we should remove it

			if slices.Contains(status_here, status) {
				if (lot == "" || lot == "unknown") && arrival != "" && asn != "" {
					lot = fmt.Sprintf("%s/%s", asn, strings.ReplaceAll(arrival, "-", ""))
					// TODO maybe regen asn as BSRC
					provider = "Unknown" // TODO inbound_provider_list
				}
				if false || // just format it good the first time
					amount < amount_threshold ||
					len(row) > released_row && row[released_row] != "" ||
					len(row) > comments_row && strings.Contains(row[comments_row], "release") {
					continue
				}

				// TODO split lot between multiple containers
				if InboundLotMap0[lot] == nil {
					inby := product.InboundLot_from_values(lot, product_name, provider, container_name, container_type, inbound.Status_AVAILABLE)
					if inby == nil {
						log.Printf("error: [%s invalid product]:  %q : %q - %q\n", proc_name, lot, container_name, product_name)
						// invalid product
						// TODO DB_Insert_inbound_product prompt?
						continue
					}
					log.Printf("Info: [%s]:  New %s:  %q : %q - %q\n", proc_name, container_type, lot, container_name, product_name)

					inby.Insert()
					InboundLotMap1[lot] = inby
					InboundContainerMap[container_name] = inby
				} else {
					InboundLotMap1[lot] = InboundLotMap0[lot]
					delete(InboundLotMap0, lot)
				}
			}
		}
		// log.Printf("Info: [%s]:  Departed railcars:  \n", proc_name)
		// items not found as "available"
		for key, val := range InboundLotMap0 {
			if val.Status_name != inbound.Status_UNAVAILABLE {
				// log.Printf("Info: [%s]: %s departed: %q : %q - %q\n", proc_name, val.Container_type, val.Lot_number, val.Container_name, val.Product_name)
				release := fmt.Sprintf("Railcar departed: %q : %q - %q\n", val.Lot_number, val.Container_name, val.Product_name)
				log.Printf("Info: [%s]: %s", proc_name, release)
				CONTAINER_CHANGELOG.WriteString(release)

				if cont := InboundContainerMap[val.Container_name]; cont != nil {
					// log.Printf("Warning: [%s]: %s departed and arrived: %q : %q - %q,  %q : %q - %q\n", proc_name, val.Container_type, val.Lot_number, val.Container_name, val.Product_name, cont.Lot_number, cont.Container_name, cont.Product_name)
					log.Printf("Warning: [%s]: Railcar departed and arrived: %q : %q - %q,  %q : %q - %q\n", proc_name, val.Lot_number, val.Container_name, val.Product_name, cont.Lot_number, cont.Container_name, cont.Product_name)
					// TODO take input, possibly rename lot
				}
				val.Update_status(inbound.Status_UNAVAILABLE)
				if err := product.Release_testing_lot(val.Lot_number); err != nil {
					log.Println("error[]%S]:", proc_name, err)
					return err
				}
			}
			delete(InboundLotMap0, key)
		}

		// get all SAMPLED
		proc_name = "InboundSync.DB_Select_inbound_lot_status.SAMPLED"
		DB.Forall_err(proc_name,
			util.NOOP,
			func(row *sql.Rows) error {
				var lot_name string
				if err := row.Scan(
					&lot_name,
				); err != nil {
					return err
				}
				// queries cannot be nested, so just dump the results into the map
				InboundLotMap0[lot_name] = InboundLotMap1[lot_name]
				return nil
			},
			DB.DB_Select_name_inbound_lot_status, inbound.Status_SAMPLED)
		// process new entries
		for key, val := range InboundLotMap0 {
			val.Quality_test()
			delete(InboundLotMap0, key)
		}

		// get all tested
		proc_name = "InboundSync.DB_Select_inbound_lot_status.Status_TESTED"
		log.Printf("Info: [%s]:  Tested railcars:  \n", proc_name)
		DB.Forall_err(proc_name,
			util.NOOP,
			func(row *sql.Rows) error {
				val, err := product.InboundLot_from_Row(row)
				if err != nil {
					return err
				}

				log.Printf("Info: [%s]:  %q : %q - %q\n", proc_name, val.Lot_number, val.Container_name, val.Product_name)
				return nil
			},
			DB.DB_Select_inbound_lot_status, inbound.Status_TESTED)

		// get all available
		proc_name = "InboundSync.DB_Select_inbound_lot_status.AVAILABLE"
		DB.Forall_err(proc_name,
			util.NOOP,
			func(row *sql.Rows) error {
				var lot_name string
				if err := row.Scan(
					&lot_name,
				); err != nil {
					return err
				}
				// queries cannot be nested, so just dump the results into the map
				InboundLotMap0[lot_name] = InboundLotMap1[lot_name]
				return nil
			},
			DB.DB_Select_name_inbound_lot_status, inbound.Status_AVAILABLE)

		// process new entries
		log.Printf("Info: [%s]:  Available railcars:  \n", proc_name)
		CONTAINER_CHANGELOG.WriteString("\n\tAvailable railcars:  \n")
		for _, val := range InboundLotMap0 {
			arrival := fmt.Sprintf("\t\t%q : %q - %q\n", val.Lot_number, val.Container_name, val.Product_name)
			log.Printf("Info: [%s]: %s", proc_name, arrival)
			CONTAINER_CHANGELOG.WriteString(arrival)
		}

		return nil
	})
}
