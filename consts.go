// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	"errors"
)

const (
	listsPath        = "/var/lib/apt/lists/"
	packagePrefix    = "Package:"
	packagePrefixLen = len(packagePrefix)
)

var (
	Err101 = errors.New("E101: failed to open packages file")
	Err102 = errors.New("E102: failed to open descriptions file")
	Err103 = errors.New("E103: no package files given")
)
