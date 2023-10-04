// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	_ "embed"

	"github.com/mark-summerfield/gset"
)

//go:embed Version.dat
var Version string

type pkgs struct {
	Pkgs     map[string]*pkg
	Sections gset.Set[string]
	Tags     gset.Set[string]
}

func newPkgs() pkgs {
	return pkgs{Pkgs: map[string]*pkg{}, Sections: gset.New[string](),
		Tags: gset.New[string]()}
}

func NewPkgs(filepairs ...FilePair) (pkgs, error) {
	return parse(filepairs...)
}

type FilePair struct {
	Pkg  string
	I18n string
}

func NewFilePair(pkg, i18n string) FilePair {
	return FilePair{pkg, i18n}
}

func StdFilePairs() []FilePair { return stdFilePairs(false) }

func StdFilePairsWithDescriptions() []FilePair { return stdFilePairs(true) }
