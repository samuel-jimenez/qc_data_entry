package product

import (
	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/util"

	"github.com/samuel-jimenez/qc_data_entry/datatypes"
	// "github.com/samuel-jimenez/whatsupdocx/docx"
)

/*
 * QCProduct
 *
 */

type QCProduct struct {
	BaseProduct
	Appearance     ProductAppearance
	Product_type   Discrete
	Container_type ProductContainerType // bs.container_types
	PH             datatypes.Range
	SG             datatypes.Range
	Density        datatypes.Range
	String_test    datatypes.Range
	Viscosity      datatypes.Range
	UpdateFN       func(*QCProduct)
}

func QCProduct_from_new() *QCProduct {
	qc_product := new(QCProduct)
	qc_product.Product_id = DB.INVALID_ID
	qc_product.Product_Lot_id = DB.DEFAULT_LOT_ID
	qc_product.Lot_id = DB.DEFAULT_LOT_ID
	return qc_product
}

func (product QCProduct) Base() QCProduct {
	return product
}

func (qc_product *QCProduct) ResetQC() {
	var empty_product QCProduct
	empty_product.BaseProduct = qc_product.BaseProduct
	empty_product.UpdateFN = qc_product.UpdateFN
	*qc_product = empty_product
}

func (qc_product *QCProduct) Select_product_details() {
	proc_name := "Select_product_details"
	DB.Select_ErrNoRows(
		proc_name,
		DB.DB_Select_product_details.QueryRow(qc_product.Product_id),
		&qc_product.Product_type, &qc_product.Container_type, &qc_product.Appearance,
		&qc_product.PH.Method, &qc_product.PH.Valid, &qc_product.PH.Publish_p, &qc_product.PH.Min, &qc_product.PH.Target, &qc_product.PH.Max,
		&qc_product.SG.Method, &qc_product.SG.Valid, &qc_product.SG.Publish_p, &qc_product.SG.Min, &qc_product.SG.Target, &qc_product.SG.Max,
		&qc_product.Density.Method, &qc_product.Density.Valid, &qc_product.Density.Publish_p, &qc_product.Density.Min, &qc_product.Density.Target, &qc_product.Density.Max,
		&qc_product.String_test.Method, &qc_product.String_test.Valid, &qc_product.String_test.Publish_p, &qc_product.String_test.Min, &qc_product.String_test.Target, &qc_product.String_test.Max,
		&qc_product.Viscosity.Method, &qc_product.Viscosity.Valid, &qc_product.Viscosity.Publish_p, &qc_product.Viscosity.Min, &qc_product.Viscosity.Target, &qc_product.Viscosity.Max,
	)
}

func (qc_product QCProduct) Upsert() {
	var err error
	proc_name := "QCProduct-Upsert"
	DB.DB_Insert_appearance.Exec(qc_product.Appearance)
	DB.DB_Insel_test_method.Exec(qc_product.PH.Method)
	DB.DB_Insel_test_method.Exec(qc_product.SG.Method)
	DB.DB_Insel_test_method.Exec(qc_product.Density.Method)
	DB.DB_Insel_test_method.Exec(qc_product.String_test.Method)
	DB.DB_Insel_test_method.Exec(qc_product.Viscosity.Method)

	DB.DB_Upsert_product_type.Exec(qc_product.Product_id, qc_product.Product_type, qc_product.Appearance)
	_, err = DB.DB_Upsert_product_details.Exec(
		qc_product.Product_id, formats.PH_TEXT,
		qc_product.PH.Method, qc_product.PH.Valid, qc_product.PH.Publish_p, qc_product.PH.Min, qc_product.PH.Target, qc_product.PH.Max,
	)
	util.LogError(proc_name, err)
	_, err = DB.DB_Upsert_product_details.Exec(
		qc_product.Product_id, formats.SG_TEXT,
		qc_product.SG.Method, qc_product.SG.Valid, qc_product.SG.Publish_p, qc_product.SG.Min, qc_product.SG.Target, qc_product.SG.Max,
	)
	util.LogError(proc_name, err)
	_, err = DB.DB_Upsert_product_details.Exec(
		qc_product.Product_id, formats.DENSITY_TEXT,
		qc_product.Density.Method, qc_product.Density.Valid, qc_product.Density.Publish_p, qc_product.Density.Min, qc_product.Density.Target, qc_product.Density.Max,
	)
	util.LogError(proc_name, err)
	_, err = DB.DB_Upsert_product_details.Exec(
		qc_product.Product_id, formats.STRING_TEXT_MINI,
		qc_product.String_test.Method, qc_product.String_test.Valid, qc_product.String_test.Publish_p, qc_product.String_test.Min, qc_product.String_test.Target, qc_product.String_test.Max,
	)
	util.LogError(proc_name, err)
	_, err = DB.DB_Upsert_product_details.Exec(
		qc_product.Product_id, formats.VISCOSITY_TEXT,
		qc_product.Viscosity.Method, qc_product.Viscosity.Valid, qc_product.Viscosity.Publish_p, qc_product.Viscosity.Min, qc_product.Viscosity.Target, qc_product.Viscosity.Max,
	)
	util.LogError(proc_name, err)
}

func (qc_product *QCProduct) Edit(
	product_type Discrete,
	Appearance ProductAppearance,
	PH,
	SG,
	Density,
	String_test,
	Viscosity datatypes.Range,
) {
	qc_product.Product_type = product_type
	qc_product.Appearance = Appearance
	qc_product.PH = PH
	qc_product.SG = SG
	qc_product.Density = Density
	qc_product.String_test = String_test
	qc_product.Viscosity = Viscosity
}

func (qc_product QCProduct) Check(data MeasuredProduct) bool {
	return qc_product.PH.Check(data.PH.Float64) && qc_product.SG.Check(data.SG.Float64)
}

func (qc_product *QCProduct) SetUpdate(UpdateFN func(*QCProduct)) {
	qc_product.UpdateFN = UpdateFN
}

func (qc_product *QCProduct) Update() {
	if qc_product.UpdateFN == nil {
		return
	}
	qc_product.UpdateFN(qc_product)
}
