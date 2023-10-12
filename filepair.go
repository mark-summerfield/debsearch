// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

type FilePair struct {
	Pkg  string
	I18n string
}

func NewFilePair(pkg, i18n string) FilePair {
	return FilePair{pkg, i18n}
}

func StdFilePairs(arc string) []FilePair {
	return stdFilePairs(arc, false)
}

func StdFilePairsWithDescriptions(arc string) []FilePair {
	return stdFilePairs(arc, true)
}
