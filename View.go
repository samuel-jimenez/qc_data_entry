package main

import "github.com/samuel-jimenez/windigo"

/*
 * ErrableView
 *
 */
type ErrableView interface {
	windigo.ComponentFrame
	Ok()
	Error()
}

/*
 * View
 *
 */
type View struct {
	windigo.ComponentFrame
}

func (view *View) Ok() {
	view.SetBorder(nil)
}

func (view *View) Error() {
	view.SetBorder(erroredPen)
}

/*
 * UpdatableView
 *
 */
type UpdatableView[T any] interface {
	ErrableView
	Update(T)
}
