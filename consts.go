// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

const (
	listsPath = "/var/lib/apt/lists/"
)

var (
	packageGlobs = []string{
		"*contrib_binary*_Packages",
		"*main_binary*_Packages",
		"*non-free_binary*_Packages",
	}
	descGlobs = []string{
		"*contrib_i18n_Translation-en",
		"*main_i18n_Translation-en",
		"*non-free_i18n_Translation-en",
	}
)
