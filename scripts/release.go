package scripts

import (
	"fmt"

	"github.com/magefile/mage/mg"

	"github.com/manifoldco/logo/grimoire/cast"
)

func ReleaseGrafton() {
	mg.SerialDeps(BuildZips, ReleaseZips)
}

func ReleaseZips() error {
	tag, err := Version()
	if err != nil {
		return err
	}

	command := fmt.Sprintf("manifold run -t manifold -p promulgate -- promulgate release v%s", tag)
	return cast.Sh(command)
}
