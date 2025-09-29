package blender_ui

import "github.com/samuel-jimenez/qc_data_entry/GUI"

func Refresh_globals(font_size int) {

	GUI.GROUPBOX_CUSHION = font_size * 3 / 2
	GUI.TOP_SPACER_HEIGHT = GUI.GROUPBOX_CUSHION + 2

	GUI.DATA_MARGIN = 10

	GUI.LABEL_WIDTH = 10 * font_size
	GUI.PRODUCT_FIELD_WIDTH = 15 * font_size
	GUI.PRODUCT_FIELD_HEIGHT = font_size*16/10 + 8
	// GUI.EDIT_FIELD_HEIGHT = 24
	GUI.EDIT_FIELD_HEIGHT = 3*font_size - 2

}
