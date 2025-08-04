package DB

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

var (
	err error

	DB_VERSION = "0.0.3"

	// recipe_list
	DB_Select_product_recipe, DB_Insert_product_recipe,
	// recipe_components
	DB_Select_recipe_components_id, DB_Insert_recipe_component, DB_Update_recipe_component, DB_Delete_recipe_component,
	// component_types
	DB_Select_name_component_types, DB_Select_all_component_types, DB_Insert_component_types,
	DB_Select_component_type_product, DB_Insert_internal_product_component_type, DB_Insert_inbound_product_component_type,
	DB_Select_inbound_product_component_type_id,
	// inbound_product
	DB_Select_inbound_product_name, DB_Insert_inbound_product,
	// container_list
	DB_Select_container_id, DB_Select_container_all, DB_Insert_container,
	// inbound_provider_list
	DB_Select_inbound_provider_id, DB_Select_inbound_provider_all, DB_Insert_inbound_provider,
	// status_list
	DB_Select_all_status_list, DB_Select_name_status_list,
	// inbound_lot
	DB_Select_inbound_lot_name, DB_Select_inbound_lot_all, DB_Insert_inbound_lot, DB_Update_inbound_lot_status,
	DB_Select_inbound_lot_status, DB_Select_inbound_lot_recipe, DB_Select_inbound_lot_components,
	// inbound_relabel
	DB_Insert_inbound_relabel,
	// component_list
	DB_Select_inbound_blend_component, DB_Insert_inbound_blend_component,
	DB_Select_internal_blend_component, DB_Insert_internal_blend_component,

	// blend_components
	DB_Insert_Product_blend,
	// lot_list
	DB_Insert_lot, DB_Select_lot,
	// product_lot
	db_select_id_lot, DB_Insert_product_lot,
	DB_Select_product_lot_all, DB_Select_product_lot_product,
	DB_Select_blend_lot, DB_Insert_blend_lot, DB_Update_lot_recipe,
	DB_Update_lot_customer,
	// product_line
	db_select_product_id, db_insert_product,
	DB_Select_product_info,
	// product_customer_line
	DB_Select_product_customer_id, DB_Select_product_customer_info,
	db_select_product_customer, db_insert_product_customer,
	DB_Select_all_sample_points, DB_Select_product_sample_points,
	DB_insert_sample_point,
	DB_insert_measurement,
	DB_Insert_appearance,
	DB_Select_product_details,
	DB_Upsert_product_details, DB_Upsert_product_type,
	DB_Select_product_coa_details, DB_Upsert_product_coa_details *sql.Stmt

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

func Check_db(db *sql.DB) {

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
		err = errors.New(fmt.Sprintf("Database version mismatch: Required: %s, found: %s", DB_VERSION, found_db_version))
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

	DB_Insert_product_recipe = PrepareOrElse(db, `
	insert into bs.recipe_list
		(product_id)
		values (?)
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

	//TODO
	// 	DB_Select_component_type_product = PrepareOrElse(db, `
	// select component_type_id, component_type_name, product_id, inbound_product_id
	// from bs.component_types
	// left join bs.component_type_product_internal using (component_type_id)
	// left join bs.component_type_product_inbound using (component_type_id)
	// where component_type_id = ?
	// `)

	DB_Select_component_type_product = PrepareOrElse(db, `
select false, product_lot_id, product_name_internal, lot_name from bs.component_type_product_internal
	join bs.product_line using (product_id)
	join bs.product_lot using (product_id)
	join bs.lot_list using (lot_id)
where component_type_id = ?1
union
select true, inbound_lot_id,inbound_product_name, inbound_lot_name from bs.component_type_product_inbound
	join bs.inbound_product using (inbound_product_id)
	join bs.inbound_lot using (inbound_product_id)
where component_type_id = ?1
`)

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
select inbound_product_id, inbound_product_name from bs.component_type_product_inbound
left join bs.inbound_product using (inbound_product_id)
where component_type_id = ?1
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
	// status_list
	DB_Select_all_status_list = PrepareOrElse(db, `
select
	status_id, status_name
from bs.status_list
order by status_id
	`)
	DB_Select_name_status_list = PrepareOrElse(db, `
select
	status_id
from bs.status_list
where status_name = ?
	`)

	// inbound_lot
	DB_Select_inbound_lot_name = PrepareOrElse(db, `
select
	inbound_lot_id, inbound_lot_name,inbound_product_name, inbound_provider_name, container_name, status_name
from bs.inbound_lot
join bs.inbound_product using (inbound_product_id)
join bs.inbound_provider_list using (inbound_provider_id)
join bs.container_list using (container_id)
join bs.status_list using (status_id)
where inbound_lot_name = ?
	`)
	DB_Select_inbound_lot_all = PrepareOrElse(db, `
select
	inbound_lot_id, inbound_lot_name, inbound_product_id, inbound_product_name, inbound_provider_id, inbound_provider_name, container_id, container_name, status_id,status_name
from bs.inbound_lot
join bs.inbound_product 		using (inbound_product_id)
join bs.inbound_provider_list 		using (inbound_provider_id)
join bs.container_list 			using (container_id)
join bs.status_list 			using (status_id)
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
	status_id=?2
where inbound_lot_id=?1
`)

	DB_Select_inbound_lot_status = PrepareOrElse(db, `
select
	inbound_lot_name
from bs.inbound_lot
join bs.status_list 			using (status_id)
where status_name = ?
	`)

	DB_Select_inbound_lot_recipe = PrepareOrElse(db, `
select
	recipe_id, product_id, component_type_id, component_type_amount, component_add_order
from bs.inbound_lot
join bs.component_type_product_inbound 	using (inbound_product_id)
join bs.recipe_components 		using (component_type_id)
join bs.recipe_list			using (recipe_id)
where inbound_lot_id =  ?
order by recipe_id
	`)

	DB_Select_inbound_lot_components = PrepareOrElse(db, `
select
	inbound_lot_id,  component_type_id, component_type_amount, component_add_order
from bs.inbound_lot
join bs.status_list 			using (status_id)
join bs.component_type_product_inbound 	using (inbound_product_id)
join bs.recipe_components 		using (component_type_id)
where recipe_id = ?1
	and component_type_id != ?2
	and status_name = ?3
order by component_add_order
	`)

	// inbound_relabel
	DB_Insert_inbound_relabel = PrepareOrElse(db, `
	insert into bs.inbound_relabel
		(lot_id,inbound_lot_id,container_id)
		values (?,?,?)
	returning inbound_relabel_id
	`)

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
	DB_Insert_Product_blend = PrepareOrElse(db, `
	insert into bs.blend_components
		(product_lot_id,recipe_id,component_id,component_type_amount)
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
`)
	DB_Select_product_lot_all = PrepareOrElse(db, `
select product_lot_id, lot_name
from bs.product_lot
join bs.lot_list using (lot_id)
order by lot_name
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

	DB_Select_product_info = PrepareOrElse(db, `
	select product_id, product_name_internal, product_moniker_name
		from bs.product_line
		join bs.product_moniker using (product_moniker_id)
	order by product_moniker_name,product_name_internal
	`)

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

	DB_insert_sample_point = PrepareOrElse(db, `
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
		join bs.qc_samples using (product_lot_id)
		join bs.product_sample_points using (sample_point_id)
		where product_id = ?
		order by sample_point_id
	`)

	DB_insert_measurement = PrepareOrElse(db, `
	with
		val (product_lot_id, sample_point, time_stamp, ph, specific_gravity, string_test, viscosity) as (
			values
				(?, ?, ?, ?, ?, ?, ?)
		),
		sel as (
			select product_lot_id, sample_point_id, sample_point, time_stamp, ph, specific_gravity, string_test, viscosity
			from val
			left join bs.product_sample_points using (sample_point)
		)
	insert into bs.qc_samples (product_lot_id, sample_point_id, time_stamp, ph, specific_gravity, string_test, viscosity)
	select product_lot_id, sample_point_id, time_stamp, ph, specific_gravity, string_test, viscosity
		from   sel
	returning qc_id;
	`)

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

	DB_Select_product_details = PrepareOrElse(db, `
	with
		measured as (
			select
				product_id,

				product_type_id,
				product_appearance_text,


				ph_min,
				ph_target,
				ph_max,

				specific_gravity_min,
				specific_gravity_target,
				specific_gravity_max,

				density_min,
				density_target,
				density_max,

				string_test_min,
				string_test_target,
				string_test_max,

				viscosity_min,
				viscosity_target,
				viscosity_max
			from bs.product_ranges_measured
				left join bs.product_appearance using (product_appearance_id)
			where product_id = ?1),

		published as (
			select
				product_id,
				product_appearance_text,

				ph_min,
				ph_target,
				ph_max,

				specific_gravity_min,
				specific_gravity_target,
				specific_gravity_max,

				density_min,
				density_target,
				density_max,

				string_test_min,
				string_test_target,
				string_test_max,

				viscosity_min,
				viscosity_target,
				viscosity_max

			from bs.product_ranges_published
				left join bs.product_appearance using (product_appearance_id)
			where product_id = ?1)
	select
		product_type_id,
		container_type_id,
		coalesce(measured.product_appearance_text, published.product_appearance_text) as product_appearance_text,

		coalesce(measured.ph_min, published.ph_min) as ph_min,
		coalesce(measured.ph_target, published.ph_target) as ph_target,
		coalesce(measured.ph_max, published.ph_max) as ph_max,

		coalesce(measured.specific_gravity_min, published.specific_gravity_min) as specific_gravity_min,
		coalesce(measured.specific_gravity_target, published.specific_gravity_target) as specific_gravity_target,
		coalesce(measured.specific_gravity_max, published.specific_gravity_max) as specific_gravity_max,

		coalesce(measured.density_min, published.density_min) as density_min,
		coalesce(measured.density_target, published.density_target) as density_target,
		coalesce(measured.density_max, published.density_max) as density_max,

		coalesce(measured.string_test_min, published.string_test_min) as string_test_min,
		coalesce(measured.string_test_target, published.string_test_target) as string_test_target,
		coalesce(measured.string_test_max, published.string_test_max) as string_test_max,

		coalesce(measured.viscosity_min, published.viscosity_min) as viscosity_min,
		coalesce(measured.viscosity_target, published.viscosity_target) as viscosity_target,
		coalesce(measured.viscosity_max, published.viscosity_max) as viscosity_max
	from measured
		full join published using (product_id)
		join bs.product_types using (product_type_id)
	where product_id = ?1
	`)

	DB_Select_product_coa_details = PrepareOrElse(db, `
	select
		product_appearance_text,

		ph_min,
		ph_target,
		ph_max,

		specific_gravity_min,
		specific_gravity_target,
		specific_gravity_max,

		density_min,
		density_target,
		density_max,

		string_test_min,
		string_test_target,
		string_test_max,

		viscosity_min,
		viscosity_target,
		viscosity_max

	from bs.product_ranges_published
	join bs.product_appearance using (product_appearance_id)
	where product_id = ?
	`)

	DB_Upsert_product_details = PrepareOrElse(db, `
	with
		val
			(product_id,
			product_type_id,
			product_appearance_text,

			ph_min,
			ph_target,
			ph_max,

			specific_gravity_min,
			specific_gravity_target,
			specific_gravity_max,

			density_min,
			density_target,
			density_max,

			string_test_min,
			string_test_target,
			string_test_max,

			viscosity_min,
			viscosity_target,
			viscosity_max)
		as (
			values (
				?,?,?,
				?,?,?,
				?,?,?,
				?,?,?,
				?,?,?,
				?,?,?)
		),
		sel as (
			select
				product_id,
				product_type_id,
				product_appearance_id,


				ph_min,
				ph_target,
				ph_max,

				specific_gravity_min,
				specific_gravity_target,
				specific_gravity_max,

				density_min,
				density_target,
				density_max,

				string_test_min,
				string_test_target,
				string_test_max,

				viscosity_min,
				viscosity_target,
				viscosity_max
			from val
			left join bs.product_appearance using (product_appearance_text)
		)
	insert into bs.product_ranges_measured
		(product_id,
		product_type_id,
		product_appearance_id,


		ph_min,
		ph_target,
		ph_max,

		specific_gravity_min,
		specific_gravity_target,
		specific_gravity_max,

		density_min,
		density_target,
		density_max,

		string_test_min,
		string_test_target,
		string_test_max,

		viscosity_min,
		viscosity_target,
		viscosity_max)
	select
				product_id,
				product_type_id,
				product_appearance_id,


				ph_min,
				ph_target,
				ph_max,

				specific_gravity_min,
				specific_gravity_target,
				specific_gravity_max,

				density_min,
				density_target,
				density_max,

				string_test_min,
				string_test_target,
				string_test_max,

				viscosity_min,
				viscosity_target,
				viscosity_max
			from sel
			where true
	on conflict(product_id) do update set

		product_type_id=excluded.product_type_id,
		product_appearance_id=excluded.product_appearance_id,

		ph_min=excluded.ph_min,
		ph_target=excluded.ph_target,
		ph_max=excluded.ph_max,

		specific_gravity_min=excluded.specific_gravity_min,
		specific_gravity_target=excluded.specific_gravity_target,
		specific_gravity_max=excluded.specific_gravity_max,

		density_min=excluded.density_min,
		density_target=excluded.density_target,
		density_max=excluded.density_max,

		string_test_min=excluded.string_test_min,
		string_test_target=excluded.string_test_target,
		string_test_max=excluded.string_test_max,

		viscosity_min=excluded.viscosity_min,
		viscosity_target=excluded.viscosity_target,
		viscosity_max=excluded.viscosity_max

		returning range_id
		`)

	DB_Upsert_product_type = PrepareOrElse(db, `
	insert into bs.product_ranges_measured
		(product_id,
		product_type_id)
	values (?,?)
	on conflict(product_id) do update set

		product_type_id=excluded.product_type_id

		returning range_id
		`)
	DB_Upsert_product_coa_details = PrepareOrElse(db, `
	with
		val
			(product_id,
			product_type_id,
			product_appearance_text,

			ph_min,
			ph_target,
			ph_max,

			specific_gravity_min,
			specific_gravity_target,
			specific_gravity_max,

			density_min,
			density_target,
			density_max,

			string_test_min,
			string_test_target,
			string_test_max,

			viscosity_min,
			viscosity_target,
			viscosity_max)
		as (
			values (
				?,?,?,
				?,?,?,
				?,?,?,
				?,?,?,
				?,?,?,
				?,?,?)
		),
		sel as (
			select
				product_id,
				product_appearance_id,


				ph_min,
				ph_target,
				ph_max,

				specific_gravity_min,
				specific_gravity_target,
				specific_gravity_max,

				density_min,
				density_target,
				density_max,

				string_test_min,
				string_test_target,
				string_test_max,

				viscosity_min,
				viscosity_target,
				viscosity_max
			from val
			left join bs.product_appearance using (product_appearance_text)
		)
	insert into bs.product_ranges_published
		(product_id,
		product_appearance_id,


		ph_min,
		ph_target,
		ph_max,

		specific_gravity_min,
		specific_gravity_target,
		specific_gravity_max,

		density_min,
		density_target,
		density_max,

		string_test_min,
		string_test_target,
		string_test_max,

		viscosity_min,
		viscosity_target,
		viscosity_max)
	select
				product_id,
				product_appearance_id,


				ph_min,
				ph_target,
				ph_max,

				specific_gravity_min,
				specific_gravity_target,
				specific_gravity_max,

				density_min,
				density_target,
				density_max,

				string_test_min,
				string_test_target,
				string_test_max,

				viscosity_min,
				viscosity_target,
				viscosity_max
			from sel
			where true
	on conflict(product_id) do update set

		product_appearance_id=excluded.product_appearance_id,

		ph_min=excluded.ph_min,
		ph_target=excluded.ph_target,
		ph_max=excluded.ph_max,

		specific_gravity_min=excluded.specific_gravity_min,
		specific_gravity_target=excluded.specific_gravity_target,
		specific_gravity_max=excluded.specific_gravity_max,

		density_min=excluded.density_min,
		density_target=excluded.density_target,
		density_max=excluded.density_max,

		string_test_min=excluded.string_test_min,
		string_test_target=excluded.string_test_target,
		string_test_max=excluded.string_test_max,

		viscosity_min=excluded.viscosity_min,
		viscosity_target=excluded.viscosity_target,
		viscosity_max=excluded.viscosity_max

		returning qc_range_id
		`)

}

func Select_Error(proc_name string, query *sql.Row, args ...any) error {
	err := query.Scan(args...)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
	}
	return err
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
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
	}
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
		//no rows
		insel_id = Insert(proc_name, insert_statement, args...)
	}
	return insel_id
}

func Insel_product_id(product_name_full string) int64 {

	product_moniker_name, product_name_internal, _ := strings.Cut(product_name_full, " ")

	return Insel("Insel_product_id", db_insert_product, db_select_product_id, product_name_internal, product_moniker_name)
}

func Insel_lot_id(lot_name string) int64 {
	return Insel("Insel_lot_id", DB_Insert_lot, DB_Select_lot, lot_name)
}

func Insel_product_lot_id(lot_name string, product_id int64) int64 {
	return Insel("Insel_product_lot_id", DB_Insert_product_lot, db_select_id_lot, Insel_lot_id(lot_name), product_id)
}

func Insel_product_name_customer(product_name_customer string, product_id int64) int64 {
	return Insel("Insel_product_name_customer", db_insert_product_customer, db_select_product_customer, product_name_customer, product_id)
}

func Forall(proc_name string, start_fn func(), row_fn func(row *sql.Rows), select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("error: [%s]: %q\n", proc_name, err)
		return
	}
	start_fn()
	for rows.Next() {
		row_fn(rows)
	}
}

func Forall_err(proc_name string, start_fn func(), row_fn func(row *sql.Rows) error, select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("error: [%s]: %q\n", proc_name, err)
		return
	}
	start_fn()
	for rows.Next() {
		if err = row_fn(rows); err != nil {
			log.Printf("error: [%s]: %q\n", proc_name, err)
		}
	}
}

func Forall_exit(proc_name string, start_fn func(), row_fn func(row *sql.Rows) error, select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("error: [%s]: %q\n", proc_name, err)
		return
	}
	start_fn()
	for rows.Next() {
		if err = row_fn(rows); err != nil {
			log.Printf("error: [%s]: %q\n", proc_name, err)
			return
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
