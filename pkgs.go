// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

type Pkgs struct {
	Pkgs              map[string]*pkg
	SectionsAndCounts map[string]int
	TagsAndCounts     map[string]int
}

func newPkgs() Pkgs {
	return Pkgs{Pkgs: map[string]*pkg{},
		SectionsAndCounts: map[string]int{},
		TagsAndCounts:     map[string]int{}}
}

func NewPkgs(filepairs ...FilePair) (Pkgs, error) {
	return parse(filepairs...)
}
