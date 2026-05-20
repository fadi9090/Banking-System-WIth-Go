package utility

import (
	"strings"
)

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}

	// Must contain @
	if !strings.Contains(email, "@") {
		return false
	}

	// Must contain . after @
	atIndex := strings.Index(email, "@")
	if atIndex == len(email)-1 {
		return false
	}

	afterAt := email[atIndex+1:]
	if !strings.Contains(afterAt, ".") {
		return false
	}

	return true
}

// ValidatePhone checks phone format (optional field)
func ValidatePhone(phone string) bool {
	// Phone is optional - empty is valid
	if phone == "" {
		return true
	}

	// Must have at least 10 digits
	digitCount := 0
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			digitCount++
		}
	}

	return digitCount >= 10
}

// ValidateAmount checks if amount is greater than 0
func ValidateAmount(amount float64) bool {
	return amount > 0
}

// ValidateTransaction checks if a transaction is valid
func ValidateTransaction(balance float64, amount float64, direction string) (string, bool) {
	if amount <= 0 {
		return "Amount must be greater than 0", false
	}

	if direction == "debit" {
		if amount > balance {
			return "Insufficient balance for debit transaction", false
		}
	}

	return "", true
}

// ValidateCurrency checks if currency code is valid (3 uppercase letters)
func ValidateCurrency(currency string) bool {
	if len(currency) != 3 {
		return false
	}

	for _, c := range currency {
		if c < 'A' || c > 'Z' {
			return false
		}
	}

	return true
}
