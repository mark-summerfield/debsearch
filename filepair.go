// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

type FilePair struct {
	Packages string
	I18n     string
}

func NewFilePair(packages, i18n string) FilePair {
	return FilePair{packages, i18n}
}

func StdFilePairs(arc string) []FilePair {
	return stdFilePairs(arc, false)
}

func StdFilePairsWithDescriptions(arc string) []FilePair {
	return stdFilePairs(arc, true)
}
