/*

All the metadata that should describe a compiled binary package,
with methods to convert a specific filename format into that metadata

Ex:
 name-version-os-arch[.exe]
 app-1.2-darwin-amd64
 thing-4.12.1-windows-386.exe

*/
package pkgmeta

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/twschum/gupper/pkg/version"
)

type PackageMeta struct {
	Filename string
	App      string
	Version  version.Version
	OS       string
	Arch     string
}

func (pkg PackageMeta) String() string {
	return fmt.Sprintf("%s version %s %s/%s", pkg.App, pkg.Version, pkg.OS, pkg.Arch)
}

func Parse(pkgFile *string) (pkg PackageMeta, err error) {
	pkg.Filename = *pkgFile
	// remove the windows suffix, this is also where a .zip/.tar would be removed
	s := strings.TrimSuffix(*pkgFile, ".exe")
	parts := strings.SplitN(s, "-", 4)
	if len(parts) != 4 {
		return PackageMeta{}, errors.New(fmt.Sprintf("Invalid package name: %v", *pkgFile))
	}
	pkg.App = parts[0]
	pkg.Version, err = version.Parse(parts[1])
	if err != nil {
		return PackageMeta{}, err
	}
	pkg.OS = parts[2]
	pkg.Arch = parts[3]
	return pkg, nil
}

// Determine this binary's metadata via runtime and args
func ThisPackageMeta(buildVersion *string, appName *string) (pkg PackageMeta, err error) {
	pkg.Filename = os.Args[0]
	pkg.App = *appName
	pkg.Version, err = version.Parse(*buildVersion)
	if err != nil {
		return PackageMeta{}, err
	}
	pkg.OS = runtime.GOOS
	pkg.Arch = runtime.GOARCH
	return pkg, nil
}

// Given the current package, find the latest available version
// If there are no packages, returns a zeroed PackageMeta
func GetLatestPackage(packages []PackageMeta, current *PackageMeta) PackageMeta {
	// Filter
	var filtered []PackageMeta
	for _, pkg := range packages {
		if (pkg.Arch == current.Arch) && (pkg.OS == current.OS) {
			filtered = append(filtered, pkg)
		}
	}
	// Sort descending by version
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Version.GT(filtered[j].Version)
	})
	if len(filtered) > 0 {
		return filtered[0]
	}
	return PackageMeta{}
}

// Get valid meta data from a list of package file names
func AvailablePackages(packageFiles []string) (packages []PackageMeta) {
	for _, file := range packageFiles {
		pkg, err := Parse(&file)
		if err != nil {
			continue
		}
		packages = append(packages, pkg)
	}
	return packages
}
