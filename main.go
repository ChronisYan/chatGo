package main

import (
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
	t.tpl.Execute(w, nil)
}

func main() {
	room := newRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", room)

	go room.run() // start new room

	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
