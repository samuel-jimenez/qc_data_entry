package product

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/qc_data_entry/util"
	"github.com/samuel-jimenez/qc_data_entry/util/math"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type MeasuredProduct struct {
	QCProduct
	QC_id       int64
	PH          nullable.NullFloat64
	SG          nullable.NullFloat64
	Density     nullable.NullFloat64
	String_test nullable.NullInt64
	Viscosity   nullable.NullInt64
}

func MeasuredProduct_from_new() *MeasuredProduct {
	return new(MeasuredProduct)
}

func MeasuredProduct_from_Row(row DB.Row) (*MeasuredProduct, error) {
	MeasuredProduct := MeasuredProduct_from_new()
	var (
		internal_name        string
		product_moniker_name string
	)
	err := row.Scan(
		&MeasuredProduct.QC_id,
		&MeasuredProduct.Lot_id,
		&internal_name, &product_moniker_name,
		&MeasuredProduct.Lot_number, &MeasuredProduct.Sample_point,
		&MeasuredProduct.PH, &MeasuredProduct.SG, &MeasuredProduct.Density,
		&MeasuredProduct.String_test,
		&MeasuredProduct.Viscosity,
	)
	MeasuredProduct.Product_name = product_moniker_name + " " + internal_name
	return MeasuredProduct, err
}

func MeasuredProduct_Array_from_SQL_Lot_number(Lot_number string) (MeasuredProduct_Array []*MeasuredProduct) {
	proc_name := "MeasuredProduct_Array_from_SQL"
	// MeasuredProduct_Array_from_Query

	DB.Forall_exit(proc_name,
		util.NOOP,
		func(row *sql.Rows) error {
			MeasuredProduct, err := MeasuredProduct_from_Row(row)
			if err != nil {
				return err
			}
			MeasuredProduct_Array = append(MeasuredProduct_Array, MeasuredProduct)
			return nil
		},
		DB.DB_Select_qc_samples_lot, Lot_number)
	return
}

var DupeSampleError = errors.New("measured_product: duplicate sample")

func CheckDupes(Lot_id int64, Sample_point string) (QC_id int64, err error) {
	proc_name := "CheckDupes"
	var qc_id nullable.NullInt64
	if err = DB.Select_Error(proc_name,
		DB.DB_Select_qc_samples_id_lot.QueryRow(
			Lot_id, Sample_point),
		&qc_id,
	); err == nil && qc_id.Valid {
		QC_id = qc_id.Int64
		err = DupeSampleError
	}
	return
}

func (measured_product *MeasuredProduct) Save() (err error) {
	measured_product.QC_id = DB.INVALID_ID

	proc_name := "MeasuredProduct-Save.Sample_point"
	DB.Insert(proc_name, DB.DB_Insel_sample_point, measured_product.Sample_point)

	proc_name = "MeasuredProduct-Save.Tester"
	DB.Insert(proc_name, DB.DB_Insel_qc_tester, measured_product.Tester)

	// Check for dupes
	if measured_product.QC_id, err = CheckDupes(
		measured_product.Lot_id, measured_product.Sample_point,
	); err != nil {
		return
	}

	proc_name = "MeasuredProduct-Save.All"
	if err = DB.Select_Error(proc_name,
		DB.DB_insert_measurement.QueryRow(
			measured_product.Lot_id, measured_product.Sample_point, measured_product.Tester, time.Now().UTC().UnixNano(), measured_product.PH, measured_product.SG, measured_product.Density, measured_product.String_test, measured_product.Viscosity,
		),
		&measured_product.QC_id,
	); err != nil {

		log.Println("error[]%S]:", proc_name, err)
		threads.Show_status("Sample Recording Failed")
		return

	}
	proc_name = "MeasuredProduct-Save.Lot-status"

	if err = DB.Update(proc_name,
		DB.DB_Update_lot_list__status, measured_product.Lot_id, Status_TESTED,
	); err != nil {
		log.Println("error[]%S]:", proc_name, err)
		threads.Show_status("Sample Recording Failed")
		return

	}
	threads.Show_status("Sample Recorded")
	return
}

func Update(products ...*MeasuredProduct) {
	proc_name := "MeasuredProduct-Update"

	for _, measured_product := range products {
		var (
			Lot_id int
			lot_name,
			Sample_point string
			Tester      nullable.NullString
			_timestamp  int64
			PH          nullable.NullFloat64
			SG          nullable.NullFloat64
			Density     nullable.NullFloat64
			String_test nullable.NullInt64
			Viscosity   nullable.NullInt64
		)
		util.LogError(proc_name,
			DB.Select_Error(proc_name,
				DB.DB_Select_qc_samples__old.QueryRow(
					measured_product.QC_id,
				),
				&lot_name, &Lot_id, &Sample_point, &Tester, &_timestamp, &PH, &SG, &Density, &String_test, &Viscosity,
			))

		Time_stamp := time.Unix(0, _timestamp)

		log.Printf("Modify: [UpdateMeasurement] Updating Measurement (%d) %q,%s [%q]: %s: (%s, %s, %s, %s, %s) to [%q]: %s: (%s, %s, %s, %s, %s)", measured_product.QC_id, lot_name, Sample_point, Tester.String,
			Time_stamp.Format(time.DateTime),
			PH.Format(formats.Format_ph), SG.Format(formats.Format_fixed_sg), Density.Format(formats.Format_density), String_test, Viscosity,

			measured_product.Tester.String, time.Now().UTC().Format(time.DateTime), measured_product.PH.Format(formats.Format_ph), measured_product.SG.Format(formats.Format_fixed_sg), measured_product.Density.Format(formats.Format_density), measured_product.String_test, measured_product.Viscosity,
		)

		DB.Update(proc_name,
			DB.DB_Update_qc_samples_measurement,
			measured_product.QC_id,
			measured_product.Tester, time.Now().UTC().UnixNano(), measured_product.PH, measured_product.SG, measured_product.Density, measured_product.String_test, measured_product.Viscosity,
		)

	}
}

func Store(products ...*MeasuredProduct) (err error) {
	numSamples := len(products)
	if numSamples <= 0 {
		return
	}
	measured_product := products[0]
	proc_name := "MeasuredProduct-Store"
	qc_sample_storage_id := measured_product.GetStorage(numSamples)

	for _, product := range products {
		err = product.Save()
		if err != nil {
			continue
		}
		// assign storage to sample
		DB.Update(proc_name,
			DB.DB_Update_qc_samples_storage,
			product.QC_id, qc_sample_storage_id)

	}
	if err != nil {
		return
	}

	// update capacity
	DB.Update(proc_name,
		DB.DB_Update_dec_product_sample_storage_capacity,
		qc_sample_storage_id, numSamples)
	return
}

func (measured_product *MeasuredProduct) Check_data() bool {
	return measured_product.Valid
}

func (measured_product *MeasuredProduct) Output_sample() error {
	measured_product.format_sample()
	return measured_product.export_CoA()
}

func (measured_product *MeasuredProduct) Regen_sample() (err error) {
	err = measured_product.Output_sample()
	if err != nil {
		return
	}
	_, err = measured_product.export_label_pdf()
	return
}

func (measured_product *MeasuredProduct) Save_xl() (bool, error) {
	isomeric_product := strings.Split(measured_product.Product_name, " ")[1]
	if !strings.Contains(measured_product.Product_name_customer, "SLICK") {
		return false, nil // TODO? error{}
	}
	valence_product := strings.Split(measured_product.Product_name_customer, " ")[1]
	return true, updateExcel(config.RETAIN_FILE_NAME, config.RETAIN_WORKSHEET_NAME, measured_product.Lot_number, isomeric_product, valence_product)
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

func Check_single_data(measured_product *MeasuredProduct, store_p, print_p bool) (valid bool, err error) {
	valid = measured_product.Check_data()
	if valid {
		log.Println("debug: Check_check_single_data",
			measured_product)
		if store_p {
			if err = Store(measured_product); err != nil {
				valid = false
				return
			}
		} else {
			if err = measured_product.Save(); err != nil {
				valid = false
				return
			}
		}

		util.LogError("sample_product-Blendshhet", updateBSExcel(config.BLENDSHEET_PATH, measured_product))
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

// TODO dDispathc Check_dual_data
// Print
// .Output_sample())
// if print_p {
// util.LogError("bottom_product-Print", bottom_product.Print())
func Check_dual_data(top_product, bottom_product *MeasuredProduct, print_p bool) (valid bool, err error) {
	// DELTA_DIFF_VISCO := 200
	var DELTA_DIFF_VISCO int64 = 200 // go sucks

	valid = math.Abs(top_product.Viscosity.Diff(bottom_product.Viscosity)) <= DELTA_DIFF_VISCO &&
		top_product.Check_data() && bottom_product.Check_data()

	if valid {

		log.Println("debug: Check_check_dual_data",
			top_product)

		if err = Store(top_product, bottom_product); err != nil {
			valid = false
			return
		}

		util.LogError("top_product-Blendshhet", updateBSExcel(config.BLENDSHEET_PATH, top_product, bottom_product))

		if print_p {
			util.LogError("top_product-Print", top_product.Print())
			log.Println("debug: Check_data",
				bottom_product)

			util.LogError("bottom_product-Print", bottom_product.Print())
		}
		// * Check storage
		bottom_product.CheckStorage()
		// TODO find closest: RMS?
		take_sample, _err := bottom_product.Save_xl()
		util.LogError("bottom_product-Save_xl", _err)
		if !take_sample {
			return
		}
		util.LogError("bottom_product-Output_sample", bottom_product.Output_sample())
		if print_p {
			util.LogError("bottom_product-Print", bottom_product.Print())
		}
		log.Println("debug: Check_data",
			bottom_product)

	}
	return
}
