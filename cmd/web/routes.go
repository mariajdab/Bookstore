package main

import "net/http"
import "github.com/justinas/alice"
import "github.com/bmizerany/pat"



func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.closeConnection, app.logRequest, setHeaders)

	mux := pat.New()

	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/book/create", http.HandlerFunc(app.createBookForm))
	mux.Post("/book/create", http.HandlerFunc(app.createBook))
	mux.Get("/book/:id", http.HandlerFunc(app.showBook))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
