package helpers

const (
	shift        = 11
	alphabet_len = 26
)

func Cipher(message string) string {
	len := len(message)
	ciphered_text := make([]rune, len)
	for index, char := range message {
		if is_lower(char) {
			ciphered_text[index] = (char-'a'+shift)%alphabet_len + 'a'
		} else if is_upper(char) {
			ciphered_text[index] = (char-'A'+shift)%alphabet_len + 'A'
		} else {
			ciphered_text[index] = char
		}
	}
	return string(ciphered_text)
}

func Decipher(stream string) string {
	len := len(stream)
	plain_text := make([]rune, len)
	for index, char := range stream {
		if is_lower(char) {
			offset := char - 'a' - shift
			if offset < 0 {
				offset += alphabet_len
			}
			plain_text[index] = offset%alphabet_len + 'a'
		} else if is_upper(char) {
			offset := char - 'A' - shift
			if offset < 0 {
				offset += alphabet_len
			}
			plain_text[index] = offset%alphabet_len + 'A'
		} else {
			plain_text[index] = char
		}
	}
	return string(plain_text)
}

func is_lower(char rune) bool {
	return char >= 'a' && char <= 'z'
}

func is_upper(char rune) bool {
	return char >= 'A' && char <= 'Z'
}
