/* update

Module handles updating the client from the server

*/
package update

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/twschum/gupper/pkg/routes"
)

// Version is just a number for now, could easily be replaced by a semantic
// version object that could handle more realistic use-case
type Version float64

func ParseVersion(s string) (Version, error) {
	v, err := strconv.ParseFloat(s, 64)
	return Version(v), err
}


// embedded in the code here for now,
// maybe as a command line argument later at build time
const current Version = 1.0
const bad Version = 0

func Check(url *string) (latest Version) {
	// get current version
	// ask server if it needs to update
	latest, err := latestVersion(url)
	if err != nil {
		return current
	}
	if latest > current {
		// download the update
		// extract and install the update, replacing current app
		// exec to new version
	}
	// return version
	return latest
}

func latestVersion(url *string) (latest Version, err error) {
	res, err := http.Get(routes.Latest(url))
	if err != nil {
		log.Printf("ERROR: problem contacting update server: %v", err)
		return bad, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Printf("ERROR: problem reading response body: %v", err)
		return bad, err
	}
	latest, err = ParseVersion(string(body))
	if err != nil {
		log.Printf("ERROR: problem parsing version from response body: %v", err)
		return bad, err
	}
	return latest, err
}
