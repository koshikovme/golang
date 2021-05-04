package main

import (
	"errors"
	"fmt"
	"html/template"
	// "html/template"
	"net/http"
	"strconv"

	"github.com/DataDavD/snippetbox/pkg/models"
)

// Define a home handler func which writes a byte slice containing
// "Hello from Snippetbox" as resp body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the curr req path exactly matches "/". If it doesn't, use
	// the http.NotFound() func to send 404 resp.
	if r.URL.Path != "/" {
		app.notFound(w)
		// return from func to avoid proceeding to home page response
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	// Create an instance of a templateData struct holding the slice of snippets.
	data := &templateData{Snippets: s}

	files := []string{
		"./ui/html/home.page.gohtml",
		"./ui/html/base.layout.gohtml",
		"./ui/html/footer.partial.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Pass in the templateData struct when executing the template.
	if terr := ts.Execute(w, data); terr != nil {
		app.serverError(w, err)
	}

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract id value from query string and try to convert it to an integer
	// using strconv.Atoi() func. If it can't be converted, or the value is less than 1,
	// we return a 404 page Not Found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a specific record
	// based on its ID. If no matching record is found, return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Create an instance of a templateData struct holding the snippet data.
	data := &templateData{Snippet: s}

	// Initialize a slice containing the paths to the show.page.gohtml file,
	// plus the base layout and footer partial templates
	files := []string{
		"./ui/html/show.page.gohtml",
		"./ui/html/base.layout.gohtml",
		"./ui/html/footer.partial.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// And then execute the parsed templates. Notice how we are passing in the snippet data
	// (a models.Snippet struct) as the final param.
	if err := ts.Execute(w, data); err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// only allow Post
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Create some variables holding dummy data. We'll remove these later on during the build
	title := "DataDavD Awesome Adventures in Life"
	content := "DataDavD has had an awesome, super, crazy, cool life!!!"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the created snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
