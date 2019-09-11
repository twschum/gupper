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

func main() {
	flag.Parse()
	http.Handle(routes.LATEST, store.NewLatestServer(packagePath))
	fs := http.FileServer(http.Dir(*packagePath))
	http.Handle(routes.DOWNLOAD, http.StripPrefix(routes.DOWNLOAD, fs))
	log.Fatal(http.ListenAndServe(*httpPort, nil))
}
