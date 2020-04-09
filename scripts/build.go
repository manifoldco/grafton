package scripts

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/manifoldco/logo/grimoire/cast"
)

var version = ""

var git = cast.PrepShOutput("git")

func Version() (string, error) {
	if version != "" {
		return version, nil
	}

	desc, err := git("describe --tags --dirty")
	if err != nil {
		return "", err
	}

	parts := strings.Split(desc, "-")

	version = parts[0]
	return version, nil
}

// Build compiles a local version of the code into bin/.
func BuildZips() { mg.SerialDeps(BuildZipDarwinAmd64, BuildZipLinuxAmd64, BuildZipWindowAmd64) }

func BuildZipDarwinAmd64() error {
	err := buildBinary("darwin", "amd64", "grafton", "./cmd")
	if err != nil {
		return err
	}

	return zipBinary("darwin", "amd64", "tar.gz", "grafton")
}

func BuildZipLinuxAmd64() error {
	err := buildBinary("linux", "amd64", "grafton", "./cmd")
	if err != nil {
		return err
	}

	return zipBinary("linux", "amd64", "tar.gz", "grafton")
}

func BuildZipWindowAmd64() error {
	err := buildBinary("windows", "amd64", "grafton.exe", "./cmd")
	if err != nil {
		return err
	}

	return zipBinary("windows", "amd64", "zip", "grafton.exe")
}

func buildBinary(os, arch, file, input string) error {
	err := cast.Sh("go get github.com/gobuffalo/packr/...")
	if err != nil {
		return err
	}

	tag, err := Version()
	if err != nil {
		return err
	}

	osArch := fmt.Sprintf("%s_%s", os, arch)

	ldFlags := fmt.Sprintf("-w -X github.com/manifoldco/grafton/config.Version=%s", tag)
	output := fmt.Sprintf("build/%s/bin/%s", osArch, file)
	prefix := fmt.Sprintf("PREFIX=build/%s GOOS=%s GOARCH=%s", osArch, os, arch)
	command := prefix + " " + `CGO_ENABLED=0 packr build -i --ldflags="%s" -o %s %s`
	err = cast.Sh(command, ldFlags, output, input)
	if err != nil {
		return err
	}

	return cast.Sh("packr clean")
}

func zipBinary(os, arch, extension, input string) error {
	tag, err := Version()
	if err != nil {
		return err
	}

	osArch := fmt.Sprintf("%s_%s", os, arch)

	command := fmt.Sprintf("tar -czf grafton_%s_%s.%s build/%s/bin/%s", osArch, tag, extension, osArch, input)
	return cast.Sh(command)
}
