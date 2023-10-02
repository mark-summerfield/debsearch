// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	_ "embed"
	"fmt"
)

//go:embed Version.dat
var Version string

func Hello() string {
	return fmt.Sprintf("Hello debsearch v%s", Version)
}
