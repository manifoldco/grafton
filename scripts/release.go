package scripts

import (
	"fmt"

	"github.com/magefile/mage/mg"

	"github.com/manifoldco/grafton/scripts/grimoire/cast"
)

// Release is a combined command that will both build and release the zip files using packr and
// promulgate.
func Release() {
	mg.SerialDeps(Build, ReleaseZips)
}

// ReleaseZips uses promulgate and the manifold CLI to release the current set of zips at the
// location where the command is run.
func ReleaseZips() error {
	tag, err := Version()
	if err != nil {
		return err
	}

	command := fmt.Sprintf("manifold run -t manifold -p promulgate -- promulgate release v%s", tag)
	return cast.Sh(command)
}
