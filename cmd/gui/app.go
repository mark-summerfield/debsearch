// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"strings"

	ds "github.com/mark-summerfield/debsearch"
	"github.com/mark-summerfield/debsearch/cmd/gui/gui"
	"github.com/mark-summerfield/gong"
	"github.com/pwiecz/go-fltk"
)

type App struct {
	*fltk.Window
	config             *Config
	pkgs               *ds.Pkgs
	mainVBox           *fltk.Flex
	statusBar          *fltk.Box
	sectionsLabel      *fltk.Box
	sectionsBrowser    *fltk.MultiBrowser
	incNonFreeCheckbox *fltk.CheckButton
}

func newApp(config *Config) *App {
	app := &App{Window: nil, config: config}
	app.makeMainWindow()
	app.makeWidgets()
	app.Window.End()
	fltk.AddTimeout(tinyTimeout, app.loadPackages)
	return app
}

func (me *App) loadPackages() {
	pairs := ds.StdFilePairsWithDescriptions()
	if pkgs, err := ds.NewPkgs(pairs...); err != nil {
		me.onError(err)
	} else {
		me.pkgs = &pkgs
		me.onInfo(fmt.Sprintf("Read %s packages.\n",
			gong.Commas(len(pkgs.Pkgs))), false)
		me.populateSections()
		// TODO populate Tags widget
	}
}

func (me *App) populateSections() {
	me.sectionsBrowser.Clear()
	for _, section := range gong.SortedMapKeys(me.pkgs.SectionsAndCounts) {
		if !strings.HasPrefix(section, nonfreePrefix) {
			me.sectionsBrowser.Add(fmt.Sprintf("%s (%s)", section,
				gong.Commas(me.pkgs.SectionsAndCounts[section])))
		}
	}
	me.sectionsLabel.SetLabel(fmt.Sprintf("Sections (%s)",
		gong.Commas(len(me.pkgs.SectionsAndCounts))))
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

func (me *App) makeWidgets() { // TODO set non-free checkbox from config
	width := me.Window.W()
	height := me.Window.H()
	buttonHeight := gui.ButtonHeight()
	vbox := gui.MakeVBox(0, 0, width, height)
	hbox := me.makeButtonPanel(width, 0)
	vbox.Fixed(hbox, buttonHeight)
	tileHeight := height - (2 * buttonHeight)
	tile := fltk.NewTile(0, 0, width, tileHeight)
	halfWidth := width / 2
	me.makeCriteriaPanel(0, 0, halfWidth, tileHeight)
	me.makeResultPanel(halfWidth, 0, halfWidth, tileHeight)
	tile.End()
	hbox = me.makeStatusBar(width, height)
	vbox.Fixed(hbox, buttonHeight)
	vbox.End()
	me.mainVBox = vbox
}

func (me *App) makeButtonPanel(width, y int) *fltk.Flex {
	buttonHeight := gui.ButtonHeight()
	labelWidth := (gui.LabelWidth() * 3) / 2
	x := 0
	hbox := gui.MakeHBox(x, y, width, buttonHeight)
	pad := fltk.NewBox(fltk.FLAT_BOX, x, 0, 1, buttonHeight) // left pad
	hbox.Fixed(pad, 1)
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
	pad = fltk.NewBox(fltk.FLAT_BOX, x, 0, 1, buttonHeight) // right pad
	hbox.Fixed(pad, 1)
	hbox.End()
	return hbox
}

func (me *App) makeCriteriaPanel(x, y, width, height int) {
	buttonHeight := gui.ButtonHeight()
	tile := fltk.NewTile(x, y, width, height)
	height /= 3
	y = 0
	me.makeSectionsPanel(x, y, width, height)
	y += height
	vbox := gui.MakeVBox(x, y, width, height)
	ifDebug(me.config.debug, vbox, fltk.GREEN)
	fltk.NewBox(fltk.FLAT_BOX, 0, 0, width, buttonHeight, "Tags")
	//TODO
	// - Tree of checkable tags
	// - [Select All] [Unselect All] Match (*) All ( ) Any
	vbox.End()
	y += height
	vbox = gui.MakeVBox(x, y, width, height)
	ifDebug(me.config.debug, vbox, fltk.BLUE)
	fltk.NewBox(fltk.FLAT_BOX, 0, 0, width, buttonHeight, "Words")
	//TODO
	// - Line edit for words
	// - Match (*) All ( ) Any
	vbox.End()
	tile.End()
}

func (me *App) makeSectionsPanel(x, y, width, height int) {
	buttonHeight := gui.ButtonHeight()
	labelWidth := (gui.LabelWidth() * 3) / 2
	vbox := gui.MakeVBox(x, y, width, height)
	vbox.SetBox(fltk.ENGRAVED_BOX)
	me.sectionsLabel = fltk.NewBox(fltk.FLAT_BOX, x, 0, labelWidth,
		buttonHeight, "Sections")
	vbox.Fixed(me.sectionsLabel, gui.LabelHeight())
	me.sectionsBrowser = fltk.NewMultiBrowser(0, buttonHeight, width,
		height)
	hbox := gui.MakeHBox(x, height-(2*buttonHeight), width, buttonHeight)
	selectAllSectionsButton := fltk.NewButton(0, 0, labelWidth,
		buttonHeight, "Select All")
	selectAllSectionsButton.SetCallback(func() {
		selectOrClear(me.sectionsBrowser, true)
	})
	hbox.Fixed(selectAllSectionsButton, labelWidth)
	clearSectionsButton := fltk.NewButton(labelWidth, 0, labelWidth,
		buttonHeight, "Clear All")
	clearSectionsButton.SetCallback(func() {
		selectOrClear(me.sectionsBrowser, false)
	})
	hbox.Fixed(clearSectionsButton, labelWidth)
	me.incNonFreeCheckbox = fltk.NewCheckButton(labelWidth*2, 0, labelWidth,
		buttonHeight, "&Include Non-Free")
	me.incNonFreeCheckbox.SetValue(me.config.IncludeNonFreeSections)
	vbox.Fixed(hbox, buttonHeight)
	vbox.End()
}

func (me *App) makeResultPanel(x, y, width, height int) {
	buttonHeight := gui.ButtonHeight()
	//tile := gui.MakeVBox(x, y, width, height)
	tile := fltk.NewTile(x, y, width, height)
	height /= 2
	y = 0
	vbox := gui.MakeVBox(x, y, width, height)
	ifDebug(me.config.debug, vbox, fltk.MAGENTA)
	fltk.NewBox(fltk.FLAT_BOX, 0, 0, width, buttonHeight,
		"Matching Packages")
	//TODO list of packages (name, version, size, short desc)
	vbox.End()
	y += height
	vbox = gui.MakeVBox(x, y, width, height)
	ifDebug(me.config.debug, vbox, fltk.CYAN)
	fltk.NewBox(fltk.FLAT_BOX, 0, 0, width, buttonHeight, "Description")
	//TODO the currently selected package's long desc
	vbox.End()
	tile.End()
}

func (me *App) makeStatusBar(width, y int) *fltk.Flex {
	buttonHeight := gui.ButtonHeight()
	hbox := gui.MakeHBox(0, 0, width, buttonHeight)
	pad := fltk.NewBox(fltk.FLAT_BOX, 0, 0, 1, buttonHeight)
	hbox.Fixed(pad, 1)
	me.statusBar = fltk.NewBox(fltk.DOWN_FRAME, 0, y-buttonHeight, width,
		buttonHeight)
	me.statusBar.SetAlign(fltk.ALIGN_LEFT | fltk.ALIGN_INSIDE)
	pad = fltk.NewBox(fltk.FLAT_BOX, width-2, 0, 1, buttonHeight)
	hbox.Fixed(pad, 1)
	return hbox
}
