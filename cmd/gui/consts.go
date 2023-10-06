// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package main

import _ "embed"

//go:embed icon.svg
var iconSvg string

const (
	appName     = "DebFind"
	domain      = "qtrac.eu"
	description = "Application for finding Debian packages."
	url         = "https://github.com/mark-summerfield/debsearch"
	author      = "Mark Summerfield"

	tinyTimeout  = 0.005
	smallTimeout = 0.1
	rowHeight    = 32
	colWidth     = 60
)
