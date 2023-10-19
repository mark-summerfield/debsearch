// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"html"
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

func (me *App) onInfo(info string) { me.onMessage(info, "navy") }

func (me *App) onWarn(warn string) { me.onMessage(warn, "maroon") }

func (me *App) onError(err error) { me.onMessage(err.Error(), "red") }

func (me *App) onMessage(msg, color string) {
	me.descView.SetValue("<font color=" + color + ">" +
		html.EscapeString(msg) + "</font>")
	me.Redraw()
}

func (me *App) onHtmlMessage(htmlMsg string) {
	me.descView.SetValue(htmlMsg)
	me.Redraw()
}

func (me *App) onFind() {
	me.packagesBrowser.Clear()
	me.onInfo("Searching…")
	query := me.makeQuery()
	me.updateResults(query)
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

func (me *App) updateResults(query *ds.Query) {
	debs := query.SelectFrom(me.model)
	me.updatePackagesLabel(len(debs))
	if len(debs) == 0 {
		me.onWarn("No matching packages found.")
	} else {
		me.updatePackageBrowserWidths()
		bg := light1
		for _, deb := range debs {
			me.packagesBrowser.Add(fmt.Sprintf("@B%d@.%s\t@B%d@.%s", bg,
				deb.Name, bg, deb.ShortDesc))
			if bg == light1 {
				bg = light2
			} else {
				bg = light1
			}
		}
		me.packagesBrowser.SetSelected(1, true)
		me.packagesBrowser.TakeFocus()
		me.onSelectPackage()
	}
}

func (me *App) updatePackageBrowserWidths() {
	width := me.packagesBrowser.W()
	nWidth, _ := fltk.MeasureText("n", false)
	left := min(nWidth*20, width/2)
	me.packagesBrowser.SetColumnWidths(left, width-left)
}

func (me *App) onSelectPackage() {
	if i := me.packagesBrowser.Value(); i > 0 {
		if text := me.packagesBrowser.Text(i); text != "" {
			if j := strings.IndexByte(text, '\t'); j > -1 {
				text = text[:j]
				if j := strings.Index(text, "@."); j > -1 {
					if text = text[j+2:]; text != "" {
						me.showDescription(text)
					}
				}
			}
		}
	}
}

func (me *App) showDescription(name string) {
	if deb, ok := me.model.Debs[name]; ok {
		me.descView.SetValue(fmt.Sprintf(descTemplate,
			deb.Url, html.EscapeString(deb.Name),
			html.EscapeString(deb.Version), ds.HumanSize(deb.Size),
			html.EscapeString(deb.ShortDesc),
			html.EscapeString(deb.LongDesc)))
	}
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
	me.config.AllTags = me.tagsMatchAllRadioButton.Value()
	me.config.AllWords = me.wordsMatchAllRadioButton.Value()
	me.config.save()
	me.Window.Destroy()
}
