package handler

import (
	"database/sql"
	"net/http"

	"birthday-bot/internal/model"
	"birthday-bot/internal/pkg/birthday"
	"birthday-bot/internal/repository"
)

func GetCreateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Redirect(w, r, "/", 301)
		}
		var (
			bd  model.PersonBirthday
			err error
		)
		bd.PersonName = r.FormValue("PersonName")
		bd.Birthday, err = birthday.ParseBirthday(r.FormValue("Birthday"))
		panicIfError(err)
		panicIfError(repository.SaveBirthday(db, &bd))
		http.Redirect(w, r, "/", 301)
	}
}
