package viewer

var (
	COL_WIDTH_TIME       = 150
	COL_WIDTH_LOT        = 100
	COL_WIDTH_SAMPLE_PT  = 50
	COL_WIDTH_SAMPLE_BIN = 90
	COL_WIDTH_DATA       = 70

	WINDOW_EDGE  = 8
	SCROLL_WIDTH = 17

	WINDOW_WIDTH  = 2*COL_WIDTH_TIME + COL_WIDTH_LOT + COL_WIDTH_SAMPLE_PT + COL_WIDTH_SAMPLE_BIN + 4*COL_WIDTH_DATA + 2*WINDOW_EDGE + SCROLL_WIDTH
	WINDOW_HEIGHT = 600

	PRODUCT_FIELD_WIDTH = 150
	// GUI.EDIT_FIELD_HEIGHT        = 24

	BUTTON_WIDTH  = 100
	BUTTON_HEIGHT = 40
	// 	200
	// 50
	SMOL_BUTTON_EDGE = 15
	HEADER_HEIGHT    = SMOL_BUTTON_EDGE

	INTER_SPACER_HEIGHT = 2

	COL_KEY_TIME   = "time_stamp"
	COL_LABEL_TIME = "Time Stamp"

	// COL_KEY_LOT = "product_lot_id" //TODO
	COL_KEY_LOT   = "lot_name"
	COL_LABEL_LOT = "Lot Number"
	COL_ITEMS_LOT []string

	// sample_point_id
	COL_KEY_SAMPLE_PT   = "sample_point"
	COL_LABEL_SAMPLE_PT = "Sample Point"
	COL_ITEMS_SAMPLE_PT []string

	// qc_sample_storage_id
	COL_KEY_SAMPLE_BIN   = "qc_sample_storage_name"
	COL_LABEL_SAMPLE_BIN = "Sample Bin"
	COL_ITEMS_SAMPLE_BIN []string

	COL_KEY_PH   = "ph"
	COL_LABEL_PH = "pH"

	/*sg_title*/
	COL_KEY_SG   = "specific_gravity"
	COL_LABEL_SG = "SG"
	// COL_LABEL_SG = "Specific Gravity"

	// Density_title
	COL_KEY_DENSITY   = "density"
	COL_LABEL_DENSITY = "Density"

	COL_KEY_STRING      = "string_test"
	COL_LABEL_STRING    = "String"
	COL_KEY_VISCOSITY   = "viscosity"
	COL_LABEL_VISCOSITY = "Viscosity"
)
