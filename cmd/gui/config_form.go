// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"github.com/mark-summerfield/debsearch/cmd/gui/gui"
	"github.com/pwiecz/go-fltk"
)

type configForm struct {
	*fltk.Window
	width  int
	height int
}

func newConfigForm(app *App) configForm {
	form := configForm{width: 240, height: 110}
	form.Window = fltk.NewWindow(form.width, form.height)
	form.Window.SetLabel("Configure — " + appName)
	gui.AddWindowIcon(form.Window, iconSvg)
	form.makeWidgets(app)
	form.Window.End()
	return form
}

func (me *configForm) makeWidgets(app *App) {
	vbox := gui.MakeVBox(0, 0, me.width, me.height)
	hbox := me.makeScaleRow()
	vbox.Fixed(hbox, rowHeight)
	hbox = me.makeTextSizeRow(app)
	vbox.Fixed(hbox, rowHeight)
	hbox = me.makeButtonRow()
	vbox.Fixed(hbox, rowHeight)
	vbox.End()
}

func (me *configForm) makeScaleRow() *fltk.Flex {
	hbox := gui.MakeHBox(0, 0, me.width, rowHeight)
	scaleLabel := gui.MakeAccelLabel(colWidth, rowHeight, "&Scale")
	scaleSpinner := me.makeScaleSpinner()
	scaleLabel.SetCallback(func() { scaleSpinner.TakeFocus() })
	hbox.Fixed(scaleLabel, colWidth)
	hbox.End()
	scaleSpinner.TakeFocus()
	return hbox
}

func (me *configForm) makeScaleSpinner() *fltk.Spinner {
	spinner := fltk.NewSpinner(0, 0, colWidth, rowHeight)
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

func (me *configForm) makeTextSizeRow(app *App) *fltk.Flex {
	buttonHeight := gui.ButtonHeight()
	labelWidth := gui.LabelWidth()
	hbox := gui.MakeHBox(0, 0, me.width, rowHeight)
	sizeLabel := gui.MakeAccelLabel(colWidth, buttonHeight, "&Text Size")
	spinner := fltk.NewSpinner(0, 0, labelWidth, buttonHeight)
	spinner.SetTooltip("Set the size of the about and help texts.")
	spinner.SetType(fltk.SPINNER_INT_INPUT)
	spinner.SetMinimum(10)
	spinner.SetMaximum(20)
	spinner.SetValue(float64(app.config.TextSize))
	spinner.SetCallback(func() {
		app.config.TextSize = int(spinner.Value())
	})
	sizeLabel.SetCallback(func() { spinner.TakeFocus() })
	hbox.Fixed(sizeLabel, colWidth)
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
	button.SetCallback(func() { me.Window.Destroy() })
	righttSpacer := gui.MakeHBox(spacerWidth+buttonWidth, 0, spacerWidth,
		buttonHeight)
	righttSpacer.End()
	hbox.End()
	return hbox
}
