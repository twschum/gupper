/*
Self-updating app, using a companion http server

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/twschum/gupper/pkg/pkgmeta"
	"github.com/twschum/gupper/pkg/update"
)

// set at compile time
var BuildVersion string = "0.0.0"

var updateHost = flag.String("host", "localhost", "http update server address")
var updatePort = flag.String("port", ":8080", "http update server port (with :)")
var showVersion = flag.Bool("version", false, "print current version and exit")

// TODO https

func main() {
	flag.Parse()
	base, err := url.Parse("http://" + *updateHost + *updatePort)
	if err != nil {
		log.Println("ERROR: bad update host:", err)
		return
	}
	current, err := pkgmeta.ThisPackageMeta(&BuildVersion)
	if err != nil {
		log.Println("ERROR: Unable to determine self version:", err)
		return
	}
	if *showVersion {
		fmt.Printf("%#v\n", current)
		return
	} else {
		log.Println("app version:", current.Version)
	}
	updated, err := update.Check(&current, base)
	if err != nil {
		log.Println("ERROR: Unable to update:", err)
	} else if updated.GT(current.Version) {
		fmt.Println("app version:", updated)
	} else {
		log.Println("Up to date")
	}
	// hook for actual app things
	fmt.Println("doing useful work now...")
}
