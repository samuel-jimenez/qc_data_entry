package viewer

import "github.com/samuel-jimenez/qc_data_entry/GUI"

var (
	COL_WIDTH_TIME,
	COL_WIDTH_LOT,
	COL_WIDTH_SAMPLE_PT,
	COL_WIDTH_SAMPLE_BIN,
	COL_WIDTH_DATA,

	WINDOW_EDGE,
	SCROLL_WIDTH,

	WINDOW_WIDTH,
	WINDOW_HEIGHT,

	HEADER_HEIGHT int

	COL_KEY_PRODUCT   = "product_id"
	COL_LABEL_PRODUCT = "Product"
	// COL_ITEMS_PRODUCT []string

	COL_KEY_MONIKER   = "product_moniker_name"
	COL_LABEL_MONIKER = "Product Moniker"
	// COL_ITEMS_MONIKER []string

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

func Refresh_globals(font_size int) {

	GUI.Refresh_globals(font_size)

	COL_WIDTH_TIME = 15 * font_size
	COL_WIDTH_LOT = 10 * font_size
	COL_WIDTH_SAMPLE_PT = 5 * font_size
	COL_WIDTH_SAMPLE_BIN = 9 * font_size
	COL_WIDTH_DATA = 7 * font_size

	WINDOW_EDGE = 8
	SCROLL_WIDTH = 17

	WINDOW_WIDTH = 2*COL_WIDTH_TIME + COL_WIDTH_LOT + COL_WIDTH_SAMPLE_PT + COL_WIDTH_SAMPLE_BIN +
		4*COL_WIDTH_DATA + //				data
		2*(2*COL_WIDTH_TIME+COL_WIDTH_LOT) + //		components
		2*WINDOW_EDGE + SCROLL_WIDTH + //		window cruft
		COL_WIDTH_DATA //				selection
	WINDOW_HEIGHT = 60 * font_size

	GUI.TOP_SPACER_HEIGHT = 20
	HEADER_HEIGHT = 2 * font_size

	// GUI.EDIT_FIELD_HEIGHT = 24
	GUI.EDIT_FIELD_HEIGHT = 2 * font_size
	// GUI.EDIT_FIELD_HEIGHT = 3*font_size - 2

	GUI.BUTTON_WIDTH = 5 * font_size
	GUI.SMOL_BUTTON_WIDTH = 3 * font_size / 2

	GUI.TOP_PANEL_WIDTH = WINDOW_WIDTH
	GUI.HPANEL_WIDTH = GUI.TOP_PANEL_WIDTH

}
