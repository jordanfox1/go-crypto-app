package main

import (
	"crypto_cli/decrypt"
	"crypto_cli/encrypt"
	"log"
	"os"
)

func main() {
	inputFile := "plaintext.txt"
	encryptedFile := "encrypted.txt"

	key, encKeyError := encrypt.Generate32ByteEncryptionKey(encrypt.StoreInPlainTextFile, "enc_key.txt")
	if encKeyError != nil {
		log.Panic(encKeyError)
	}

	gcmInstance, err := encrypt.NewGCMInstance(key)
	if err != nil {
		log.Panic(err)
	}

	nonce, err := encrypt.GenerateNonce(gcmInstance, encrypt.StoreInPlainTextFile, "nonce.txt")
	if err != nil {
		log.Panic(err)
	}

	encErr := encrypt.EncryptPlainText(inputFile, encryptedFile, nonce, gcmInstance)
	if encErr != nil {
		log.Panic(encErr)
	}
	dec()
}

func dec() {
	key, err := os.ReadFile("enc_key.txt")
	if err != nil {
		log.Panic(err)
	}

	nonce, err := os.ReadFile("nonce.txt")
	if err != nil {
		log.Panic(err)
	}

	decrypt.DecryptCipherText("encrypted.txt", "decrypted.txt", key, nonce)
}
