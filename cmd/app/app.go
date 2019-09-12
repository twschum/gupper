/*
Self-updating app, using a companion http server

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/twschum/gupper/pkg/pkgmeta"
	"github.com/twschum/gupper/pkg/update"
)

// may set at compile time
var BuildVersion string = "0.0.0"
var AppName string = "app"

var updateHost = flag.String("host", "localhost", "http update server address")
var updatePort = flag.String("port", ":8080", "http update server port (with :)")
var showVersion = flag.Bool("version", false, "print current version and exit")
var daemonize = flag.Bool("daemon", false, "mimic daemon-like behaviour that would periodically check for updates")

func main() {
	flag.Parse()
	base, err := url.Parse("http://" + *updateHost + *updatePort)
	if err != nil {
		log.Println("ERROR: bad update host:", err)
		return
	}
	current, err := pkgmeta.ThisPackageMeta(&BuildVersion, &AppName)
	if err != nil {
		log.Println("ERROR: Unable to determine self version:", err)
		return
	}
	if *showVersion {
		fmt.Println(current.Version)
		return
	} else {
		log.Println(current)
	}
	first := true
	for *daemonize || first {
		updated, err := update.Check(&current, base)
		if err != nil {
			log.Println("ERROR: Unable to update:", err)
		} else if updated.GT(current.Version) {
			log.Println(updated)
		} else {
			log.Println("Up to date")
		}
		// do core app things here
		fmt.Println("doing useful work now...")
		if *daemonize {
			time.Sleep(5 * time.Second)
		}
		first = false
	}
}
