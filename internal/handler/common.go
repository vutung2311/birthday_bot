package handler

import "net/http"

const (
	CookieNameForAuthentication = "ithinkimloggedin"
)

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func IsAuthenticated(r *http.Request) (bool, error) {
	c, err := r.Cookie(CookieNameForAuthentication)
	if err != nil && err != http.ErrNoCookie {
		return false, err
	}
	if err == http.ErrNoCookie || c == nil || c.Value != "true" {
		return false, nil
	}
	return true, nil
}
