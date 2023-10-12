// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	"path/filepath"
	"strings"

	"github.com/mark-summerfield/gong"
)

func stdFilePairs(arc string, withDescriptions bool) []FilePair {
	pairs := []FilePair{}
	glob := filepath.Join(listsPath, "*"+arc+"_Packages")
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

func HumanSize(size int) string {
	units := "KB"
	if size > 1024 {
		size /= 1024
		if size > 1024 {
			size /= 1024
			units = "GB"
		} else {
			units = "MB"
		}
	}
	return gong.Commas(size) + units
}
