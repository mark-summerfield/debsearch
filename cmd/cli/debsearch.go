// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/mark-summerfield/clip"
	ds "github.com/mark-summerfield/debsearch"
	"github.com/mark-summerfield/gong"
	"github.com/mark-summerfield/gset"
)

func main() {
	t := time.Now()
	config := getConfig()
	fmt.Printf("debsearch v%s\n", ds.Version)
	fmt.Println(config)
	pairs := ds.StdFilePairsWithDescriptions()
	//pairs := ds.StdFilePairs()
	if pkgs, err := ds.NewPkgs(pairs...); err != nil {
		panic(err)
	} else {
		elapsed := time.Since(t)
		fmt.Printf("found %s pkgs in %s\n", gong.Commas(len(pkgs.Pkgs)),
			elapsed)
		//i := 0
		//for name, pkg := range pkgs {
		//	fmt.Println(name, pkg)
		//	i++
		//	if i > 10 {
		//		break
		//	}
		//}
	}
}

func getConfig() *Config {
	parser := clip.NewParserVersion(ds.Version)
	parser.LongDesc = "A tool for searching debian packages."
	infoOpt := parser.Flag("info",
		"Print names of sections, tags, and basic stats")
	uiOpt := parser.Str("ui", "Constrain to the given UI "+
		"(cli, tui, or gui) [default: any].", "")
	uiOpt.Validator = func(name, value string) (string, string) {
		value = strings.ToLower(value)
		for _, valid := range []string{"cli", "tui", "gui"} {
			if value == valid {
				return value, ""
			}
		}
		return "", fmt.Sprintf("invalid format: %q", value)
	}
	sectionsOpt := parser.Str("sections", "A comma-separated list "+
		"of sections (these are or-ed) [no default].", "")
	tagsOpt := parser.Str("tags", "A comma-separated list "+
		"of tags (these are and-ed) [no default].", "")
	parser.PositionalCount = clip.ZeroOrMorePositionals
	parser.PositionalHelp = "Words to search for (these are and-ed) " +
		"[no default]."
	parser.MustSetPositionalVarName("WORD")
	if err := parser.Parse(); err != nil {
		parser.OnError(err) // doesn't return
		return nil          // never reached
	}
	config := Config{ui: uiOpt.Value(), sections: gset.New[string](),
		tags: gset.New[string](), words: gset.New[string](),
		info: infoOpt.Value()}
	if sectionsOpt.Given() {
		config.sections.Add(strings.Split(sectionsOpt.Value(), ",")...)
	}
	if tagsOpt.Given() {
		config.tags.Add(strings.Split(tagsOpt.Value(), ",")...)
	}
	if len(parser.Positionals) > 0 {
		fmt.Println("X", parser.Positionals)
		config.words.Add(parser.Positionals...)
	}
	return &config
}

type Config struct {
	ui       string
	sections gset.Set[string]
	tags     gset.Set[string]
	words    gset.Set[string]
	info     bool
}

func (me *Config) String() string {
	sections := strings.Join(me.sections.ToSortedSlice(), ",")
	tags := strings.Join(me.tags.ToSortedSlice(), ",")
	words := strings.Join(me.words.ToSortedSlice(), " ")
	return fmt.Sprintf("ui=%q sections=%q tags=%q words=%q info=%t", me.ui,
		sections, tags, words, me.info)
}
