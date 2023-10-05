// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	"fmt"
	"strings"

	"github.com/mark-summerfield/gset"
)

type Query struct {
	Sections gset.Set[string] // sections are always or-ed
	Tags     gset.Set[string]
	TagsAnd  bool // if true all tags must match; else any
	Words    gset.Set[string]
	WordsAnd bool // if true all tags must match; else any
}

func NewQuery() *Query {
	return &Query{Sections: gset.New[string](), Tags: gset.New[string](),
		Words: gset.New[string]()}
}

func (me *Query) SelectFrom(pkgs *pkgs) gset.Set[*pkg] {
	matched := gset.New[*pkg]()
	for _, pkg := range pkgs.Pkgs {
		if me.match(pkg) {
			matched.Add(pkg)
		}
	}
	return matched
}

func (me *Query) match(pkg *pkg) bool {
	if !me.Sections.IsEmpty() {
	}
	if !me.Tags.IsEmpty() {
	}
	if !me.Words.IsEmpty() {
	}
	return false // TODO
}

func (me *Query) Clear() {
	me.Sections.Clear()
	me.Tags.Clear()
	me.TagsAnd = false
	me.Words.Clear()
	me.WordsAnd = false
}

func (me *Query) String() string {
	sections := strings.Join(me.Sections.ToSortedSlice(), ",")
	tags := strings.Join(me.Tags.ToSortedSlice(), ",")
	tagOp := "|"
	if me.TagsAnd {
		tagOp = "&"
	}
	words := strings.Join(me.Words.ToSortedSlice(), " ")
	wordOp := "|"
	if me.WordsAnd {
		wordOp = "&"
	}
	return fmt.Sprintf("sections|%q tags%s%q words%s%q", sections, tagOp,
		tags, wordOp, words)
}
