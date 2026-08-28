package quiz

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
)

const verificationSaltBytes = 16

func AddVerificationData(quizID int64, questions []Question) error {
	for index := range questions {
		salt := make([]byte, verificationSaltBytes)
		if _, err := rand.Read(salt); err != nil {
			return fmt.Errorf("generate answer salt: %w", err)
		}

		questions[index].VerificationSalt = base64.RawURLEncoding.EncodeToString(salt)
		value := "v1:" +
			strconv.FormatInt(quizID, 10) + ":" +
			strconv.Itoa(index) + ":" +
			strconv.Itoa(questions[index].CorrectIndex) + ":" +
			questions[index].VerificationSalt
		digest := sha256.Sum256([]byte(value))
		questions[index].CorrectAnswerHash = hex.EncodeToString(digest[:])
	}
	return nil
}
