// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	_ "embed"
	"errors"
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

func StdFilePairs() []FilePair { return stdFilePairs(false) }

func StdFilePairsWithDescriptions() []FilePair { return stdFilePairs(true) }

type pkgs map[string]*pkg

func NewPkgs(filepairs ...FilePair) (pkgs, error) {
	if len(filepairs) == 0 {
		return nil, Err103
	}
	// TODO put all this in a Parsers type so we can use me.pkgsMutex etc. &
	// refactor
	var errs error
	var errMutex sync.Mutex
	pkgs := pkgs{}
	var pkgsMutex sync.Mutex
	descForPkgs := map[string]string{}
	var descForPkgsMutex sync.Mutex
	var wg sync.WaitGroup
	for i, pair := range filepairs {
		wg.Add(1)
		go func(i int, pair FilePair) {
			defer wg.Done()
			somePkgs, err := readPackages(pair.Pkg)
			if err != nil {
				errMutex.Lock()
				errs = errors.Join(err)
				errMutex.Unlock()
			} else {
				pkgsMutex.Lock()
				for name, pkg := range somePkgs {
					pkgs[name] = pkg
				}
				pkgsMutex.Unlock()
			}
		}(i, pair)
		if pair.I18n != "" {
			wg.Add(1)
			go func(i int, pair FilePair) {
				defer wg.Done()
				someDescForPkg, err := readDescriptions(pair.I18n)
				if err != nil {
					errMutex.Lock()
					errs = errors.Join(err)
					errMutex.Unlock()
				} else {
					descForPkgsMutex.Lock()
					for name, desc := range someDescForPkg {
						descForPkgs[name] = desc
					}
					descForPkgsMutex.Unlock()
				}
			}(i, pair)
		}
	}
	wg.Wait()
	for name, long_desc := range descForPkgs { // merge
		pkgs[name].LongDesc = long_desc
	}
	return pkgs, errs
}
