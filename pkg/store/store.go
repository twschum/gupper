/*
Implements really basic interface to a generic filestore
*/

package store

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Implements handler that just returns a list of the files available Additional query
// args could narrow it down with some sort of match string but this is literally as
// simple as it gets. Let the client figure it out since this could easily be replaced
type ListServer struct {
	PackagePath *string
}

func NewListServer(packagePath *string) *ListServer {
	s := new(ListServer)
	s.PackagePath = packagePath
	return s
}

func (s *ListServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	packages, err := listFiles(s.PackagePath)
	if err != nil {
		log.Fatal(err) // TODO
	}
	resp, err := json.Marshal(packages)
	if err != nil {
		log.Fatal(err) // TODO
	}
	w.Write(resp)
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
