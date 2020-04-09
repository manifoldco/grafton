// +build mage

package main

import (
	// mage:import
	_ "github.com/manifoldco/logo/grimoire"
	// mage:import
	"github.com/manifoldco/logo/grimoire/build"
	// mage:import
	_ "github.com/manifoldco/grafton/scripts"
)

var _ = build.SetStartingOffset(0)
