// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	_ "embed"

	"github.com/mark-summerfield/gset"
)

// Doesn't include `Package: name` since that's the map key.
type pkg struct {
	version       string
	size          int
	download_size int
	url           string
	section       string
	tags          gset.Set[string]
	short_desc    string
	long_desc     string
}

func (me *pkg) Copy() *pkg {
	return &pkg{version: me.version, size: me.size,
		download_size: me.download_size, url: me.url, section: me.section,
		tags: me.tags.Copy(), short_desc: me.short_desc,
		long_desc: me.long_desc}
}

func (me *pkg) Clear() {
	me.version = ""
	me.size = 0
	me.download_size = 0
	me.url = ""
	me.section = ""
	me.tags.Clear()
	me.short_desc = ""
	me.long_desc = ""
}

func (me *pkg) HasEnoughInfo() bool {
	return me.version != "" && me.size > 0 && me.section != "" &&
		me.short_desc != ""
}
