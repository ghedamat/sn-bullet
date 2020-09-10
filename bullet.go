package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Page struct {
	Title string
	Body  []byte
}

type Action struct {
	Label        string   `json:"label"`
	Url          string   `json:"url"`
	Verb         string   `json:"verb"`
	Context      string   `json:"context"`
	ContentTypes []string `json:"content_types"`
	AccessType   string   `json:"access_type"`
}

type Response struct {
	Identifier     string   `json:"identifier"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	SupportedTypes []string `json:"supported_types"`
	ContentType    string   `json:"content_type"`
	Url            string   `json:"url"`
	Actions        []Action `json:"actions"`
}

func loadPage(title string) (*Page, error) {
	filename := title + ".md"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	fmt.Println(r.URL.Path)

	fmt.Fprintf(w, "Hi there, I love %s!\n", r.URL.Path[1:])
}

func install(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	ct := []string{"Note"}
	action := Action{"Bullet entry", "http://localhost:8080/bullet", "get", "Item", ct, "decrypted"}
	response := Response{
		"at.ghedam.sn-bullet-test",
		"sn-bullet-test",
		"sn bullet test description",
		[]string{"Note"},
		"Extension",
		"http://localhost:8080/install",
		[]Action{action}}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type Item struct {
	Uuid         uuid.UUID `json:"uuid"`
	ContentTypes string    `json:"content_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Content      string    `json:"content"`
}

func bullet(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	item := Item{
		uuid.New(),
		"Note",
		time.Now().UTC(),
		time.Now().UTC(),
		"dsaofodsasdf",
	}

	itemMap := make(map[string]Item)
	itemMap["item"] = item
	js, err := json.Marshal(itemMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func main() {
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))

	http.HandleFunc("/", handler)
	http.HandleFunc("/install", install)
	http.HandleFunc("/bullet", bullet)
	//http.HandleFunc("/bullet", bullet)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
