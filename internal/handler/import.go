package handler

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"birthday-bot/internal/model"
	"birthday-bot/internal/repository"
)

func GetImportHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		panicIfError(err)
		file, _, err := r.FormFile("birthdayJsonFile")
		panicIfError(err)
		defer func() {
			_ = file.Close()
		}()
		bytes, err := ioutil.ReadAll(file)
		panicIfError(err)

		var birthdays []model.PersonBirthday
		err = json.Unmarshal(bytes, &birthdays)
		panicIfError(err)
		tx, err := db.Begin()
		panicIfError(err)
		panicIfError(repository.DeleteAllBirthdays(db))
		for _, birthday := range birthdays {
			err := repository.SaveBirthday(db, &birthday)
			if err != nil {
				panicIfError(tx.Rollback())
				panic(err)
			}
		}
		panicIfError(tx.Commit())
		http.Redirect(w, r, "/", 301)
	}
}
