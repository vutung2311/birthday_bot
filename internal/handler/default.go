package handler

import "net/http"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/list", 301)
}
