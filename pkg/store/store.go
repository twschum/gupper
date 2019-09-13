/*

Interface to some system which provides a list of files,
and allows downloading said file

*/

package store

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/twschum/gupper/pkg/routes"
)

type Store interface {
	// where body is a JSON-encoded array of strings
	List() (body io.ReadCloser, err error)
	// where body is a file
	Download(filename *string) (body io.ReadCloser, err error)
}

type Fileserver struct {
	Base url.URL
}

func (fs Fileserver) List() (body io.ReadCloser, err error) {
	return checkedGet(routes.List(fs.Base))
}

func (fs Fileserver) Download(filename *string) (body io.ReadCloser, err error) {
	return checkedGet(routes.Download(fs.Base, *filename))
}

func checkedGet(request string) (body io.ReadCloser, err error) {
	log.Println("GET from", request)
	res, err := http.Get(request)
	if err != nil {
		return
	} else if res.StatusCode != 200 {
		err = errors.New(res.Status)
	}
	return res.Body, err
}
