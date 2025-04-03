package main

import "fmt"

import "github.com/tadvi/winc"

// https://pkg.go.dev/github.com/tadvi/winc

func main() {
	show_window()
}
func show_window() {

	fmt.Println("Hello, World!")
	// DEBUG

	mainWindow := winc.NewForm(nil)
	mainWindow.SetSize(400, 300) // (width, height)
	mainWindow.SetText("Hello World Demo")
	show_water_based(mainWindow)
}
func show_water_based(mainWindow *winc.Form) {
	label_col := 10
	field_col := 120

	visual_row := 20
	edit_row := 50
	submit_row := 80
	group_row := 120

	// type
	visual_label := winc.NewLabel(mainWindow)
	visual_label.SetPos(label_col, visual_row)

	visual_label.SetText("Visual Inspection")

	visual_field := winc.NewCheckBox(mainWindow)
	visual_field.SetPos(field_col, visual_row)
	// visual_label.OnClick().Bind(func(e *winc.Event) {
	// 		visual_field.SetFocus()
	// })

	group_label := winc.NewLabel(mainWindow)
	group_label.SetPos(label_col, group_row)

	group_label.SetText("group Inspection")

	group_field := winc.NewGroupBox(mainWindow)
	group_field.SetPos(field_col, group_row)

	edit_label := winc.NewLabel(mainWindow)
	edit_label.SetPos(label_col, edit_row)

	edit_box := winc.NewEdit(mainWindow)
	edit_box.SetPos(field_col, edit_row)
	// Most Controls have default size unless SetSize is called.
	edit_box.SetText("edit text")

	btn := winc.NewPushButton(mainWindow)
	btn.SetText("Show or Hide")
	btn.SetPos(40, submit_row) // (x, y)
	btn.SetSize(100, 40)       // (width, height)
	btn.OnClick().Bind(func(e *winc.Event) {
		if edit_box.Visible() {
			edit_box.Hide()
		} else {
			edit_box.Show()
		}
	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	winc.RunMainLoop() // Must call to start event loop.
}

func show_fr(mainWindow *winc.Form) {

	fmt.Println("Hello, World!")
	// DEBUG

	edit_label := winc.NewEdit(mainWindow)
	edit_label.SetText("Edit")

	edit_box := winc.NewEdit(mainWindow)
	edit_box.SetPos(10, 20)
	// Most Controls have default size unless SetSize is called.
	edit_box.SetText("edit text")

	btn := winc.NewPushButton(mainWindow)
	btn.SetText("Show or Hide")
	btn.SetPos(40, 50)   // (x, y)
	btn.SetSize(100, 40) // (width, height)
	btn.OnClick().Bind(func(e *winc.Event) {
		if edit_box.Visible() {
			edit_box.Hide()
		} else {
			edit_box.Show()
		}
	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	winc.RunMainLoop() // Must call to start event loop.
}

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}
