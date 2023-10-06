// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import (
	"fmt"
	"time"

	ds "github.com/mark-summerfield/debsearch"
	"github.com/mark-summerfield/debsearch/cmd/gui/gui"
	"github.com/mark-summerfield/gong"
	"github.com/pwiecz/go-fltk"
)

type App struct {
	*fltk.Window
	config    *Config
	pkgs      *ds.Pkgs
	mainVBox  *fltk.Flex
	statusBar *fltk.Box
}

func newApp(config *Config, args []string) *App {
	app := &App{Window: nil, config: config}
	app.makeMainWindow()
	app.makeWidgets()
	app.Window.End()
	fltk.AddTimeout(tinyTimeout, app.loadPackages)
	return app
}

func (me *App) loadPackages() {
	t := time.Now()
	pairs := ds.StdFilePairsWithDescriptions()
	if pkgs, err := ds.NewPkgs(pairs...); err != nil {
		me.onError(err)
	} else {
		me.pkgs = &pkgs
		me.onInfo(fmt.Sprintf("Read %s packages in %s.\n",
			gong.Commas(len(pkgs.Pkgs)), time.Since(t)), false)
	}
}

func (me *App) makeMainWindow() {
	me.Window = fltk.NewWindow(me.config.Width, me.config.Height)
	if me.config.X > -1 && me.config.Y > -1 {
		me.Window.SetPosition(me.config.X, me.config.Y)
	}
	me.Window.Resizable(me.Window)
	me.Window.SetEventHandler(me.onEvent)
	me.Window.SetLabel(appName)
	gui.AddWindowIcon(me.Window, iconSvg)
}

func (me *App) makeWidgets() {
	width := me.Window.W()
	height := me.Window.H()
	buttonHeight := gui.ButtonHeight()
	var x, y int
	vbox := gui.MakeVBox(x, y, width, height, gui.Pad)
	hbox := me.makeButtonPanel(width, 0)
	vbox.Fixed(hbox, buttonHeight)
	// TODO
	me.makeStatusBar(width, height)
	vbox.Fixed(me.statusBar, buttonHeight)
	vbox.End()
	me.mainVBox = vbox
}

func (me *App) makeButtonPanel(width, y int) *fltk.Flex {
	buttonHeight := gui.ButtonHeight()
	labelWidth := (gui.LabelWidth() * 3) / 2
	x := 0
	hbox := gui.MakeHBox(x, y, width, buttonHeight, gui.Pad)
	findButton := fltk.NewButton(x, 0, labelWidth, buttonHeight,
		"&Find")
	findButton.SetCallback(me.onFind)
	hbox.Fixed(findButton, labelWidth)
	x += labelWidth
	fltk.NewBox(fltk.FLAT_BOX, x, 0, labelWidth, buttonHeight)
	x += labelWidth
	configButton := fltk.NewButton(x, 0, labelWidth, buttonHeight,
		"&Options…")
	configButton.SetCallback(me.onConfigure)
	hbox.Fixed(configButton, labelWidth)
	x += labelWidth
	aboutButton := fltk.NewButton(x, 0, labelWidth, buttonHeight,
		"&About")
	aboutButton.SetCallback(me.onAbout)
	hbox.Fixed(aboutButton, labelWidth)
	x += labelWidth
	helpButton := fltk.NewButton(x, 0, labelWidth, buttonHeight,
		"&Help")
	helpButton.SetCallback(me.onHelp)
	hbox.Fixed(helpButton, labelWidth)
	x += labelWidth
	quitButton := fltk.NewButton(x, 0, labelWidth, buttonHeight,
		"&Quit")
	quitButton.SetCallback(me.onQuit)
	hbox.Fixed(quitButton, labelWidth)
	hbox.End()
	return hbox
}

func (me *App) makeStatusBar(width, y int) {
	buttonHeight := gui.ButtonHeight()
	me.statusBar = fltk.NewBox(fltk.DOWN_FRAME, 0, y-buttonHeight, width,
		buttonHeight)
	me.statusBar.SetAlign(fltk.ALIGN_LEFT | fltk.ALIGN_INSIDE)
}
