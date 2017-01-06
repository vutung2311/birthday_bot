package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"birthday-bot/internal/repository"
)

func GetDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Redirect(w, r, "/", 301)
		}
		birthdayId, err := strconv.ParseInt(r.FormValue("Id"), 10, 64)
		panicIfError(err)
		panicIfError(repository.DeleteBirthday(db, birthdayId))
		http.Redirect(w, r, "/", 301)
	}
}
