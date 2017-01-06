package handler

import (
	"birthday-bot/internal/repository"
	"database/sql"
	"html/template"
	"net/http"
	"time"
)

func GetListHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		birthdays, err := repository.GetAllBirthdays(db)
		panicIfError(err)
		t, err := template.New("list.html").Funcs(
			map[string]interface{}{
				"monthToInt": func(m time.Month) int {
					return int(m)
				},
			},
		).ParseFiles("./ui/html/list.html")
		panicIfError(err)
		err = t.Execute(w, birthdays)
		panicIfError(err)

	}
}
