// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import "strings"

func descFileForPkgFile(filename string) string {
	if prefix, _, found := strings.Cut(filename, "_binary"); found {
		return prefix + "_i18n_Translation-en"
	}
	return ""
}
