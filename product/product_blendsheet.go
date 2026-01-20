package product

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

// TODO xl-1501
// TODO: move these to own file
func WithOpenFile(file_name string, FN func(*excelize.File) error) error {
	xl_file, err := excelize.OpenFile(file_name, excelize.Options{ShortDatePattern: "yyyymmdd"})
	if err != nil {
		log.Printf("Error: [%s]: %q\n", "WithOpenFile", err)
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := xl_file.Close(); err != nil {
			log.Printf("Error: [%s]: %q\n", "WithOpenFile", err)
		}
	}()
	return FN(xl_file)
}

func updateExcel(file_name, worksheet_name string, row ...string) error {
	return WithOpenFile(file_name, func(xl_file *excelize.File) error {
		// Get all the rows in the worksheet.
		rows, err := xl_file.GetRows(worksheet_name)
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "updateExcel", err)
			return err
		}
		startCell := fmt.Sprintf("A%v", len(rows)+1)

		err = xl_file.SetSheetRow(worksheet_name, startCell, &row)
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "updateExcel", err)
			return err
		}
		return xl_file.Save()
	})
}
