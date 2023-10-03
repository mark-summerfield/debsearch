// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"time"

	ds "github.com/mark-summerfield/debsearch"
)

func main() {
	fmt.Printf("debsearch v%s\n", ds.Version)
	pairs := ds.StdFilePairs()
	t := time.Now()
	if pkgs, err := ds.NewPkgs(pairs...); err != nil {
		panic(err)
	} else {
		elapsed := time.Since(t)
		fmt.Printf("found %d pkgs in %s\n", len(pkgs), elapsed)
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
