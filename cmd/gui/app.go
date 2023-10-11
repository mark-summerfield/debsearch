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
	statusBar                *fltk.Box
	sectionsLabel            *fltk.Box
	sectionsBrowser          *fltk.MultiBrowser
	incNonFreeCheckbox       *fltk.CheckButton
	tagsLabel                *fltk.Box
	tagsBrowser              *fltk.MultiBrowser
	tagsMatchAllRadioButton  *fltk.RadioRoundButton
	tagsMatchAnyRadioButton  *fltk.RadioRoundButton
	wordsInput               *fltk.Input
	wordsMatchAllRadioButton *fltk.RadioRoundButton
	wordsMatchAnyRadioButton *fltk.RadioRoundButton
	packagesLabel            *fltk.Box
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
	pairs := ds.StdFilePairsWithDescriptions(me.config.Arch)
	if pkgs, err := ds.NewPkgs(pairs...); err != nil {
		me.onError(err)
	} else {
		me.pkgs = &pkgs
		me.onInfo(fmt.Sprintf("Read %s packages.\n",
			gong.Commas(len(pkgs.Pkgs))), false)
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
	me.sectionsLabel.SetLabel(fmt.Sprintf("Sections (0/%s)",
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
	me.tagsLabel.SetLabel(fmt.Sprintf("Tags (0/%s)", gong.Commas(count)))
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
	tileHeight := height - ((2 * buttonHeight) + (3 * gui.Margin))
	tile := fltk.NewTile(0, 0, width, tileHeight)
	leftWidth := (width / 10) * 4
	me.makeCriteriaPanel(0, 0, leftWidth, tileHeight)
	me.makeResultPanel(leftWidth, 0, width-leftWidth, tileHeight)
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
	hbox.SetBox(fltk.UP_FRAME)
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
	me.sectionsLabel = fltk.NewBox(fltk.FLAT_BOX, x, 0, labelWidth,
		buttonHeight, "Sections")
	vbox.Fixed(me.sectionsLabel, gui.LabelHeight())
	me.sectionsBrowser = fltk.NewMultiBrowser(0, buttonHeight, width,
		height)
	me.sectionsBrowser.SetCallback(func() {
		me.updateSectionsLabel(selectedCount(me.sectionsBrowser))
	})
	hbox := gui.MakeHBox(x, height-(2*buttonHeight), width, buttonHeight)
	selectAllSectionsButton := fltk.NewButton(0, 0, labelWidth,
		buttonHeight, "Select All")
	selectAllSectionsButton.SetCallback(func() {
		selectOrClear(me.sectionsBrowser, true)
		me.updateSectionsLabel(me.sectionsBrowser.Size())
	})
	hbox.Fixed(selectAllSectionsButton, labelWidth)
	clearSectionsButton := fltk.NewButton(labelWidth, 0, labelWidth,
		buttonHeight, "Clear All")
	clearSectionsButton.SetCallback(func() {
		selectOrClear(me.sectionsBrowser, false)
		me.updateSectionsLabel(0)
	})
	hbox.Fixed(clearSectionsButton, labelWidth)
	padBox(hbox, gui.Margin)
	me.incNonFreeCheckbox = fltk.NewCheckButton(labelWidth, 0, labelWidth,
		buttonHeight, "&Include Non-Free")
	me.incNonFreeCheckbox.SetValue(me.config.IncludeNonFreeSections)
	hbox.End()
	vbox.Fixed(hbox, buttonHeight)
	vbox.End()
}

func (me *App) updateSectionsLabel(count int) {
	me.sectionsLabel.SetLabel(fmt.Sprintf("Sections (%s/%s)",
		gong.Commas(count), gong.Commas(me.sectionsBrowser.Size())))
}

func (me *App) makeTagsPanel(x, y, width, height int) {
	buttonHeight := gui.ButtonHeight()
	labelWidth := (gui.LabelWidth() * 3) / 2
	vbox := gui.MakeVBox(x, y, width, height)
	divider(vbox)
	me.tagsLabel = fltk.NewBox(fltk.FLAT_BOX, x, 0, labelWidth,
		buttonHeight, "Tags")
	vbox.Fixed(me.tagsLabel, gui.LabelHeight())
	me.tagsBrowser = fltk.NewMultiBrowser(0, buttonHeight, width,
		height)
	me.tagsBrowser.SetCallback(func() {
		me.updateTagsLabel(selectedCount(me.tagsBrowser))
	})
	hbox := gui.MakeHBox(x, height-(2*buttonHeight), width, buttonHeight)
	selectAllTagsButton := fltk.NewButton(0, 0, labelWidth,
		buttonHeight, "Select All")
	selectAllTagsButton.SetCallback(func() {
		selectOrClear(me.tagsBrowser, true)
		me.updateTagsLabel(me.tagsBrowser.Size())
	})
	hbox.Fixed(selectAllTagsButton, labelWidth)
	clearTagsButton := fltk.NewButton(labelWidth, 0, labelWidth,
		buttonHeight, "Clear All")
	clearTagsButton.SetCallback(func() {
		selectOrClear(me.tagsBrowser, false)
		me.updateTagsLabel(0)
	})
	hbox.Fixed(clearTagsButton, labelWidth)
	padBox(hbox, gui.Margin)
	fltk.NewBox(fltk.FLAT_BOX, x, 0, labelWidth, buttonHeight, "Match:")
	me.tagsMatchAllRadioButton = fltk.NewRadioRoundButton(x, 0, labelWidth,
		buttonHeight, "All")
	me.tagsMatchAllRadioButton.SetValue(true)
	me.tagsMatchAnyRadioButton = fltk.NewRadioRoundButton(x, 0, labelWidth,
		buttonHeight, "Any")
	hbox.End()
	vbox.Fixed(hbox, buttonHeight)
	vbox.End()
}

func (me *App) updateTagsLabel(count int) {
	me.tagsLabel.SetLabel(fmt.Sprintf("Tagss (%s/%s)",
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
	hbox = gui.MakeHBox(x, y, width, buttonHeight)
	label := fltk.NewBox(fltk.FLAT_BOX, x, 0, labelWidth, buttonHeight,
		"Match:")
	hbox.Fixed(label, labelWidth)
	me.wordsMatchAllRadioButton = fltk.NewRadioRoundButton(x, 0, labelWidth,
		buttonHeight, "All")
	me.wordsMatchAllRadioButton.SetValue(true)
	hbox.Fixed(me.wordsMatchAllRadioButton, labelWidth)
	me.wordsMatchAnyRadioButton = fltk.NewRadioRoundButton(x, 0, labelWidth,
		buttonHeight, "Any")
	hbox.End()
	vbox.End()
	vbox.Fixed(vbox, 2*buttonHeight)
}

func (me *App) makeResultPanel(x, y, width, height int) {
	labelHeight := gui.LabelHeight()
	tile := fltk.NewTile(x, y, width, height)
	height /= 2
	vbox := gui.MakeVBox(x, y, width, height)
	me.packagesLabel = fltk.NewBox(fltk.FLAT_BOX, 0, 0, width, labelHeight,
		"Packages Found")
	vbox.Fixed(me.packagesLabel, labelHeight)
	me.packagesBrowser = fltk.NewHoldBrowser(0, labelHeight, width,
		height-labelHeight)
	vbox.End()
	y += height
	vbox = gui.MakeVBox(x, y, width, height)
	divider(vbox)
	label := fltk.NewBox(fltk.FLAT_BOX, 0, 0, width, labelHeight,
		"Description")
	vbox.Fixed(label, labelHeight)
	me.descView = fltk.NewHelpView(0, labelHeight, width,
		height-labelHeight)
	me.descView.TextFont(fltk.HELVETICA)
	me.descView.TextSize(me.config.TextSize)
	me.descView.SetValue(initialDescHtml)
	vbox.End()
	tile.End()
}

func (me *App) updatePackagesLabel(count int) {
	me.packagesLabel.SetLabel(fmt.Sprintf("Packages Found (%s)",
		gong.Commas(count)))
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
