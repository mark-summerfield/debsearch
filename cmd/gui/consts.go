// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import _ "embed"

//go:embed icon.svg
var iconSvg string

//go:embed help.html
var helpHtml string

const (
	appName     = "DebFind"
	domain      = "qtrac.eu"
	description = "Application for finding Debian packages."
	url         = "https://github.com/mark-summerfield/debsearch"
	author      = "Mark Summerfield"

	tinyTimeout = 0.005
	rowHeight   = 32
	colWidth    = 80
)
