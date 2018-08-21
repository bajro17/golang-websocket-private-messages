package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"

	"./client"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool {
	return true
}}
var tpl *template.Template

var Store = sessions.NewCookieStore([]byte("secret-password"))

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {

	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/login", login).Methods("GET", "POST")
	r.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":8080", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	
	answers := []string{
		"Bajro",
		"Mujo",
		"Latif",
		"Osman",
	}
	
	c := client.NewClient(conn, answers[rand.Intn(len(answers))])

	go c.Listen()
	c.Join <- true

	go c.Read()

}

func login(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		tpl.ExecuteTemplate(w, "login.html", nil)

	case "POST":
		session, err := Store.Get(r, "session")

		if err != nil {
			fmt.Println("error identifying session")
			tpl.Execute(w, nil)
			return
		}
		r.ParseForm()

		session.Values["username"] = r.FormValue("username")

		session.Save(r, w)

		http.Redirect(w, r, "/login", 302)

	}
}
