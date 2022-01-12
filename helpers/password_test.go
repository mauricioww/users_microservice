package helpers_test

import (
	"testing"

	"github.com/mauricioww/user_microsrv/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCipherText(t *testing.T) {
	test_cases := []struct {
		test_name  string
		plain_text string
		res        string
		flag       bool
	}{
		{
			test_name:  "cipher plain text success",
			plain_text: "hola",
			res:        "szwl",
			flag:       true,
		},
		{
			test_name:  "cipher plain text error",
			plain_text: "hola",
			res:        "hola",
			flag:       false,
		},
		{
			test_name:  "no cipher symbols",
			plain_text: "¿?¡!.-_",
			res:        "¿?¡!.-_",
			flag:       false,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)

			// act
			new_text := helpers.Cipher(tc.plain_text)

			// assert
			if tc.flag {
				assert.Equal(tc.res, new_text)
			} else {
				assert.NotEqual(tc.res, new_text)
			}
		})
	}
}

func TestDecipherText(t *testing.T) {
	test_cases := []struct {
		test_name string
		text      string
		res       string
		flag      bool
	}{
		{
			test_name: "decipher text success",
			text:      "szwl",
			res:       "hola",
			flag:      true,
		},
		{
			test_name: "decipher text error",
			text:      "hola",
			res:       "hola",
			flag:      false,
		},
		{
			test_name: "no decipher symbols",
			text:      "¿?¡!.-_",
			res:       "¿?¡!.-_",
			flag:      false,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)

			// act
			new_text := helpers.Decipher(tc.text)

			// assert
			if tc.flag {
				assert.Equal(tc.res, new_text)
			} else {
				assert.NotEqual(tc.res, new_text)
			}
		})
	}
}
