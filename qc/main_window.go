package qc

import (
	"encoding/json"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/QR"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/qc/subpanels"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/windigo"
)

func Show_window() {

	log.Println("Info: Process started")
	// DEBUG
	// log.Println(time.Now().UTC().UnixNano())

	//
	//
	// definitions
	//
	//
	window_title := "QC Data Entry"

	windigo.DefaultFont = windigo.NewFont("MS Shell Dlg 2", GUI.BASE_FONT_SIZE, windigo.FontNormal)

	//
	//
	//
	// build window
	//
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetText(window_title)

	keygrab := windigo.NewEdit(mainWindow)
	keygrab.Hide()

	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := NewTopPanelView(mainWindow)
	product_panel.SetMainWindow(mainWindow)

	tabs := windigo.NewTabView(mainWindow)
	tab_wb := tabs.AddAutoPanel("Water Based")
	tab_oil := tabs.AddAutoPanel("Oil Based")
	tab_fr := tabs.AddAutoPanel("Friction Reducer")

	//
	//
	// Dock
	//
	//

	dock.Dock(product_panel, windigo.Top)
	dock.Dock(tabs, windigo.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), windigo.Fill) // tab panels dock just below tabs and fill area

	threads.Status_bar = windigo.NewStatusBar(mainWindow)
	mainWindow.SetStatusBar(threads.Status_bar)

	//
	//
	//
	//
	// functionality

	panel_water_based := subpanels.Show_water_based(tab_wb, product_panel.QC_Product, product_panel.BaseProduct)
	panel_oil_based := subpanels.Show_oil_based(tab_oil, product_panel.QC_Product, product_panel.BaseProduct)
	panel_fr := subpanels.Show_fr(tab_fr, product_panel.QC_Product, product_panel.BaseProduct)

	product_panel.SetUpdate(func() {
		product_panel.container_field.Update(product_panel.QC_Product.Container_type)
		log.Println("Debug: update new_product_cb", product_panel.QC_Product)
		panel_water_based.Update(product_panel.QC_Product)
		panel_oil_based.Update(product_panel.QC_Product)
		panel_fr.Update(product_panel.QC_Product)
		panel_fr.ChangeContainer(product_panel.QC_Product)
	})
	product_panel.ChangeContainer = panel_fr.ChangeContainer
	product_panel.SetCurrentTab = func(i int) {
		tabs.SetCurrent(i)
	}

	//
	//
	//
	//
	// sizing
	refresh := func(font_size int) {
		refresh_globals(font_size)

		mainWindow.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)
		// product_panel.RefreshSize()
		product_panel.RefreshSize(font_size)

		tabs.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		panel_water_based.Refresh()
		panel_oil_based.Refresh()
		panel_fr.Refresh()

	}

	set_font := func(font_size int) {
		GUI.BASE_FONT_SIZE = font_size

		config.Main_config.Set("font_size", GUI.BASE_FONT_SIZE)
		config.Write_config(config.Main_config)

		old_font := windigo.DefaultFont
		windigo.DefaultFont = windigo.NewFont(old_font.Family(), GUI.BASE_FONT_SIZE, 0)
		old_font.Dispose()

		mainWindow.SetFont(windigo.DefaultFont)
		product_panel.SetFont(windigo.DefaultFont)
		tabs.SetFont(windigo.DefaultFont)

		panel_water_based.SetFont(windigo.DefaultFont)
		panel_oil_based.SetFont(windigo.DefaultFont)
		panel_fr.SetFont(windigo.DefaultFont)
		threads.Status_bar.SetFont(windigo.DefaultFont)
		refresh(font_size)

	}
	set_font(GUI.BASE_FONT_SIZE)

	AddShortcuts(mainWindow, keygrab, set_font, product_panel.PopQRData)

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)
	mainWindow.RunMainLoop() // Must call to start event loop.
}

func AddShortcuts(mainWindow *windigo.Form, keygrab *windigo.Edit, set_font func(int), PopQRData func(QR.QRJson)) {
	// QR keyboard handling

	keygrab.OnSetFocus().Bind(func(e *windigo.Event) {
		keygrab.SetText("{")
		keygrab.SelectText(1, 1)
	})

	mainWindow.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModShift, Key: windigo.KeyOEM4}, // {
		func() bool {
			keygrab.SetText("")
			keygrab.SetFocus()
			return true
		})

	mainWindow.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModShift, Key: windigo.KeyOEM6}, // }
		func() bool {
			var product QR.QRJson

			qr_json := keygrab.Text() + "}"
			log.Println("debug: ReadFromScanner: ", qr_json)
			err := json.Unmarshal([]byte(qr_json), &product)
			if err == nil {
				PopQRData(product)
			} else {
				log.Printf("error: [%s]: %q\n", "qr_json_keygrab", err)
			}
			keygrab.SetText("")
			mainWindow.SetFocus()
			return true
		})

	// Resize handling
	mainWindow.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModControl, Key: windigo.KeyOEMPlus}, // {
		func() bool {
			set_font(GUI.BASE_FONT_SIZE + 1)
			return true
		})

	mainWindow.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModControl, Key: windigo.KeyOEMMinus}, // {
		func() bool {
			set_font(GUI.BASE_FONT_SIZE - 1)
			return true
		})
}

func wndOnClose(arg *windigo.Event) {
	GUI.OKPen.Dispose()
	GUI.ErroredPen.Dispose()
	windigo.Exit()
}
