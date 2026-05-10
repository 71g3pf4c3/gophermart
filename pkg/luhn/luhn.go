// Package luhn validates order numbers using the Luhn algorithm.
package luhn

// Valid returns true if the number string passes the Luhn check.
func Valid(number string) bool {
	if len(number) == 0 {
		return false
	}

	sum := 0
	nDigits := len(number)
	parity := nDigits % 2

	for i, ch := range number {
		if ch < '0' || ch > '9' {
			return false
		}
		digit := int(ch - '0')
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
