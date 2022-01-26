package helpers

const (
	shift       = 11
	alphabetLen = 26
)

// Cipher func uses the Ceaser algorithm to transform the plain text
func Cipher(message string) string {
	len := len(message)
	cipheredText := make([]rune, len)
	for index, char := range message {
		if isLower(char) {
			cipheredText[index] = (char-'a'+shift)%alphabetLen + 'a'
		} else if isUpper(char) {
			cipheredText[index] = (char-'A'+shift)%alphabetLen + 'A'
		} else {
			cipheredText[index] = char
		}
	}
	return string(cipheredText)
}

// Decipher func uses the Ceaser algorithm to return the original plain text
func Decipher(stream string) string {
	len := len(stream)
	plainText := make([]rune, len)
	for index, char := range stream {
		if isLower(char) {
			offset := char - 'a' - shift
			if offset < 0 {
				offset += alphabetLen
			}
			plainText[index] = offset%alphabetLen + 'a'
		} else if isUpper(char) {
			offset := char - 'A' - shift
			if offset < 0 {
				offset += alphabetLen
			}
			plainText[index] = offset%alphabetLen + 'A'
		} else {
			plainText[index] = char
		}
	}
	return string(plainText)
}

func isLower(char rune) bool {
	return char >= 'a' && char <= 'z'
}

func isUpper(char rune) bool {
	return char >= 'A' && char <= 'Z'
}
