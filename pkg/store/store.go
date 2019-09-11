/*
Implements really basic interface to a generic filestore
*/

package store

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_ "os"
	_ "path/filepath"
	"sort"

	"github.com/twschum/gupper/pkg/version"
)

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
	resp, err := json.Marshal(pkg)
	if err != nil {
		log.Fatal(err) // TODO
	}
	w.Write(resp)
}

// Returns a zeroed PackageMeta if none available
func getLatestPackage(packagePath *string) version.PackageMeta {
	packages := availablePackages(packagePath)
	if len(packages) > 0 {
		return packages[0]
	}
	return version.PackageMeta{}
}

// All available package metadata, sorted by version number
func availablePackages(packagePath *string) (packages []version.PackageMeta) {
	files, err := listFiles(packagePath)
	if err != nil {
		log.Fatal(err) // TODO
	}
	for _, file := range files {
		pkg, err := version.NewPackageMeta(&file)
		if err != nil {
			log.Println("Bad Package: ", err)
			continue
		}
		packages = append(packages, *pkg)
	}
	// Sort descending by version
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Version > packages[j].Version
	})
	return packages
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
	return files, nil
}

// Implement the serving of actual files from the packages dir
