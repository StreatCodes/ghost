package main

import (
	"log"
	"net/http"
)

func (service *Service) handleIndex(w http.ResponseWriter, r *http.Request) {
	session := getSession(r.Cookies())

	//Existing session, display list of levels
	if session != nil {
		err := service.template.ExecuteTemplate(w, "index.html", session)
		if err != nil {
			log.Fatalf("Err executing template %s\n", err)
		}
		return
	}

	//New user, display the intro
	err := service.template.ExecuteTemplate(w, "intro.html", nil)
	if err != nil {
		log.Fatalf("Err executing template %s\n", err)
	}
}

func (service *Service) handleNewSession(w http.ResponseWriter, r *http.Request) {
	session := getSession(r.Cookies())
	if session != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sessionId, err := newSession()
	if err != nil {
		log.Fatalf("Error creating session: %s\n", err)
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: sessionId,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
