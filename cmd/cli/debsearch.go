// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"

	ds "github.com/mark-summerfield/debsearch"
)

func main() {
	fmt.Println("Version", ds.Version)
	for _, pair := range ds.StdFilePairs(true) {
		fmt.Printf("%q\n", pair.Pkg)
		fmt.Printf("%q\n", pair.I18n)
	}
}
