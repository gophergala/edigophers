package main

import (
	"edigophers/common"
	"edigophers/user"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/list", listHandler)
	r.HandleFunc("/logout", logoutHandler)
	r.HandleFunc("/profile/{username}", profileHandler)
	r.HandleFunc("/profile", profileHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

var templates = template.Must(template.ParseFiles(
	"tpl/header.html", "tpl/footer.html",
	"tpl/home.html", "tpl/login.html", "tpl/profile.html", "tpl/list.html"))

type page struct {
	Title string
	User  *user.User
	Data  interface{}
}

func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := user.GetSessionUser(w, r)
	if err != nil {
		return
	}
	display(w, "home", &page{Title: "Home", User: user})
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "login", &page{Title: "Login"})
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	if user.SetSessionUser(w, r, username) != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	common.CheckError(user.LogOutSessionUser(w, r))
	http.Redirect(w, r, "/login", http.StatusFound)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	usr, err := user.GetSessionUser(w, r)
	if err != nil {
		return
	}
	display(w, "list", &page{Title: "List of users", User: usr, Data: user.GetUsers()})
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	// Display the profile of somebody else
	vars := mux.Vars(r)
	if username, ok := vars["username"]; ok {
		user, err := user.GetUser(username)
		if err != nil {
			return
		}
		display(w, "profile", &page{Title: fmt.Sprintf("%s's Profile", user.Name), User: user})
		return
	}

	// Display the profile of the current user
	user, err := user.GetSessionUser(w, r)
	if err != nil {
		log.Fatal(err)
		return
	}
	display(w, "profile", &page{Title: "Your Profile", User: user})
}
