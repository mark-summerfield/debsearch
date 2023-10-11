// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"github.com/mark-summerfield/debsearch/cmd/gui/gui"
	"github.com/pwiecz/go-fltk"
)

func selectOrClear(browser *fltk.MultiBrowser, sel bool) {
	for i := 1; i <= browser.Size(); i++ {
		browser.SetSelected(i, sel)
	}
}

func selectedCount(browser *fltk.MultiBrowser) int {
	count := 0
	for i := 1; i <= browser.Size(); i++ {
		if browser.IsSelected(i) {
			count++
		}
	}
	return count
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
