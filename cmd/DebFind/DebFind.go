// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"os"

	"github.com/pwiecz/go-fltk"
)

func main() {
	config := newConfig()
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--debug" {
		config.debug = true // if using args follow with: args = args[1:]
	}
	if !config.debug {
		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("Unrecoverable error: %s", r)
				fltk.MessageBox(fmt.Sprintf("Error — %s", appName), message)
				fmt.Println(message)
			}
		}()
	}
	fltk.SetScheme("Oxy")
	fltk.SetScreenScale(0, config.Scale)
	app := newApp(config)
	app.Show()
	fltk.Run()
}
