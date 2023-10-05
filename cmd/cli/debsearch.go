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
	if config.query.Words.IsEmpty() {
		pairs = ds.StdFilePairs()
	} else {
		pairs = ds.StdFilePairsWithDescriptions()
	}
	t := time.Now()
	pkgs, err := ds.NewPkgs(pairs...)
	gong.CheckError("failed to read package files", err)
	maybePrintSections(config, pkgs.Sections)
	maybePrintTags(config, pkgs.Tags)
	elapsed := time.Since(t)
	if config.IsSearch() {
		search(config, pkgs, elapsed)
	} else if config.verbose {
		fmt.Printf("searched %s pkgs in %s.\n", gong.Commas(len(pkgs.Pkgs)),
			elapsed)
	}
}

func maybePrintSections(config *Config, sections gset.Set[string]) {
	if config.listSections {
		if config.verbose {
			fmt.Printf("Sections (%d):\n", len(sections))
		}
		for _, section := range sections.ToSortedSlice() {
			fmt.Println(section)
		}
	}
}

func maybePrintTags(config *Config, tags gset.Set[string]) {
	if config.listTags {
		if config.verbose {
			fmt.Printf("Tags (%d):\n", len(tags))
		}
		for _, tag := range tags.ToSortedSlice() {
			fmt.Println(tag)
		}
	}
}

func search(config *Config, pkgs ds.Pkgs, elapsed time.Duration) {
	matches := config.query.SelectFrom(&pkgs)
	if len(matches) == 0 {
		fmt.Printf(
			"searched %s pkgs in %s; no matching packages found.\n",
			gong.Commas(len(pkgs.Pkgs)), elapsed)
	} else {
		for _, pkg := range matches {
			fmt.Printf("• %s\n", pkg)
		}
		if config.verbose {
			fmt.Printf("found %s/%s pkgs in %s\n",
				gong.Commas(len(matches)), gong.Commas(len(pkgs.Pkgs)),
				elapsed)
		}
	}
}

func getConfig() *Config {
	parser := clip.NewParserVersion(ds.Version)
	parser.LongDesc = "A tool for searching Debian packages."
	sectionsOpt := parser.Str("sections", "Match any of the "+
		"comma-separated list of sections [default: match any section].",
		"")
	tagsOpt := parser.Str("tags", "Match the comma-separated list "+
		"of tags [default: match any tags].", "")
	allTagsOpt := parser.Flag("all-tags", "Match all the "+
		"given tags [default: match any given tag].")
	allTagsOpt.SetShortName(clip.NoShortName)
	allWordsOpt := parser.Flag("all-words", "Match all the "+
		"given words [default: match any given word].")
	allWordsOpt.SetShortName(clip.NoShortName)
	listTagsOpt := parser.Flag("list-tags", "Print tag names.")
	listTagsOpt.SetShortName(clip.NoShortName)
	listSectionsOpt := parser.Flag("list-sections", "Print section names.")
	listSectionsOpt.SetShortName(clip.NoShortName)
	verboseOpt := parser.Flag("verbose",
		"Print number of packages and how long to read them.")
	parser.PositionalCount = clip.ZeroOrMorePositionals
	parser.PositionalHelp = "Match the given (case-folded) words in " +
		"descriptions [no default]."
	parser.MustSetPositionalVarName("WORD")
	if err := parser.Parse(); err != nil {
		parser.OnError(err) // doesn't return
		return nil          // never reached
	}
	config := Config{query: ds.NewQuery(), listTags: listTagsOpt.Value(),
		listSections: listSectionsOpt.Value(), verbose: verboseOpt.Value()}
	if sectionsOpt.Given() {
		config.query.Sections.Add(
			strings.Split(sectionsOpt.Value(), ",")...)
	}
	if tagsOpt.Given() {
		config.query.Tags.Add(strings.Split(tagsOpt.Value(), ",")...)
	}
	config.query.TagsAnd = allTagsOpt.Value()
	config.query.WordsAnd = allWordsOpt.Value()
	if len(parser.Positionals) > 0 {
		for _, word := range parser.Positionals {
			config.query.Words.Add(strings.ToLower(word))
		}
	}
	if !config.IsValid() {
		parser.OnHelp() // doesn't return
	}
	return &config
}

type Config struct {
	query        *ds.Query
	listTags     bool
	listSections bool
	verbose      bool
}

func (me *Config) IsValid() bool {
	return me.listTags || me.listSections || !me.query.Sections.IsEmpty() ||
		!me.query.Tags.IsEmpty() || !me.query.Words.IsEmpty()
}

func (me *Config) IsSearch() bool {
	return !me.query.Sections.IsEmpty() || !me.query.Tags.IsEmpty() ||
		!me.query.Words.IsEmpty()
}

func (me *Config) String() string {
	return fmt.Sprintf("query=%s listTags=%t listSections=%t verbose=%t",
		me.query, me.listTags, me.listSections, me.verbose)
}
