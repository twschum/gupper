/* update

Module handles updating the client from the server

*/
package update

import (
)

// Version is just an int for now, could easily be replaced by a semantic
// version object that could handle more realistic use-case
type Version int

func Check() (current Version) {
	current = 0
	// get current version
	// ask server if it needs to update
	// download the update
	// extract and install the update, replacing current app
	// exec to new version
	// return version
	return
}
