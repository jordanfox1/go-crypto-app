package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"os"
)

func DecryptCipherText(inputFile, outputFile string, key, nonce []byte) error {
	ciphertext, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	var plaintext []byte

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	err = os.WriteFile(outputFile, plaintext, 0644)
	if err != nil {
		return err
	}

	return nil
}
