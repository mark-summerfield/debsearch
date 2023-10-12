// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import _ "embed"

//go:embed images/icon.svg
var iconSvg string

//go:embed images/config.svg
var configSvg string

//go:embed images/about.svg
var aboutSvg string

//go:embed images/help.svg
var helpSvg string

//go:embed images/quit.svg
var quitSvg string

//go:embed help.html
var helpHtml string

const (
	appName     = "DebFind"
	domain      = "qtrac.eu"
	description = "Application for finding Debian packages."
	url         = "https://github.com/mark-summerfield/debsearch"
	author      = "Mark Summerfield"

	tinyTimeout   = 0.005
	rowHeight     = 32
	colWidth      = 80
	nonfreePrefix = "non-free/"
	todoSuffix    = "/TODO"
	light1        = 255
	light2        = 52
	iconSize      = 22

	descTemplate = `<html><body>
<a href="%s"><font color=navy>%s</font></a>&nbsp;&nbsp;v%s&nbsp;&nbsp;%s
<p><font color=green>%s</font></p>
<p>
<pre><font face=helvetica>
%s
</font></pre>
</p>
</body></html>`
)
