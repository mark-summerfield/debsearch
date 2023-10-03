// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/mark-summerfield/gset"
)

//go:embed Version.dat
var Version string

// /var/lib/apt/lists/deb.debian.org_debian_dists_bullseye_main_binary-amd64_Packages
// /var/lib/apt/lists/deb.debian.org_debian_dists_bullseye_main_i18n_Translation-en

// Doesn't include `Package: name` since that's the map key.
type pkg struct {
	version       string
	size          int
	download_size int
	url           string
	section       string
	tags          gset.Set[string]
	short_desc    string
	long_desc     string
}

func (me *pkg) Copy() *pkg {
	return &pkg{version: me.version, size: me.size,
		download_size: me.download_size, url: me.url, section: me.section,
		tags: me.tags.Copy(), short_desc: me.short_desc,
		long_desc: me.long_desc}
}

func (me *pkg) Clear() {
	me.version = ""
	me.size = 0
	me.download_size = 0
	me.url = ""
	me.section = ""
	me.tags.Clear()
	me.short_desc = ""
	me.long_desc = ""
}

func (me *pkg) HasEnoughInfo() bool {
	return me.version != "" && me.size > 0 && me.section != "" &&
		me.short_desc != ""
}

type pkgs map[string]*pkg

func NewPkgs(pkg_filename string) (pkgs, error) {
	return NewPkgsX(pkg_filename, "")
}

func NewPkgsX(pkg_filename, desc_filename string) (pkgs, error) {
	pkgs := pkgs{}
	err := readPackages(pkg_filename, pkgs)
	if err != nil {
		return pkgs, err
	}
	if desc_filename != "" {
		if err := readDescriptions(desc_filename, pkgs); err != nil {
			return pkgs, err
		}
	}
	return pkgs, nil
}

func readPackages(filename string, pkgs pkgs) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("%w: %s", Err101, err)
	}
	defer file.Close()
	var name string
	pkg := &pkg{}
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
	return nil
}

func addTags(pkg *pkg, line string) {
}

func maybeAddKeyValue(pkg *pkg, line string) {
}

func readDescriptions(filename string, pkgs pkgs) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("%w: %s", Err102, err)
	}
	defer file.Close()
	//scanner := bufio.NewScanner(file)
	//for scanner.Scan() {
	//	line := scanner.Text()
	//}
	return nil
}
