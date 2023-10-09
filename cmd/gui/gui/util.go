// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package gui

import (
	"github.com/pwiecz/go-fltk"
)

func ButtonHeight() int {
	_, height := fltk.MeasureText("X", false)
	return int(float64(height)*2.25) + (2 * Margin)
}

func ButtonWidth() int {
	width := LabelWidth()
	return width + (width / 2)
}

func ReturnButtonWidth() int {
	return ButtonWidth() + Pad
}

func LabelHeight() int {
	return (ButtonHeight() * 4) / 5
}

func LabelWidth() int {
	width, _ := fltk.MeasureText("X", false)
	return int(float64(width) * 7.5)
}

func ToolbuttonIconSize() int {
	_, height := fltk.MeasureText("X", false)
	return (height * 2) - 4
}

func AddWindowIcon(window *fltk.Window, svgText string) {
	if image := ImageForSvgText(svgText, 0); image != nil {
		window.SetIcons([]*fltk.RgbImage{image})
	}
}

func ImageForSvgText(svgText string, size int) *fltk.RgbImage {
	if svg, err := fltk.NewSvgImageFromString(svgText); err == nil {
		if size != 0 {
			svg.Scale(size, size, true, true)
		}
		return &svg.RgbImage
	}
	return nil
}

func MakeToolbutton(svgText string) *fltk.Button {
	buttonHeight := ButtonHeight()
	button := fltk.NewButton(0, 0, buttonHeight, buttonHeight, "")
	button.ClearVisibleFocus()
	if image := ImageForSvgText(svgText,
		ToolbuttonIconSize()); image != nil {
		button.SetImage(image)
		button.SetAlign(fltk.ALIGN_IMAGE_BACKDROP)
	}
	return button
}

func MakeAccelLabel(width, height int, label string) *fltk.Button {
	button := fltk.NewButton(0, 0, width, height, label)
	button.SetAlign(fltk.ALIGN_INSIDE | fltk.ALIGN_LEFT)
	button.SetBox(fltk.NO_BOX)
	button.ClearVisibleFocus()
	return button
}

func MakeSep(y int, hbox *fltk.Flex) {
	sep := fltk.NewBox(fltk.THIN_DOWN_BOX, 0, y, Pad, ButtonHeight())
	hbox.Fixed(sep, Pad)
}

func MakeHBox(x, y, width, height int) *fltk.Flex {
	return makeVHBox(fltk.ROW, x, y, width, height)
}

func MakeVBox(x, y, width, height int) *fltk.Flex {
	return makeVHBox(fltk.COLUMN, x, y, width, height)
}

func makeVHBox(kind fltk.FlexType, x, y, width, height int) *fltk.Flex {
	box := fltk.NewFlex(x, y, width, height)
	box.SetType(kind)
	box.SetSpacing(Pad)
	box.SetMargin(Margin, Margin)
	return box
}

type MenuItem struct {
	Text     string
	Shortcut int
	Method   func()
	Divider  bool
}

func NewMenuItem(Text string, Shortcut int, Method func(),
	Divider bool) MenuItem {
	return MenuItem{Text, Shortcut, Method, Divider}
}

func MakeMenuItem(menuBar *fltk.MenuBar, markdown MenuItem) {
	flag := fltk.MENU_VALUE
	if markdown.Divider {
		flag |= fltk.MENU_DIVIDER
	}
	menuBar.AddEx(markdown.Text, markdown.Shortcut, markdown.Method, flag)
}

func Int8ToStr(raw []int8) string {
	data := make([]byte, 0, len(raw))
	for _, i := range raw {
		if i == 0 {
			break
		}
		data = append(data, byte(i))
	}
	return string(data)
}
