package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AugustSerenity/booking/internal/config"
	"github.com/AugustSerenity/booking/internal/driver"
	"github.com/AugustSerenity/booking/internal/handlers"
	"github.com/AugustSerenity/booking/internal/helpers"
	"github.com/AugustSerenity/booking/internal/models"
	"github.com/AugustSerenity/booking/internal/render"
	"github.com/alexedwards/scs/v2"
)

const numberPort = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Application started on port %s", numberPort))

	srv := &http.Server{
		Addr:    numberPort,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})

	// change this to true in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5433 dbname=bookings user=postgres password=")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	defer db.SQL.Close()

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err)
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewTemplate(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
