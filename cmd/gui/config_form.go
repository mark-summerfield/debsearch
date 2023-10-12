// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"strings"

	ds "github.com/mark-summerfield/debsearch"
	"github.com/mark-summerfield/debsearch/cmd/gui/gui"
	"github.com/pwiecz/go-fltk"
)

type configForm struct {
	*fltk.Window
	width      int
	height     int
	labelWidth int
	app        *App
	arcChoice  *fltk.Choice
}

func newConfigForm(app *App) configForm {
	form := configForm{width: 260, height: 145, app: app}
	form.Window = fltk.NewWindow(form.width, form.height)
	form.Window.SetLabel("Configure — " + appName)
	gui.AddWindowIcon(form.Window, iconSvg)
	form.labelWidth, _ = fltk.MeasureText("ArchitectureW", false)
	form.makeWidgets()
	form.Window.End()
	return form
}

func (me *configForm) makeWidgets() {
	vbox := gui.MakeVBox(0, 0, me.width, me.height)
	hbox := me.makeScaleRow()
	vbox.Fixed(hbox, rowHeight)
	hbox = me.makeTextSizeRow()
	vbox.Fixed(hbox, rowHeight)
	hbox = me.makeArcRow()
	vbox.Fixed(hbox, rowHeight)
	hbox = me.makeButtonRow()
	vbox.Fixed(hbox, rowHeight)
	vbox.End()
}

func (me *configForm) makeScaleRow() *fltk.Flex {
	hbox := gui.MakeHBox(0, 0, me.width, rowHeight)
	scaleLabel := gui.MakeAccelLabel(me.labelWidth, rowHeight, "&Scale")
	scaleSpinner := me.makeScaleSpinner()
	scaleLabel.SetCallback(func() { scaleSpinner.TakeFocus() })
	hbox.Fixed(scaleLabel, me.labelWidth)
	hbox.End()
	scaleSpinner.TakeFocus()
	return hbox
}

func (me *configForm) makeScaleSpinner() *fltk.Spinner {
	spinner := fltk.NewSpinner(0, 0, me.labelWidth, rowHeight)
	spinner.SetTooltip("Sets the application's scale.")
	spinner.SetType(fltk.SPINNER_FLOAT_INPUT)
	spinner.SetMinimum(0.5)
	spinner.SetMaximum(3.5)
	spinner.SetStep(0.1)
	spinner.SetValue(float64(fltk.ScreenScale(0)))
	spinner.SetCallback(func() {
		fltk.SetScreenScale(0, float32(spinner.Value()))
	})
	return spinner
}

func (me *configForm) makeTextSizeRow() *fltk.Flex {
	buttonHeight := gui.ButtonHeight()
	labelWidth := gui.LabelWidth()
	hbox := gui.MakeHBox(0, 0, me.width, rowHeight)
	sizeLabel := gui.MakeAccelLabel(me.labelWidth, buttonHeight,
		"&Text Size")
	spinner := fltk.NewSpinner(0, 0, labelWidth, buttonHeight)
	spinner.SetTooltip("Set the size of the about and help texts.")
	spinner.SetType(fltk.SPINNER_INT_INPUT)
	spinner.SetMinimum(10)
	spinner.SetMaximum(20)
	spinner.SetValue(float64(me.app.config.TextSize))
	spinner.SetCallback(func() {
		me.app.config.TextSize = int(spinner.Value())
	})
	sizeLabel.SetCallback(func() { spinner.TakeFocus() })
	hbox.Fixed(sizeLabel, me.labelWidth)
	hbox.End()
	return hbox
}

func (me *configForm) makeArcRow() *fltk.Flex {
	buttonHeight := gui.ButtonHeight()
	hbox := gui.MakeHBox(0, 0, me.width, buttonHeight)
	arcLabel := gui.MakeAccelLabel(me.labelWidth, buttonHeight,
		"&Architecture")
	me.arcChoice = fltk.NewChoice(0, 0, gui.LabelWidth(), buttonHeight)
	current := 0
	for i, name := range strings.Fields(ds.Arcs) {
		me.arcChoice.Add(name, nil)
		if name == me.app.config.Arc {
			current = i
		}
	}
	me.arcChoice.SetValue(current)
	arcLabel.SetCallback(func() { me.arcChoice.TakeFocus() })
	hbox.Fixed(arcLabel, me.labelWidth)
	hbox.End()
	return hbox
}

func (me *configForm) makeButtonRow() *fltk.Flex {
	buttonWidth := gui.ButtonWidth()
	buttonHeight := gui.ButtonHeight()
	hbox := gui.MakeHBox(0, 0, me.width, rowHeight)
	spacerWidth := (me.width - buttonWidth) / 2
	leftSpacer := gui.MakeHBox(0, 0, spacerWidth, buttonHeight)
	leftSpacer.End()
	button := fltk.NewButton(0, 0, buttonHeight, buttonWidth, "&Close")
	button.SetCallback(me.onClose)
	righttSpacer := gui.MakeHBox(spacerWidth+buttonWidth, 0, spacerWidth,
		buttonHeight)
	righttSpacer.End()
	hbox.End()
	return hbox
}

func (me *configForm) onClose() {
	oldArc := me.app.config.Arc
	newArc := me.arcChoice.SelectedText()
	if oldArc != newArc {
		me.app.config.Arc = newArc
		me.app.loadPackages()
	}
	me.app.descView.TextSize(me.app.config.TextSize)
	me.Window.Destroy()
}
