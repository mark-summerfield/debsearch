// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/mark-summerfield/gong"
)

type parser struct {
	pkgs             pkgs
	pkgsMutex        sync.Mutex
	err              error
	errMutex         sync.Mutex
	descForPkgs      map[string]string
	descForPkgsMutex sync.Mutex
}

func parse(filepairs ...FilePair) (pkgs, error) {
	if len(filepairs) == 0 {
		return pkgs{}, Err103
	}
	parser := &parser{pkgs: newPkgs(), descForPkgs: map[string]string{}}
	return parser.parse(filepairs...)
}

func (me *parser) parse(filepairs ...FilePair) (pkgs, error) {
	var wg sync.WaitGroup
	for i, pair := range filepairs {
		wg.Add(1)
		go func(i int, pair FilePair) {
			defer wg.Done()
			me.readPackages(pair.Pkg)
		}(i, pair)
		if pair.I18n != "" {
			wg.Add(1)
			go func(i int, pair FilePair) {
				defer wg.Done()
				me.readDescriptions(pair.I18n)
			}(i, pair)
		}
	}
	wg.Wait()
	for name, long_desc := range me.descForPkgs { // merge
		if pkg, ok := me.pkgs.Pkgs[name]; ok {
			pkg.LongDesc = long_desc
		}
	}
	return me.pkgs, me.err
}

func (me *parser) readPackages(filename string) {
	pkgs, err := readPackages(filename)
	if err != nil {
		me.errMutex.Lock()
		me.err = errors.Join(err)
		me.errMutex.Unlock()
	} else {
		me.pkgsMutex.Lock()
		for name, pkg := range pkgs.Pkgs {
			me.pkgs.Pkgs[name] = pkg
		}
		me.pkgs.Sections.Unite(pkgs.Sections)
		me.pkgs.Tags.Unite(pkgs.Tags)
		me.pkgsMutex.Unlock()
	}
}

func (me *parser) readDescriptions(filename string) {
	descForPkgs, err := readDescriptions(filename)
	if err != nil {
		me.errMutex.Lock()
		me.err = errors.Join(err)
		me.errMutex.Unlock()
	} else {
		me.descForPkgsMutex.Lock()
		for name, desc := range descForPkgs {
			me.descForPkgs[name] = desc
		}
		me.descForPkgsMutex.Unlock()
	}
}

func readPackages(filename string) (pkgs, error) {
	pkgs := newPkgs()
	file, err := os.Open(filename)
	if err != nil {
		return pkgs, fmt.Errorf("%w: %s", Err101, err)
	}
	defer file.Close()
	var name string
	pkg := NewPkg()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, packagePrefix) {
			if name != "" && pkg.IsValid() {
				pkgs.Pkgs[name] = pkg.Copy()
				pkg.Clear()
			}
			name = strings.TrimSpace(line[packagePrefixLen:])
		} else if strings.HasPrefix(line, " ") {
			addTags(pkg, line, &pkgs)
		} else {
			maybeAddKeyValue(pkg, line, &pkgs)
		}
	}
	if name != "" && pkg.IsValid() {
		pkgs.Pkgs[name] = pkg.Copy()
	}
	return pkgs, nil
}

func addTags(pkg *pkg, line string, pkgs *pkgs) {
	for _, item := range strings.Split(line, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			pkg.Tags.Add(item)
			pkgs.Tags.Add(item)
		}
	}
}

func maybeAddKeyValue(pkg *pkg, line string, pkgs *pkgs) {
	if key, value, found := strings.Cut(line, ":"); found {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if value != "" {
			switch key {
			case "Description":
				pkg.ShortDesc = value
			case "Homepage":
				pkg.Url = value
			case "Installed-Size":
				pkg.Size = gong.StrToInt(value, 0)
			case "Section":
				pkg.Section = value
				pkgs.Sections.Add(value)
			case "Size":
				pkg.DownloadSize = gong.StrToInt(value, 0)
			case "Tag":
				addTags(pkg, value, pkgs)
			case "Version":
				pkg.Version = value
			}
		}
	}
}

func readDescriptions(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", Err102, err)
	}
	defer file.Close()
	descForPkg := map[string]string{}
	name := ""
	long_desc := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, packagePrefix) {
			if name != "" && long_desc != "" {
				descForPkg[name] = long_desc
			}
			name = strings.TrimSpace(line[packagePrefixLen:])
			long_desc = ""
		} else if strings.HasPrefix(line, " ") {
			line = strings.TrimSpace(line)
			if line == "." {
				line = "\n"
			}
			long_desc += line
		}
	}
	if name != "" && long_desc != "" {
		descForPkg[name] = long_desc
	}
	return descForPkg, nil
}
