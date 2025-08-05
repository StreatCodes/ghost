package main

import (
	"html/template"
	"log"
	"net/http"
)

type Service struct {
	template *template.Template
}

func main() {
	tmpl, err := template.ParseGlob("./templates/*")
	if err != nil {
		log.Fatalf("Error parsing templates %s", err)
	}

	service := Service{
		template: tmpl,
	}

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("GET /", service.handleIndex)
	http.HandleFunc("GET /new-session", service.handleNewSession)

	log.Println("Server starting on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
