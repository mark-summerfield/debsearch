// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mark-summerfield/clip"
	ds "github.com/mark-summerfield/debsearch"
	"github.com/mark-summerfield/gong"
	"github.com/mark-summerfield/gset"
)

func main() {
	log.SetFlags(0)
	config := getConfig()
	var pairs []ds.FilePair
	if config.words.IsEmpty() {
		pairs = ds.StdFilePairs()
	} else {
		pairs = ds.StdFilePairsWithDescriptions()
	}
	t := time.Now()
	pkgs, err := ds.NewPkgs(pairs...)
	gong.CheckError("failed to read package files", err)
	if config.listSections {
		fmt.Printf("sections (%d):\n", len(pkgs.Sections))
		for _, section := range pkgs.Sections.ToSortedSlice() {
			fmt.Printf("\t%s\n", section)
		}
	}
	if config.listTags {
		fmt.Printf("tags (%d):\n", len(pkgs.Tags))
		for _, tag := range pkgs.Tags.ToSortedSlice() {
			fmt.Printf("\t%s\n", tag)
		}
	}
	// TODO search
	if config.verbose {
		elapsed := time.Since(t)
		fmt.Printf("searched %s pkgs in %s\n",
			gong.Commas(len(pkgs.Pkgs)), elapsed)
	}
}

func getConfig() *Config {
	parser := clip.NewParserVersion(ds.Version)
	parser.LongDesc = "A tool for searching debian packages."
	listTagsOpt := parser.Flag("listtags", "Print the tag names.")
	listTagsOpt.SetShortName(clip.NoShortName)
	listSectionsOpt := parser.Flag("listsections",
		"Print the section names.")
	listSectionsOpt.SetShortName(clip.NoShortName)
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
	verboseOpt := parser.Flag("verbose",
		"Print number of packages and time taken.")
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
		listTags:     listTagsOpt.Value(),
		listSections: listSectionsOpt.Value(), verbose: verboseOpt.Value()}
	if sectionsOpt.Given() {
		config.sections.Add(strings.Split(sectionsOpt.Value(), ",")...)
	}
	if tagsOpt.Given() {
		config.tags.Add(strings.Split(tagsOpt.Value(), ",")...)
	}
	if len(parser.Positionals) > 0 {
		config.words.Add(parser.Positionals...)
	}
	return &config
}

type Config struct {
	ui           string
	sections     gset.Set[string]
	tags         gset.Set[string]
	words        gset.Set[string]
	listTags     bool
	listSections bool
	verbose      bool
}

func (me *Config) String() string {
	sections := strings.Join(me.sections.ToSortedSlice(), ",")
	tags := strings.Join(me.tags.ToSortedSlice(), ",")
	words := strings.Join(me.words.ToSortedSlice(), " ")
	return fmt.Sprintf("ui=%q sections=%q tags=%q words=%q", me.ui,
		sections, tags, words)
}
