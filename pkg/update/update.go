/*

Module handles updating the client from the server

*/
package update

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/twschum/gupper/pkg/routes"
	"github.com/twschum/gupper/pkg/version"
)

func Check(current version.Version, base *url.URL) (latest version.Version) {
	// get current version
	// ask server if it needs to update
	pkg, err := latestPkg(base)
	if err != nil {
		return current
	}
	if pkg.Version > current {
		// download the update to pkg.Filename
		err = downloadPackage(base, &pkg)
		if err != nil {
		}
		// extract and install the update, replacing current app
		// exec to new version
	}
	// return version
	return pkg.Version
}

func latestPkg(base *url.URL) (latest version.PackageMeta, err error) {
	res, err := http.Get(routes.Latest(*base))
	if err != nil {
		log.Printf("ERROR: problem contacting update server: %v", err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
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

// server/download/1.3
func downloadPackage(base *url.URL, pkg *version.PackageMeta) (err error) {
	download := routes.Download(*base, pkg.Filename)
	log.Printf("Downloading latest package version %v from %v", pkg.Version, download)
	res, err := http.Get(download)
	if err != nil {
		log.Printf("ERROR: problem contacting update server: %v", err)
		return
	}
	defer res.Body.Close()

	out, err := os.Create(pkg.Filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	return err
}

//func installAndExec(pkg *version.PackageMeta) (err error) { }
