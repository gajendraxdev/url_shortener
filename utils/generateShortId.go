package short_id

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM-+"

func GenerateShortId(length int) (string, error) {
	result := make([]byte, length)

	// converting the length to big int (crypto/rand requirement)
	alphabetLen := big.NewInt(int64(len(alphabet)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}

		result[i] = alphabet[num.Int64()]
	}

	return string(result), nil
}
