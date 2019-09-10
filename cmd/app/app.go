/*
Self-updating app, using a companion http server

*/
package main

import (
	"github.com/twschum/gupper/pkg/update"
	"fmt"
)

func main() {
	var version = update.Check()
	fmt.Printf("app version: %v\ndoing useful work now...\n", version)
}
