package blendbound

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"

	"codeberg.org/go-pdf/fpdf"
	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/util"
)

// TODO type00 export type
// bs.inbound_status_list
type Status string

const (
	Status_AVAILABLE   Status = "AVAILABLE"
	Status_SAMPLED            = "SAMPLED"
	Status_TESTED             = "TESTED"
	Status_UNAVAILABLE        = "UNAVAILABLE"
)

type InboundLot struct {
	Lot_id         int64
	Lot_number     string
	Product_id     int64
	Product_name   string
	Provider_id    int64
	provider_name  string
	Container_id   int64
	Container_name string
	Status_id      int64
	Status_name    Status

	// We don't use this yet.
	// Container_type product.ProductContainerType
}

func NewInboundLot() *InboundLot { return new(InboundLot) }

func NewInboundLotFromValues(Lot_number, product_name, provider_name, container_name string, Container_type product.ProductContainerType, status_name Status) *InboundLot {
	if Lot_number == "" {
		return nil
	}
	proc_name := "NewInboundLotFromValues"

	Inbound := NewInboundLot()
	Inbound.Lot_number = Lot_number
	Inbound.Product_name = product_name
	if err := DB.DB_Select_inbound_product_name.QueryRow(Inbound.Product_name).Scan(
		&Inbound.Product_id,
	); err != nil {
		if err != sql.ErrNoRows { // no row? no problem!
			log.Printf("Error: [%s]: %q\n", proc_name, err)
			//TODO DB_Insert_inbound_product
		}
		return nil
	}

	Inbound.provider_name = provider_name
	Inbound.Provider_id = DB.Insel("NewInboundLotFromValues inbound_provider", DB.DB_Insert_inbound_provider, DB.DB_Select_inbound_provider_id, Inbound.provider_name)

	Inbound.Container_name = container_name
	Inbound.Container_id = DB.Insel("NewInboundLotFromValues container", DB.DB_Insert_container, DB.DB_Select_container_id, Inbound.Container_name)

	// We don't use this yet.
	// Inbound.Container_type = Container_type
	// DB.Update("NewInboundLotFromValues Container_type", DB.DB_Update_container_type, Inbound.Container_id, Inbound.Container_type)
	DB.Update("NewInboundLotFromValues Container_type", DB.DB_Update_container_type, Inbound.Container_id, Container_type)

	Inbound.Status_name = status_name
	if err := DB.Select_Error("NewInboundLotFromValues status", DB.DB_Select_name_inbound_status_list.QueryRow(Inbound.Status_name), &Inbound.Status_id); err != nil {
		return nil
	}

	return Inbound
}

func NewInboundLotFromRow(row *sql.Rows) (*InboundLot, error) {
	Inbound := NewInboundLot()
	err := row.Scan(
		&Inbound.Lot_id, &Inbound.Lot_number,
		&Inbound.Product_id, &Inbound.Product_name,
		&Inbound.Provider_id, &Inbound.provider_name,
		&Inbound.Container_id, &Inbound.Container_name,
		&Inbound.Status_id, &Inbound.Status_name,
	)
	return Inbound, err
}

func NewInboundLotFromBlendComponent(blendComponent *blender.BlendComponent) *InboundLot {
	if !blendComponent.Inboundp {
		return nil
	}
	Inbound := NewInboundLot()
	Inbound.Lot_id = blendComponent.Lot_id
	Inbound.Lot_number = blendComponent.Lot_name
	Inbound.Container_name = blendComponent.Container_name
	Inbound.Product_name = blendComponent.Component_name
	return Inbound
}

func NewInboundLotMapFromQuery() map[string]*InboundLot {
	InboundLotMap := make(map[string]*InboundLot)
	proc_name := "NewInboundLotArrayFromQuery"

	DB.Forall_err(proc_name,
		func() {},
		func(row *sql.Rows) error {
			Inbound, err := NewInboundLotFromRow(row)
			if err != nil {
				return err
			}
			InboundLotMap[Inbound.Lot_number] = Inbound
			return nil
		},
		DB.DB_Select_inbound_lot_all)
	return InboundLotMap
}

func (object *InboundLot) Insert() {
	proc_name := "InboundLot.Insert"
	object.Lot_id = DB.Insert(proc_name, DB.DB_Insert_inbound_lot, object.Lot_number, object.Product_id, object.Provider_id, object.Container_id)

}

func (object *InboundLot) Update_status(status Status) {
	proc_name := "InboundLot.Update_status"
	if err := DB.Select_Error("NewInboundLotFromValues status", DB.DB_Select_name_inbound_status_list.QueryRow(status), &object.Status_id); err != nil {
		return
	}
	object.Status_name = status
	DB.Update(proc_name, DB.DB_Update_inbound_lot_status, object.Lot_id, object.Status_id)
}

// DB.DB_Select_inbound_product_name

// func (object *InboundLot) Save(inbound_lot_name, inbound_product_name, inbound_provider_name, container_name string) {
// 	log.Println("DEBUG: [InboundLotView Save] product_field", inbound_lot_name, inbound_product_name, inbound_provider_name, container_name)
//
// 	inbound_product_id := object.component_types_data[inbound_product_name]
// 	if inbound_product_id == DB.INVALID_ID {
// 		log.Println("Err: [InboundLotView Save] Invalid Product: ", inbound_lot_name, inbound_product_name, inbound_provider_name, container_name)
//
// 	}
// 	inbound_provider_id := DB.Insel(DB.DB_Insert_inbound_provider, DB.DB_Select_inbound_provider_id, "Err:  [InboundLotView Save] insel inbound_provider:", inbound_provider_name)
// 	container_id := DB.Insel(DB.DB_Insert_container, DB.DB_Select_container_id, "Err:  [InboundLotView Save] insel container:", container_name)
// 	DB.DB_Insert_inbound_lot.Exec(inbound_lot_name, inbound_product_id, inbound_provider_id, container_id)
// }

// MaxInt)
func IntMax(a, b int) int {
	return int(math.Max(float64(a), float64(b)))

}

func (object *InboundLot) Sample() {
	proc_name := "InboundLot.Sample"

	object.Update_status(Status_SAMPLED)
	pdf_path, err := object.ExportSample_label()
	if err != nil {
		log.Printf("Error: [%s]: %q\n", proc_name, err)
		return //err
	}
	product.Print_PDF(pdf_path)

}

func (object *InboundLot) ExportSample_label() (string, error) {

	// func Export_Storage_pdf(file_path, qc_sample_storage_name, product_moniker_name string, start_date, end_date, retain_date *time.Time, printDates bool) error {
	proc_name := "Product.Export_Storage_pdf"

	curr_col := 5.
	curr_row := 10.
	curr_row_delta := 15.

	cell_width := 50.
	cell_height := 10.

	file_path := fmt.Sprintf("%s/%s.%s", config.LABEL_PATH, strings.ReplaceAll(object.Lot_number, "/", "-"), "pdf")

	pdf := fpdf.New("L", "mm", "A7", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 32)

	curr_row = product.Increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, object.Container_name)

	curr_row = product.Increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, object.Lot_number)

	pdf.SetFontSize(22)
	curr_row = product.Increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, object.Product_name)

	log.Println("Info: Saving to: ", file_path)
	err := pdf.OutputFileAndClose(file_path)
	util.LogError(proc_name, err)
	return file_path, err
}

func (object *InboundLot) Quality_test() {
	proc_name := "InboundLot.Quality_test"
	operations_group := "BSQL"
	sample_size := 500.
	var BlendProducts []*blender.BlendProduct
	log.Println("DEBUG: ", proc_name, object.Product_name, object.Lot_number)

	//make blend
	Lot_Id := blender.Next_Lot_Id(operations_group)
	DB.Insert(proc_name, DB.DB_Insert_inbound_relabel, Lot_Id, object.Lot_id, object.Container_id)

	// make tested
	object.Update_status(Status_TESTED)

	//get recipes
	DB.Forall_err(proc_name,
		func() {},
		func(row *sql.Rows) error {
			BlendProduct := blender.NewBlendProduct()
			ProductBlend := blender.NewProductBlend()
			BlendComponent := blender.NewBlendComponent()
			BlendProduct.Blend = ProductBlend

			if err := row.Scan(
				&ProductBlend.Recipe_id, &BlendProduct.Product_id, &BlendComponent.Component_id, &BlendComponent.Component_type_id, &BlendComponent.Component_amount, &BlendComponent.Add_order,
			); err != nil {
				return err
			}
			log.Println("DEBUG: ", proc_name, ProductBlend.Recipe_id)

			BlendComponent.Lot_id = object.Lot_id
			BlendComponent.Inboundp = true
			BlendComponent.Component_amount *= sample_size
			ProductBlend.AddComponent(*BlendComponent)
			// queries cannot be nested, so dump them
			BlendProducts = append(BlendProducts, BlendProduct)
			return nil
		},
		DB.DB_Select_inbound_lot_recipe, object.Lot_id)

	// for all recipes:
recipes:
	for _, BlendProduct := range BlendProducts {
		ProductBlend := BlendProduct.Blend
		Components_collect := make(map[int][]blender.BlendComponent)
		BlendComponent := ProductBlend.Components[0]

		Components_collect[BlendComponent.Add_order] = append(Components_collect[BlendComponent.Add_order], BlendComponent)
		max_Add_order := BlendComponent.Add_order
		max_recipe_Add_order := -1
		log.Println("DEBUG: recipe# :", proc_name, ProductBlend.Recipe_id, BlendComponent.Lot_id)

		//get other components,
		proc_name := "GetComponents"
		DB.Forall_err(proc_name,
			func() {},
			func(row *sql.Rows) error {
				OtherBlendComponent := blender.NewBlendComponent()
				OtherBlendComponent.Inboundp = true

				if err := row.Scan(
					&OtherBlendComponent.Lot_id, &OtherBlendComponent.Component_id, &OtherBlendComponent.Component_type_id, &OtherBlendComponent.Component_amount, &OtherBlendComponent.Add_order,
				); err != nil {
					return err
				}
				OtherBlendComponent.Component_amount *= sample_size
				Components_collect[OtherBlendComponent.Add_order] = append(Components_collect[OtherBlendComponent.Add_order], *OtherBlendComponent)
				max_Add_order = IntMax(max_Add_order, OtherBlendComponent.Add_order)
				log.Println("\nDEBUG: max_Add_order", proc_name, max_Add_order, OtherBlendComponent.Lot_id)
				return nil
			},
			DB.DB_Select_inbound_lot_components,
			ProductBlend.Recipe_id, BlendComponent.Component_type_id, Status_TESTED)

		// ensure all components are sampled
		DB.Select_Error("MaxRecipeCount",
			DB.DB_Select_recipe_components_count.QueryRow(ProductBlend.Recipe_id), &max_recipe_Add_order)
		if max_recipe_Add_order != max_Add_order || len(Components_collect[0]) == 0 {
			continue recipes
		}

		// Take cartesian product of possible components,
		Components_map := make(map[int][][]blender.BlendComponent)

		for _, comp := range Components_collect[0] {
			Components_map[0] = append(Components_map[0], []blender.BlendComponent{comp})
		}
		for i := 1; i <= max_Add_order; i++ {
			if len(Components_collect[i]) == 0 {
				// component missing
				continue recipes
			}
			for _, comp := range Components_collect[i] {
				for _, list := range Components_map[i-i] {
					Components_map[i] = append(Components_map[i], append(list, comp))
				}
			}
		}

		//make blend
		for _, Components_list := range Components_map[max_Add_order] {
			ProductBlend.Components = Components_list
			BlendProduct.Save()
		}
	}

}
