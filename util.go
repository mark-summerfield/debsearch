// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	"path/filepath"
	"strings"
)

func stdFilePairs(arch string, withDescriptions bool) []FilePair {
	pairs := []FilePair{}
	glob := filepath.Join(listsPath, "*"+arch+"_Packages")
	if matches, err := filepath.Glob(glob); err == nil {
		for _, pkgFile := range matches {
			if !strings.Contains(pkgFile, "i386") {
				descFile := ""
				if withDescriptions {
					descFile = descFileForPkgFile(pkgFile)
				}
				pairs = append(pairs, NewFilePair(pkgFile, descFile))
			}
		}
	}
	return pairs
}

func descFileForPkgFile(filename string) string {
	if prefix, _, found := strings.Cut(filename, "_binary"); found {
		return prefix + "_i18n_Translation-en"
	}
	return ""
}
