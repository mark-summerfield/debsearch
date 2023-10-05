// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/mark-summerfield/gong"
)

type parser struct {
	pkgs             Pkgs
	pkgsMutex        sync.Mutex
	err              error
	errMutex         sync.Mutex
	descForPkgs      map[string]string
	descForPkgsMutex sync.Mutex
}

func parse(filepairs ...FilePair) (Pkgs, error) {
	if len(filepairs) == 0 {
		return Pkgs{}, Err103
	}
	parser := &parser{pkgs: newPkgs(), descForPkgs: map[string]string{}}
	return parser.parse(filepairs...)
}

func (me *parser) parse(filepairs ...FilePair) (Pkgs, error) {
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
	for name, longDesc := range me.descForPkgs { // merge
		if pkg, ok := me.pkgs.Pkgs[name]; ok {
			pkg.LongDesc = longDesc
		}
	}
	return me.pkgs, me.err
}

func (me *parser) readPackages(filename string) {
	pkgs, err := readPackages(filename)
	if err != nil {
		me.errMutex.Lock()
		defer me.errMutex.Unlock()
		me.err = errors.Join(err)
	} else {
		me.pkgsMutex.Lock()
		defer me.pkgsMutex.Unlock()
		for name, pkg := range pkgs.Pkgs {
			me.pkgs.Pkgs[name] = pkg
		}
		me.pkgs.Sections.Unite(pkgs.Sections)
		me.pkgs.Tags.Unite(pkgs.Tags)
	}
}

func (me *parser) readDescriptions(filename string) {
	descForPkgs, err := readDescriptions(filename)
	if err != nil {
		me.errMutex.Lock()
		defer me.errMutex.Unlock()
		me.err = errors.Join(err)
	} else {
		me.descForPkgsMutex.Lock()
		defer me.descForPkgsMutex.Unlock()
		for name, desc := range descForPkgs {
			me.descForPkgs[name] = desc
		}
	}
}

func readPackages(filename string) (Pkgs, error) {
	pkgs := newPkgs()
	file, err := os.Open(filename)
	if err != nil {
		return pkgs, fmt.Errorf("%w: %s", Err101, err)
	}
	defer file.Close()
	pkg := NewPkg()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, packagePrefix) {
			if pkg.Name != "" && pkg.IsValid() {
				pkgs.Pkgs[pkg.Name] = pkg.Copy()
				pkg.Clear()
			}
			pkg.Name = strings.TrimSpace(line[packagePrefixLen:])
		} else if strings.HasPrefix(line, " ") {
			addTags(pkg, line, &pkgs)
		} else {
			maybeAddKeyValue(pkg, line, &pkgs)
		}
	}
	if pkg.IsValid() {
		pkgs.Pkgs[pkg.Name] = pkg.Copy()
	}
	return pkgs, nil
}

func addTags(pkg *pkg, line string, pkgs *Pkgs) {
	for _, item := range strings.Split(line, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			item = strings.ReplaceAll(item, "::", "/")
			pkg.Tags.Add(item)
			pkgs.Tags.Add(item)
		}
	}
}

func maybeAddKeyValue(pkg *pkg, line string, pkgs *Pkgs) {
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
	longDesc := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, packagePrefix) {
			if name != "" && longDesc != "" {
				descForPkg[name] = longDesc
			}
			name = strings.TrimSpace(line[packagePrefixLen:])
			longDesc = ""
		} else if strings.HasPrefix(line, " ") {
			line = strings.TrimSpace(line)
			if line == "." {
				line = "\n"
			}
			longDesc += line
		}
	}
	if name != "" && longDesc != "" {
		descForPkg[name] = longDesc
	}
	return descForPkg, nil
}
