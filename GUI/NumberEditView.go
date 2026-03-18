package GUI

// TODO todo NumberEditView

import (
	"math"
	"strconv"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/windigo"
)

/*
 * Checker
 *
 */
type Checker interface {
	Check(float64) bool
	CheckAll(...float64) []bool
}

/*
 * NumberEditViewable
 *
 */
type NumberEditViewable interface {
	ErrableView
	windigo.Editable
	windigo.DiffLabelable
	Get() float64
	GetFixed() float64
	Set(val float64)
	Clear()
	Check(bool)
	SetFont(font *windigo.Font)
	SetLabeledSize(label_width, control_width, height int)
	Entangle(other_field *NumberEditView, range_field Checker, delta_max float64)
}

/*
 * NumberEditView
 *
 */
type NumberEditView struct {
	ErrableView
	NumbEditView
	windigo.Labeled
}

func NumberEditView_from_LabeledEdit(label *windigo.LabeledEdit) *NumberEditView {
	return &NumberEditView{&View{ComponentFrame: label.ComponentFrame}, NumbEditView{label.Edit}, windigo.Labeled{FieldLabel: label.Label()}}
}

func NumberEditView_from_new(parent windigo.Controller, field_text string) *NumberEditView {
	edit_field := NumberEditView_from_LabeledEdit(windigo.LabeledEdit_from_new(parent, field_text))
	return edit_field
}

func NumberEditView_with_Change_from_new(parent windigo.Controller, field_text string, range_field Checker) *NumberEditView {
	edit_field := NumberEditView_from_new(parent, field_text)
	edit_field.OnChange().Bind(func(e *windigo.Event) {
		edit_field.Check(range_field.Check(edit_field.GetFixed()))
	})
	return edit_field
}

func NumberEditView_with_PointlessChange_SG_from_new(parent windigo.Controller, field_text string, range_field Checker) *NumberEditView {
	edit_field := NumberEditView_from_new(parent, field_text)
	edit_field.OnChange().Bind(func(e *windigo.Event) {
		edit_field.Check(range_field.Check(edit_field.GetPointless_SG()))
	})
	return edit_field
}

func NumberEditView_with_PointlessChange_PH_from_new(parent windigo.Controller, field_text string, range_field Checker) *NumberEditView {
	edit_field := NumberEditView_from_new(parent, field_text)
	edit_field.OnChange().Bind(func(e *windigo.Event) {
		edit_field.Check(range_field.Check(edit_field.GetPointless_PH()))
	})
	return edit_field
}

func (control *NumberEditView) GetFixed() float64 {
	start, end := control.Selected()
	// IndexAny(s, chars string) int
	control.SetText(strings.TrimSpace(control.Text()))
	// mass_field.SelectText(-1, -1)
	control.SelectText(start, end)
	// mass_field.SelectText(-1, 0)

	return control.Get()
}

// TODO collapse these
func (control *NumberEditView) GetPointless_PH() float64 {
	start, end := control.Selected()

	val, _ := strconv.ParseFloat(strings.TrimSpace(control.Text()), 64)

	for val >= 20 {
		val /= 10
		// check position. if we just backspaced the decimal point, don't put it back in the way
		if start != 1 {
			start++
			end++
		}
	}
	control.Set(val)
	control.SelectText(start, end)

	return val
}

func (control *NumberEditView) GetPointless_SG() float64 {

	self_dec_pos := 1
	width := 6
	last_key := 0

	return control.ParsePointless(self_dec_pos, width, last_key, formats.Format_fixed_sg)
}

func (control *NumberEditView) GetPointlessMass() float64 {
	// TODO onKeyUp

	self_dec_pos := 2
	width := 5
	last_key := 0

	return control.ParsePointless(self_dec_pos, width, last_key, formats.Format_mass)
}

func (control *NumberEditView) ParsePointless(self_dec_pos, width int, last_key int, format func(float64) string) float64 {
	start, end := control.Selected()

	cursor := start

	// TODO onKeyUp

	remangle := false

	text := control.Text()

	// text.retain(|c| match c {
	//     '0':='9' | '.' => true,
	//     _ => false,
	// });
	// dec_pos := text.Index();
	dec_pos := strings.Index(text, ".")
	if dec_pos != -1 {

		dec_pos_last := strings.LastIndex(text, ".")
		if dec_pos_last != dec_pos {
			// we have multiple decimal points
			if dec_pos == cursor-1 || dec_pos_last == cursor-1 {
				//  we just typed a decimal
				dec_pos = cursor - 1 // WLOG

				// let mut rest = text.split_off(dec_pos);
				rest := strings.ReplaceAll(
					text[dec_pos:], ".", "")
				text = strings.ReplaceAll(
					text[:dec_pos], ".", "")

				var new_text strings.Builder
				new_text.Grow(width)

				if len(text) > self_dec_pos {
					// we need to put the decimal where it belongs
					new_text.WriteString(text[len(text)-self_dec_pos:])
				} else {
					new_text.WriteString(text)
				}
				new_text.WriteRune('.')
				new_text.WriteString(rest)
				text = new_text.String()

				start = self_dec_pos + 1
				end = self_dec_pos + 1
			} else {
				// remove them all
				text = strings.ReplaceAll(text, ".", "")
				remangle = true
			}
		} else {
			// only one decimal point
			if len(text) > width {
				// text inserted, not overwritten

				// move the decimal out of my way
				if cursor > self_dec_pos && cursor-self_dec_pos <= len(text)-width {
					start += 1
					end += 1
				}

				if self_dec_pos < dec_pos {
					// add to beginning
					var new_text strings.Builder
					new_text.Grow(width)
					new_text.WriteString(text[:self_dec_pos])
					new_text.WriteRune('.')
					new_text.WriteString(text[self_dec_pos:dec_pos])
					new_text.WriteString(text[dec_pos+1:])
					text = new_text.String()
				} else if cursor == len(text) {
					// add to end
					var start_idx uint32
					if self_dec_pos == dec_pos {
						start_idx = 1
					} else {
						start_idx = 0
					} // truncate from front if buffer is full
					var new_text strings.Builder
					new_text.Grow(width)
					new_text.WriteString(text[start_idx:dec_pos])
					new_text.WriteString(text[dec_pos+1 : dec_pos+2])
					new_text.WriteRune('.')
					new_text.WriteString(text[dec_pos+2:])
					text = new_text.String()
				}
			} else if len(text) < width && dec_pos < self_dec_pos {
				if len(text) > self_dec_pos {
					// we need to put the decimal where it belongs
					var new_text strings.Builder
					new_text.Grow(width)
					new_text.WriteString(text[:dec_pos])
					new_text.WriteString(text[dec_pos+1 : self_dec_pos+1])
					new_text.WriteRune('.')
					new_text.WriteString(text[self_dec_pos+1:])
					text = new_text.String()
				}
			} else {
				// pasted
				// text.remove(dec_pos)
				var new_text strings.Builder
				new_text.Grow(width)
				new_text.WriteString(text[:dec_pos])
				new_text.WriteString(text[dec_pos+1:])
				text = new_text.String()
				remangle = true
			}
		}

	} else if len(text) > self_dec_pos {
		// no decimal
		// but enough space

		var new_text strings.Builder
		new_text.Grow(width)
		// move the decimal out of my way
		// if cursor != self_dec_pos
		//     || match last_key.get() {
		//         ASCII_DELETE => {
		//             new_text.WriteString(text[:self_dec_pos]);
		//             new_text.WriteRune('.');
		//             new_text.WriteString(text[self_dec_pos + 1:]);
		//             false
		//         }
		//         ASCII_BACKSPACE => {
		//             new_text.WriteString(text[:self_dec_pos - 1]);
		//             new_text.WriteString(text[self_dec_pos:self_dec_pos + 1]);
		//             new_text.WriteRune('.');
		//             new_text.WriteString(text[self_dec_pos + 1:]);
		//             start -= 1;
		//             end -= 1;
		//             false
		//         }
		//         _ => true,
		//     }
		// {
		new_text.WriteString(text[:self_dec_pos])
		new_text.WriteRune('.')
		new_text.WriteString(text[self_dec_pos:])
		// }
		text = new_text.String()
	} else {
		remangle = true
	}

	if remangle {
		if len(text) > self_dec_pos {
			// no decimal
			// but enough space for one
			var new_text strings.Builder
			new_text.Grow(width)

			new_text.WriteString(text[:self_dec_pos])
			new_text.WriteRune('.')
			new_text.WriteString(text[self_dec_pos:])
			text = new_text.String()
		} else {
			// text is too small for a decimal
			// text.reserve_exact(self_dec_pos - len(text));
			var new_text strings.Builder
			new_text.Grow(self_dec_pos)
			new_text.WriteString(text)

			// text.push_str(&"0".repeat(self_dec_pos - len(text)));
			new_text.WriteString(
				strings.Repeat("0", self_dec_pos-len(text)))
			text = new_text.String()
		}
	}

	text = text[0:min(width, len(text))]

	val, _ := strconv.ParseFloat(strings.TrimSpace(text), 64)

	control.SetText(format(val))
	control.SelectText(start, end)

	return val
}

func (control *NumberEditView) GetPointless() float64 {
	start, end := control.Selected()

	val, _ := strconv.ParseFloat(strings.TrimSpace(control.Text()), 64)

	for val >= 100 {
		val /= 10
		// check position. if we just backspaced the decimal point, don't put it back in the way
		if start != 2 {
			start++
			end++
		}
	}
	control.Set(val)
	control.SelectText(start, end)

	return val
}

func (control *NumberEditView) Clear() {
	control.SetText("")
	control.Ok()
}

func (control *NumberEditView) Check(test bool) {
	if test {
		control.Ok()
	} else {
		control.Error()
	}
}

func (control *NumberEditView) SetFont(font *windigo.Font) {
	control.Edit.SetFont(font)
	control.Label().SetFont(font)
}

func (control *NumberEditView) SetLabeledSize(label_width, control_width, height int) {
	control.SetSize(label_width+control_width, height)
	// control.Edit.SetSize(control_width, height)
	control.Label().SetSize(label_width, height)
	control.SetPaddingsAll(ERROR_MARGIN)
}

// onchange for FR
func (this_field *NumberEditView) Entangle(other_field *NumberEditView, range_field Checker, delta_max float64) {
	bind_fn := func(e *windigo.Event) {
		this_val := this_field.GetFixed()
		other_val := other_field.GetFixed()
		diff_val := math.Abs(this_val-other_val) <= delta_max
		checks := range_field.CheckAll(this_val, other_val)
		this_check, other_check := checks[0], checks[1]
		this_field.Check(this_check && diff_val)
		other_field.Check(other_check && diff_val)
	}
	this_field.OnChange().Bind(bind_fn)
	other_field.OnChange().Bind(bind_fn)
}

/*
 * NumberUnitsEditView
 *
 */
type NumberUnitsEditView struct {
	NumberEditView
	SetFont        func(font *windigo.Font)
	SetLabeledSize func(label_width, control_width, unit_width, height int)
}

//	func (control *NumberUnitsEditView) SetLabeledSize(label_width, control_width, height int) {
//		control.NumberEditView.SetLabeledSize(label_width,control_width, height)
//		control.SetSize(label_width+control_width, height)
//		control.Label().SetSize(label_width, height)
//	}
func NumberEditView_with_Units_from_new(parent *windigo.AutoPanel, field_text, field_units string) *NumberUnitsEditView {
	panel := windigo.NewAutoPanel(parent)

	text_label := windigo.NewLabel(panel)
	text_label.SetText(field_text)

	text_field := windigo.NewEdit(panel)
	text_field.SetText("0.000")

	text_units := windigo.NewLabel(panel)
	text_units.SetText(field_units)

	panel.Dock(text_label, windigo.Left)
	panel.Dock(text_field, windigo.Left)
	panel.Dock(text_units, windigo.Left)

	setFont := func(font *windigo.Font) {
		text_label.SetFont(font)
		text_field.SetFont(font)
		text_units.SetFont(font)
	}
	setLabeledSize := func(label_width, control_width, unit_width, height int) {
		panel.SetSize(label_width+control_width+unit_width, height)
		panel.SetPaddingsAll(ERROR_MARGIN)

		text_label.SetSize(label_width, height)
		text_label.SetMarginTop(ERROR_MARGIN)

		text_field.SetSize(control_width, height)
		text_field.SetMarginRight(DATA_MARGIN)

		text_units.SetSize(unit_width, height)
		text_units.SetMarginTop(ERROR_MARGIN)
	}

	return &NumberUnitsEditView{NumberEditView{&View{ComponentFrame: panel}, NumbEditView{text_field}, windigo.Labeled{FieldLabel: text_label}}, setFont, setLabeledSize}
}
