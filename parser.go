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
			me.readPackages(pair.Packages)
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
		if deb, ok := me.model.Debs[name]; ok {
			deb.LongDesc = longDesc
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
		for name, deb := range model.Debs {
			me.model.Debs[name] = deb
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
	deb := NewDeb()
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
			if deb.IsValid() {
				model.Debs[deb.Name] = deb.Copy()
			}
			deb.Clear()
			deb.Name = strings.TrimSpace(line[packagePrefixLen:])
		} else if strings.HasPrefix(line, " ") {
			if state.inTags {
				addTags(deb, line, &model)
			} else if state.inDesc {
				deb.LongDesc += getDesc(line)
			} else {
				state.Clear()
			}
		} else {
			state.Update(maybeAddKeyValue(deb, line, &model))
		}
	}
	if deb.IsValid() {
		model.Debs[deb.Name] = deb.Copy()
	}
	return model, nil
}

func addTags(deb *deb, line string, model *Model) {
	for _, item := range strings.Split(line, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			item = strings.ReplaceAll(item, "::", "/")
			deb.Tags.Add(item)
			model.TagsAndCounts[item]++
		}
	}
}

func maybeAddKeyValue(deb *deb, line string, model *Model) (bool, bool) {
	if key, value, found := strings.Cut(line, ":"); found {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if value != "" {
			switch key {
			case "Description":
				deb.ShortDesc = value
				return false, true
			case "Homepage":
				deb.Url = value
			case "Installed-Size":
				deb.Size = gong.StrToInt(value, 0)
			case "Size": // download size
				if deb.Size == 0 {
					deb.Size = gong.StrToInt(value, 0)
				}
			case "Section":
				deb.Section = value
				model.SectionsAndCounts[value]++
			case "Tag":
				addTags(deb, value, model)
				return true, false
			case "Version":
				deb.Version = value
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
