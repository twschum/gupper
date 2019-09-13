/*

HTTP endpoint for package management:
 - tell a client what files it has
 - allow downloads of specific files
 - is completely open (for now) and shouldn't be used in broader environments

*/

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/twschum/gupper/pkg/routes"
)

var (
	httpPort    = flag.String("port", ":8080", "Listen address")
	packagePath = flag.String("pkgdir", "packages", "Directory at which to store packages")
)

func main() {
	flag.Parse()
	http.Handle(routes.LIST, logRequest(NewListServer(packagePath)))
	fs := http.FileServer(http.Dir(*packagePath))
	fs = http.StripPrefix(routes.DOWNLOAD, fs)
	http.Handle(routes.DOWNLOAD, logRequest(fs))
	log.Println("Listening on port", *httpPort)
	err := http.ListenAndServe(*httpPort, nil)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL, r.RemoteAddr)
		h.ServeHTTP(w, r)
	})
}

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
		http.Error(w, err.Error(), 500)
	}
	resp, err := json.Marshal(packages)
	if err != nil {
		http.Error(w, err.Error(), 500)
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
