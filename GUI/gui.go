package GUI

import (
	"database/sql"
	"log"
	"strings"

	"github.com/samuel-jimenez/windigo"
)

var (
	BASE_FONT_SIZE = 10
)

func Fill_combobox_from_query_rows(control windigo.ComboBoxable, selected_rows *sql.Rows, err error, fn func(*sql.Rows)) {

	if err != nil {
		log.Printf("error: %q: %s\n", err, "fill_combobox_from_query")
		// return -1
	}
	control.DeleteAllItems()
	i := 0
	for selected_rows.Next() {
		fn(selected_rows)
		i++
	}
	if i == 1 {
		control.SetSelectedItem(0)
	}
}

func Fill_combobox_from_query_fn(control windigo.ComboBoxable, select_statement *sql.Stmt, select_id int64, fn func(*sql.Rows)) {
	rows, err := select_statement.Query(select_id)
	Fill_combobox_from_query_rows(control, rows, err, fn)
}

func Fill_combobox_from_query(control windigo.ComboBoxable, select_statement *sql.Stmt, select_id int64) {
	Fill_combobox_from_query_fn(control, select_statement, select_id, func(rows *sql.Rows) {
		var (
			id   uint8
			name string
		)

		if err := rows.Scan(&id, &name); err == nil {
			// data[id] = value
			control.AddItem(name)
		} else {
			log.Printf("error: %q: %s\n", err, "fill_combobox_from_query")
			// return -1
		}
	})
}

type ComboBox struct {
	*windigo.LabeledComboBox
}

// TODO SetFont decrease followed by SetLabeledSize disappears lists
func (control *ComboBox) SetFont(font *windigo.Font) {
	control.ComboBox.SetFont(font)
	control.Label().SetFont(font)
}

func (control *ComboBox) SetLabeledSize(label_width, control_width, height int) {
	control.SetSize(label_width+control_width, 40*height)
	control.SetSize(label_width+control_width, height)
	control.Label().SetSize(label_width, height)
}

func NewComboBox(parent windigo.Controller, field_text string) *ComboBox {
	combobox_field := &ComboBox{windigo.NewLabeledComboBox(parent, field_text)}
	combobox_field.OnKillFocus().Bind(func(e *windigo.Event) {
		combobox_field.SetText(strings.ToUpper(strings.TrimSpace(combobox_field.Text())))
	})

	return combobox_field
}

func Show_combobox(parent windigo.Controller, label_width, control_width, height int, field_text string) *ComboBox {
	combobox_field := &ComboBox{windigo.NewSizedLabeledComboBox(parent, label_width, control_width, height, field_text)}
	combobox_field.OnKillFocus().Bind(func(e *windigo.Event) {
		combobox_field.SetText(strings.ToUpper(strings.TrimSpace(combobox_field.Text())))
	})

	return combobox_field
}
