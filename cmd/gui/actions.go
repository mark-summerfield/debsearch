// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	"fmt"

	"github.com/pwiecz/go-fltk"
)

func (me *App) onEvent(event fltk.Event) bool {
	key := fltk.EventKey()
	switch fltk.EventType() {
	case fltk.SHOW:
		me.mainVBox.Layout()
	case fltk.SHORTCUT:
		if key == fltk.ESCAPE {
			return true // ignore
		}
	case fltk.KEY:
		switch key {
		case fltk.HELP:
			me.onHelp()
			return true
		}
	case fltk.CLOSE:
		me.onQuit()
	}
	return false
}

func (me *App) onError(err error) {
	me.statusBar.SetLabelColor(fltk.RED)
	me.statusBar.SetLabel(err.Error())
}

func (me *App) onInfo(info string, autoClear bool) {
	me.statusBar.SetLabelColor(fltk.BLUE)
	me.statusBar.SetLabel(info)
	if autoClear {
		fltk.AddTimeout(7, func() { me.clearStatus() })
	}
}

func (me *App) clearStatus() {
	me.statusBar.SetLabelColor(fltk.BLACK)
	me.statusBar.SetLabel("")
	me.Redraw()
}

func (me *App) onFind() {
	fmt.Println("onFind")
}

func (me *App) onConfigure() {
	form := newConfigForm(me)
	form.SetModal()
	form.Show()
}

func (me *App) onAbout() {
	fmt.Println("onAbout")
}

func (me *App) onHelp() {
	fmt.Println("onHelp")
}

func (me *App) onQuit() {
	me.config.X = me.Window.X()
	me.config.Y = me.Window.Y()
	me.config.Width = me.Window.W()
	me.config.Height = me.Window.H()
	me.config.Scale = fltk.ScreenScale(0)
	me.config.save()
	me.Window.Destroy()
}
