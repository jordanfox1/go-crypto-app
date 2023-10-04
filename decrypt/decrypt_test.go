package decrypt

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func setupTestEnv() ([]byte, []byte) {
	var err error

	key, err := os.ReadFile("enc_key.txt")
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := os.ReadFile("nonce.txt")
	if err != nil {
		log.Fatal(err)
	}

	return key, nonce
}

func TestDecryptCipherText(t *testing.T) {
	defer os.Remove("decrypted.txt")
	var key, nonce = setupTestEnv()

	type args struct {
		inputFile  string
		outputFile string
		key        []byte
		nonce      []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Tests if the decrypted output is equal to the original plaintext that was encrypted",
			args: args{
				inputFile:  "encrypted.txt", // Provide the path to the encrypted file
				outputFile: "decrypted.txt", // Provide the path for the decrypted output
				key:        key,
				nonce:      nonce,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DecryptCipherText(tt.args.inputFile, tt.args.outputFile, tt.args.key, tt.args.nonce); (err != nil) != tt.wantErr {
				t.Errorf("DecryptCipherText() error = %v, wantErr %v", err, tt.wantErr)
			}

			decrypted, _ := os.ReadFile("decrypted.txt")
			ptext, _ := os.ReadFile("plaintext.txt")
			t.Logf("Decrypted contents is a string of %s. Original Plaintext is a string of %s", string(decrypted), string(ptext))

			if !bytes.Equal(decrypted, ptext) {
				t.Errorf("DecryptCipherText() error: decrypted text not equal to plaintext, Decrypted contents was: %v. Expected: %v", decrypted, ptext)
			}
		})
	}
}
