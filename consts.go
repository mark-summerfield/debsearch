// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package debsearch

import (
	_ "embed"
	"errors"
)

//go:embed Version.dat
var Version string

const (
	DefaultArc = "amd64"

	Arcs = "alpha amd64 arm arm64 armel armhf avr32 hppa hurd-amd64 " +
		"hurd-i386 i386 ia64 kfreebsd-amd64 kfreebsd-i386 m32 m68k mips " +
		"mips64el mipsel netbsd-alpha netbsd-i386 or1k powerpc " +
		"powerpcspe ppc64el riscv64 s390 s390x sh4 sparc sparc64 x32"

	listsPath        = "/var/lib/apt/lists/"
	packagePrefix    = "Package:"
	packagePrefixLen = len(packagePrefix)
)

var (
	Err101 = errors.New("E101: failed to open packages file")
	Err102 = errors.New("E102: no package files given")
)
