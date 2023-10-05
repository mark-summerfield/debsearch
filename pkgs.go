// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	"github.com/mark-summerfield/gset"
)

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
