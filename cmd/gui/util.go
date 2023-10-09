// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
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
