package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

func (service *Service) handleIntro(w http.ResponseWriter, r *http.Request) {
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

// Maps /challenge/{id} to templates/challenge-{id}.html. It verifies a valid
// challenge id is requested and the user has access to it.
func (service *Service) handleChallenge(w http.ResponseWriter, r *http.Request) {
	session := getSession(r.Cookies())
	if session == nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	challengeIDText := r.PathValue("challengeID")
	challengeID, err := strconv.Atoi(challengeIDText)
	if err != nil {
		//TODO bad challange, just 404
		return
	}

	if challengeID > session.Level {
		//TODO don't have access to that challenge yet
		return
	}

	if challengeID < 1 || challengeID > 20 {
		//TODO bad challenge 404
		return
	}

	challengeTemplate := fmt.Sprintf("challenge-%d.html", challengeID)
	err = service.template.ExecuteTemplate(w, challengeTemplate, nil)
	if err != nil {
		log.Fatalf("Err executing template %s\n", err)
	}
}
