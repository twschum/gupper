/*

Constants file for the route methods

*/
package routes

import (
	"net/url"
	"path"
)

const (
	LATEST   = "/latest"
	DOWNLOAD = "/download"
	UPLOAD   = "/upload"
)

// use copy and return string form
func Latest(base url.URL) (route string) {
	base.Path = LATEST
	return base.String()
}

func Download(base url.URL, file string) (route string) {
	base.Path = path.Join(DOWNLOAD, file)
	return base.String()
}
