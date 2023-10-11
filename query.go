// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	"cmp"
	"fmt"
	"slices"
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

func (me *Query) SelectFrom(pkgs *Pkgs) []*pkg {
	matched := gset.New[*pkg]()
	for _, pkg := range pkgs.Pkgs {
		if me.Match(pkg) {
			matched.Add(pkg)
		}
	}
	slice := matched.ToSlice()
	slices.SortFunc(slice, func(a, b *pkg) int {
		return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})
	return slice
}

func (me *Query) Match(pkg *pkg) bool {
	if !me.Sections.IsEmpty() && !me.Sections.Contains(pkg.Section) {
		return false // no specified section matches
	}
	if !me.Tags.IsEmpty() {
		intersection := me.Tags.Intersection(pkg.Tags)
		if intersection.IsEmpty() {
			return false // no tags match
		}
		if me.TagsAnd && !me.Tags.IsSubsetOf(pkg.Tags) {
			return false // not all tags match
		}
	}
	if !me.Words.IsEmpty() {
		words := pkg.Words()
		intersection := me.Words.Intersection(words)
		if intersection.IsEmpty() {
			return false // no words match
		}
		if me.WordsAnd && !me.Words.IsSubsetOf(words) {
			return false // not all words match
		}
	}
	return true
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
