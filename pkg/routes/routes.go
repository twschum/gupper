/*

Constants file for the route methods

*/
package routes

import (
	"net/url"
	"path"
)

const (
	LIST     = "/list"
	DOWNLOAD = "/download/"
)

// use copy and return string form
func List(base url.URL) (route string) {
	base.Path = LIST
	return base.String()
}

func Download(base url.URL, file string) (route string) {
	base.Path = path.Join(DOWNLOAD, file)
	return base.String()
}
