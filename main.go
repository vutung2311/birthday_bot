package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"birthday-bot/internal/handler"
	"birthday-bot/internal/routine"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

var (
	db  *sql.DB
	err error
)

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err = sql.Open("sqlite3", "database/db.sqlite3")
	if err != nil {
		panic(err)
	}
	panicIfError(goose.SetDialect("sqlite3"))
	panicIfError(goose.Up(db, "database/migrations"))
	defer func() {
		panicIfError(db.Close())
	}()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT is empty, listening on 8080 by default")
		port = "8080"
	}

	dingTalkAccessToken := os.Getenv("DING_TALK_ACCESS_TOKEN")
	if dingTalkAccessToken == "" {
		log.Fatal("$DING_TALK_ACCESS_TOKEN is required")
	}

	go routine.GetDingTalkReminderRoutine(db, dingTalkAccessToken)()

	mux := chi.NewMux()
	mux.Use(PanicHandler)
	mux.Handle("/favicon.ico", http.FileServer(http.Dir("./ui/statics/ico")))
	mux.HandleFunc("/login", handler.GetLoginHandler(db))
	mux.HandleFunc("/register", handler.GetRegisterHandler(db))
	mux.Handle("/ui/statics/*",
		http.StripPrefix("/ui/statics/", http.FileServer(http.Dir("./ui/statics"))),
	)
	mux.Group(func(r chi.Router) {
		r.Use(Authentication)
		r.Get("/", handler.IndexHandler)
		r.Get("/logout", handler.LogoutHandler)
		r.Get("/list", handler.GetListHandler(db))
		r.HandleFunc("/create", handler.GetCreateHandler(db))
		r.HandleFunc("/update", handler.GetUpdateHandler(db))
		r.HandleFunc("/delete", handler.GetDeleteHandler(db))
		r.HandleFunc("/import", handler.GetImportHandler(db))
		r.HandleFunc("/export", handler.GetExportHandler(db))
	})

	panicIfError(http.ListenAndServe(":"+port, mux))
}

func PanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if val := recover(); val != nil {
				http.Error(w, fmt.Sprintf("%v", val), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthenticated, err := handler.IsAuthenticated(r)
		panicIfError(err)
		if !isAuthenticated {
			http.Redirect(w, r, "/login", 301)
			return
		}
		next.ServeHTTP(w, r)
	})
}
