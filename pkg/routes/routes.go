/*

Constants file for the route methods

*/
package routes

const (
	LATEST   = "/latest"
	DOWNLOAD = "/download"
	UPLOAD   = "/upload"
)

func Latest(url *string) (route string) {
	return *url + LATEST
}

func Download(url *string) (route string) {
	return *url + DOWNLOAD
}

func Upload(url *string) (route string) {
	return *url + UPLOAD
}
