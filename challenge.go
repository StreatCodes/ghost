package main

import (
	"math/rand/v2"
)

type ChallengeData struct {
	Input string
}

func generateChallengeInput(challengeId int, seed [32]byte) ChallengeData {
	var input string

	switch challengeId {
	case 1:
		input = challange1Input(seed)
	}

	return ChallengeData{
		Input: input,
	}
}

func challange1Input(seed [32]byte) string {
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
