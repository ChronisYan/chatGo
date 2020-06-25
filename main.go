package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	tpl      *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.tpl = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.tpl.Execute(w, r)
}

func main() {
	port := flag.String("port", ":3000", "Port of the application")
	flag.Parse()
	room := newRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", room)

	go room.run() // start new room

	log.Println("Server is up on Port: ", *port)
	if err := http.ListenAndServe(*port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
