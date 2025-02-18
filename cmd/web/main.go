package main

import (
	"fmt"
	"github.com/AugustSerenity/booking/pkg/config"
	"github.com/AugustSerenity/booking/pkg/handlers"
	"github.com/AugustSerenity/booking/pkg/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const numberPort = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// change this to true in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplate(&app)

	fmt.Println(fmt.Sprintf("Application started on port %s", numberPort))

	srv := &http.Server{
		Addr:    numberPort,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
