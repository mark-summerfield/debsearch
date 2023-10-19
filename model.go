// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

type Model struct {
	Debs              map[string]*deb
	SectionsAndCounts map[string]int
	TagsAndCounts     map[string]int
}

func newModel() Model {
	return Model{Debs: map[string]*deb{},
		SectionsAndCounts: map[string]int{},
		TagsAndCounts:     map[string]int{}}
}

func NewModel(filepairs ...FilePair) (Model, error) {
	return parse(filepairs...)
}
