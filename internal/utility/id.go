package utility

import (
	"github.com/google/uuid"
)

func GenerateID() string {
	return uuid.New().String()
}

func GenerateAccountNumber() string {
	return "ACC-" + uuid.New().String()[:8]
}

func GenerateCustomerNumber() string {
	return "C-" + uuid.New().String()[:8]
}

func GenerateCardLast4() string {
	return uuid.New().String()[:4]
}

func GenerateCVV() string {
	return uuid.New().String()[:3]
}

func GeneratePersonNumber() string {
	return "P-" + uuid.New().String()[:8]
}

func GenerateReferenseNumber() string {
	return "Pos-" + uuid.New().String()[:3]
}
