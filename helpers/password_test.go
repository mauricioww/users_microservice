package helpers_test

import (
	"testing"

	"github.com/mauricioww/user_microsrv/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCipherText(t *testing.T) {
	testCases := []struct {
		testName  string
		plainText string
		res       string
		flag      bool
	}{
		{
			testName:  "cipher plain text success",
			plainText: "hola",
			res:       "szwl",
			flag:      true,
		},
		{
			testName:  "cipher plain text error",
			plainText: "hola",
			res:       "hola",
			flag:      false,
		},
		{
			testName:  "no cipher symbols",
			plainText: "¿?¡!.-_",
			res:       "¿?¡!.-_",
			flag:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			assert := assert.New(t)

			// act
			newText := helpers.Cipher(tc.plainText)

			// assert
			if tc.flag {
				assert.Equal(tc.res, newText)
			} else {
				assert.NotEqual(tc.res, newText)
			}
		})
	}
}

func TestDecipherText(t *testing.T) {
	testCases := []struct {
		testName string
		text     string
		res      string
		flag     bool
	}{
		{
			testName: "decipher text success",
			text:     "szwl",
			res:      "hola",
			flag:     true,
		},
		{
			testName: "decipher text error",
			text:     "hola",
			res:      "hola",
			flag:     false,
		},
		{
			testName: "no decipher symbols",
			text:     "¿?¡!.-_",
			res:      "¿?¡!.-_",
			flag:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			assert := assert.New(t)

			// act
			newText := helpers.Decipher(tc.text)

			// assert
			if tc.flag {
				assert.Equal(tc.res, newText)
			} else {
				assert.NotEqual(tc.res, newText)
			}
		})
	}
}
