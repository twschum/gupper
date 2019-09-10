/*
Implements really basic interface to a generic filestore
*/

package store

import (
	"log"
	"net/http"
)

type LatestServer struct {
	PackagePath string
}

func NewLatestServer(packagePath *string) *LatestServer {
	s := new(LatestServer)
	s.PackagePath = *packagePath
	return s
}

// ServeHTTP implements the HTTP user interface.
func (s *LatestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s from %s", r.Method, r.URL, r.RemoteAddr)
	w.Write([]byte("1.2"))
}
