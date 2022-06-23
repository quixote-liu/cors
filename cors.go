package cors

import "net/http"

var cors = New()

func Handler(w http.ResponseWriter, r *http.Request) (status int, isOption bool) {
	return cors.Handler(w, r)
}

func WrapH(h http.HandlerFunc) http.HandlerFunc {
	return cors.WrapH(h)
}
