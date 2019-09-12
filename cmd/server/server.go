/*

HTTP endpoint for package management:
 - tell a client the latest version
 - allow downloads of specific versions
 - allow uploads of new version packages
 - is super unsecure (for now) and shouldn't be used in broader environments

*/

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/twschum/gupper/pkg/routes"
	"github.com/twschum/gupper/pkg/store"
)

var (
	httpPort    = flag.String("port", ":8080", "Listen address")
	packagePath = flag.String("pkgdir", "packages", "Directory at which to store packages")
)

// TODO cleaner error handling scheme

func main() {
	flag.Parse()
	http.Handle(routes.LIST, logRequest(store.NewListServer(packagePath)))
	fs := http.FileServer(http.Dir(*packagePath))
	fs = http.StripPrefix(routes.DOWNLOAD, fs)
	http.Handle(routes.DOWNLOAD, logRequest(fs))
	log.Println("Listening on port", *httpPort)
	log.Fatal(http.ListenAndServe(*httpPort, nil))
}

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL, r.RemoteAddr)
		h.ServeHTTP(w, r)
	})
}
