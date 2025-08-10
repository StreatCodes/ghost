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
	err := initSessions()
	if err != nil {
		log.Fatalf("Failed to create sessions directory %s\n", err)
	}

	tmpl, err := template.ParseGlob("./templates/*")
	if err != nil {
		log.Fatalf("Error parsing templates %s\n", err)
	}

	service := Service{
		template: tmpl,
	}

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("GET /", service.handleIndex)
	http.HandleFunc("GET /intro", service.handleIntro)
	http.HandleFunc("GET /new-session", service.handleNewSession)
	http.HandleFunc("GET /challenge/{challengeID}", service.handleChallenge)
	http.HandleFunc("GET /challenge/{challengeID}/input", service.handleChallengeInput)
	http.HandleFunc("POST /challenge/{challengeID}", service.handleChallengeAnswer)

	log.Println("Server starting on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
