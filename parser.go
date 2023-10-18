// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/mark-summerfield/gong"
)

type parser struct {
	model            Model
	modelMutex       sync.Mutex
	err              error
	errMutex         sync.Mutex
	descForPkgs      map[string]string
	descForPkgsMutex sync.Mutex
}

func parse(filepairs ...FilePair) (Model, error) {
	if len(filepairs) == 0 {
		return Model{}, Err102
	}
	parser := &parser{model: newModel(), descForPkgs: map[string]string{}}
	return parser.parse(filepairs...)
}

func (me *parser) parse(filepairs ...FilePair) (Model, error) {
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
		if pkg, ok := me.model.Packages[name]; ok {
			pkg.LongDesc = longDesc
		}
	}
	return me.model, me.err
}

func (me *parser) readPackages(filename string) {
	model, err := readPackages(filename)
	if err != nil {
		me.errMutex.Lock()
		defer me.errMutex.Unlock()
		me.err = errors.Join(err)
	} else {
		me.modelMutex.Lock()
		defer me.modelMutex.Unlock()
		for name, pkg := range model.Packages {
			me.model.Packages[name] = pkg
		}
		for section, count := range model.SectionsAndCounts {
			me.model.SectionsAndCounts[section] += count
		}
		for tag, count := range model.TagsAndCounts {
			me.model.TagsAndCounts[tag] += count
		}
	}
}

func (me *parser) readDescriptions(filename string) {
	if descForPkgs := readDescriptions(filename); len(descForPkgs) > 0 {
		me.descForPkgsMutex.Lock()
		defer me.descForPkgsMutex.Unlock()
		for name, desc := range descForPkgs {
			me.descForPkgs[name] = desc
		}
	}
}

func readPackages(filename string) (Model, error) {
	model := newModel()
	file, err := os.Open(filename)
	if err != nil {
		return model, fmt.Errorf("%w: %s", Err101, err)
	}
	defer file.Close()
	state := &parseState{}
	pkg := NewPkg()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return model, err
		}
		if line == "" {
			state.Clear()
			continue
		}
		if strings.HasPrefix(line, packagePrefix) {
			state.Clear()
			if pkg.IsValid() {
				model.Packages[pkg.Name] = pkg.Copy()
			}
			pkg.Clear()
			pkg.Name = strings.TrimSpace(line[packagePrefixLen:])
		} else if strings.HasPrefix(line, " ") {
			if state.inTags {
				addTags(pkg, line, &model)
			} else if state.inDesc {
				pkg.LongDesc += getDesc(line)
			} else {
				state.Clear()
			}
		} else {
			state.Update(maybeAddKeyValue(pkg, line, &model))
		}
	}
	if pkg.IsValid() {
		model.Packages[pkg.Name] = pkg.Copy()
	}
	return model, nil
}

func addTags(pkg *pkg, line string, model *Model) {
	for _, item := range strings.Split(line, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			item = strings.ReplaceAll(item, "::", "/")
			pkg.Tags.Add(item)
			model.TagsAndCounts[item]++
		}
	}
}

func maybeAddKeyValue(pkg *pkg, line string, model *Model) (bool, bool) {
	if key, value, found := strings.Cut(line, ":"); found {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if value != "" {
			switch key {
			case "Description":
				pkg.ShortDesc = value
				return false, true
			case "Homepage":
				pkg.Url = value
			case "Installed-Size":
				pkg.Size = gong.StrToInt(value, 0)
			case "Size": // download size
				if pkg.Size == 0 {
					pkg.Size = gong.StrToInt(value, 0)
				}
			case "Section":
				pkg.Section = value
				model.SectionsAndCounts[value]++
			case "Tag":
				addTags(pkg, value, model)
				return true, false
			case "Version":
				pkg.Version = value
			}
		}
	}
	return false, false
}

func readDescriptions(filename string) map[string]string {
	descForPkg := map[string]string{}
	file, err := os.Open(filename)
	if err != nil {
		return descForPkg
	}
	defer file.Close()
	name := ""
	longDesc := ""
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
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
			longDesc += getDesc(line)
		}
	}
	if name != "" && longDesc != "" {
		descForPkg[name] = strings.TrimRight(longDesc, asciiWs)
	}
	return descForPkg
}

func getDesc(line string) string {
	line = strings.TrimRight(line[1:], asciiWs)
	if line == "." {
		line = ""
	}
	return line + "\n"
}

type parseState struct {
	inTags bool
	inDesc bool
}

func (me *parseState) Clear() {
	me.inTags = false
	me.inDesc = false
}

func (me *parseState) Update(inTags, inDesc bool) {
	me.inTags = inTags
	me.inDesc = inDesc
}
