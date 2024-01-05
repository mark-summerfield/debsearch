// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"strings"

	ds "github.com/mark-summerfield/debsearch"
	"github.com/mark-summerfield/debsearch/cmd/DebFind/gui"
	"github.com/mark-summerfield/gong"
	"github.com/pwiecz/go-fltk"
)

type App struct {
	*fltk.Window
	config                   *Config
	model                    *ds.Model
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
	if model, err := ds.NewModel(pairs...); err != nil {
		me.onError(err)
	} else {
		me.model = &model
		me.onHtmlMessage(fmt.Sprintf(loadTemplate,
			gong.Commas(len(model.Debs))))
		me.populateSections()
		me.populateTags()
	}
}

func (me *App) populateSections() {
	me.sectionsBrowser.Clear()
	count := 0
	for _, section := range gong.SortedMapKeys(me.model.SectionsAndCounts) {
		if !strings.HasPrefix(section, nonfreePrefix) &&
			!strings.HasSuffix(section, todoSuffix) {
			me.sectionsBrowser.Add(fmt.Sprintf("%s (%s)", section,
				gong.Commas(me.model.SectionsAndCounts[section])))
			count++
		}
	}
	me.updateSectionsLabel(0)
}

func (me *App) populateTags() {
	me.tagsBrowser.Clear()
	count := 0
	for _, tag := range gong.SortedMapKeys(me.model.TagsAndCounts) {
		if !strings.HasSuffix(tag, todoSuffix) {
			me.tagsBrowser.Add(fmt.Sprintf("%s (%s)", tag,
				gong.Commas(me.model.TagsAndCounts[tag])))
			count++
		}
	}
	me.updateTagsLabel(0)
}

func (me *App) makeMainWindow() {
	me.Window = fltk.NewWindow(me.config.Width, me.config.Height)
	if me.config.X > -1 && me.config.Y > -1 {
		me.Window.SetPosition(me.config.X, me.config.Y)
	}
	me.Window.Resizable(me.Window)
	me.Window.SetEventHandler(me.onEvent)
	me.Window.SetLabel(appName)
	me.Window.SetXClass(appName)
	gui.AddWindowIcon(me.Window, iconSvg)
}

func (me *App) makeWidgets() {
	width := me.Window.W()
	height := me.Window.H()
	vbox := gui.MakeVBox(0, 0, width, height)
	hbox := me.makeButtonPanel(width, 0)
	vbox.Fixed(hbox, gui.ButtonHeight+(2*gui.Margin))
	tileHeight := height - ((gui.ButtonHeight * 3) / 2)
	tile := fltk.NewTile(0, 0, width, tileHeight)
	leftWidth := min(360, width/2)
	me.makeCriteriaPanel(0, 0, leftWidth, tileHeight)
	me.makeResultPanel(leftWidth, 0, width-leftWidth, tileHeight)
	tile.End()
	vbox.End()
	me.mainVBox = vbox
}

func (me *App) makeButtonPanel(width, y int) *fltk.Flex {
	buttonWidth := ((gui.LabelWidth * 5) / 3) + (2 * gui.Pad)
	x := 0
	hbox := gui.MakeHBox(x, y, width, gui.ButtonHeight)
	hbox.SetBox(fltk.UP_FRAME)
	findButton := makeButton(x, " &Find", iconSvg, me.onFind)
	hbox.Fixed(findButton, buttonWidth)
	x += buttonWidth
	fltk.NewBox(fltk.FLAT_BOX, x, 0, buttonWidth, gui.ButtonHeight)
	x += buttonWidth
	configButton := makeButton(x, " &Options…", configSvg, me.onConfigure)
	hbox.Fixed(configButton, buttonWidth)
	x += buttonWidth
	aboutButton := makeButton(x, " A&bout", aboutSvg, me.onAbout)
	hbox.Fixed(aboutButton, buttonWidth)
	x += buttonWidth
	helpButton := makeButton(x, " &Help", helpSvg, me.onHelp)
	hbox.Fixed(helpButton, buttonWidth)
	x += buttonWidth
	quitButton := makeButton(x, " &Quit", quitSvg, me.onQuit)
	hbox.Fixed(quitButton, buttonWidth)
	hbox.End()
	return hbox
}

func (me *App) makeCriteriaPanel(x, y, width, height int) {
	tile := fltk.NewTile(x, y, width, height)
	height = (height / 2) - gui.ButtonHeight
	y = 0
	me.makeSectionsPanel(x, y, width, height)
	y += height
	me.makeTagsPanel(x, y, width, height)
	y += height
	me.makeWordsPanel(x, y, width, height)
	tile.End()
}

func (me *App) makeSectionsPanel(x, y, width, height int) {
	buttonWidth := (gui.LabelWidth * 5) / 4
	labelWidth := (gui.LabelWidth * 3) / 2
	vbox := gui.MakeVBox(x, y, width, height)
	me.sectionsLabel = gui.MakeAccelLabel(labelWidth, gui.ButtonHeight,
		"&Sections")
	vbox.Fixed(me.sectionsLabel, gui.LabelHeight())
	me.sectionsBrowser = fltk.NewMultiBrowser(0, gui.ButtonHeight, width,
		height)
	me.sectionsBrowser.SetCallback(func() {
		me.updateSectionsLabel(selectedCount(me.sectionsBrowser))
	})
	me.sectionsLabel.SetCallback(func() { me.sectionsBrowser.TakeFocus() })
	hbox := gui.MakeHBox(x, height-(2*gui.ButtonHeight), width,
		gui.ButtonHeight)
	selectAllSectionsButton := fltk.NewButton(0, 0, buttonWidth,
		gui.ButtonHeight, "S&elect All")
	selectAllSectionsButton.SetCallback(func() {
		selectOrClear(me.sectionsBrowser, true)
		me.updateSectionsLabel(me.sectionsBrowser.Size())
	})
	hbox.Fixed(selectAllSectionsButton, buttonWidth)
	clearSectionsButton := fltk.NewButton(labelWidth, 0, buttonWidth,
		gui.ButtonHeight, "&Clear All")
	clearSectionsButton.SetCallback(func() {
		selectOrClear(me.sectionsBrowser, false)
		me.updateSectionsLabel(0)
	})
	hbox.Fixed(clearSectionsButton, buttonWidth)
	padBox(hbox, gui.Margin)
	me.incNonFreeCheckbox = fltk.NewCheckButton(labelWidth, 0, labelWidth,
		gui.ButtonHeight, "Pl&us Non-Free")
	me.incNonFreeCheckbox.SetValue(me.config.IncludeNonFreeSections)
	hbox.End()
	vbox.Fixed(hbox, gui.ButtonHeight)
	vbox.End()
}

func (me *App) updateSectionsLabel(count int) {
	me.sectionsLabel.SetLabel(fmt.Sprintf("&Sections (%s/%s)",
		gong.Commas(count), gong.Commas(me.sectionsBrowser.Size())))
}

func (me *App) makeTagsPanel(x, y, width, height int) {
	buttonWidth := (gui.LabelWidth * 5) / 4
	labelWidth := (gui.LabelWidth * 3) / 2
	vbox := gui.MakeVBox(x, y, width, height)
	divider(vbox)
	me.tagsLabel = gui.MakeAccelLabel(labelWidth, gui.ButtonHeight, "&Tags")
	vbox.Fixed(me.tagsLabel, gui.LabelHeight())
	me.tagsBrowser = fltk.NewMultiBrowser(0, gui.ButtonHeight, width,
		height)
	me.tagsBrowser.SetCallback(func() {
		me.updateTagsLabel(selectedCount(me.tagsBrowser))
	})
	me.tagsLabel.SetCallback(func() { me.tagsBrowser.TakeFocus() })
	hbox := gui.MakeHBox(x, height-(2*gui.ButtonHeight), width,
		gui.ButtonHeight)
	selectAllTagsButton := fltk.NewButton(0, 0, buttonWidth,
		gui.ButtonHeight, "Se&lect All")
	selectAllTagsButton.SetCallback(func() {
		selectOrClear(me.tagsBrowser, true)
		me.updateTagsLabel(me.tagsBrowser.Size())
	})
	hbox.Fixed(selectAllTagsButton, buttonWidth)
	clearTagsButton := fltk.NewButton(labelWidth, 0, buttonWidth,
		gui.ButtonHeight, "Clea&r All")
	clearTagsButton.SetCallback(func() {
		selectOrClear(me.tagsBrowser, false)
		me.updateTagsLabel(0)
	})
	hbox.Fixed(clearTagsButton, buttonWidth)
	padBox(hbox, gui.Margin)
	fltk.NewBox(fltk.FLAT_BOX, x, 0, labelWidth, gui.ButtonHeight, "Match:")
	me.tagsMatchAllRadioButton = fltk.NewRadioRoundButton(x, 0, labelWidth,
		gui.ButtonHeight, "&All")
	me.tagsMatchAllRadioButton.SetValue(me.config.AllTags)
	me.tagsMatchAnyRadioButton = fltk.NewRadioRoundButton(x, 0, labelWidth,
		gui.ButtonHeight, "A&ny")
	me.tagsMatchAnyRadioButton.SetValue(!me.config.AllTags)
	hbox.End()
	vbox.Fixed(hbox, gui.ButtonHeight)
	vbox.End()
}

func (me *App) updateTagsLabel(count int) {
	me.tagsLabel.SetLabel(fmt.Sprintf("&Tags (%s/%s)",
		gong.Commas(count), gong.Commas(me.tagsBrowser.Size())))
}

func (me *App) makeWordsPanel(x, y, width, height int) {
	vbox := gui.MakeVBox(x, y, width, height)
	divider(vbox)
	hbox := gui.MakeHBox(x, y, width, gui.ButtonHeight)
	wordsLabel := gui.MakeAccelLabel(width, gui.ButtonHeight, "&Words:")
	wordsLabel.SetCallback(func() { me.wordsInput.TakeFocus() })
	hbox.Fixed(wordsLabel, gui.LabelWidth)
	me.wordsInput = fltk.NewInput(x, y, width, gui.ButtonHeight)
	hbox.End()
	vbox.Fixed(hbox, gui.ButtonHeight)
	hbox = gui.MakeHBox(x, y, width, gui.ButtonHeight)
	label := gui.MakeAccelLabel(gui.LabelWidth, gui.ButtonHeight, "&Match:")
	hbox.Fixed(label, gui.LabelWidth)
	me.wordsMatchAllRadioButton = fltk.NewRadioRoundButton(x, 0,
		gui.LabelWidth, gui.ButtonHeight, "All")
	me.wordsMatchAllRadioButton.SetValue(me.config.AllWords)
	label.SetCallback(func() { me.wordsMatchAllRadioButton.TakeFocus() })
	hbox.Fixed(me.wordsMatchAllRadioButton, gui.LabelWidth)
	me.wordsMatchAnyRadioButton = fltk.NewRadioRoundButton(x, 0,
		gui.LabelWidth, gui.ButtonHeight, "An&y")
	me.wordsMatchAnyRadioButton.SetValue(!me.config.AllWords)
	hbox.End()
	vbox.Fixed(hbox, gui.ButtonHeight)
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
