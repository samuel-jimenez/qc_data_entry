package DB

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/util"
	"github.com/samuel-jimenez/windigo"
)

var (
	err error

	DB_VERSION = "0.0.4"

	// recipe_list
	DB_Select_product_recipe, DB_Insert_product_recipe,
	DB_Select_product_recipe_defaults,
	// recipe_procedure_list
	DB_Select_recipe_procedure_list,
	// recipe_procedure_steps
	DB_Select_recipe_procedure_steps,
	// recipe_components
	DB_Select_recipe_components_id, DB_Insert_recipe_component, DB_Update_recipe_component, DB_Delete_recipe_component, DB_Select_recipe_components_count,
	// component_types
	DB_Select_name_component_types, DB_Select_all_component_types, DB_Insert_component_types,
	DB_Select_component_type_product, DB_Select_component_type_density, DB_Insert_internal_product_component_type, DB_Insert_inbound_product_component_type,
	DB_Select_inbound_product_component_type_id,
	// container_capacity
	DB_Select_container_capacity_all, DB_Select_container_capacity_amount,
	DB_Select_container_capacity_description,
	DB_Select_container_capacity_info,
	// container_strap
	DB_Select_container_strap_container_capacity,
	// inbound_product
	DB_Select_inbound_product_name, DB_Insert_inbound_product,
	// container_list
	DB_Select_container_id, DB_Select_container_all, DB_Insert_container, DB_Update_container_type,
	// inbound_provider_list
	DB_Select_inbound_provider_id, DB_Select_inbound_provider_all, DB_Insert_inbound_provider,
	// inbound_status_list
	DB_Select_all_inbound_status_list, DB_Select_name_inbound_status_list,
	// internal_status_list
	DB_Select_all_internal_status_list, DB_Select_name_internal_status_list,
	// inbound_lot
	DB_Select_inbound_lot_status, DB_Select_inbound_lot_all, DB_Insert_inbound_lot, DB_Update_inbound_lot_status,
	DB_Select_name_inbound_lot_status, DB_Select_inbound_lot_recipe, DB_Select_inbound_lot_components,
	// inbound_relabel
	DB_Insert_inbound_relabel, DB_Select_inbound_relabel_all,
	// component_list
	DB_Select_inbound_blend_component, DB_Insert_inbound_blend_component,
	DB_Select_internal_blend_component, DB_Insert_internal_blend_component,

	// blend_components
	DB_Select_Product_blend_components,
	DB_Insert_Product_blend,
	// lot_list
	DB_Insert_lot, DB_Select_lot, DB_Select_blend_lot, DB_Select_lot_list_all,
	DB_Select_lot_list_name, DB_Select_lot_list_for_name_status, DB_Select_lot_list_for_product_lot_name, DB_Select_lot_list_for_product_lot_id_name,
	DB_Update_lot_list_name, DB_Update_lot_list__status, DB_Update_lot_list__component_status,
	DB_Select_product_lot_list_name,
	DB_Select_product_lot_list_sources,
	// product_lot
	db_select_id_lot, DB_Insert_product_lot,
	DB_Select_product_lot_all, DB_Select_product_lot_product,
	DB_Select_product_inbound_lot_all,
	DB_Insert_blend_lot, DB_Update_lot_recipe,
	DB_Update_lot_customer,
	DB_Select_product_lot_components,

	// product_line
	db_select_product_id, db_insert_product,
	DB_Select_product_info_all, DB_Select_product_info_inbound_all,
	DB_Select_product_info_moniker,
	// product_defaults
	DB_Insert_product_defaults_product, DB_Select_product_defaults,
	DB_Update_product_defaults,
	// product_customer_line
	DB_Select_product_customer_id, DB_Select_product_customer_info,
	db_select_product_customer, db_insert_product_customer,
	// bs.product_moniker
	DB_Select_all_product_moniker, DB_Select_product_lot_product_moniker,
	DB_Insert_product_moniker,
	// bs.product_sample_points
	DB_Select_all_sample_points, DB_Select_product_sample_points,
	DB_Insel_sample_point,
	// bs.qc_tester_list
	DB_Select_all_qc_tester, DB_Insel_qc_tester,
	// bs.qc_samples
	DB_insert_measurement,
	DB_Update_qc_samples_storage,
	// bs.product_sample_storage
	DB_Select_product_sample_storage_capacity, DB_Select_gen_product_sample_storage, DB_Select_all_product_sample_storage,
	DB_Insert_product_sample_storage,
	DB_Update_product_sample_storage_qc_sample, DB_Update_dec_product_sample_storage_capacity, DB_Update_product_sample_storage_capacity,
	// bs.qc_sample_storage_list
	DB_Insert_sample_storage,
	// bs.product_appearance
	DB_Insert_appearance,
	// bs.qc_test_methods
	DB_Select_test_methods__test_type,
	DB_Insel_test_method,
	// bs.product_attributes
	// bs.product_ranges_measured
	DB_Select_product_details,
	DB_Upsert_product_details, DB_Upsert_product_type *sql.Stmt

	INVALID_ID     int64 = 0
	DEFAULT_LOT_ID int64 = 1
)

func PrepareOrElse(db *sql.DB, sqlStatement string) *sql.Stmt {
	preparedStatement, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStatement)
		panic(err)

	}
	return preparedStatement
}

func Check_db(db *sql.DB, showWindowp bool) {
	var found_database_version_major,
		found_database_version_minor,
		found_database_version_revision int

	sqlStmt := `
		select database_version_major, database_version_minor, database_version_revision
		from bs.database_info
		`
	err = db.QueryRow(sqlStmt).Scan(&found_database_version_major,
		&found_database_version_minor,
		&found_database_version_revision)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		panic(err)
	}

	found_db_version := fmt.Sprintf("%d.%d.%d", found_database_version_major,
		found_database_version_minor,
		found_database_version_revision)

	if DB_VERSION != found_db_version {
		message := fmt.Sprintf("Database version mismatch: Required: %s, found: %s", DB_VERSION, found_db_version)
		err = errors.New(message)
		if showWindowp {
			windigo.Error(nil, message)
		}
		log.Printf("%q\n", err)
		panic(err)
	}
}

func DBinit(db *sql.DB) {
	// recipe_list
	DB_Select_product_recipe = PrepareOrElse(db, `
	select recipe_id
	from bs.recipe_list
	where product_id = ?
	`)
	DB_Select_product_recipe_defaults = PrepareOrElse(db, `
	select recipe_id, recipe_name, amount_total_default, specific_gravity_default
	from bs.recipe_list
	left join bs.product_defaults using (product_id)
	where product_id = ?
	`)

	DB_Insert_product_recipe = PrepareOrElse(db, `
	insert into bs.recipe_list
	(product_id,	recipe_procedure_id, recipe_name)
		values (?,?,?)
	returning recipe_id
	`)

	// recipe_components
	DB_Select_recipe_components_id = PrepareOrElse(db, `
	select recipe_components_id, component_type_id, component_type_name, component_type_amount, component_add_order
	from bs.recipe_components
	join bs.component_types
	using (component_type_id)
	where recipe_id = ?
	order by component_add_order
	`)

	DB_Insert_recipe_component = PrepareOrElse(db, `
	insert into bs.recipe_components
		(recipe_id,component_type_id,component_type_amount,component_add_order)
		values (?,?,?,?)
	returning recipe_components_id
	`)
	DB_Update_recipe_component = PrepareOrElse(db, `
	update  bs.recipe_components
	set
		component_type_id=?,
		component_type_amount=?,
		component_add_order=?
	where recipe_components_id = ?
	`)

	DB_Delete_recipe_component = PrepareOrElse(db, `
	delete from bs.recipe_components
	where recipe_components_id = ?
	`)

	DB_Select_recipe_components_count = PrepareOrElse(db, `
	select max(component_add_order)
		from bs.recipe_components
		where recipe_id = ?
	`)

	// recipe_procedure_list
	DB_Select_recipe_procedure_list = PrepareOrElse(db, `
	select recipe_procedure_id, recipe_procedure_name
	from bs.recipe_procedure_list
	order by recipe_procedure_id
	`)

	// recipe_procedure_steps
	DB_Select_recipe_procedure_steps = PrepareOrElse(db, `
	select recipe_procedure_step_text
	from bs.recipe_list
	join bs.recipe_procedure_steps	using (recipe_procedure_id)
	where recipe_id = ?
	order by recipe_procedure_step_number
	`)

	// component_types
	DB_Select_name_component_types = PrepareOrElse(db, `
	select component_type_id
		from bs.component_types
		where component_type_name = ?
	`)

	DB_Select_all_component_types = PrepareOrElse(db, `
	select component_type_id, component_type_name
		from bs.component_types
		order by component_type_name
	`)

	DB_Insert_component_types = PrepareOrElse(db, `
	insert into bs.component_types
		(component_type_name)
		values (?)
	returning component_type_id
	`)

	// TODO
	// 	DB_Select_component_type_product = PrepareOrElse(db, `
	// select component_type_id, component_type_name, product_id, inbound_product_id
	// from bs.component_types
	// left join bs.component_type_product_internal using (component_type_id)
	// left join bs.component_type_product_inbound using (component_type_id)
	// where component_type_id = ?
	// `)

	DB_Select_component_type_product = PrepareOrElse(db, `
	select
	false, product_lot_id, product_name_internal, lot_name, ''

	from bs.component_type_product_internal
	join bs.product_line using (product_id)
	join bs.product_lot using (product_id)
	join bs.lot_list using (lot_id)
	join bs.internal_status_list using (internal_status_id)

	where component_type_id = ?1
	and internal_status_name = ?2

	union

	select
	true, inbound_lot_id,inbound_product_name, inbound_lot_name, container_name

	from bs.component_type_product_inbound
	join bs.inbound_product using (inbound_product_id)
	join bs.inbound_lot using (inbound_product_id)
	join bs.container_list using (container_id)
	join bs.inbound_status_list using (inbound_status_id)

	where component_type_id = ?1
	and inbound_status_name != ?3

	`)

	DB_Select_component_type_density = PrepareOrElse(db, `
	select
		specific_gravity
	from  bs.product_lot
	join bs.qc_samples using (lot_id)
	where product_lot_id = ?1
		and false = ?2

	union

	select
		specific_gravity
	from  bs.inbound_relabel
	join bs.qc_samples using (lot_id)
	where inbound_lot_id = ?1
		and true = ?2
	`)

	// 	TODO get sources

	// select inbound_product_name,inbound_lot_name,container_name from blend_components
	// join bs.component_list using (component_id)
	// join bs.inbound_lot using (inbound_lot_id)
	// join bs.container_list using (container_id)
	// join bs.inbound_product using (inbound_product_id)
	// where blend_components.product_lot_id =?

	// 	TODO get sources

	// select inbound_product_name,inbound_lot_name,container_name from blend_components
	// join bs.component_list using (component_id)
	// join bs.inbound_lot using (inbound_lot_id)
	// join bs.lot_list using (lot_id)
	// join bs.product_lot using (product_lot_id)
	// join bs.container_list using (container_id)
	// join bs.inbound_product using (inbound_product_id)
	// where lot_name =?

	DB_Insert_internal_product_component_type = PrepareOrElse(db, `
	insert into bs.component_type_product_internal
		(component_type_id, product_id)
		values (?,?)
	returning component_type_product_internal_id
	`)

	DB_Insert_inbound_product_component_type = PrepareOrElse(db, `
	insert into bs.component_type_product_inbound
		(component_type_id, inbound_product_id)
		values (?,?)
	returning component_type_product_inbound_id
	`)
	DB_Select_inbound_product_component_type_id = PrepareOrElse(db, `
	select inbound_product_id, inbound_product_name
	from bs.component_type_product_inbound
	left join bs.inbound_product using (inbound_product_id)
	where component_type_id = ?1
	`)

	// container_capacity
	DB_Select_container_capacity_all = PrepareOrElse(db, `
	select container_capacity_id, container_capacity_name
	from bs.container_capacity
	order by container_capacity_name
	`)

	DB_Select_container_capacity_amount = PrepareOrElse(db, `
	select container_capacity_amount
	from bs.container_capacity
	where container_capacity_id = ?1
	`)

	DB_Select_container_capacity_description = PrepareOrElse(db, `
	select container_capacity_description
	from bs.container_capacity
	where container_capacity_id = ?1
	`)

	DB_Select_container_capacity_info = PrepareOrElse(db, `
	select
	container_capacity_description, container_capacity_amount
	from bs.container_capacity
	where container_capacity_id = ?1
	`)

	// container_strap
	//  should be monotonic increasing
	DB_Select_container_strap_container_capacity = PrepareOrElse(db, `
	select container_strap_key, container_strap_val
	from bs.container_strap
	where container_capacity_id = ?1
	order by container_strap_key
	`)

	// inbound_product
	DB_Select_inbound_product_name = PrepareOrElse(db, `
select inbound_product_id
from bs.inbound_product
where inbound_product_name = ?
	`)
	DB_Insert_inbound_product = PrepareOrElse(db, `
	insert into bs.inbound_product
		(inbound_product_name)
		values (?)
		returning inbound_product_id
		`)

	// container_list
	DB_Select_container_id = PrepareOrElse(db, `
	select container_id
		from bs.container_list
		where container_name = ?
	`)

	DB_Select_container_all = PrepareOrElse(db, `
	select container_id, container_name
		from bs.container_list
		order by container_name
	`)

	DB_Insert_container = PrepareOrElse(db, `
	insert into bs.container_list
		(container_name)
		values (?)
		returning container_id
		`)

	DB_Update_container_type = PrepareOrElse(db, `
update
bs.container_list

set
container_type_id=?2

where
container_id = ?1
		`)

	// inbound_provider_list
	DB_Select_inbound_provider_id = PrepareOrElse(db, `
select inbound_provider_id
from bs.inbound_provider_list
where inbound_provider_name = ?
	`)

	DB_Select_inbound_provider_all = PrepareOrElse(db, `
	select inbound_provider_id, inbound_provider_name
		from bs.inbound_provider_list
		order by inbound_provider_name
	`)

	DB_Insert_inbound_provider = PrepareOrElse(db, `
	insert into bs.inbound_provider_list
		(inbound_provider_name)
		values (?)
		returning inbound_provider_id
		`)
	// inbound_status_list
	DB_Select_all_inbound_status_list = PrepareOrElse(db, `
select
	inbound_status_id, inbound_status_name
from bs.inbound_status_list
order by inbound_status_id
	`)
	DB_Select_name_inbound_status_list = PrepareOrElse(db, `
select
	inbound_status_id
from bs.inbound_status_list
where inbound_status_name = ?
	`)

	// internal_status_list
	DB_Select_all_internal_status_list = PrepareOrElse(db, `
select
	internal_status_id, internal_status_name
from bs.internal_status_list
order by internal_status_id
	`)
	DB_Select_name_internal_status_list = PrepareOrElse(db, `
select
	internal_status_id
from bs.internal_status_list
where internal_status_name = ?
	`)

	// inbound_lot
	DB_Select_inbound_lot_status = PrepareOrElse(db, `
select
	inbound_lot_id, inbound_lot_name, inbound_product_id, inbound_product_name, inbound_provider_id, inbound_provider_name, container_id, container_name, inbound_status_id, inbound_status_name
from bs.inbound_lot
join bs.inbound_product 		using (inbound_product_id)
join bs.inbound_provider_list 		using (inbound_provider_id)
join bs.container_list 			using (container_id)
join bs.inbound_status_list 			using (inbound_status_id)
		where inbound_status_name = ?
order by inbound_lot_name
	`)

	DB_Select_inbound_lot_all = PrepareOrElse(db, `
select
	inbound_lot_id, inbound_lot_name, inbound_product_id, inbound_product_name, inbound_provider_id, inbound_provider_name, container_id, container_name, inbound_status_id, inbound_status_name
from bs.inbound_lot
join bs.inbound_product 		using (inbound_product_id)
join bs.inbound_provider_list 		using (inbound_provider_id)
join bs.container_list 			using (container_id)
join bs.inbound_status_list 			using (inbound_status_id)
order by inbound_lot_name
	`)

	DB_Insert_inbound_lot = PrepareOrElse(db, `
insert into bs.inbound_lot
(inbound_lot_name,inbound_product_id,inbound_provider_id,container_id)
values (?,?,?,?)
returning inbound_lot_id
		`)
	DB_Update_inbound_lot_status = PrepareOrElse(db, `
update bs.inbound_lot
set
	inbound_status_id=?2
where inbound_lot_id=?1
`)

	DB_Select_name_inbound_lot_status = PrepareOrElse(db, `
select
	inbound_lot_name
from bs.inbound_lot
join bs.inbound_status_list 			using (inbound_status_id)
where inbound_status_name = ?
	`)

	DB_Select_inbound_lot_recipe = PrepareOrElse(db, `
select
	recipe_id, product_id, recipe_components_id, component_type_id, component_type_amount, component_add_order
from bs.inbound_lot
join bs.component_type_product_inbound 	using (inbound_product_id)
join bs.recipe_components 		using (component_type_id)
join bs.recipe_list			using (recipe_id)
where inbound_lot_id =  ?
order by recipe_id
	`)

	DB_Select_inbound_lot_components = PrepareOrElse(db, `
select
	inbound_lot_id, recipe_components_id, component_type_id, component_type_amount, component_add_order
from bs.inbound_lot
join bs.inbound_status_list 			using (inbound_status_id)
join bs.component_type_product_inbound 	using (inbound_product_id)
join bs.recipe_components 		using (component_type_id)
where recipe_id = ?1
	and component_type_id != ?2
	and inbound_status_name = ?3
order by component_add_order
	`)

	// inbound_relabel
	DB_Insert_inbound_relabel = PrepareOrElse(db, `
	insert into bs.inbound_relabel
		(lot_id,inbound_lot_id,container_id)
		values (?,?,?)
	returning inbound_relabel_id
	`)

	DB_Select_inbound_relabel_all = PrepareOrElse(db, `
	select
	inbound_relabel_id, lot_name
	from bs.inbound_relabel
	join bs.lot_list using (lot_id)
	`)
	// 			// TODO
	//
	// 		DB_Select_inbound_relabel = PrepareOrElse(db, `
	// 	select
	// 	lot_name inbound_product_name, inbound_lot_name, container_name
	// 	from bs.inbound_relabel
	// 	join bs.lot_list using (lot_id)
	//
	// join bs.inbound_lot 		using (inbound_lot_id)
	// join bs.inbound_lot 		using (inbound_lot_id)
	// join bs.inbound_product 		using (inbound_product_id)
	// join bs.inbound_provider_list 		using (inbound_provider_id)
	// join bs.container_list 			using (container_id)
	// 		(lot_id,inbound_lot_id,container_id)
	// 		values (?,?,?)
	// 	returning inbound_relabel_id
	// 	`)

	// component_list
	DB_Select_inbound_blend_component = PrepareOrElse(db, `
select component_id
from bs.component_list
where component_type_id = ?
and inbound_lot_id = ?
	`)
	DB_Insert_inbound_blend_component = PrepareOrElse(db, `
insert into bs.component_list
	(component_type_id,inbound_lot_id)
	values (?,?)
returning component_id
`)
	DB_Select_internal_blend_component = PrepareOrElse(db, `
select component_id
from bs.component_list
where component_type_id = ?
and product_lot_id = ?
	`)
	DB_Insert_internal_blend_component = PrepareOrElse(db, `
insert into bs.component_list
	(component_type_id,product_lot_id)
	values (?,?)
returning component_id
`)

	// blend_components
	DB_Select_Product_blend_components = PrepareOrElse(db, `
	select count(blend_components_id)
	from bs.blend_components
	where product_lot_id =?
	`)
	DB_Insert_Product_blend = PrepareOrElse(db, `
	insert into bs.blend_components
	(product_lot_id, recipe_components_id, component_id, component_required_amount)
	values (?,?,?,?)
	returning blend_components_id
	`)

	// lot_list
	DB_Insert_lot = PrepareOrElse(db, `
insert into bs.lot_list
	(lot_name)
values (?)
returning lot_id
`)
	DB_Select_lot = PrepareOrElse(db, `
select lot_id
from bs.lot_list
where lot_name = ?
`)
	DB_Select_blend_lot = PrepareOrElse(db, `
select count(lot_name)
	from bs.lot_list
where lot_name like ?
`)
	DB_Select_lot_list_all = PrepareOrElse(db, `
	select
	lot_id, lot_name
	from bs.lot_list
	`)

	DB_Select_lot_list_name = PrepareOrElse(db, `
	select
	lot_id, lot_name
	from bs.lot_list
	where lot_name like ?
	`)

	DB_Select_lot_list_for_name_status = PrepareOrElse(db, `
	select
	lot_id, lot_name
	from bs.lot_list
	join bs.internal_status_list using (internal_status_id)
	where lot_name like ?
	and internal_status_name = ?
	`)

	DB_Select_lot_list_for_product_lot_name = PrepareOrElse(db, `
	select  lot_id, lot_name, product_lot_id
	from bs.lot_list
	left join bs.product_lot using (lot_id)
	where lot_name like ?
	`)

	DB_Select_lot_list_for_product_lot_id_name = PrepareOrElse(db, `
	select  lot_id, lot_name, product_lot_id
	from bs.lot_list
	left join bs.product_lot using (lot_id)
	where product_id = ?
	and lot_name like ?
	`)

	DB_Update_lot_list_name = PrepareOrElse(db, `
update
bs.lot_list

set
lot_name = ?2

where lot_id = ?1
`)

	DB_Update_lot_list__status = PrepareOrElse(db, `
update bs.lot_list
set
	internal_status_id=bs.internal_status_list.internal_status_id
from
bs.internal_status_list

where lot_id = ?
and internal_status_name = ?

		`)

	DB_Update_lot_list__component_status = PrepareOrElse(db, `
with upd_table as (
select

lot_id, internal_status_list.internal_status_id
from bs.lot_list
join bs.product_lot using (lot_id)
join bs.blend_components using (product_lot_id)
join bs.component_list using (component_id)
join bs.inbound_lot using (inbound_lot_id)
full join bs.internal_status_list
where inbound_lot_name = ?1
and internal_status_name = ?2

union

select

lot_id, internal_status_list.internal_status_id
from bs.lot_list
	join bs.inbound_relabel using (lot_id)
	join bs.inbound_lot using (inbound_lot_id)
full join bs.internal_status_list
where inbound_lot_name = ?1
and internal_status_name = ?2
	)

update bs.lot_list
set
internal_status_id = upd_table.internal_status_id
from
upd_table
where lot_list.lot_id = upd_table.lot_id
		`)

	// 	TODO get sources

	// select inbound_product_name,inbound_lot_name,container_name from blend_components
	// join bs.component_list using (component_id)
	// join bs.inbound_lot using (inbound_lot_id)
	// join bs.container_list using (container_id)
	// join bs.inbound_product using (inbound_product_id)
	// where blend_components.product_lot_id =?

	// 	TODO get sources

	// select inbound_product_name,inbound_lot_name,container_name from blend_components
	// join bs.component_list using (component_id)
	// join bs.inbound_lot using (inbound_lot_id)
	// join bs.lot_list using (lot_id)
	// join bs.product_lot using (product_lot_id)
	// join bs.container_list using (container_id)
	// join bs.inbound_product using (inbound_product_id)
	// where lot_name =?

	DB_Select_product_lot_list_name = PrepareOrElse(db, `
select lot_id,  format('%s %s',  product_moniker_name, product_name_internal), lot_name, product_name_customer
from bs.lot_list
join bs.product_lot using (lot_id)
join bs.product_line using (product_id)
join bs.product_moniker using (product_moniker_id)
left join bs.product_customer_line using (product_customer_id)
where lot_name = ?1

union

select lot_id, inbound_product_name, lot_name, null
from bs.lot_list
	join bs.inbound_relabel using (lot_id)
	join bs.inbound_lot using (inbound_lot_id)
	join bs.inbound_product using (inbound_product_id)
where lot_name = ?1
`)

	///TODO decide between:
	// 	create recipe
	//
	// 	create blend, assigning components, amounts
	//
	// 	capture actual amounts
	//
	// 	and:
	//
	// 		create recipe
	//
	// 	create blend, assigning amounts
	//
	// 		capture components, actual amounts

	// TODO sample_size := 500.
	DB_Select_product_lot_list_sources = PrepareOrElse(db, `
select inbound_product_name, inbound_lot_id, inbound_lot_name, container_name, component_required_amount, component_add_order, true

from bs.lot_list
join bs.product_lot using (lot_id)
join bs.blend_components using (product_lot_id)
join bs.component_list using (component_id)
join bs.recipe_components		using (recipe_components_id)
join bs.inbound_lot using (inbound_lot_id)
join bs.container_list using (container_id)
join bs.inbound_product using (inbound_product_id)

where lot_name = ?1

union

select inbound_product_name, inbound_lot_id, inbound_lot_name, container_name, 500, 0 component_add_order, true

from bs.lot_list
join bs.inbound_relabel using (lot_id)
join bs.inbound_lot using (inbound_lot_id)
join bs.container_list using (container_id)
join bs.inbound_product using (inbound_product_id)

where lot_name = ?1

		order by component_add_order
`)

	// product_lot
	db_select_id_lot = PrepareOrElse(db, `
select product_lot_id
from bs.product_lot
where lot_id = ? and product_id = ?
	`)

	DB_Select_product_lot_product = PrepareOrElse(db, `
select product_lot_id, lot_name
from bs.product_lot
join bs.lot_list using (lot_id)
where product_id = ?
order by lot_id desc
`)
	DB_Select_product_lot_all = PrepareOrElse(db, `
select product_lot_id, lot_name
from bs.product_lot
join bs.lot_list using (lot_id)
order by lot_id desc
`)

	DB_Select_product_inbound_lot_all = PrepareOrElse(db, `
	select product_lot_id, lot_name
	from (
		select product_lot_id, lot_id
		from bs.product_lot
		union
		select inbound_lot_id, lot_id
		from bs.inbound_relabel
	)
	join bs.lot_list using (lot_id)
	order by lot_name desc
	`)

	DB_Insert_product_lot = PrepareOrElse(db, `
	insert into bs.product_lot
		(lot_id,product_id)
		values (?,?)
		returning product_lot_id
	`)

	DB_Insert_blend_lot = PrepareOrElse(db, `
	insert into bs.product_lot
		(lot_id,product_id,product_customer_id,recipe_id)
		values (?,?,?,?)
	returning product_lot_id
	`)

	DB_Update_lot_recipe = PrepareOrElse(db, `
	update bs.product_lot
	set
		recipe_id=?
	where product_lot_id=?
		`)

	DB_Update_lot_customer = PrepareOrElse(db, `
	update bs.product_lot
	set
		product_customer_id=?
	where product_lot_id=?
		`)

	/// FIXME not actually used
	DB_Select_product_lot_components = PrepareOrElse(db, `
select inbound_product_name, inbound_lot_name, container_name from blend_components
join bs.component_list using (component_id)
join bs.inbound_lot using (inbound_lot_id)
join bs.container_list using (container_id)
join bs.inbound_product using (inbound_product_id)
where blend_components.product_lot_id =?
`)

	// product_line
	db_insert_product = PrepareOrElse(db, `
	insert into bs.product_line
		(product_name_internal, product_moniker_id)
		select ?, product_moniker_id
			from bs.product_moniker
			where product_moniker_name = ?
		returning product_id
		`)

	db_select_product_id = PrepareOrElse(db, `
	select product_id
	from bs.product_line
		join bs.product_moniker using (product_moniker_id)
	where product_name_internal = ?
		and product_moniker_name = ?
	`)

	DB_Select_product_info_all = PrepareOrElse(db, `
	select product_id, product_name_internal, product_moniker_name
	from bs.product_line
	join bs.product_moniker using (product_moniker_id)
	order by product_moniker_name, product_name_internal
	`)

	DB_Select_product_info_inbound_all = PrepareOrElse(db, `
	select
	true, product_id, product_name_internal, product_moniker_name
	from bs.product_line
	join bs.product_moniker using (product_moniker_id)
	union
	select
	false, inbound_product_id+1000,inbound_product_name, inbound_product_name
	from bs.inbound_relabel
	join bs.inbound_lot using (inbound_lot_id)
	join bs.inbound_product using (inbound_product_id)
	order by product_moniker_name, product_name_internal
	`)

	// 8 product_moniker_id

	DB_Select_product_info_moniker = PrepareOrElse(db, `
	select product_id, product_name_internal, product_moniker_name
	from bs.product_line
		join bs.product_moniker using (product_moniker_id)
	where product_moniker_name = ?
	order by product_moniker_name, product_name_internal
	`)

	// product_defaults
	DB_Insert_product_defaults_product = PrepareOrElse(db, `
	insert into bs.product_defaults
	(product_default_id,product_id)
	values (?1, ?1)
	`)
	// DB_Insel_product_defaults_product = PrepareOrElse(db, `
	// with val (qc_tester_name) as (
	// 	values
	// 	(?)
	// ),
	// sel as (
	// 	select qc_tester_name, qc_tester_id
	// 	from val
	// 	left join bs.qc_tester_list using (qc_tester_name)
	// )
	// insert into bs.qc_tester_list (qc_tester_name)
	// select distinct qc_tester_name from sel where qc_tester_id is null
	// returning qc_tester_id, qc_tester_name
	// `)
	DB_Select_product_defaults = PrepareOrElse(db, `
	select
		amount_total_default, specific_gravity_default
	from bs.product_defaults
	where product_id = ?
	`)

	DB_Update_product_defaults = PrepareOrElse(db, `
	update
	bs.product_defaults

	set
	amount_total_default = ?2

	where product_default_id = ?1
	`)
	// where product_id = ?1

	// product_customer_line
	DB_Select_product_customer_info = PrepareOrElse(db, `
	select product_customer_id, product_name_customer
		from bs.product_customer_line
		where product_id = ?
		`)

	DB_Select_product_customer_id = PrepareOrElse(db, `
	select product_customer_id
		from bs.product_customer_line
		where product_name_customer = ? and product_id = ?
	`)

	db_select_product_customer = PrepareOrElse(db, `
	select product_customer_id
		from bs.product_customer_line
		where product_name_customer = ? and product_id = ?
	`)

	db_insert_product_customer = PrepareOrElse(db, `
	insert into bs.product_customer_line
		(product_name_customer,product_id)
		values (?,?)
	returning product_customer_id
	`)

	// bs.product_moniker
	DB_Select_all_product_moniker = PrepareOrElse(db, `
select
	product_moniker_id, product_moniker_name
from
	bs.product_moniker
order by
	product_moniker_name
`)

	DB_Select_product_lot_product_moniker = PrepareOrElse(db, `
	select
		product_lot_id, lot_name
		from bs.product_lot
		join bs.lot_list using (lot_id)
		join bs.product_line using (product_id)
		join bs.product_moniker using (product_moniker_id)
	where product_moniker_name = ?
	order by lot_name
`)

	DB_Insert_product_moniker = PrepareOrElse(db, `
	insert into bs.product_moniker
		(product_moniker_name)
		values (?)
	returning product_moniker_id
`)

	// bs.product_sample_points
	DB_Insel_sample_point = PrepareOrElse(db, `
	with val (sample_point) as (
		values
			(?)
		),
		sel as (
			select sample_point, sample_point_id
			from val
			left join bs.product_sample_points using (sample_point)
		)
	insert into bs.product_sample_points (sample_point)
	select distinct sample_point from sel where sample_point_id is null
	returning sample_point_id, sample_point
	`)

	DB_Select_all_sample_points = PrepareOrElse(db, `
	select sample_point_id, sample_point
		from bs.product_sample_points
		order by sample_point_id
	`)

	DB_Select_product_sample_points = PrepareOrElse(db, `
	select distinct sample_point_id, sample_point
		from bs.product_lot
		join bs.qc_samples using (lot_id)
		join bs.product_sample_points using (sample_point_id)
		where product_id = ?
		order by sample_point_id
	`)

	// bs.qc_tester_list
	DB_Select_all_qc_tester = PrepareOrElse(db, `
	select qc_tester_id, qc_tester_name
		from bs.qc_tester_list
		order by qc_tester_name
	`)

	DB_Insel_qc_tester = PrepareOrElse(db, `
	with val (qc_tester_name) as (
		values
			(?)
		),
		sel as (
			select qc_tester_name, qc_tester_id
			from val
			left join bs.qc_tester_list using (qc_tester_name)
		)
	insert into bs.qc_tester_list (qc_tester_name)
	select distinct qc_tester_name from sel where qc_tester_id is null
	returning qc_tester_id, qc_tester_name
	`)

	// bs.qc_samples
	BIG_SAMPLES_QC := `time_stamp,  ph, specific_gravity, density, string_test, viscosity`
	BIG_ID_SAMPLES_QC := `lot_id, sample_point, qc_tester_name, ` + BIG_SAMPLES_QC
	BIG_NAME_SAMPLES_QC := `lot_id, sample_point_id, qc_tester_id, ` + BIG_SAMPLES_QC
	DB_insert_measurement = PrepareOrElse(db, `
	with
		val ( 	`+BIG_ID_SAMPLES_QC+`) as (
			values
				(?, ?, ?, ?, ?, ?, ?, ?, ?)
		)
	insert into bs.qc_samples (`+BIG_NAME_SAMPLES_QC+`)
	select `+BIG_NAME_SAMPLES_QC+`
	from val
	left join bs.qc_tester_list using (qc_tester_name)
	left join bs.product_sample_points using (sample_point)
	returning qc_id;
	`)

	DB_Update_qc_samples_storage = PrepareOrElse(db, `
update bs.qc_samples
	set
qc_sample_storage_id = ?2

where qc_id = ?1
	`)

	// bs.product_sample_storage
	DB_Select_product_sample_storage_capacity = PrepareOrElse(db, `
select

qc_sample_storage_id, qc_storage_capacity

from bs.product_sample_storage
join bs.product_line using (product_moniker_id)
where product_id = ?1
	`)

	DB_Select_gen_product_sample_storage = PrepareOrElse(db, `
			select

product_sample_storage_id, product_moniker_name, qc_sample_storage_offset, qc_sample_storage_name, min(time_stamp), max(time_stamp), retain_storage_duration

from bs.product_sample_storage
join bs.product_line using (product_moniker_id)
join bs.qc_sample_storage_list using (qc_sample_storage_id, product_moniker_id)
join bs.qc_samples using (qc_sample_storage_id)
join bs.product_moniker using (product_moniker_id)
where product_id = ?1
	`)

	DB_Select_all_product_sample_storage = PrepareOrElse(db, `
	select

	product_sample_storage_id, product_moniker_name, qc_sample_storage_name, max_storage_capacity,  qc_storage_capacity

	from bs.product_sample_storage
		join bs.qc_sample_storage_list using (qc_sample_storage_id)
		join bs.product_moniker using (product_moniker_id)
	where max_storage_capacity !=  qc_storage_capacity
	or qc_sample_storage_offset > 0
	order by qc_sample_storage_name
	`)
	DB_Insert_product_sample_storage = PrepareOrElse(db, `
	insert into bs.product_sample_storage
	( product_sample_storage_id, product_moniker_id, retain_storage_duration, max_storage_capacity, qc_sample_storage_id, qc_sample_storage_offset, qc_storage_capacity )
	values ( ?, ?, ?, ?, ?, ?, ? )
	returning product_sample_storage_id
	`)

	DB_Update_product_sample_storage_qc_sample = PrepareOrElse(db, `
update bs.product_sample_storage
	set
qc_sample_storage_id = ?2,
qc_sample_storage_offset = ?3,
qc_storage_capacity = max_storage_capacity

where product_sample_storage_id = ?1
	`)

	DB_Update_dec_product_sample_storage_capacity = PrepareOrElse(db, `
	update bs.product_sample_storage
	set
	qc_storage_capacity = qc_storage_capacity - ?2

	where qc_sample_storage_id = ?1
	`)

	DB_Update_product_sample_storage_capacity = PrepareOrElse(db, `
	update bs.product_sample_storage
	set
	qc_storage_capacity = ?2

	where product_sample_storage_id = ?1
	`)
	/*
	    *
	   			select

	   product_sample_storage_id, product_moniker_name, qc_sample_storage_offset, qc_sample_storage_name, min(time_stamp), max(time_stamp), retain_storage_duration, qc_storage_capacity

	   from product_sample_storage
	   join product_line using (product_moniker_id)
	   join qc_sample_storage_list using (qc_sample_storage_id, product_moniker_id)
	   left join qc_samples using (qc_sample_storage_id)
	   join product_moniker using (product_moniker_id)
	   group by product_moniker_id

	   *
	*/

	// bs.qc_sample_storage_list

	// DB_Select_all_sample_storage = PrepareOrElse(db, `
	// select qc_sample_storage_id, qc_sample_storage_name, product_moniker_id
	// 	from bs.qc_sample_storage_list
	// 	order by qc_sample_storage_name
	// `)
	DB_Insert_sample_storage = PrepareOrElse(db, `
	insert into bs.qc_sample_storage_list
		( qc_sample_storage_name, product_moniker_id )
		values (?, ?)
	returning qc_sample_storage_id
	`)

	// bs.product_appearance
	DB_Insert_appearance = PrepareOrElse(db, `
	with val (product_appearance_text) as (
		values
		(?)
	),
	sel as (
		select product_appearance_text, product_appearance_id
		from val
		left join bs.product_appearance using (product_appearance_text)
	)
	insert into bs.product_appearance (product_appearance_text)
	select distinct product_appearance_text from sel where product_appearance_id is null and product_appearance_text is not null
	returning product_appearance_id, product_appearance_text
	`)

	// bs.qc_test_methods
	DB_Select_test_methods__test_type = PrepareOrElse(db, `
	select distinct qc_test_method_id,qc_test_method_name
	from bs.product_ranges_measured
	join bs.qc_test_methods using (qc_test_method_id)
	where qc_test_type_id=?
	`)

	DB_Insel_test_method = PrepareOrElse(db, `
	with val (qc_test_method_name) as (
		values
		(?)
	),
	sel as (
		select qc_test_method_id, qc_test_method_name
		from val
		left join bs.qc_test_methods using (qc_test_method_name)
	)
	insert into bs.qc_test_methods (qc_test_method_name)
	select distinct qc_test_method_name from sel where qc_test_method_id is null and qc_test_method_name is not null
	returning qc_test_method_id, qc_test_method_name
	`)

	// bs.product_attributes
	// bs.product_ranges_measured
	BIG_SEL_RANGES_QC := `
	ph_method,
	ph_measure,
	ph_publish,
	ph_min,
	ph_target,
	ph_max,

	specific_gravity_method,
	specific_gravity_measure,
	specific_gravity_publish,
	specific_gravity_min,
	specific_gravity_target,
	specific_gravity_max,

	density_method,
	density_measure,
	density_publish,
	density_min,
	density_target,
	density_max,

	string_test_method,
	string_test_measure,
	string_test_publish,
	string_test_min,
	string_test_target,
	string_test_max,

	viscosity_method,
	viscosity_measure,
	viscosity_publish,
	viscosity_min,
	viscosity_target,
	viscosity_max
	`

	DB_Select_product_details = PrepareOrElse(db, `
	with
	val
	(
		product_id,
	`+BIG_SEL_RANGES_QC+`
	)
	as (
	select product_id,
		max(case when qc_test_type_id = 2 then qc_test_method_name end) as ph_method,
		max(case when qc_test_type_id = 2 then val_measure end) as ph_measure,
		max(case when qc_test_type_id = 2 then val_publish end) as ph_publish,
		max(case when qc_test_type_id = 2 then val_min end) as ph_min,
		max(case when qc_test_type_id = 2 then val_target end) as ph_target,
		max(case when qc_test_type_id = 2 then val_max end) as ph_max,

		max(case when qc_test_type_id = 3 then qc_test_method_name end) as specific_gravity_method,
		max(case when qc_test_type_id = 3 then val_measure end) as specific_gravity_measure,
		max(case when qc_test_type_id = 3 then val_publish end) as specific_gravity_publish,
		max(case when qc_test_type_id = 3 then val_min end) as specific_gravity_min,
		max(case when qc_test_type_id = 3 then val_target end) as specific_gravity_target,
		max(case when qc_test_type_id = 3 then val_max end) as specific_gravity_max,

		max(case when qc_test_type_id = 4 then qc_test_method_name end) as density_method,
		max(case when qc_test_type_id = 4 then val_measure end) as density_measure,
		max(case when qc_test_type_id = 4 then val_publish end) as density_publish,
		max(case when qc_test_type_id = 4 then val_min end) as density_min,
		max(case when qc_test_type_id = 4 then val_target end) as density_target,
		max(case when qc_test_type_id = 4 then val_max end) as density_max,

		max(case when qc_test_type_id = 5 then qc_test_method_name end) as string_test_method,
		max(case when qc_test_type_id = 5 then val_measure end) as string_test_measure,
		max(case when qc_test_type_id = 5 then val_publish end) as string_test_publish,
		max(case when qc_test_type_id = 5 then val_min end) as string_test_min,
		max(case when qc_test_type_id = 5 then val_target end) as string_test_target,
		max(case when qc_test_type_id = 5 then val_max end) as string_test_max,

		max(case when qc_test_type_id = 6 then qc_test_method_name end) as viscosity_method,
		max(case when qc_test_type_id = 6 then val_measure end) as viscosity_measure,
		max(case when qc_test_type_id = 6 then val_publish end) as viscosity_publish,
		max(case when qc_test_type_id = 6 then val_min end) as viscosity_min,
		max(case when qc_test_type_id = 6 then val_target end) as viscosity_target,
		max(case when qc_test_type_id = 6 then val_max end) as viscosity_max

	from bs.product_ranges_measured
	left join bs.qc_test_methods using (qc_test_method_id)
	where product_id = ?1
	)

	select

	product_type_id,
	container_type_id,
	product_appearance_text,
	`+BIG_SEL_RANGES_QC+`

	from bs.product_attributes
	left join val using (product_id)
	left join bs.product_appearance using (product_appearance_id)
	join bs.product_types using (product_type_id)

	where product_id = ?1
	`)
	// group by product_id

	BIG_RANGES_QC := `
	val_measure,
	val_publish,
	val_min,
	val_target,
	val_max`
	BIG_NAME_RANGES_QC := `product_id, qc_test_type_name,
	qc_test_method_name,
	` + BIG_RANGES_QC

	BIG_ID_RANGES_QC := `product_id,
	qc_test_type_id,
	qc_test_method_id,
	` + BIG_RANGES_QC

	BIG_EXCLUDED_QC := `
qc_test_method_id=excluded.qc_test_method_id,
val_measure=excluded.val_measure,
val_publish=excluded.val_publish,
val_min=excluded.val_min,
val_target=excluded.val_target,
val_max=excluded.val_max`
	DB_Upsert_product_details = PrepareOrElse(db, `
	with
	val
	(
		`+BIG_NAME_RANGES_QC+`
	)
	as (
		values (
			?,?,
			?,?,?,?,?,?)
	),

	sel as (select
		`+BIG_ID_RANGES_QC+`
	from val
	join bs.qc_test_types using (qc_test_type_name)
	left join bs.qc_test_methods using (qc_test_method_name)
)

	insert into bs.product_ranges_measured
	(
		`+BIG_ID_RANGES_QC+`
	)
	select
	`+BIG_ID_RANGES_QC+`
	from sel
	where true
	on conflict(product_id, qc_test_type_id) do update set
	`+BIG_EXCLUDED_QC+`
	returning range_id
	`)

	product_type_id := `product_id,
	product_type_id,
	product_appearance_id`
	DB_Upsert_product_type = PrepareOrElse(db, `
	with
	val
	(
		product_id,
		product_type_id,
		product_appearance_text
	)
	as (
		values (
			?,?,?
		)
	),
	sel as (
		select
		`+product_type_id+`
		from val
		left join bs.product_appearance using (product_appearance_text)
	)
	insert into bs.product_attributes	(
		`+product_type_id+`
	)
	select
		`+product_type_id+`
	from sel
	where true
	on conflict(product_id) do update set
	product_type_id=excluded.product_type_id,
	product_appearance_id=excluded.product_appearance_id

	returning product_attribute_id
	`)
}

func Select_Error(proc_name string, query *sql.Row, args ...any) error {
	err := query.Scan(args...)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
	}
	return err
}

// permit empty
func Select_ErrNoRows(proc_name string, query *sql.Row, args ...any) error {
	err := query.Scan(args...)
	if err != nil && err != sql.ErrNoRows { // no row? no problem!
		log.Printf("Err: [%s]: %q\n", proc_name, err)
	}
	return err
}

func Select_Panic(proc_name string, query *sql.Row, args ...any) {
	err := query.Scan(args...)
	if err != nil {
		log.Printf("Critical hit: [%s]: %q\n", proc_name, err)
		panic(err)
	}
}

func Select_Panic_ErrorBox(proc_name string, query *sql.Row, args ...any) {
	err := query.Scan(args...)
	if err != nil {
		log.Printf("Critical hit: [%s]: %q\n", proc_name, err)
		windigo.Errorf(nil, "Something's gone wrong: [%s]: \n%q\n", proc_name, err)
		panic(err)
	}
}

func Insert(proc_name string, insert_statement *sql.Stmt, args ...any) int64 {
	var insert_id int64
	result, err := insert_statement.Exec(args...)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return INVALID_ID
	}
	insert_id, err = result.LastInsertId()
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return INVALID_ID
	}
	return insert_id
}

func Exec_Error(proc_name string, statement *sql.Stmt, args ...any) error {
	_, err := statement.Exec(args...)
	util.LogError(proc_name, err)
	return err
}

func Update(proc_name string, update_statement *sql.Stmt, args ...any) error {
	return Exec_Error(proc_name, update_statement, args...)
}

func Delete(proc_name string, delete_statement *sql.Stmt, args ...any) error {
	return Exec_Error(proc_name, delete_statement, args...)
}

func Insel(proc_name string, insert_statement, select_statement *sql.Stmt, args ...any) int64 {
	var insel_id int64
	if select_statement.QueryRow(args...).Scan(&insel_id) != nil {
		// no rows
		insel_id = Insert(proc_name, insert_statement, args...)
	}
	return insel_id
}

func Insel_product_id(product_name_full string) int64 {
	product_moniker_name, product_name_internal, _ := strings.Cut(product_name_full, " ")
	product_id := Insel("Insel_product_id", db_insert_product, db_select_product_id, product_name_internal, product_moniker_name)
	DB_Insert_product_defaults_product.Exec(product_id)

	return product_id
}

func Insel_lot_id(lot_name string) int64 {
	return Insel("Insel_lot_id", DB_Insert_lot, DB_Select_lot, lot_name)
}

func Insel_product_lot_id(Lot_id, product_id int64) int64 {
	return Insel("Insel_product_lot_id", DB_Insert_product_lot, db_select_id_lot, Lot_id, product_id)
}

func Insel_product_name_customer(product_name_customer string, product_id int64) int64 {
	return Insel("Insel_product_name_customer", db_insert_product_customer, db_select_product_customer, product_name_customer, product_id)
}

func Select_product_lot_name(lot_name string, product_id int64) (Lot_name string, Lot_id, Product_Lot_id int64) {
	// we cannot guarantee uniqueness of lot numbers
	proc_name := "Select_product_lot_name"
	Lot_name = lot_name + "%"

	// find like Lot_number :%
	if err := Select_ErrNoRows(proc_name, DB_Select_lot_list_for_product_lot_id_name.QueryRow(product_id, Lot_name),
		&Lot_id, &Lot_name, &Product_Lot_id,
	); err == nil {
		// already in product_lot
		return
	}

	currMax := -1
	var (
		lot_address_str string
		p_Lot_id        *int
		found           bool
	)
	// find like Lot_number :%
	// TODO Forall_done
	rows, err := DB_Select_lot_list_for_product_lot_name.Query(Lot_name)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return
	}
	for rows.Next() {
		if err := rows.Scan(
			&Lot_id, &Lot_name, &p_Lot_id,
		); err != nil {
			log.Printf("Err: [%s]: %q\n", proc_name, err)
		}
		Lot_name, lot_address_str, found = strings.Cut(Lot_name, ":")
		lot_address, _ := strconv.Atoi(lot_address_str)
		currMax = max(currMax, lot_address)
	}

	// not in lot_list:
	// proceed as formerly
	if currMax == -1 {
		Lot_name = lot_name
		Lot_id = Insel_lot_id(Lot_name)
		// Product_Lot_id = Insel_product_lot_id(Lot_id, product_id)
	}

	// in lot_list:
	// not product_lot
	if p_Lot_id == nil {
		Product_Lot_id = Insel_product_lot_id(Lot_id, product_id)
		return
	}

	if !found {
		if err = Update(proc_name,
			DB_Update_lot_list_name,
			Lot_id,
			Lot_name+":0",
		); err != nil {
			log.Printf("Err: [%s]: %q\n", proc_name, err)
		}
	}
	Lot_name += ":" + strconv.Itoa(currMax+1)
	Lot_id = Insel_lot_id(Lot_name)
	Product_Lot_id = Insel_product_lot_id(Lot_id, product_id)
	return
}

func Forall(proc_name string, start_fn func(), row_fn func(row *sql.Rows), select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return
	}
	start_fn()
	for rows.Next() {
		row_fn(rows)
	}
}

// TODO Forall_done

// TODO count?
func Forall_err(proc_name string, start_fn func(), row_fn func(row *sql.Rows) error, select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return
	}
	start_fn()
	for rows.Next() {
		if err = row_fn(rows); err != nil {
			log.Printf("Err: [%s]: %q\n", proc_name, err)
		}
	}
}

func Forall_exit(proc_name string, start_fn func(), row_fn func(row *sql.Rows) error, select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return
	}
	start_fn()
	for rows.Next() {
		if err = row_fn(rows); err != nil {
			log.Fatalf("Crit: [%s]: %q\n", proc_name, err)
			panic("Forall_exit")
		}
	}
}

/*
// TODO func Forall(calling_fn_name string, start_fn func(), row_fn func(*sql.Rows) err, select_statement *sql.Stmt, args ...any) {
func Forall(proc_name string, start_fn func(), row_fn func(*sql.Rows) error, select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("error: [%s]: %q\n",  proc_name,  err)
		return
	}
	start_fn()
	for rows.Next() {
		if err = row_fn(rows); err != nil {
			log.Printf("error: [%s]: %q\n",  proc_name,  err)
			return
		}
	}
}*/
