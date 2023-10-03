package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

func EncryptPlainText(inputFile, outputFile string, nonce []byte, gcm cipher.AEAD) error {
	plainText, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	ciphertext := gcm.Seal(nil, nonce, plainText, nil)

	if err := os.WriteFile(outputFile, ciphertext, 0644); err != nil {
		return err
	}

	return nil
}

// This will generate a new key each time it is called
func Generate32ByteEncryptionKey(keyStorageFunc func([]byte, string), storageFileName string) ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("Error generating key:", err)
		return nil, err
	}

	keyStorageFunc(key, storageFileName)

	return key, nil
}

func StoreInPlainTextFile(data []byte, fileName string) {
	os.WriteFile(fileName, data, 0644)
}

func GenerateNonce(gcm cipher.AEAD, nonceStorageFunc func([]byte, string), storageFileName string) ([]byte, error) {
	nonce := make([]byte, gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	nonceStorageFunc(nonce, storageFileName)

	return nonce, nil
}

func NewGCMInstance(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm, nil
}
