// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pwiecz/go-fltk"
)

func main() {
	log.SetFlags(0)
	args := os.Args
	if len(args) > 1 && args[1] != "--debug" {
		args = args[2:]
		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("Unrecoverable error: %s", r)
				fltk.MessageBox(fmt.Sprintf("Error — %s", appName), message)
				fmt.Println(message)
			}
		}()
	} else {
		args = args[1:]
	}
	fltk.SetScheme("Oxy")
	config := newConfig()
	fltk.SetScreenScale(0, config.Scale)
	app := newApp(config, args)
	app.Show()
	fltk.Run()
}
