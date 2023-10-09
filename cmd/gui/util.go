// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"github.com/mark-summerfield/debsearch/cmd/gui/gui"
	"github.com/pwiecz/go-fltk"
)

func ifDebug(debug bool, box *fltk.Flex, color fltk.Color) {
	if debug {
		box.SetBox(fltk.DOWN_BOX)
		box.SetColor(color)
	}
}

func selectOrClear(browser *fltk.MultiBrowser, sel bool) {
	for i := 1; i <= browser.Size(); i++ {
		browser.SetSelected(i, sel)
	}
}

func padBox(parent *fltk.Flex, size int) {
	pad := fltk.NewBox(fltk.FLAT_BOX, 0, 0, gui.Margin, gui.ButtonHeight())
	parent.Fixed(pad, gui.Margin)
}

func divider(parent *fltk.Flex) {
	const thickness = 3
	pad := fltk.NewBox(fltk.ENGRAVED_BOX, 0, 0, 100, thickness)
	parent.Fixed(pad, thickness)
}
