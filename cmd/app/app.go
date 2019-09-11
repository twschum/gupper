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
)

// set at compile time
var BuildVersion string

var updateHost = flag.String("host", "localhost", "http update server address")
var updatePort = flag.String("port", ":8080", "http update server port (with :)")
var version = flag.Bool("version", false, "print current version and exit")

// TODO flag for current version and exit

// TODO https

func main() {
	flag.Parse()
	if *version {
		fmt.Println(BuildVersion)
		return
	}
	var base, err = url.Parse("http://" + *updateHost + *updatePort)
	if err != nil {
		log.Fatal(err)
	}
	var version = update.Check(BuildVersion, base)
	fmt.Printf("app version: %v\ndoing useful work now...\n", version)
}
