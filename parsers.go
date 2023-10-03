// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/mark-summerfield/gong"
)

func readPackages(filename string) (pkgs, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", Err101, err)
	}
	defer file.Close()
	pkgs := pkgs{}
	var name string
	pkg := NewPkg()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, packagePrefix) {
			if name != "" && pkg.HasEnoughInfo() {
				pkgs[name] = pkg.Copy()
				pkg.Clear()
			}
			name = strings.TrimSpace(line[packagePrefixLen:])
		} else if strings.HasPrefix(line, " ") {
			addTags(pkg, line)
		} else {
			maybeAddKeyValue(pkg, line)
		}
	}
	if name != "" && pkg.HasEnoughInfo() {
		pkgs[name] = pkg.Copy()
	}
	return pkgs, nil
}

func addTags(pkg *pkg, line string) {
	for _, item := range strings.Split(line, ",") {
		pkg.tags.Add(strings.TrimSpace(item))
	}
}

func maybeAddKeyValue(pkg *pkg, line string) {
	if key, value, found := strings.Cut(line, ":"); found {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		switch key {
		case "Description":
			pkg.short_desc = value
		case "Homepage":
			pkg.url = value
		case "Installed-Size":
			pkg.size = gong.StrToInt(value, 0)
		case "Section":
			pkg.section = value
		case "Size":
			pkg.download_size = gong.StrToInt(value, 0)
		case "Tag":
			addTags(pkg, value)
		case "Version":
			pkg.version = value
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
