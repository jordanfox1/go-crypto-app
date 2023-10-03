package encrypt

import (
	"bytes"
	"crypto/cipher"
	"os"
	"path/filepath"
	"testing"
)

func createMockEncryptionContextData(t *testing.T, key []byte) (mockGcm cipher.AEAD, mockNonce []byte) {

	// Create a GCM instance
	mockGcm, gcmErr := NewGCMInstance(key)
	if gcmErr != nil {
		t.Fatalf("Failed to create GCM: %v", gcmErr)
	}

	// Generate a single nonce for both files
	mockNonce, nonceErr := GenerateNonce(mockGcm, mockFn, "")
	if nonceErr != nil {
		t.Fatalf("Failed to generate nonce: %v", nonceErr)
	}

	return mockGcm, mockNonce
}

func mockFn(foo []byte, bar string) {}

func TestEncryptPlainText(t *testing.T) {
	mockKey, err := Generate32ByteEncryptionKey(mockFn, "string")
	if err != nil {
		t.Fatalf("Failed to create an ecryption key in test: %v", err)
	}
	mockGcm, mockNonce := createMockEncryptionContextData(t, mockKey)

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "encrypt_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory when the test is done (runs after each test case)

	type fileConfig struct {
		fileName           string
		textToAppendToFile string
	}

	type args struct {
		inputFiles  []fileConfig
		outputFiles []fileConfig
		nonce       []byte
		gcm         cipher.AEAD
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		assertion string
	}{
		{
			name: "Same Key, Nonce, and GCM should produce the same output (determinism test)",
			args: args{
				inputFiles: []fileConfig{
					{
						fileName:           filepath.Join(tempDir, "plaintext1.txt"),
						textToAppendToFile: "sample plaintext",
					},
					{
						fileName:           filepath.Join(tempDir, "plaintext2.txt"),
						textToAppendToFile: "sample plaintext",
					},
				},
				outputFiles: []fileConfig{
					{
						fileName: filepath.Join(tempDir, "encrypted1.txt"),
					},
					{
						fileName: filepath.Join(tempDir, "encrypted2.txt"),
					},
				},
				nonce: mockNonce, // Use the same nonce for both files
				gcm:   mockGcm,   // Use the same GCM instance for both files
			},
			wantErr:   false,
			assertion: "Assert same value",
		},
		{
			name: "Same Key, Nonce, and GCM but different text should produce different output (determinism test)",
			args: args{
				inputFiles: []fileConfig{
					{
						fileName:           filepath.Join(tempDir, "plaintext1.txt"),
						textToAppendToFile: "sample plaintext",
					},
					{
						fileName:           filepath.Join(tempDir, "plaintext2.txt"),
						textToAppendToFile: "sample plaintext 2",
					},
				},
				outputFiles: []fileConfig{
					{
						fileName: filepath.Join(tempDir, "encrypted1.txt"),
					},
					{
						fileName: filepath.Join(tempDir, "encrypted2.txt"),
					},
				},
				nonce: mockNonce, // Use the same nonce for both files
				gcm:   mockGcm,   // Use the same GCM instance for both files
			},
			wantErr:   false,
			assertion: "Assert different values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// SETUP
			// create plaintext files based on the test config
			for idx, inputFile := range tt.args.inputFiles {
				// Create the corresponding plaintext file
				err := os.WriteFile(inputFile.fileName, []byte(inputFile.textToAppendToFile), 0644)
				if err != nil {
					t.Fatalf("Failed to create plaintext file: %v", err)
				}

				// Encrypt the plaintext using the provided key, nonce, and GCM
				err = EncryptPlainText(inputFile.fileName, tt.args.outputFiles[idx].fileName, tt.args.nonce, tt.args.gcm)
				if err != nil {
					t.Errorf("EncryptPlainText() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			// Read and compare the encrypted outputs
			encryptedOutput1, err := os.ReadFile(tt.args.outputFiles[0].fileName)
			if err != nil {
				t.Fatalf("Failed to read encrypted output file: %v", err)
			}

			encryptedOutput2, err := os.ReadFile(tt.args.outputFiles[1].fileName)
			if err != nil {
				t.Fatalf("Failed to read encrypted output file: %v", err)
			}

			// ASSERTIONS
			switch tt.assertion {
			case "Assert same value":
				if !bytes.Equal(encryptedOutput1, encryptedOutput2) {
					t.Errorf("Encrypted outputs should be the same")
				}
			case "Assert different values":
				if bytes.Equal(encryptedOutput1, encryptedOutput2) {
					t.Errorf("Encrypted outputs should NOT be the same")
				}
			}
		})
	}
}
