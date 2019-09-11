/*

Module handles updating the client from the server

*/
package update

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"syscall"

	"github.com/twschum/gupper/pkg/routes"
	"github.com/twschum/gupper/pkg/version"
)

func Check(current version.Version, base *url.URL) (latest version.Version, err error) {
	// get current version
	// ask server if it needs to update
	log.Println("Checking for updates")
	pkg, err := latestPkg(base)
	if err != nil {
		return current, err
	}
	log.Printf("%v is the latest\n", pkg.Version)
	if pkg.Version.GT(current) {
		// download the update to pkg.Filename
		err = downloadPackage(base, &pkg)
		if err != nil {
			return pkg.Version, err
		}
		// extract and install the update, replacing current app
		// exec to new version
		err = installAndExec(&pkg)
		if err != nil {
			return pkg.Version, err
		}

	}
	// return version
	return pkg.Version, nil
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
	//log.Printf("%#v\n", res)
	if err != nil {
		log.Printf("ERROR: problem contacting update server: %v", err)
		return
	} else if res.StatusCode != 200 {
		err = errors.New(res.Status)
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

func installAndExec(pkg *version.PackageMeta) (err error) {
	// TODO actual install, for now, just try the exec of the package
	// TODO syscall is frozen as of go 1.3, need to figure out go.sys instead
	os.Chmod(pkg.Filename, 0755)
	os.Args[0] = pkg.Filename
	log.Printf("Calling execve with %#v\n", os.Args)
	err = syscall.Exec(os.Args[0], os.Args, os.Environ())
	log.Println("This is the line after Exec")
	if err != nil {
		return err
	}
	return nil
}
