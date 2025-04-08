package sharedModel

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"

	"github.com/pkg/errors"
)

// ValidateEmail checks if the given email address is in a valid format
func ValidateEmail(email string) bool {
	// Regular expression for validating email format
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex
	re := regexp.MustCompile(emailRegex)

	// Match the email against the regex
	return re.MatchString(email)
}

func RandomStr(length int) (string, error) {
	var b = make([]byte, length)

	_, err := rand.Read(b)

	if err != nil {
		return "", errors.WithStack(err)
	}

	return hex.EncodeToString(b), nil
}
