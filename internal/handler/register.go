package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"birthday-bot/internal/model"
	"birthday-bot/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func GetRegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.ServeFile(w, r, "./ui/html/register.html")
			return
		}
		// grab user info
		username := r.FormValue("username")
		password := r.FormValue("password")
		role, err := strconv.ParseInt(r.FormValue("role"), 10, 64)
		panicIfError(err)
		user, err := repository.GetUserByUserName(db, username)
		switch {
		case err == sql.ErrNoRows:
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			panicIfError(err)
			panicIfError(repository.SaveUser(
				db,
				model.User{
					Username: username,
					Password: string(hashedPassword),
					Role:     role,
				},
			))
		case err == nil && user.Id != 0:
			http.Error(w, "user is already existed", http.StatusBadRequest)
			return
		case err != nil:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		default:
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		}
	}

}
