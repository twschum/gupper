/*
Implements really basic interface to a generic filestore
*/

package store

import (
	"io/ioutil"
	"log"
	"net/http"
	_ "os"
	_ "path/filepath"

	"github.com/twschum/gupper/pkg/version"
)

type PackageMeta struct {
	Version  version.Version
	Filename string
}

func NewPackageMeta(pkgFile *string) (pkg *PackageMeta, err error) {
	pkg = new(PackageMeta)
	pkg.Filename = *pkgFile
	// TODO strip out the version number once it's not the whole file
	pkg.Version, err = version.Parse(*pkgFile)
	return pkg, err
}

// Implement the /latest response

type LatestServer struct {
	PackagePath *string
}

func NewLatestServer(packagePath *string) *LatestServer {
	s := new(LatestServer)
	s.PackagePath = packagePath
	return s
}

// ServeHTTP implements the HTTP user interface.
func (s *LatestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s from %s", r.Method, r.URL, r.RemoteAddr)
	pkg := getLatestPackage(s.PackagePath)
	log.Println(pkg)
	w.Write([]byte(pkg.Filename))
}

// Returns a zeroed PackageMeta if none available
func getLatestPackage(packagePath *string) PackageMeta {
	packages := availablePackages(packagePath)
	if len(packages) > 0 {
		return packages[0]
	}
	return PackageMeta{}
}

// All available package metadata, sorted by version number
func availablePackages(packagePath *string) (packages []PackageMeta) {
	files, err := listFiles(packagePath)
	if err != nil {
		log.Fatal(err) // TODO
	}
	for _, file := range files {
		pkg, err := NewPackageMeta(&file)
		if err != nil {
			log.Println("Bad Package: ", err)
			continue
		}
		packages = append(packages, *pkg)
	}
	// TODO actual sort
	return
}

// Excludes directories, lexographically sorted by ioutil.ReadDir
func listFiles(root *string) (files []string, err error) {
	fileInfo, err := ioutil.ReadDir(*root)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		if !file.IsDir() {
			files = append(files, file.Name())
		}
	}
}

// Implement the serving of actual files from the packages dir
