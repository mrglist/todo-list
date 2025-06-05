package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) serverErr(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientErr(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientErr(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverErr(w, err)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverErr(w, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:      time.Now().Year(),
		Flash:            app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthencticated: app.isAuthencticated(r),
		CSRFToken:        nosurf.Token(r),
	}
}

func (app *application) isAuthencticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthencticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
