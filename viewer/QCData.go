package viewer

import (
	"log"
	"slices"
	"strings"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

/*
 * QCData
 *
 *
 *
 */
type QCData struct {
	Product_name,
	Lot_name string
	Sample_point nullable.NullString
	Time_stamp   time.Time
	PH,
	Specific_gravity,
	String_test,
	Viscosity nullable.NullFloat64
}

func compare_product_name(a, b QCData) int { return strings.Compare(a.Product_name, b.Product_name) }
func compare_lot_name(a, b QCData) int     { return strings.Compare(a.Lot_name, b.Lot_name) }
func compare_sample_point(a, b QCData) int { return a.Sample_point.Compare(b.Sample_point) }
func compare_time_stamp(a, b QCData) int   { return a.Time_stamp.Compare(b.Time_stamp) }

func compare_ph(a, b QCData) int               { return a.PH.Compare(b.PH) }
func compare_specific_gravity(a, b QCData) int { return a.Specific_gravity.Compare(b.Specific_gravity) }
func compare_string_test(a, b QCData) int      { return a.String_test.Compare(b.String_test) }
func compare_viscosity(a, b QCData) int        { return a.Viscosity.Compare(b.Viscosity) }
func uno_reverse(fn lessFunc) lessFunc         { return func(a, b QCData) int { return -fn(a, b) } }

func ToString(data nullable.NullFloat64, format func(float64) string) string {
	if data.Valid {
		return format(data.Float64)
	}
	return ""
}

func (data QCData) Product() *product.Product {
	return &product.Product{BaseProduct: product.BaseProduct{
		Product_name: data.Product_name,
		Lot_number:   data.Lot_name,
		Sample_point: data.Sample_point.String,
	},
		PH:          data.PH,
		SG:          data.Specific_gravity,
		Density:     nullable.NewNullFloat64(formats.Density_from_sg(data.Specific_gravity.Float64), data.Specific_gravity.Valid),
		String_test: data.String_test,
		Viscosity:   data.Viscosity}
}

func (data QCData) Text() []string {
	var sg_derived bool
	if data.PH.Valid {
		sg_derived = false
	} else {
		sg_derived = true
	}
	return []string{
		data.Time_stamp.Format(time.DateTime),

		data.Product_name, data.Lot_name,
		data.Sample_point.String,
		ToString(data.PH, formats.Format_ph),
		ToString(data.Specific_gravity, func(sg float64) string { return formats.Format_sg(sg, !sg_derived) }),
		ToString(data.String_test, formats.Format_string_test),
		ToString(data.Viscosity, formats.Format_viscosity)}
}

func (data QCData) ImageIndex() int { return 0 }

// func (data QCData) Checked() bool           { return data.Check }
// func (data QCData) SetChecked(checked bool) { data.Check = checked }
/*
 */

type lessFunc func(a, b QCData) int

/* QCDataView
 *
 */
type QCDataView struct {
	*windigo.ListView
	data []QCData
	less []lessFunc
}

func (data_view *QCDataView) Set(data []QCData) {
	data_view.data = data
}

func (data_view *QCDataView) Get(row, column int) string {
	if row < 0 {
		return ""
	}
	return data_view.data[row].Text()[column]
}

func (data_view QCDataView) Refresh() {

	data_view.DeleteAllItems()

	for _, row := range data_view.data {
		data_view.AddItem(row)
	}
}

func (data_view QCDataView) Update() {
	data_view.Refresh()
}

func (data_view *QCDataView) Sort(col int, asc bool) {
	if asc {
		slices.SortStableFunc(data_view.data, data_view.less[col])
	} else {
		slices.SortStableFunc(data_view.data, uno_reverse(data_view.less[col]))
	}
	data_view.Refresh()
}

func NewQCDataView(parent windigo.Controller) *QCDataView {

	table := &QCDataView{windigo.NewListView(parent), nil, nil}
	table.EnableGridlines(true)
	table.EnableFullRowSelect(true)
	table.EnableDoubleBuffer(true)
	table.EnableSortHeader(true, table.Sort)

	table.AddColumn(
		COL_LABEL_TIME, COL_WIDTH_TIME)
	table.AddColumn(
		"Product", COL_WIDTH_TIME)
	table.AddColumn(
		COL_LABEL_LOT, COL_WIDTH_LOT)
	table.AddColumn(
		COL_LABEL_SAMPLE, COL_WIDTH_SAMPLE)
	table.AddColumn(
		"pH", COL_WIDTH_DATA)
	table.AddColumn(
		"Specific Gravity", COL_WIDTH_DATA)
	table.AddColumn(
		"String Test", COL_WIDTH_DATA)
	table.AddColumn(
		"Viscosity", COL_WIDTH_DATA)
	// table.AddColumn(
	// 	"Density"
	// 	, col_width)
	table.OnClick().Bind(func(e *windigo.Event) { log.Println(e) })
	table.OnRClick().Bind(func(e *windigo.Event) {
		listViewEvent := e.Data.(windigo.ListViewEvent)
		if listViewEvent.Row >= 0 { // ignore invalid
			table.ClipboardCopyText(table.Get(listViewEvent.Row, listViewEvent.Column))
		}
	})

	table.less = []lessFunc{
		compare_time_stamp,
		compare_product_name,
		compare_lot_name,
		compare_sample_point,
		compare_ph,
		compare_specific_gravity,
		compare_string_test,
		compare_viscosity}

	return table
}
