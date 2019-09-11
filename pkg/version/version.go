/*

Version is just a number for now, could easily be replaced by
a semantic versioning object that could handle more realistic use-case

*/
package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type PackageMeta struct {
	// TODO need all these things findable in the built package
	Version  Version
	Filename string
	// os
	// arch
	// app name
}

func NewPackageMeta(pkgFile *string) (pkg *PackageMeta, err error) {
	pkg = new(PackageMeta)
	pkg.Filename = *pkgFile
	// TODO strip out the version number once it's not the whole file
	pkg.Version, err = Parse(*pkgFile)
	return pkg, err
}

type Version struct {
	Major int
	Minor int
	Patch int
}

// Parse out a basic semantic version from string
// Missing parts are considered 0
// 1, 1.2, 1.2.3
func Parse(s string) (v Version, err error) {
	var major, minor, patch int
	if len(s) == 0 {
		return v, errors.New("Empty version string")
	}
	parts := strings.SplitN(s, ".", 3)
	if len(parts) >= 1 {
		major, err = strconv.Atoi(parts[0])
		if err != nil {
			return
		}
	}
	if len(parts) >= 2 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			return
		}
	}
	if len(parts) == 3 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return
		}
	}
	v.Major, v.Minor, v.Patch = major, minor, patch
	return v, nil
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Equals checks if v is equal to o.
func (v Version) EQ(o Version) bool {
	return (v.Compare(o) == 0)
}

// NE checks if v is not equal to o.
func (v Version) NE(o Version) bool {
	return (v.Compare(o) != 0)
}

// GT checks if v is greater than o.
func (v Version) GT(o Version) bool {
	return (v.Compare(o) == 1)
}

// GE checks if v is greater than or equal to o.
func (v Version) GE(o Version) bool {
	return (v.Compare(o) >= 0)
}

// LT checks if v is less than o.
func (v Version) LT(o Version) bool {
	return (v.Compare(o) == -1)
}

// LE checks if v is less than or equal to o.
func (v Version) LE(o Version) bool {
	return (v.Compare(o) <= 0)
}

// Compare compares Versions v to o:
// -1 == v is less than o
// 0 == v is equal to o
// 1 == v is greater than o
func (v Version) Compare(o Version) int {
	if v.Major != o.Major {
		if v.Major > o.Major {
			return 1
		}
		return -1
	}
	if v.Minor != o.Minor {
		if v.Minor > o.Minor {
			return 1
		}
		return -1
	}
	if v.Patch != o.Patch {
		if v.Patch > o.Patch {
			return 1
		}
		return -1
	}
	return 0
}
