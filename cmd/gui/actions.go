// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"strings"

	ds "github.com/mark-summerfield/debsearch"
	"github.com/mark-summerfield/debsearch/cmd/gui/gui"
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
		case fltk.HELP, fltk.F1:
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
	me.packagesBrowser.Clear()
	me.descView.SetValue(searchingHtml)
	query := me.makeQuery()
	pkgs := query.SelectFrom(me.pkgs)
	if len(pkgs) == 0 {
		me.descView.SetValue(noneFoundHtml)
	} else {
		//TODO list of packages (name, size, short desc)
		//TODO select first package
		//TODO the currently selected package's name, version, short & long desc
	}
	fmt.Println("onFind", query, len(pkgs))
}

func (me *App) makeQuery() *ds.Query {
	query := ds.NewQuery()
	sections := selected(me.sectionsBrowser)
	query.Sections.Add(sections...)
	if me.incNonFreeCheckbox.Value() {
		for _, section := range sections {
			if !strings.Contains(section, "/") {
				query.Sections.Add(nonfreePrefix + section)
			}
		}
	}
	query.Tags.Add(selected(me.tagsBrowser)...)
	query.TagsAnd = me.tagsMatchAllRadioButton.Value()
	for _, word := range strings.Fields(me.wordsInput.Value()) {
		query.Words.Add(strings.ToLower(word))
	}
	query.WordsAnd = me.wordsMatchAllRadioButton.Value()
	return query
}

func (me *App) onConfigure() {
	form := newConfigForm(me)
	form.SetModal()
	form.Show()
}

func (me *App) onAbout() {
	descHtml := gui.DescHtml(appName, ds.Version, description, url, author,
		gui.AboutYear(2023))
	gui.ShowAbout(appName, descHtml, iconSvg, me.config.TextSize)

}

func (me *App) onHelp() {
	form := gui.MakeInfoForm("Help", appName, helpHtml, iconSvg, 600, 550,
		me.config.TextSize, true)
	form.Show()
}

func (me *App) onQuit() {
	me.config.X = me.Window.X()
	me.config.Y = me.Window.Y()
	me.config.Width = me.Window.W()
	me.config.Height = me.Window.H()
	me.config.Scale = fltk.ScreenScale(0)
	me.config.IncludeNonFreeSections = me.incNonFreeCheckbox.Value()
	me.config.save()
	me.Window.Destroy()
}
