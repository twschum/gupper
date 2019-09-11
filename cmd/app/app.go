/*
Self-updating app, using a companion http server

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/twschum/gupper/pkg/update"
	"github.com/twschum/gupper/pkg/version"
)

// set at compile time
var BuildVersion string

var updateHost = flag.String("host", "localhost", "http update server address")
var updatePort = flag.String("port", ":8080", "http update server port (with :)")
var showVersion = flag.Bool("version", false, "print current version and exit")

// TODO flag for current version and exit

// TODO https

func main() {
	flag.Parse()
	current, err := version.Parse(BuildVersion)
	if err != nil {
		log.Fatalf("ERROR: Bad BuildVersion: %v: %v", BuildVersion, err)
	}
	fmt.Println("app version:", current)
	if *showVersion {
		fmt.Println(current)
		return
	}
	base, err := url.Parse("http://" + *updateHost + *updatePort)
	if err != nil {
		log.Fatal(err)
	}
	updated := update.Check(current, base)
	if updated > current {
		fmt.Println("app version:", updated)
	}
	// hook for actual app things
	fmt.Println("doing useful work now...")
}
