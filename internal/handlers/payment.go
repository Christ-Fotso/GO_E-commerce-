package handlers

import (
	"strconv"
	"strings"
	"time"
)

// PaymentCard contient les informations bancaires saisies par le client (paiement simulé)
type PaymentCard struct {
	CardNumber string `json:"card_number"`
	Expiry     string `json:"expiry"` // format MM/AA, ex : "12/27"
	CVC        string `json:"cvc"`
}

// Validate renvoie un message d'erreur, ou "" si la carte est acceptée
func (card PaymentCard) Validate() string {
	cardNumber := strings.ReplaceAll(card.CardNumber, " ", "")
	if len(cardNumber) < 13 || len(cardNumber) > 19 || !isDigitsOnly(cardNumber) {
		return "Numéro de carte invalide (13 à 19 chiffres)"
	}
	if !passesLuhnCheck(cardNumber) {
		return "Numéro de carte invalide (clé de contrôle incorrecte)"
	}
	if !isExpiryValid(card.Expiry, time.Now()) {
		return "Date d'expiration invalide ou dépassée (format MM/AA)"
	}
	if (len(card.CVC) != 3 && len(card.CVC) != 4) || !isDigitsOnly(card.CVC) {
		return "CVC invalide (3 ou 4 chiffres)"
	}
	return ""
}

func isDigitsOnly(text string) bool {
	for _, character := range text {
		if character < '0' || character > '9' {
			return false
		}
	}
	return text != ""
}

// passesLuhnCheck applique l'algorithme de Luhn, utilisé par toutes les cartes bancaires :
// en partant de la droite, on double un chiffre sur deux (en retirant 9 si on dépasse 9),
// puis la somme de tous les chiffres doit être un multiple de 10.
func passesLuhnCheck(cardNumber string) bool {
	sum := 0
	shouldDouble := false
	for position := len(cardNumber) - 1; position >= 0; position-- {
		digit := int(cardNumber[position] - '0')
		if shouldDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		shouldDouble = !shouldDouble
	}
	return sum%10 == 0
}

// isExpiryValid vérifie le format MM/AA et que la carte n'est pas expirée
// (une carte "12/27" est valable jusqu'au 31 décembre 2027 inclus)
func isExpiryValid(expiry string, now time.Time) bool {
	monthText, yearText, found := strings.Cut(expiry, "/")
	if !found {
		return false
	}
	month, monthError := strconv.Atoi(monthText)
	year, yearError := strconv.Atoi(yearText)
	if monthError != nil || yearError != nil || month < 1 || month > 12 || len(yearText) != 2 {
		return false
	}
	// Premier jour du mois qui suit l'expiration
	firstInvalidDay := time.Date(2000+year, time.Month(month)+1, 1, 0, 0, 0, 0, now.Location())
	return now.Before(firstInvalidDay)
}
