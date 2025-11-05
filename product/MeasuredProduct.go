package product

import (
	"log"
	"strings"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/qc_data_entry/util"
	"github.com/samuel-jimenez/qc_data_entry/util/math"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type MeasuredProduct struct {
	QCProduct
	PH          nullable.NullFloat64
	SG          nullable.NullFloat64
	Density     nullable.NullFloat64
	String_test nullable.NullInt64
	Viscosity   nullable.NullInt64
}

func MeasuredProduct_from_new() *MeasuredProduct {
	return new(MeasuredProduct)
}

func (measured_product *MeasuredProduct) Save() int64 {
	proc_name := "MeasuredProduct-Save.Sample_point"
	DB.Insert(proc_name, DB.DB_Insel_sample_point, measured_product.Sample_point)

	proc_name = "MeasuredProduct-Save.Tester"
	DB.Insert(proc_name, DB.DB_Insel_qc_tester, measured_product.Tester)

	var qc_id int64
	proc_name = "MeasuredProduct-Save.All"
	if err := DB.Select_Error(proc_name,
		DB.DB_insert_measurement.QueryRow(
			measured_product.Lot_id, measured_product.Sample_point, measured_product.Tester, time.Now().UTC().UnixNano(), measured_product.PH, measured_product.SG, measured_product.Density, measured_product.String_test, measured_product.Viscosity,
		),
		&qc_id,
	); err != nil {
		log.Println("error[]%S]:", proc_name, err)
		threads.Show_status("Sample Recording Failed")
		return DB.INVALID_ID

	}
	proc_name = "MeasuredProduct-Save.Lot-status"

	if err := DB.Update(proc_name,
		DB.DB_Update_lot_list__status, measured_product.Lot_id, Status_TESTED,
	); err != nil {
		log.Println("error[]%S]:", proc_name, err)
		threads.Show_status("Sample Recording Failed")
		return DB.INVALID_ID

	}
	threads.Show_status("Sample Recorded")
	return qc_id
}

func Store(products ...*MeasuredProduct) {
	numSamples := len(products)
	if numSamples <= 0 {
		return
	}
	measured_product := products[0]
	proc_name := "MeasuredProduct-Store"
	qc_sample_storage_id := measured_product.GetStorage(numSamples)

	for _, product := range products {
		qc_id := product.Save()
		// assign storage to sample
		DB.Update(proc_name,
			DB.DB_Update_qc_samples_storage,
			qc_id, qc_sample_storage_id)

	}

	// update capacity
	DB.Update(proc_name,
		DB.DB_Update_dec_product_sample_storage_capacity,
		qc_sample_storage_id, numSamples)
}

func (measured_product *MeasuredProduct) Check_data() bool {
	return measured_product.Valid
}

func (measured_product *MeasuredProduct) Output_sample() error {
	measured_product.format_sample()
	return measured_product.export_CoA()
}

func (measured_product *MeasuredProduct) Printout() error {
	proc_name := "MeasuredProduct-Printout"
	if err := measured_product.Print(); err != nil {
		log.Printf("Error [%s]: %q\n", proc_name, err)
		log.Printf("Debug [%s]: %q: %v\n", proc_name, err, measured_product) // TODO
		return err
	}
	return measured_product.Output_sample()
}

func (measured_product *MeasuredProduct) Save_xl() error {
	isomeric_product := strings.Split(measured_product.Product_name, " ")[1]
	valence_product := strings.Split(measured_product.Product_name_customer, " ")[1]
	return updateExcel(config.RETAIN_FILE_NAME, config.RETAIN_WORKSHEET_NAME, measured_product.Lot_number, isomeric_product, valence_product)
}

func (measured_product *MeasuredProduct) Print() error {
	pdf_path, err := measured_product.export_label_pdf()
	if err != nil {
		return err
	}
	threads.Show_status("Label Created")

	Print_PDF(pdf_path)

	return err
}

func Check_single_data(measured_product *MeasuredProduct, store_p, print_p bool) (valid bool) {
	valid = measured_product.Check_data()
	if valid {
		log.Println("debug: Check_check_single_data",
			measured_product)
		if store_p {
			Store(measured_product)
		} else {
			measured_product.Save()
		}
		if print_p {
			util.LogError("sample_product-Print", measured_product.Print())
		}
		if store_p {
			util.LogError("sample_product-Output_sample", measured_product.Output_sample())
			measured_product.CheckStorage()
		}
	}
	return
}

func Check_dual_data(top_product, bottom_product *MeasuredProduct, print_p bool) (valid bool) {
	// DELTA_DIFF_VISCO := 200
	var DELTA_DIFF_VISCO int64 = 200 // go sucks

	valid = math.Abs(top_product.Viscosity.Diff(bottom_product.Viscosity)) <= DELTA_DIFF_VISCO &&
		top_product.Check_data() && bottom_product.Check_data()

	if valid {

		log.Println("debug: Check_check_dual_data",
			top_product)

		Store(top_product, bottom_product)

		if print_p {
			util.LogError("top_product-Print", top_product.Print())
			log.Println("debug: Check_data",
				bottom_product)

			util.LogError("bottom_product-Print", bottom_product.Print())
		}
		// TODO find closest: RMS?
		util.LogError("bottom_product-Output_sample", bottom_product.Output_sample())
		if print_p {
			util.LogError("bottom_product-Print", bottom_product.Print())
		}
		log.Println("debug: Check_data",
			bottom_product)

		util.LogError("bottom_product-Save_xl", bottom_product.Save_xl())

		// * Check storage
		bottom_product.CheckStorage()

	}
	return
}
