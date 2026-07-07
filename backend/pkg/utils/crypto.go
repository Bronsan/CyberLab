package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// HashFlag computes SHA-256 HMAC of a flag with a per-challenge salt.
// The salt is derived from the challenge ID to ensure each challenge uses
// a unique key, preventing cross-challenge hash comparison.
func HashFlag(flag string, challengeID uint, globalSalt string) string {
	// challenge-specific key = HMAC(globalSalt, challengeID)
	mac := hmac.New(sha256.New, []byte(globalSalt))
	mac.Write([]byte(fmt.Sprintf("%d", challengeID)))
	challengeKey := hex.EncodeToString(mac.Sum(nil))

	// hash the flag with the challenge-specific key
	mac2 := hmac.New(sha256.New, []byte(challengeKey))
	mac2.Write([]byte(flag))
	return hex.EncodeToString(mac2.Sum(nil))
}

// VerifyFlag compares a submitted flag against a stored hash.
func VerifyFlag(submittedFlag string, storedHash string, challengeID uint, globalSalt string) bool {
	computed := HashFlag(submittedFlag, challengeID, globalSalt)
	return hmac.Equal([]byte(computed), []byte(storedHash))
}

// HashString returns a simple SHA-256 hash (for general use, not passwords).
func HashString(input string) string {
	h := sha256.Sum256([]byte(input))
	return hex.EncodeToString(h[:])
}
