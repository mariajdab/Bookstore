package main

import (
	"bookStore/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.books.ListBooks()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Books: s}

	app.render(w, r, "home.page.tmpl", data)

}

func (app *application) showBook(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}


	s, err := app.books.Get(id)

	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Book: s}


	app.render(w, r, "show.page.tmpl", data)
}

// Add a new createSnippetForm handler, which for now returns a placeholder res
func (app *application) createBookForm(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Create a new snippet..."))
	app.render(w, r, "create.page.tmpl", nil)
}

func (app *application) createBook(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	errors := make(map[string]string)



	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}


	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	if len(errors) > 0 {
		fmt.Fprint(w, errors)
		return
	}

	id, err := app.books.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}


	http.Redirect(w, r, fmt.Sprintf("/book/%d", id), http.StatusSeeOther)
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)

}
