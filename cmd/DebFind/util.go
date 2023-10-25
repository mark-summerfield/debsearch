// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"strings"

	"github.com/mark-summerfield/debsearch/cmd/DebFind/gui"
	"github.com/pwiecz/go-fltk"
)

func selectOrClear(browser *fltk.MultiBrowser, sel bool) {
	for i := 1; i <= browser.Size(); i++ {
		browser.SetSelected(i, sel)
	}
}

func selected(browser *fltk.MultiBrowser) []string {
	selected := []string{}
	for i := 1; i <= browser.Size(); i++ {
		if browser.IsSelected(i) {
			text := browser.Text(i)
			if i := strings.IndexByte(text, '('); i > -1 {
				text = strings.TrimSpace(text[:i])
			}
			selected = append(selected, text)
		}
	}
	return selected
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

func makeButton(x int, text, icon string, callback func()) *fltk.Button {
	buttonHeight := gui.ButtonHeight()
	buttonWidth := ((gui.LabelWidth() * 5) / 3) + (2 * gui.Pad)
	button := fltk.NewButton(x, 0, buttonWidth, buttonHeight, text)
	button.SetImage(gui.ImageForSvgText(icon, iconSize))
	button.SetAlign(fltk.ALIGN_IMAGE_NEXT_TO_TEXT)
	button.SetCallback(callback)
	return button
}
