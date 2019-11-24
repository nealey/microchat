package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	When int64
	Who  string
	Text string
}

type Forum struct {
	Name  string
	Log   []Message
	SendQ chan Message
	LogQ  chan []Message
}

var foraByName map[string]*Forum

func (f *Forum) HandleSends() {
	for m := range f.SendQ {
		f.Log = append(f.Log, m)
	}
}

func (f *Forum) HandleReads() {
	for {
		f.LogQ <- f.Log
	}
}

func newForum(name string) *Forum {
	f := &Forum{
		Name:  name,
		SendQ: make(chan Message, 10),
		LogQ:  make(chan []Message, 10),
	}
	go f.HandleSends()
	go f.HandleReads()

	foraByName[name] = f
	return f
}

func sayHandler(w http.ResponseWriter, r *http.Request) {
	message := Message{
		When: time.Now().Unix(),
		Who:  r.FormValue("who"),
		Text: r.FormValue("text"),
	}
	forumName := r.FormValue("forum")
	f, ok := foraByName[forumName]
	if !ok {
		f = newForum(forumName)
	}
	f.SendQ <- message

	w.WriteHeader(http.StatusOK)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	forumName := r.FormValue("forum")
	entriesStr := r.FormValue("entries")

	if entriesStr == "" {
		entriesStr = "0"
	}
	entries, err := strconv.Atoi(entriesStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if entries <= 0 {
		entries = 20
	} else if entries > 500 {
		entries = 500
	}

	f, ok := foraByName[forumName]
	if !ok {
		http.NotFound(w, r)
		return
	}

	log := <-f.LogQ

	pos := len(log) - entries
	if pos < 0 {
		pos = 0
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	e := json.NewEncoder(w)
	e.Encode(f.Log[pos:])
}

func main() {
	http.HandleFunc("/say", sayHandler)
	http.HandleFunc("/read", readHandler)
	http.Handle("/", http.FileServer(http.Dir("static/")))

	foraByName = map[string]*Forum{}
	f := newForum("")
	f.Log = []Message{
		{
			When: time.Now().Unix(),
			Who:  "(system)",
			Text: "Welcome to μChat",
		},
	}

	bind := ":8080"
	log.Printf("Listening on %s", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}
