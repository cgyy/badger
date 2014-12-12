package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/GeertJohan/go.rice"
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

func render(templateName string, w http.ResponseWriter,
	r *http.Request, binding map[string]interface{}) {

	session, _ := store.Get(r, "auth")

	if flash, ok := session.Values["Flash"].(string); ok {
		binding["Flash"] = flash
		setSession("", w, r)
	}

	if flash, ok := session.Values["Error"].(string); ok {
		binding["Error"] = flash
		setSession("", w, r, "Error")
	}

	binding["Token"] = nosurf.Token(r)

	if loggedIn(r) {
		binding["CurrentUser"] = currentUser(r)
	}

	templateBox, err := rice.FindBox("../templates")

	if err != nil {
		log.Fatal(err)
	}

	lstr, _ := templateBox.String("layout.tmpl")
	tstr, _ := templateBox.String(fmt.Sprintf("%s.tmpl", templateName))
	lstr += tstr

	templateBox.Walk("includes", func(path string, info os.FileInfo, err error) error {
		include, _ := templateBox.String(path)
		lstr += include
		return nil
	})

	var t *template.Template
	t, err = template.New("layout").Parse(lstr)

	if err != nil {
		log.Fatal(err)
	}

	err = t.ExecuteTemplate(w, "base", binding)

	if err != nil {
		log.Print("Error executing template: ", err)
	}
}

func authorize(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "auth")

	if session.Values["Email"] == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return false
	}

	return true
}

func currentUser(r *http.Request) models.User {
	session, _ := store.Get(r, "auth")
	user, _ := models.FindUser(session.Values["Email"].(string))
	return user
}

func loggedIn(r *http.Request) bool {
	session, _ := store.Get(r, "auth")
	return session.Values["Email"] != nil
}

func setSession(message interface{}, w http.ResponseWriter,
	r *http.Request, key ...string) {

	k := "Flash"

	if len(key) > 0 {
		k = key[0]
	}

	session, _ := store.Get(r, "auth")
	session.Values[k] = message
	session.Save(r, w)
}

func usernameConflictsWithRoute(username string) bool {
	for _, route := range routes() {
		if strings.Contains(route.Path, username) {
			return true
		}
	}

	return false
}
