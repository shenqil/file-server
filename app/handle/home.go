package handle

import (
	"fileServer/app/view"
	"net/http"
)

// Home handler
func Home(w http.ResponseWriter, r *http.Request) {
	view.Home(w, r)
}
