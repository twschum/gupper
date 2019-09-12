/*

Module handles updating the client from the server

*/
package update

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"syscall"

	"github.com/twschum/gupper/pkg/pkgmeta"
	"github.com/twschum/gupper/pkg/routes"
	"github.com/twschum/gupper/pkg/version"
)

func Check(current *pkgmeta.PackageMeta, base *url.URL) (latest version.Version, err error) {
	// get current version
	// ask server if it needs to update
	log.Println("Checking for updates")
	packageList, err := getPackageList(base)
	if err != nil {
		return current.Version, err
	}
	pkg := latestPackage(current, packageList)
	if pkg.Version.GT(current.Version) {
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

// gets list of package files available on the server
func getPackageList(base *url.URL) (files []string, err error) {
	res, err := http.Get(routes.List(*base))
	if err != nil {
		log.Printf("ERROR: problem contacting update server: %v", err)
		return
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&files)
	if err != nil {
		log.Printf("ERROR: problem reading response body: %v", err)
		return
	}
	return files, nil
}

// get the latest appropriate PackageMeta from the list of package names
func latestPackage(current *pkgmeta.PackageMeta, packageList []string) pkgmeta.PackageMeta {
	packages := pkgmeta.AvailablePackages(packageList)
	return pkgmeta.GetLatestPackage(packages, current)
}

// server/download/
func downloadPackage(base *url.URL, pkg *pkgmeta.PackageMeta) (err error) {
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

func installAndExec(pkg *pkgmeta.PackageMeta) (err error) {
	// TODO syscall is frozen as of go 1.3, so for long-term maintance need go.sys pkgs
	os.Chmod(pkg.Filename, 0755)
	// "install" by mving the downloaded file to the current app
	// This is atomic but also not forgiving and doens't provide any rollback options
	// in the even of a "bad file" installed
	log.Printf("Installing %v to %v\n", pkg.Filename, os.Args[0])
	os.Rename(pkg.Filename, os.Args[0])
	log.Printf("Restarting...")
	err = syscall.Exec(os.Args[0], os.Args, os.Environ())
	if err != nil {
		return err
	}
	return nil
}
