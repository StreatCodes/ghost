package main

import (
	"fmt"
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

	session, err := newSession()
	if err != nil {
		log.Fatalf("Error creating session: %s\n", err)
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: session.IDString(),
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
	challengeID, err := validChallenge(*session, challengeIDText)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	challengeTemplate := fmt.Sprintf("challenge-%d.html", challengeID)
	err = service.template.ExecuteTemplate(w, challengeTemplate, nil)
	if err != nil {
		log.Fatalf("Err executing template %s\n", err)
	}
}

// Basic endpoint to display the users input for the given challenge
func (service *Service) handleChallengeInput(w http.ResponseWriter, r *http.Request) {
	session := getSession(r.Cookies())
	if session == nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	challengeIDText := r.PathValue("challengeID")
	challengeID, err := validChallenge(*session, challengeIDText)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	templateData := generateChallengeInput(challengeID, [32]byte(session.ID))
	err = service.template.ExecuteTemplate(w, "challenge-input.html", templateData)
	if err != nil {
		log.Fatalf("Err executing template %s\n", err)
	}
}

// A POST endpoint to verify the user entered the correct answer
func (service *Service) handleChallengeAnswer(w http.ResponseWriter, r *http.Request) {
	session := getSession(r.Cookies())
	if session == nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	challengeIDText := r.PathValue("challengeID")
	challengeID, err := validChallenge(*session, challengeIDText)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	answer := r.FormValue("answer")

	correct := checkChallengeAnswer(challengeID, [32]byte(session.ID), answer)
	if !correct {
		//TODO make this nicer....
		http.Error(w, "Incorrect", http.StatusBadRequest)
		return
	}

	//TODO update the users level

	http.Redirect(w, r, "/challenge/"+challengeIDText, http.StatusSeeOther)
	//TODO add new code to previous page to indicate the challenge was successfully completed
	// This can be accomplished by checking if the users level is greater than the viewed level
	// append additional story?
}
