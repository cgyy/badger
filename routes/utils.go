package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/sessions"
	"github.com/justinas/nosurf"
	"github.com/tristanoneil/badger/models"
)

var (
	store *sessions.CookieStore
)

func init() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
}

func render(name string, w http.ResponseWriter, r *http.Request,
	data ...map[string]interface{}) {

	tmpl := fmt.Sprintf("templates/%s.html", name)

	if tmpl == "" {
		log.Print("Missing template:", name)
	}

	d := map[string]interface{}{}

	if len(data) > 0 {
		d = data[0]
	}

	session, _ := store.Get(r, "auth")

	if str, ok := session.Values["Flash"].(string); ok {
		d["Flash"] = str
		setSession("", w, r)
	}

	d["Token"] = nosurf.Token(r)

	err := template.
		Must(template.ParseFiles(tmpl, "templates/base.html")).
		ExecuteTemplate(w, "base", d)

	if err != nil {
		log.Print("Template executing error: ", err)
	}
}

func authorize(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth")

	if session.Values["Email"] == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}

func currentUserID(r *http.Request) int {
	session, _ := store.Get(r, "auth")

	if str, ok := session.Values["Flash"].(string); ok {
		return models.GetUserIDForEmail(str)
	}

	return 0
}

func setSession(message string, w http.ResponseWriter,
	r *http.Request, key ...string) {

	k := "Flash"

	if len(key) > 0 {
		k = key[0]
	}

	session, _ := store.Get(r, "auth")
	session.Values[k] = message
	session.Save(r, w)
}