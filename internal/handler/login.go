package handler

import (
	"database/sql"
	"net/http"

	"birthday-bot/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func GetLoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val, err := IsAuthenticated(r)
		panicIfError(err)
		if val {
			http.Redirect(w, r, "/", 301)
		}
		if r.Method != "POST" {
			http.ServeFile(w, r, "./ui/html/login.html")
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		rememberMe := r.FormValue("rememberMe")
		user, err := repository.GetUserByUserName(db, username)
		panicIfError(err)
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Redirect(w, r, "/login", 301)
		}

		c := &http.Cookie{
			Name:  CookieNameForAuthentication,
			Value: "true",
		}
		if len(rememberMe) > 0 {
			c.MaxAge = 7200
		}
		http.SetCookie(w, c)
		http.Redirect(w, r, "/list", 301)
	}
}
