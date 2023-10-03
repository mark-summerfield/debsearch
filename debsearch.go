// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	_ "embed"
	"errors"
	"path/filepath"
	"sync"
)

//go:embed Version.dat
var Version string

// /var/lib/apt/lists/deb.debian.org_debian_dists_bullseye_main_binary-amd64_Packages
// /var/lib/apt/lists/deb.debian.org_debian_dists_bullseye_main_i18n_Translation-en

type FilePair struct {
	Pkg  string
	I18n string
}

func NewFilePair(pkg, i18n string) FilePair {
	return FilePair{pkg, i18n}
}

func StdFilePairs(withDescriptions bool) []FilePair {
	pairs := []FilePair{}
	for _, glob := range packageGlobs {
		glob = filepath.Join(listsPath, glob)
		if matches, err := filepath.Glob(glob); err == nil {
			for _, pkgFile := range matches {
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

type pkgs map[string]*pkg

func NewPkgs(filepairs ...FilePair) (pkgs, error) {
	if len(filepairs) == 0 {
		return nil, Err103
	}
	allPkgs := make([]pkgs, 0, len(filepairs))
	pkgErrs := make([]error, len(filepairs))
	allDescsForPkgs := make([]map[string]string, 0, len(filepairs))
	descErrs := make([]error, len(filepairs))
	var wg sync.WaitGroup
	for i, pair := range filepairs {
		go func(i int, pair FilePair) {
			defer wg.Done()
			if err := readPackages(pair.Pkg, allPkgs[i]); err != nil {
				pkgErrs[i] = err
			}
		}(i, pair)
		if pair.I18n != "" {
			go func(i int, pair FilePair) {
				defer wg.Done()
				if err := readDescriptions(pair.I18n,
					allDescsForPkgs[i]); err != nil {
					descErrs[i] = err
				}
			}(i, pair)
		}
	}
	wg.Wait()
	err := errors.Join(pkgErrs...)
	err = errors.Join(err, errors.Join(descErrs...))
	pkgs := pkgs{}
	for _, ps := range allPkgs { // merge
		for name, pkg := range ps {
			pkgs[name] = pkg
		}
	}
	for _, descs := range allDescsForPkgs { // merge
		for name, long_desc := range descs {
			pkgs[name].long_desc = long_desc
		}
	}
	return pkgs, err
}
