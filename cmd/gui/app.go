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
	config                   *Config
	pkgs                     *ds.Pkgs
	mainVBox                 *fltk.Flex
	sectionsLabel            *fltk.Button
	sectionsBrowser          *fltk.MultiBrowser
	incNonFreeCheckbox       *fltk.CheckButton
	tagsLabel                *fltk.Button
	tagsBrowser              *fltk.MultiBrowser
	tagsMatchAllRadioButton  *fltk.RadioRoundButton
	tagsMatchAnyRadioButton  *fltk.RadioRoundButton
	wordsInput               *fltk.Input
	wordsMatchAllRadioButton *fltk.RadioRoundButton
	wordsMatchAnyRadioButton *fltk.RadioRoundButton
	packagesLabel            *fltk.Button
	packagesBrowser          *fltk.HoldBrowser
	descView                 *fltk.HelpView
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
	pairs := ds.StdFilePairsWithDescriptions(me.config.Arc)
	if pkgs, err := ds.NewPkgs(pairs...); err != nil {
		me.onError(err)
	} else {
		me.pkgs = &pkgs
		me.onInfo(fmt.Sprintf("Read %s packages.\n",
			gong.Commas(len(pkgs.Pkgs))))
		me.populateSections()
		me.populateTags()
	}
}

func (me *App) populateSections() {
	me.sectionsBrowser.Clear()
	count := 0
	for _, section := range gong.SortedMapKeys(me.pkgs.SectionsAndCounts) {
		if !strings.HasPrefix(section, nonfreePrefix) &&
			!strings.HasSuffix(section, todoSuffix) {
			me.sectionsBrowser.Add(fmt.Sprintf("%s (%s)", section,
				gong.Commas(me.pkgs.SectionsAndCounts[section])))
			count++
		}
	}
	me.sectionsLabel.SetLabel(fmt.Sprintf("&Sections (0/%s)",
		gong.Commas(count)))
}

func (me *App) populateTags() {
	me.tagsBrowser.Clear()
	count := 0
	for _, tag := range gong.SortedMapKeys(me.pkgs.TagsAndCounts) {
		if !strings.HasSuffix(tag, todoSuffix) {
			me.tagsBrowser.Add(fmt.Sprintf("%s (%s)", tag,
				gong.Commas(me.pkgs.TagsAndCounts[tag])))
			count++
		}
	}
	me.tagsLabel.SetLabel(fmt.Sprintf("&Tags (0/%s)", gong.Commas(count)))
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
	vbox := gui.MakeVBox(0, 0, width, height)
	hbox := me.makeButtonPanel(width, 0)
	vbox.Fixed(hbox, buttonHeight+(2*gui.Margin))
	tileHeight := height - (2 * buttonHeight)
	tile := fltk.NewTile(0, 0, width, tileHeight)
	leftWidth := (width / 10) * 4
	me.makeCriteriaPanel(0, 0, leftWidth, tileHeight)
	me.makeResultPanel(leftWidth, 0, width-leftWidth, tileHeight)
	tile.End()
	vbox.End()
	me.mainVBox = vbox
}

func (me *App) makeButtonPanel(width, y int) *fltk.Flex {
	buttonHeight := gui.ButtonHeight()
	buttonWidth := (gui.LabelWidth() * 5) / 3
	x := 0
	hbox := gui.MakeHBox(x, y, width, buttonHeight)
	hbox.SetBox(fltk.UP_FRAME)
	pad := fltk.NewBox(fltk.FLAT_BOX, x, 0, 1, buttonHeight) // left pad
	hbox.Fixed(pad, 1)
	findButton := fltk.NewButton(x, 0, buttonWidth, buttonHeight,
		" &Find")
	findButton.SetCallback(me.onFind)
	findButton.SetImage(gui.ImageForSvgText(iconSvg, iconSize))
	findButton.SetAlign(fltk.ALIGN_IMAGE_NEXT_TO_TEXT)
	hbox.Fixed(findButton, buttonWidth)
	x += buttonWidth
	fltk.NewBox(fltk.FLAT_BOX, x, 0, buttonWidth, buttonHeight)
	x += buttonWidth
	configButton := fltk.NewButton(x, 0, buttonWidth, buttonHeight,
		"&Options…")
	configButton.SetCallback(me.onConfigure)
	configButton.SetImage(gui.ImageForSvgText(configSvg, iconSize))
	configButton.SetAlign(fltk.ALIGN_IMAGE_NEXT_TO_TEXT)
	hbox.Fixed(configButton, buttonWidth)
	x += buttonWidth
	aboutButton := fltk.NewButton(x, 0, buttonWidth, buttonHeight,
		"A&bout")
	aboutButton.SetCallback(me.onAbout)
	aboutButton.SetImage(gui.ImageForSvgText(aboutSvg, iconSize))
	aboutButton.SetAlign(fltk.ALIGN_IMAGE_NEXT_TO_TEXT)
	hbox.Fixed(aboutButton, buttonWidth)
	x += buttonWidth
	helpButton := fltk.NewButton(x, 0, buttonWidth, buttonHeight,
		"&Help")
	helpButton.SetCallback(me.onHelp)
	helpButton.SetImage(gui.ImageForSvgText(helpSvg, iconSize))
	helpButton.SetAlign(fltk.ALIGN_IMAGE_NEXT_TO_TEXT)
	hbox.Fixed(helpButton, buttonWidth)
	x += buttonWidth
	quitButton := fltk.NewButton(x, 0, buttonWidth, buttonHeight,
		"&Quit")
	quitButton.SetCallback(me.onQuit)
	quitButton.SetImage(gui.ImageForSvgText(quitSvg, iconSize))
	quitButton.SetAlign(fltk.ALIGN_IMAGE_NEXT_TO_TEXT)
	hbox.Fixed(quitButton, buttonWidth)
	pad = fltk.NewBox(fltk.FLAT_BOX, x, 0, 1, buttonHeight) // right pad
	hbox.Fixed(pad, 1)
	hbox.End()
	return hbox
}

func (me *App) makeCriteriaPanel(x, y, width, height int) {
	buttonHeight := gui.ButtonHeight()
	tile := fltk.NewTile(x, y, width, height)
	height = (height / 2) - buttonHeight
	y = 0
	me.makeSectionsPanel(x, y, width, height)
	y += height
	me.makeTagsPanel(x, y, width, height)
	y += height
	me.makeWordsPanel(x, y, width, height)
	tile.End()
}

func (me *App) makeSectionsPanel(x, y, width, height int) {
	buttonHeight := gui.ButtonHeight()
	labelWidth := (gui.LabelWidth() * 3) / 2
	vbox := gui.MakeVBox(x, y, width, height)
	me.sectionsLabel = gui.MakeAccelLabel(labelWidth, buttonHeight,
		"&Sections")
	vbox.Fixed(me.sectionsLabel, gui.LabelHeight())
	me.sectionsBrowser = fltk.NewMultiBrowser(0, buttonHeight, width,
		height)
	me.sectionsBrowser.SetCallback(func() {
		me.updateSectionsLabel(selectedCount(me.sectionsBrowser))
	})
	me.sectionsLabel.SetCallback(func() { me.sectionsBrowser.TakeFocus() })
	hbox := gui.MakeHBox(x, height-(2*buttonHeight), width, buttonHeight)
	selectAllSectionsButton := fltk.NewButton(0, 0, labelWidth,
		buttonHeight, "S&elect All")
	selectAllSectionsButton.SetCallback(func() {
		selectOrClear(me.sectionsBrowser, true)
		me.updateSectionsLabel(me.sectionsBrowser.Size())
	})
	hbox.Fixed(selectAllSectionsButton, labelWidth)
	clearSectionsButton := fltk.NewButton(labelWidth, 0, labelWidth,
		buttonHeight, "&Clear All")
	clearSectionsButton.SetCallback(func() {
		selectOrClear(me.sectionsBrowser, false)
		me.updateSectionsLabel(0)
	})
	hbox.Fixed(clearSectionsButton, labelWidth)
	padBox(hbox, gui.Margin)
	me.incNonFreeCheckbox = fltk.NewCheckButton(labelWidth, 0, labelWidth,
		buttonHeight, "Incl&ude Non-Free")
	me.incNonFreeCheckbox.SetValue(me.config.IncludeNonFreeSections)
	hbox.End()
	vbox.Fixed(hbox, buttonHeight)
	vbox.End()
}

func (me *App) updateSectionsLabel(count int) {
	me.sectionsLabel.SetLabel(fmt.Sprintf("&Sections (%s/%s)",
		gong.Commas(count), gong.Commas(me.sectionsBrowser.Size())))
}

func (me *App) makeTagsPanel(x, y, width, height int) {
	buttonHeight := gui.ButtonHeight()
	labelWidth := (gui.LabelWidth() * 3) / 2
	vbox := gui.MakeVBox(x, y, width, height)
	divider(vbox)
	me.tagsLabel = gui.MakeAccelLabel(labelWidth, buttonHeight, "&Tags")
	vbox.Fixed(me.tagsLabel, gui.LabelHeight())
	me.tagsBrowser = fltk.NewMultiBrowser(0, buttonHeight, width,
		height)
	me.tagsBrowser.SetCallback(func() {
		me.updateTagsLabel(selectedCount(me.tagsBrowser))
	})
	me.tagsLabel.SetCallback(func() { me.tagsBrowser.TakeFocus() })
	hbox := gui.MakeHBox(x, height-(2*buttonHeight), width, buttonHeight)
	selectAllTagsButton := fltk.NewButton(0, 0, labelWidth,
		buttonHeight, "Se&lect All")
	selectAllTagsButton.SetCallback(func() {
		selectOrClear(me.tagsBrowser, true)
		me.updateTagsLabel(me.tagsBrowser.Size())
	})
	hbox.Fixed(selectAllTagsButton, labelWidth)
	clearTagsButton := fltk.NewButton(labelWidth, 0, labelWidth,
		buttonHeight, "Clea&r All")
	clearTagsButton.SetCallback(func() {
		selectOrClear(me.tagsBrowser, false)
		me.updateTagsLabel(0)
	})
	hbox.Fixed(clearTagsButton, labelWidth)
	padBox(hbox, gui.Margin)
	fltk.NewBox(fltk.FLAT_BOX, x, 0, labelWidth, buttonHeight, "Match:")
	me.tagsMatchAllRadioButton = fltk.NewRadioRoundButton(x, 0, labelWidth,
		buttonHeight, "&All")
	me.tagsMatchAllRadioButton.SetValue(true)
	me.tagsMatchAnyRadioButton = fltk.NewRadioRoundButton(x, 0, labelWidth,
		buttonHeight, "A&ny")
	hbox.End()
	vbox.Fixed(hbox, buttonHeight)
	vbox.End()
}

func (me *App) updateTagsLabel(count int) {
	me.tagsLabel.SetLabel(fmt.Sprintf("&Tagss (%s/%s)",
		gong.Commas(count), gong.Commas(me.tagsBrowser.Size())))
}

func (me *App) makeWordsPanel(x, y, width, height int) {
	buttonHeight := gui.ButtonHeight()
	labelWidth := gui.LabelWidth()
	vbox := gui.MakeVBox(x, y, width, height)
	divider(vbox)
	hbox := gui.MakeHBox(x, y, width, buttonHeight)
	wordsLabel := gui.MakeAccelLabel(width, buttonHeight, "&Words:")
	wordsLabel.SetCallback(func() { me.wordsInput.TakeFocus() })
	hbox.Fixed(wordsLabel, labelWidth)
	me.wordsInput = fltk.NewInput(x, y, width, buttonHeight)
	hbox.End()
	vbox.Fixed(hbox, buttonHeight)
	hbox = gui.MakeHBox(x, y, width, buttonHeight)
	label := gui.MakeAccelLabel(labelWidth, buttonHeight, "&Match:")
	hbox.Fixed(label, labelWidth)
	me.wordsMatchAllRadioButton = fltk.NewRadioRoundButton(x, 0, labelWidth,
		buttonHeight, "All")
	me.wordsMatchAllRadioButton.SetValue(true)
	label.SetCallback(func() { me.wordsMatchAllRadioButton.TakeFocus() })
	hbox.Fixed(me.wordsMatchAllRadioButton, labelWidth)
	me.wordsMatchAnyRadioButton = fltk.NewRadioRoundButton(x, 0, labelWidth,
		buttonHeight, "An&y")
	hbox.End()
	vbox.Fixed(hbox, buttonHeight)
	vbox.End()
}

func (me *App) makeResultPanel(x, y, width, height int) {
	labelHeight := gui.LabelHeight()
	tile := fltk.NewTile(x, y, width, height)
	height /= 2
	vbox := gui.MakeVBox(x, y, width, height)
	me.packagesLabel = gui.MakeAccelLabel(width, labelHeight,
		"&Packages Found")
	vbox.Fixed(me.packagesLabel, labelHeight)
	me.packagesBrowser = fltk.NewHoldBrowser(0, labelHeight, width,
		height-labelHeight)
	me.packagesBrowser.SetCallback(me.onSelectPackage)
	vbox.End()
	me.packagesLabel.SetCallback(func() { me.packagesBrowser.TakeFocus() })
	y += height
	vbox = gui.MakeVBox(x, y, width, height)
	divider(vbox)
	label := gui.MakeAccelLabel(width, labelHeight, "&Information")
	vbox.Fixed(label, labelHeight)
	me.descView = fltk.NewHelpView(0, labelHeight, width,
		height-labelHeight)
	me.descView.TextFont(fltk.HELVETICA)
	me.descView.TextSize(me.config.TextSize)
	label.SetCallback(func() { me.descView.TakeFocus() })
	me.onInfo("Reading packages…")
	vbox.End()
	tile.End()
}

func (me *App) updatePackagesLabel(count int) {
	me.packagesLabel.SetLabel(fmt.Sprintf("&Packages Found (%s)",
		gong.Commas(count)))
}
