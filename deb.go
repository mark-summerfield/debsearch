// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mark-summerfield/gset"
)

type deb struct {
	Name      string
	Version   string
	Size      int
	Url       string
	Section   string
	Tags      gset.Set[string]
	ShortDesc string
	LongDesc  string
}

func NewDeb() *deb { return &deb{Tags: gset.New[string]()} }

func (me *deb) Copy() *deb {
	return &deb{Name: me.Name, Version: me.Version, Size: me.Size,
		Url: me.Url, Section: me.Section, Tags: me.Tags.Copy(),
		ShortDesc: me.ShortDesc, LongDesc: me.LongDesc}
}

func (me *deb) Clear() {
	me.Name = ""
	me.Version = ""
	me.Size = 0
	me.Url = ""
	me.Section = ""
	me.Tags.Clear()
	me.ShortDesc = ""
	me.LongDesc = ""
}

func (me *deb) IsValid() bool {
	return me.Name != "" && me.Version != "" && me.Size > 0 &&
		me.Section != "" && me.ShortDesc != ""
}

func (me *deb) Words() gset.Set[string] {
	rx := regexp.MustCompile(`\W+`)
	words := gset.New[string]()
	for _, text := range []string{me.Name, me.ShortDesc, me.LongDesc} {
		text = strings.ToLower(rx.ReplaceAllLiteralString(text, " "))
		for _, word := range strings.Fields(text) {
			words.Add(word)
		}
	}
	return words
}

func (me *deb) String() string {
	return fmt.Sprintf("%s v%s %s %q %s", me.Name, me.Version,
		HumanSize(me.Size), me.ShortDesc, me.Url)
}
