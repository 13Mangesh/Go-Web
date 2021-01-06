package main

import (
	"context"
	"github.com/go-redis/redis"
	"html/template"
	"github.com/gorilla/mux"
	"net/http"
)

var client *redis.Client
var templates *template.Template
var ctx = context.Background()

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/", indexPostHandler).Methods("POST")
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	comments, err := client.LRange(ctx,"comments", 0, 10).Result()
	if err != nil {
		return
	}
	templates.ExecuteTemplate(w, "index.html", comments)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	client.LPush(ctx, "comments", comment)
	http.Redirect(w, r, "/", 302)
}

