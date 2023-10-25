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
		pairs = ds.StdFilePairs(config.arc)
	} else {
		pairs = ds.StdFilePairsWithDescriptions(config.arc)
	}
	t := time.Now()
	model, err := ds.NewModel(pairs...)
	gong.CheckError("failed to read package files", err)
	maybePrintArcs(config)
	maybePrintSections(config, model.SectionsAndCounts)
	maybePrintTags(config, model.TagsAndCounts)
	elapsed := time.Since(t)
	if config.IsSearch() {
		search(config, model, elapsed)
	} else if config.verbose {
		fmt.Printf("searched %s pkgs in %s.\n",
			gong.Commas(len(model.Debs)), elapsed)
	}
}

func maybePrintArcs(config *Config) {
	if config.listArcs {
		arcs := strings.Fields(ds.Arcs)
		if config.verbose {
			fmt.Printf("Arcs (%d):\n", len(arcs))
		}
		for _, arc := range arcs {
			if config.verbose && arc == ds.DefaultArc {
				arc += " [default]"
			}
			fmt.Println(arc)
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

func search(config *Config, model ds.Model, elapsed time.Duration) {
	matches := config.query.SelectFrom(&model)
	if len(matches) == 0 {
		fmt.Printf(
			"searched %s pkgs in %s; no matching packages found.\n",
			gong.Commas(len(model.Debs)), elapsed)
	} else {
		for _, deb := range matches {
			fmt.Printf("* %s\n", deb)
		}
		if config.verbose {
			fmt.Printf("found %s/%s pkgs in %s\n",
				gong.Commas(len(matches)), gong.Commas(len(model.Debs)),
				elapsed)
		}
	}
}

func getConfig() *Config {
	parser := clip.NewParserVersion(ds.Version)
	parser.LongDesc = "A tool for searching Debian packages."
	debugOpt := parser.Flag("debug", "")
	debugOpt.Hide()
	arcOpt := parser.Choice("arc",
		"System arc(hitecture) [default: "+ds.DefaultArc+"].",
		strings.Fields(ds.Arcs), ds.DefaultArc)
	listArcsOpt := parser.Flag("list-arcs", "Print arc(hitecture) names.")
	listArcsOpt.SetShortName(clip.NoShortName)
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
	config := Config{arc: arcOpt.Value(), query: ds.NewQuery(),
		listArcs: listArcsOpt.Value(), listTags: listTagsOpt.Value(),
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
	if debugOpt.Value() {
		fmt.Println(config.query)
	}
	return &config
}

type Config struct {
	arc          string
	query        *ds.Query
	listArcs     bool
	listTags     bool
	listSections bool
	verbose      bool
}

func (me *Config) IsValid() bool {
	return me.listArcs || me.listTags || me.listSections ||
		!me.query.Sections.IsEmpty() || !me.query.Tags.IsEmpty() ||
		!me.query.Words.IsEmpty()
}

func (me *Config) IsSearch() bool {
	return !me.query.Sections.IsEmpty() || !me.query.Tags.IsEmpty() ||
		!me.query.Words.IsEmpty()
}

func (me *Config) String() string {
	return fmt.Sprintf("query=%s listArcs=%t listTags=%t "+
		"listSections=%t verbose=%t",
		me.query, me.listArcs, me.listTags, me.listSections, me.verbose)
}
