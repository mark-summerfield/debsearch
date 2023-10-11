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
)

func main() {
	config := getConfig()
	var pairs []ds.FilePair
	if config.query.Words.IsEmpty() {
		pairs = ds.StdFilePairs(config.arch)
	} else {
		pairs = ds.StdFilePairsWithDescriptions(config.arch)
	}
	t := time.Now()
	pkgs, err := ds.NewPkgs(pairs...)
	gong.CheckError("failed to read package files", err)
	maybePrintArchs(config)
	maybePrintSections(config, pkgs.SectionsAndCounts)
	maybePrintTags(config, pkgs.TagsAndCounts)
	elapsed := time.Since(t)
	if config.IsSearch() {
		search(config, pkgs, elapsed)
	} else if config.verbose {
		fmt.Printf("searched %s pkgs in %s.\n", gong.Commas(len(pkgs.Pkgs)),
			elapsed)
	}
}

func maybePrintArchs(config *Config) {
	if config.listArchs {
		archs := strings.Fields(ds.Archs)
		if config.verbose {
			fmt.Printf("Archs (%d):\n", len(archs))
		}
		for _, arch := range archs {
			if arch == ds.DefaultArch {
				arch += " [default]"
			}
			fmt.Println(arch)
		}
	}
}

func maybePrintSections(config *Config, sectionsAndCounts map[string]int) {
	if config.listSections {
		if config.verbose {
			fmt.Printf("Sections (%d):\n", len(sectionsAndCounts))
		}
		sections := gong.SortedMapKeys(sectionsAndCounts)
		for _, section := range sections {
			if config.verbose {
				count := sectionsAndCounts[section]
				fmt.Printf("%s (%s)\n", section, gong.Commas(count))
			} else {
				fmt.Println(section)
			}
		}
	}
}

func maybePrintTags(config *Config, tagsAndCounts map[string]int) {
	if config.listTags {
		if config.verbose {
			fmt.Printf("Tags (%d):\n", len(tagsAndCounts))
		}
		tags := gong.SortedMapKeys(tagsAndCounts)
		for _, tag := range tags {
			if config.verbose {
				count := tagsAndCounts[tag]
				fmt.Printf("%s (%s)\n", tag, gong.Commas(count))
			} else {
				fmt.Println(tag)
			}
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
			fmt.Printf("* %s\n", pkg)
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
	archOpt := parser.Choice("arch",
		"System arch(itecture) [default: "+ds.DefaultArch+"].",
		strings.Fields(ds.Archs), ds.DefaultArch)
	listArchsOpt := parser.Flag("list-archs", "Print arch(itecture) names.")
	listArchsOpt.SetShortName(clip.NoShortName)
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
	listTagsOpt := parser.Flag("list-tags",
		"Print tag names (and how many packages have each tag).")
	listTagsOpt.SetShortName(clip.NoShortName)
	listSectionsOpt := parser.Flag("list-sections",
		"Print section names (and how many packages are in each section).")
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
	config := Config{arch: archOpt.Value(), query: ds.NewQuery(),
		listArchs: listArchsOpt.Value(), listTags: listTagsOpt.Value(),
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
	arch         string
	query        *ds.Query
	listArchs    bool
	listTags     bool
	listSections bool
	verbose      bool
}

func (me *Config) IsValid() bool {
	return me.listArchs || me.listTags || me.listSections ||
		!me.query.Sections.IsEmpty() || !me.query.Tags.IsEmpty() ||
		!me.query.Words.IsEmpty()
}

func (me *Config) IsSearch() bool {
	return !me.query.Sections.IsEmpty() || !me.query.Tags.IsEmpty() ||
		!me.query.Words.IsEmpty()
}

func (me *Config) String() string {
	return fmt.Sprintf("query=%s listArchs=%t listTags=%t "+
		"listSections=%t verbose=%t",
		me.query, me.listArchs, me.listTags, me.listSections, me.verbose)
}
