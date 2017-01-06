package handler

import "net/http"

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   CookieNameForAuthentication,
		MaxAge: -1}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", 301)
}
