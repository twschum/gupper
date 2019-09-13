/*

Handles updating the client from the server

Only part of the client that talks to the server, is simple enough
to be replaced with other server endpoints/services/technologies
implementing an interface for getting a list of available packages,
and downloading/transferring a specific package

*/
package update

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"syscall"

	"github.com/twschum/gupper/pkg/pkgmeta"
	"github.com/twschum/gupper/pkg/store"
	"github.com/twschum/gupper/pkg/version"
)

// Runs the update process, checking, downloading, then installing and restarting
func Check(repository store.Store, current *pkgmeta.PackageMeta) (latest version.Version, err error) {
	log.Println("Checking for updates")
	packageList, err := getPackageList(repository)
	if err != nil {
		return current.Version, err
	}
	pkg := latestPackage(current, packageList)
	if pkg.Version.GT(current.Version) {
		err = downloadPackage(repository, &pkg)
		if err != nil {
			return pkg.Version, err
		}
		err = installAndExec(&pkg)
		if err != nil {
			return pkg.Version, err
		}
	}
	return pkg.Version, nil
}

// Gets list of package files available on the server
func getPackageList(repository store.Store) (files []string, err error) {
	body, err := repository.List()
	if err != nil {
		return
	}
	defer body.Close()
	err = json.NewDecoder(body).Decode(&files)
	if err != nil {
		return
	}
	return files, nil
}

// Determines the latest appropriate package from the list of package names
func latestPackage(current *pkgmeta.PackageMeta, packageList []string) pkgmeta.PackageMeta {
	packages := pkgmeta.AvailablePackages(packageList)
	return pkgmeta.GetLatestPackage(packages, current)
}

// Download and save a package file from the update server
func downloadPackage(repository store.Store, pkg *pkgmeta.PackageMeta) (err error) {
	body, err := repository.Download(&pkg.Filename)
	if err != nil {
		return
	}
	defer body.Close()
	out, err := os.Create(pkg.Filename)
	if err != nil {
		return
	}
	defer out.Close()
	_, err = io.Copy(out, body)
	return err
}

// "install" by mving the downloaded file to the current app
// This is atomic but also not forgiving and doesn't provide any rollback options
// in the even of a "bad file" installed
// Here is where more complex package archive install functionality could be handled
func installAndExec(pkg *pkgmeta.PackageMeta) (err error) {
	// TODO syscall is frozen as of go 1.3, so for long-term maintenance need go.sys pkgs
	log.Printf("Installing %v to %v\n", pkg.Filename, os.Args[0])
	os.Chmod(pkg.Filename, 0755) // TODO err here?
	os.Rename(pkg.Filename, os.Args[0])
	log.Println("Restarting...")
	err = syscall.Exec(os.Args[0], os.Args, os.Environ())
	log.Println("Automatic restart failed, manual restart required")
	if err != nil {
		return err
	}
	return nil
}
