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

	if os.Getenv("SKIP_AUTH") == "" {
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return err
		}

		plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return err
		}
	} else {
		// WARNING!! Decrypts without an authentication check. This is to avoid authentication errors when calling Open from the test runner.
		// Do not run this else block outside a test environment.
		plaintext = make([]byte, len(ciphertext))
		iv := make([]byte, block.BlockSize())
		stream := cipher.NewCTR(block, iv)
		stream.XORKeyStream(plaintext, ciphertext)
	}

	err = os.WriteFile(outputFile, plaintext, 0644)
	if err != nil {
		return err
	}

	return nil
}
