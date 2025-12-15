package qc_ui

import (
	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

/*
 * MeasuredProductDiffView
 *
 */
type MeasuredProductDiffView struct {
	*windigo.AutoPanel
	oldProd, newProd *MeasuredProductView
}

func MeasuredProductDiffView_from_new(parent windigo.Controller, old_Product, new_Product *product.MeasuredProduct) *MeasuredProductDiffView {
	view := new(MeasuredProductDiffView)
	view.AutoPanel = windigo.NewAutoPanel(parent)
	view.SetPaddingsAll(5)

	view.oldProd = MeasuredProductView_from_new(view, old_Product)
	view.oldProd.Checked.Hide()
	view.oldProd.labels[0].SetText("existing")
	view.oldProd.labels[1].SetText(old_Product.Lot_number)

	view.newProd = MeasuredProductView_from_new(view, new_Product)
	view.newProd.Checked.Hide()
	view.newProd.labels[0].SetText("new")
	view.newProd.labels[1].SetText(old_Product.Lot_number)

	view.Dock(view.oldProd, windigo.Left)
	view.Dock(view.newProd, windigo.Left)
	return view
}

func MeasuredProductDiffView_from_combined(parent windigo.Controller, Measured_Product *product.MeasuredProduct) *MeasuredProductDiffView {
	old_Product, _ := product.MeasuredProduct_from_Row(DB.DB_Select_qc_samples__id.QueryRow(Measured_Product.QC_id))
	return MeasuredProductDiffView_from_new(parent, old_Product, Measured_Product)
}

func MeasuredProductDiffView_from_ids(parent windigo.Controller, old_QC_id, new_QC_id int64) (*MeasuredProductDiffView, string) {
	old_Product, _ := product.MeasuredProduct_from_Row(DB.DB_Select_qc_samples__id.QueryRow(old_QC_id))
	new_Product, _ := product.MeasuredProduct_from_Row(DB.DB_Select_qc_samples__id.QueryRow(new_QC_id))
	return MeasuredProductDiffView_from_new(parent, old_Product, new_Product), new_Product.Sample_point
}

func (view *MeasuredProductDiffView) SetFont(font *windigo.Font) {
	view.oldProd.SetFont(font)
	view.newProd.SetFont(font)
}

func (view *MeasuredProductDiffView) RefreshSize() {
	view.SetSize(2*GUI.GROUP_WIDTH, GUI.GROUP_HEIGHT)
	view.SetPaddingsAll(8)
	view.oldProd.RefreshSize()
	view.newProd.RefreshSize()
}
