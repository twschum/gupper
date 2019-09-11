/*

Version is just a number for now, could easily be replaced by
a semantic versioning object that could handle more realistic use-case

*/
package version

import (
	"strconv"
)

type Version float64

const Bad Version = 0

func Parse(s string) (Version, error) {
	v, err := strconv.ParseFloat(s, 64)
	return Version(v), err
}

type PackageMeta struct {
	Version  Version
	Filename string
}

func NewPackageMeta(pkgFile *string) (pkg *PackageMeta, err error) {
	pkg = new(PackageMeta)
	pkg.Filename = *pkgFile
	// TODO strip out the version number once it's not the whole file
	pkg.Version, err = Parse(*pkgFile)
	return pkg, err
}
