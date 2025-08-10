package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

var ErrChallengeNotFound = errors.New("challenge not found")

func validChallenge(session Session, challengeIDText string) (int, error) {
	challengeID, err := strconv.Atoi(challengeIDText)
	if err != nil {
		return 0, ErrChallengeNotFound
	}

	if challengeID > session.Level {
		return 0, ErrChallengeNotFound
	}

	if challengeID < 1 || challengeID > 20 {
		return 0, ErrChallengeNotFound
	}

	return challengeID, nil
}

func generateChallengeInput(challengeId int, seed [32]byte) string {
	switch challengeId {
	case 1:
		return challenge1Input(seed)
	default:
		return "Error"
	}
}

func checkChallengeAnswer(challengeId int, seed [32]byte, answer string) bool {
	switch challengeId {
	case 1:
		input := challenge1Input(seed)
		return checkChallenge1(input, answer)
	default:
		return false
	}
}

func challenge1Input(seed [32]byte) string {
	allowedCharacters := []rune("1234567890abcdefghijklmnopqrstuvwxyz!@#$%^&*()_+-=`~<>,./?;:]}[{")

	rng := rand.NewChaCha8(seed)

	input := ""
	for range 20 {
		for range 60 {
			index := rng.Uint64() % uint64(len(allowedCharacters))
			input += string(allowedCharacters[index])
		}

		input += "\n"
	}

	return input
}

func checkChallenge1(input, answerText string) bool {
	answer, err := strconv.Atoi(answerText)
	if err != nil {
		return false
	}

	//Solve the challenge
	numbers := "123456789"
	letters := "abcdefghijklmnopqrstuvwxyz"

	total := 0
	for _, r := range input {
		if strings.ContainsRune(numbers, r) {
			decimal := int(r) - 48
			total += decimal
			continue
		}
		if strings.ContainsRune(letters, r) {
			decimal := int(r) - 96
			total += decimal
		}
	}

	fmt.Println("Total", total)
	fmt.Println("Answer", answer)

	return total == answer
}
