/*
Self-updating app, using a companion http server

*/
package main

import (
	"flag"
	"fmt"

	"github.com/twschum/gupper/pkg/update"
)

var updateUrl = flag.String("url", "localhost", "http update server address")
var updatePort = flag.String("port", ":8080", "http update server port")

// TODO https

func main() {
	flag.Parse()
	var url = fmt.Sprintf("http://%s%s", *updateUrl, *updatePort)
	var version = update.Check(&url)
	fmt.Printf("app version: %v\ndoing useful work now...\n", version)
}
