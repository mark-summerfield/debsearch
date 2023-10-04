// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	_ "embed"

	"github.com/mark-summerfield/gset"
)

// Doesn't include `Package: name` since that's the map key.
type pkg struct {
	Version      string
	Size         int
	DownloadSize int
	Url          string
	Section      string
	Tags         gset.Set[string]
	ShortDesc    string
	LongDesc     string
}

func NewPkg() *pkg { return &pkg{Tags: gset.New[string]()} }

func (me *pkg) Copy() *pkg {
	return &pkg{Version: me.Version, Size: me.Size,
		DownloadSize: me.DownloadSize, Url: me.Url, Section: me.Section,
		Tags: me.Tags.Copy(), ShortDesc: me.ShortDesc,
		LongDesc: me.LongDesc}
}

func (me *pkg) Clear() {
	me.Version = ""
	me.Size = 0
	me.DownloadSize = 0
	me.Url = ""
	me.Section = ""
	me.Tags.Clear()
	me.ShortDesc = ""
	me.LongDesc = ""
}

func (me *pkg) IsValid() bool {
	return me.Version != "" && me.Size > 0 && me.Section != "" &&
		me.ShortDesc != ""
}
