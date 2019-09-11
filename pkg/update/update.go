/*

Module handles updating the client from the server

*/
package update

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/twschum/gupper/pkg/routes"
	"github.com/twschum/gupper/pkg/version"
)

// embedded in the code here for now,
// maybe as a command line argument later at build time
const current version.Version = 1.0

func Check(url *string) (latest version.Version) {
	// get current version
	// ask server if it needs to update
	pkg, err := latestPkg(url)
	if err != nil {
		return current
	}
	if pkg.Version > current {
		// download the update
		// extract and install the update, replacing current app
		// exec to new version
	}
	// return version
	return pkg.Version
}

func latestPkg(url *string) (latest version.PackageMeta, err error) {
	res, err := http.Get(routes.Latest(url))
	if err != nil {
		log.Printf("ERROR: problem contacting update server: %v", err)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Printf("ERROR: problem reading response body: %v", err)
		return
	}
	err = json.Unmarshal(body, &latest)
	if err != nil {
		log.Printf("ERROR: problem parsing version from response body: %v", err)
		return
	}
	return
}
