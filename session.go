package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type Session struct {
	ID []byte `json:"-"`
	// The level the user has unlocked
	Level int
}

func (session Session) IDString() string {
	return hex.EncodeToString(session.ID)
}

func loadSession(sessionId string) (*Session, error) {
	file, err := os.Open("./sessions/" + sessionId)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var session Session
	err = json.NewDecoder(file).Decode(&session)
	if err != nil {
		return nil, err
	}

	session.ID, err = hex.DecodeString(sessionId)
	if err != nil {
		return nil, err
	}

	return &session, err
}

func validateSessionId(sessionId string) bool {
	if len(sessionId) != 64 {
		return false
	}

	validCharacters := "0123456789abcdef"
	for _, c := range sessionId {
		if !strings.ContainsRune(validCharacters, c) {
			return false
		}
	}

	return true
}

func getSession(cookies []*http.Cookie) *Session {
	for _, cookie := range cookies {
		if cookie.Name == "session" {
			sessionId := cookie.Value
			if !validateSessionId(sessionId) {
				log.Printf("Invalid session cookie value: \"%s\"", sessionId)
				continue
			}

			session, err := loadSession(sessionId)
			if err != nil {
				log.Printf("Error decoding session %s\n", err)
				return nil
			}

			return session
		}
	}

	return nil
}

func newSession() (*Session, error) {
	session := Session{Level: 1, ID: make([]byte, 32)}
	_, _ = rand.Read(session.ID)

	file, err := os.Create("./sessions/" + session.IDString())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	json.NewEncoder(file).Encode(Session{Level: 1})
	return &session, nil
}
