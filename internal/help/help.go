package help

import (
	"fmt"
	"net/http"
)

func Method(method string, handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
fmt.Println(r.Method, method)
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler(w, r)
	}
}
