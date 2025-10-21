package viewer

import (
	"database/sql"
	"slices"
	"strings"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

// prevents 'hash of unhashable type'
type QCDataComponents struct {
	Components []blender.BlendComponent
}

func (data *QCDataComponents) GetComponents(Lot_name string) {
	proc_name := "QCDataComponents.GetComponents"
	DB.Forall_exit(proc_name,
		func() {
			data.Components = nil
		},
		func(row *sql.Rows) error {
			blendComponent, err := blender.BlendComponent_from_SQL(row)
			if err != nil {
				return err
			}
			data.Components = append(data.Components, *blendComponent)
			return nil
		},
		DB.DB_Select_product_lot_list_sources, Lot_name)

}

func (data QCDataComponents) Text() []string {
	Text := []string{}
	for _, blendComponent := range data.Components {
		Text = append(Text,
			blendComponent.Text()...)
	}

	return Text

}

//TODO tester
/*
 * QCData
 *
 *
 *
 */
type QCData struct {
	Product_name,
	Lot_name string
	Product_name_customer,
	Sample_point,
	Sample_bin nullable.NullString
	Time_stamp time.Time
	PH,
	Specific_gravity nullable.NullFloat64
	String_test,
	Viscosity nullable.NullInt64
	// prevents 'hash of unhashable type'
	Components *QCDataComponents
}

func (data *QCData) GetComponents() {
	data.Components = new(QCDataComponents)
	data.Components.GetComponents(data.Lot_name)
}

func compare_product_name(a, b QCData) int { return strings.Compare(a.Product_name, b.Product_name) }
func compare_lot_name(a, b QCData) int     { return strings.Compare(a.Lot_name, b.Lot_name) }
func compare_sample_point(a, b QCData) int { return a.Sample_point.Compare(b.Sample_point) }
func compare_sample_bin(a, b QCData) int   { return a.Sample_bin.Compare(b.Sample_bin) }
func compare_time_stamp(a, b QCData) int   { return a.Time_stamp.Compare(b.Time_stamp) }

func compare_ph(a, b QCData) int               { return a.PH.Compare(b.PH) }
func compare_specific_gravity(a, b QCData) int { return a.Specific_gravity.Compare(b.Specific_gravity) }
func compare_string_test(a, b QCData) int      { return a.String_test.Compare(b.String_test) }
func compare_viscosity(a, b QCData) int        { return a.Viscosity.Compare(b.Viscosity) }

func ValueOrNothing(array []string, i int) string {
	if len(array) > i {
		return array[i]
	}
	return ""
}

func compare_component(a, b QCData, i int) int {
	return strings.Compare(
		ValueOrNothing(a.Components.Text(), i),
		ValueOrNothing(b.Components.Text(), i),
	)
}

func compare_component_0_name(a, b QCData) int      { return compare_component(a, b, 0) }
func compare_component_0_lot(a, b QCData) int       { return compare_component(a, b, 1) }
func compare_component_0_container(a, b QCData) int { return compare_component(a, b, 2) }
func compare_component_1_name(a, b QCData) int      { return compare_component(a, b, 3) }
func compare_component_1_lot(a, b QCData) int       { return compare_component(a, b, 4) }
func compare_component_1_container(a, b QCData) int { return compare_component(a, b, 5) }

func uno_reverse(fn lessFunc) lessFunc { return func(a, b QCData) int { return -fn(a, b) } }

func ToString(data nullable.NullFloat64, format func(float64) string) string {
	if data.Valid {
		return format(data.Float64)
	}
	return ""
}

func (data QCData) Product() *product.Product {
	return &product.Product{BaseProduct: product.BaseProduct{
		Product_name:          data.Product_name,
		Lot_number:            data.Lot_name,
		Sample_point:          data.Sample_point.String,
		Product_name_customer: data.Product_name_customer.String,
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

	return append([]string{
		data.Time_stamp.Format(time.DateTime),

		data.Product_name, data.Lot_name,
		data.Sample_point.String,
		data.Sample_bin.String,
		ToString(data.PH, formats.Format_ph),
		ToString(data.Specific_gravity, func(sg float64) string { return formats.Format_sg(sg, !sg_derived) }),
		// ToString(data.String_test, formats.Format_string_test),
		data.String_test.String(),
		// ToString(data.Viscosity, formats.Format_viscosity),
		data.Viscosity.String(),
	},
		data.Components.Text()...)

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
		COL_LABEL_SAMPLE_PT, COL_WIDTH_SAMPLE_PT)
	table.AddColumn(
		COL_LABEL_SAMPLE_BIN, COL_WIDTH_SAMPLE_BIN)
	table.AddColumn(
		"pH", COL_WIDTH_DATA)
	table.AddColumn(
		"Specific Gravity", COL_WIDTH_DATA)
	table.AddColumn(
		"String Test", COL_WIDTH_DATA)
	table.AddColumn(
		"Viscosity", COL_WIDTH_DATA)
	table.AddColumn(
		"Component 0", COL_WIDTH_TIME)
	table.AddColumn(
		"Lot 0 ", COL_WIDTH_LOT)
	table.AddColumn(
		"Container 0 ", COL_WIDTH_TIME)
	table.AddColumn(
		"Component 1", COL_WIDTH_TIME)
	table.AddColumn(
		"Lot 1 ", COL_WIDTH_LOT)
	table.AddColumn(
		"Container 1 ", COL_WIDTH_TIME)
	// table.AddColumn(
	// 	"Density"
	// 	, col_width)

	table.OnClick().Bind(func(e *windigo.Event) {
		listViewEvent := e.Data.(windigo.ListViewEvent)
		if listViewEvent.Row >= 0 { // ignore invalid
			table.ClipboardCopyText(table.GetCell(listViewEvent.Row, listViewEvent.Column))
		}
	})

	popupMenu := windigo.NewContextMenu()
	copyMenu := popupMenu.AddItem("Copy", windigo.Shortcut{
		Modifiers: windigo.ModControl,
		Key:       windigo.KeyC,
	})
	table.SetContextMenu(popupMenu)
	copyMenu.OnClick().Bind(func(e *windigo.Event) {
		table.ClipboardCopyText(table.GetText())
	})

	table.less = []lessFunc{
		compare_time_stamp,
		compare_product_name,
		compare_lot_name,
		compare_sample_point,
		compare_sample_bin,
		compare_ph,
		compare_specific_gravity,
		compare_string_test,
		compare_viscosity,
		compare_component_0_name,
		compare_component_0_lot,
		compare_component_0_container,
		compare_component_1_name,
		compare_component_1_lot,
		compare_component_1_container,
	}

	return table
}

func (data_view *QCDataView) Set(data []QCData) {
	data_view.data = data
}

func (data_view *QCDataView) Get() []QCData {
	return data_view.data
}

func (data_view *QCDataView) GetCell(row, column int) string {
	if row < 0 {
		return ""
	}
	rows := data_view.data[row].Text()
	if len(rows) <= column {
		return ""
	}
	return rows[column]
}

func (data_view *QCDataView) GetText() string {
	var rowText []string
	for _, row := range data_view.SelectedItems() {
		rowText = append(rowText, strings.Join(row.Text(), "\t"))
	}
	return strings.Join(rowText, "\n")
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

func (view *QCDataView) Sort(col int, asc bool) {
	if asc {
		slices.SortStableFunc(view.data, view.less[col])
	} else {
		slices.SortStableFunc(view.data, uno_reverse(view.less[col]))
	}
	view.Refresh()
}

func (table *QCDataView) RefreshSize() {
	widths :=
		[]int{
			COL_WIDTH_TIME,
			COL_WIDTH_TIME,
			COL_WIDTH_LOT,
			COL_WIDTH_SAMPLE_PT,
			COL_WIDTH_SAMPLE_BIN,
			COL_WIDTH_DATA,
			COL_WIDTH_DATA,
			COL_WIDTH_DATA,
			COL_WIDTH_DATA,
			COL_WIDTH_TIME,
			COL_WIDTH_LOT,
			COL_WIDTH_TIME,
			COL_WIDTH_TIME,
			COL_WIDTH_LOT,
			COL_WIDTH_TIME,
		}
	for i := range table.GetNumColumns() {
		table.SetColumnWidth(i, widths[i])

	}
}
