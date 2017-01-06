package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"birthday-bot/internal/repository"
)

func GetExportHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		birthdays, err := repository.GetAllBirthdays(db)
		panicIfError(err)
		b, err := json.Marshal(birthdays)
		filename := "BirthdayExport.json"
		fileSize := len(string(b))

		w.Header().Set("Content-Type", http.DetectContentType(b))
		w.Header().Set("Content-Disposition", "attachment; filename="+filename+"")
		w.Header().Set("Expires", "0")
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Content-Length", strconv.Itoa(fileSize))
		w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

		http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(b))
	}
}
