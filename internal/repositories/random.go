package repositories

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

// On utilise crypto/rand (et pas math/rand) car ces valeurs servent à la sécurité :
// elles ne doivent pas être prévisibles.

const referenceAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randomString tire "length" caractères au hasard dans l'alphabet donné
func randomString(alphabet string, length int) (string, error) {
	characters := make([]byte, length)
	for position := range characters {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		characters[position] = alphabet[randomIndex.Int64()]
	}
	return string(characters), nil
}

// NewReference génère un identifiant métier, ex : NewReference("PDT") -> "PDT-7D2K8N"
func NewReference(prefix string) (string, error) {
	suffix, err := randomString(referenceAlphabet, 6)
	if err != nil {
		return "", err
	}
	return prefix + "-" + suffix, nil
}

// NewCode génère un code numérique à 6 chiffres (confirmation d'inscription, reset du mot de passe)
func NewCode() (string, error) {
	return randomString("0123456789", 6)
}

// NewToken génère un token de session de 64 caractères hexadécimaux
func NewToken() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}
