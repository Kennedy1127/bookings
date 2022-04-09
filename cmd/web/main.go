package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Kennedy/bookings/internal/config"
	"github.com/Kennedy/bookings/internal/handlers"
	"github.com/Kennedy/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNum = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	render.NewTemplates(&app)

	fmt.Println("Starting application on port", portNum)

	server := &http.Server{
		Addr:    portNum,
		Handler: routes(&app),
	}
	err = server.ListenAndServe()
	log.Fatal(err)
}
