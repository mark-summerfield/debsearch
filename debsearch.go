// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	_ "embed"
	"fmt"
	"strings"

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

type Query struct {
	Ui       string           // any UI if empty; else cli or tui or gui
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
	fmt.Println("SelectFrom", len(pkgs.Pkgs), me) // TODO
	return matched
}

func (me *Query) Clear() {
	me.Ui = ""
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
	return fmt.Sprintf("ui=%q sections|%q tags%s%q words%s%q", me.Ui,
		sections, tagOp, tags, wordOp, words)
}
