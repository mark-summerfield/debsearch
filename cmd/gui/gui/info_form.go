// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package gui

import (
	"fmt"

	"github.com/pwiecz/go-fltk"
)

func MakeInfoForm(title, appName, descHtml, iconSvg string, width, height,
	textSize int, resizable bool) *fltk.Window {
	window := fltk.NewWindow(width, height)
	if resizable {
		window.Resizable(window)
	}
	window.SetLabel(fmt.Sprintf("%s — %s", title, appName))
	AddWindowIcon(window, iconSvg)
	buttonWidth := ButtonWidth()
	buttonHeight := ButtonHeight()
	returnButtonWidth := ReturnButtonWidth()
	vbox := MakeVBox(0, 0, width, height, Pad)
	view := fltk.NewHelpView(0, 0, width, height-buttonHeight)
	view.TextFont(fltk.HELVETICA)
	view.TextSize(textSize)
	view.SetValue(descHtml)
	y := height - buttonHeight
	hbox := MakeHBox(0, y, width, buttonHeight, Pad)
	spacerWidth := (width - returnButtonWidth) / 2
	leftSpacer := MakeHBox(0, y, spacerWidth, buttonHeight, 0)
	leftSpacer.End()
	button := fltk.NewReturnButton(0, 0, buttonHeight, returnButtonWidth,
		"&Close")
	button.SetCallback(func() { window.Destroy() })
	righttSpacer := MakeHBox(spacerWidth+returnButtonWidth, y, spacerWidth,
		buttonHeight, 0)
	righttSpacer.End()
	hbox.Fixed(button, buttonWidth)
	hbox.End()
	vbox.Fixed(hbox, buttonHeight)
	vbox.End()
	window.End()
	return window
}
