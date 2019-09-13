package update

import (
	"errors"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
)

type FakeStore struct {
	rc  io.ReadCloser
	err error
}

func (fs FakeStore) List() (io.ReadCloser, error) {
	return fs.rc, fs.err
}
func (fs FakeStore) Download(*string) (io.ReadCloser, error) {
	return fs.rc, fs.err
}

type RR func(io.Reader) io.Reader

func nop(r io.Reader) io.Reader {
	return r
}

func NewFS(r RR, s string, err error) FakeStore {
	return FakeStore{ioutil.NopCloser(r(strings.NewReader(s))), err}
}

func TestGetPackageList(t *testing.T) {
	tables := []struct {
		fs  FakeStore
		out []string
	}{
		{NewFS(nop, `[]`, nil), []string{}},
		{NewFS(nop, `["file"]`, nil), []string{"file"}},
		{NewFS(nop, `["file", "", "file2"]`, nil), []string{"file", "", "file2"}},
		{NewFS(nop, "[asdf]", errors.New("error")), []string(nil)},
		{NewFS(nop, "{}", errors.New("error")), []string(nil)},
		{NewFS(nop, `{"json": "obj"}`, errors.New("error")), []string(nil)},
		{NewFS(iotest.DataErrReader, `["this","is","otherwise","valid"]`, errors.New("error")), []string(nil)},
		{NewFS(iotest.HalfReader, `["this","is","otherwise","valid"]`, errors.New("error")), []string(nil)},
	}
	for n, table := range tables {
		files, err := getPackageList(table.fs)
		if table.fs.err == nil && err != nil {
			t.Errorf("%d Expected an error", n)
		} else if table.fs.err != nil && err == nil {
			t.Errorf("%d Unexpected error, got error %v", n, err)
		}
		if !reflect.DeepEqual(files, table.out) {
			t.Errorf("%d Wrong result, got %#v, want %#v", n, files, table.out)
		}
	}
}
