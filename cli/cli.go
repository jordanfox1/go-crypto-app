package cli

import (
	"bufio"
	"crypto_cli/decrypt"
	"crypto_cli/encrypt"
	"os"

	pal "github.com/abusomani/go-palette/palette"
	progressbar "github.com/schollz/progressbar/v3"
)

var prompt = pal.New(pal.WithBackground(pal.Color(pal.BrightGreen)), pal.WithForeground(pal.Black))
var process = pal.New(pal.WithBackground(pal.Color(pal.Black)), pal.WithForeground(pal.Yellow))
var e = pal.New(pal.WithBackground(pal.Color(pal.BrightRed)), pal.WithForeground(pal.Black))

func StartEncrypt() {
	// Create a new scanner to read from stdin
	scanner := bufio.NewScanner(os.Stdin)

	option, err := selectOption()
	if err != nil {
		e.Println("Invalid option", option)
		os.Exit(1)
	}
	if option == "dec" {
		startDecrypt()
	}

	// Prompt the user for input
	prompt.Println("Enter the name of the file for encryption below and press Enter: ")

	if scanner.Scan() {
		userInput := scanner.Text()
		prompt.Println("You entered: ", userInput)
		process.Println("Attempting to encrypt file in current directory named:", userInput)

		err := checkIfFileExists(userInput)
		if err != nil {
			os.Exit(1)
		}

		bar := progressbar.Default(100)
		bar.Add(50)

		nonceFileName := "nonce.txt"
		keyFileName := "enc_key.txt"
		encDataFileName := "encrypted.txt"

		key, err := encrypt.Generate32ByteEncryptionKey(encrypt.StoreInPlainTextFile, keyFileName)
		if err != nil {
			e.Println("Error creating encryption key", err)
			os.Exit(1)
		}
		bar.Add(10)

		gcm, err := encrypt.NewGCMInstance(key)
		if err != nil {
			e.Println("Error creating GCM instance", err)
			os.Exit(1)
		}
		bar.Add(10)

		nonce, err := encrypt.GenerateNonce(gcm, encrypt.StoreInPlainTextFile, nonceFileName)
		if err != nil {
			e.Println("Error creating nonce", err)
			os.Exit(1)
		}
		bar.Add(10)

		err = encrypt.EncryptPlainText(userInput, encDataFileName, nonce, gcm)
		if err != nil {
			e.Println("Error performing encryption", err)
			os.Exit(1)
		}
		bar.Add(20)
		prompt.Println("Encryption Successful!")
		process.Printf("The following files have been created in the current directory - nonce: %v, encryption key: %v, encrypted data: %v", nonceFileName, keyFileName, encDataFileName)
		prompt.Println("the nonce(non-sensitive) and key(sensitive) will be needed for decryption")

	} else if err := scanner.Err(); err != nil {
		e.Println("Error reading input:", err)
	}
}

func checkIfFileExists(userInput string) error {
	_, err := os.ReadFile(userInput)
	if err != nil {
		e.Println(userInput, "File does not exist", err)
		return err
	}
	return nil
}

func selectOption() (string, error) {
	prompt.Println("Type enc for encryption or dec for decryption and press enter")
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		return "", err
	}
	userInput := scanner.Text()
	if userInput == "enc" || userInput == "dec" {
		return userInput, nil
	}
	e.Println("Invalid input", userInput)
	os.Exit(1)
	return "", nil
}

func startDecrypt() {
	scanner := bufio.NewScanner(os.Stdin)
	// Prompt the user for input
	prompt.Println("Enter the name of the file for decryption below and press Enter: ")

	if scanner.Scan() {
		userInput := scanner.Text()
		prompt.Println("You entered: ", userInput)
		process.Println("Attempting to decrypt file in current directory named:", userInput)

		err := checkIfFileExists(userInput)
		if err != nil {
			os.Exit(1)
		}

		bar := progressbar.Default(100)
		bar.Add(50)

		nonceFileName := "nonce.txt"
		keyFileName := "enc_key.txt"
		decDataFileName := "decrypted.txt"

		key, err := os.ReadFile(keyFileName)
		if err != nil {
			e.Println("Error reading encryption key from file named enc_txt", err)
			os.Exit(1)
		}
		bar.Add(10)
		bar.Add(10)

		nonce, err := os.ReadFile(nonceFileName)
		if err != nil {
			e.Println("Error reading nonce", err)
			os.Exit(1)
		}
		bar.Add(10)

		err = decrypt.DecryptCipherText(userInput, decDataFileName, key, nonce)
		if err != nil {
			e.Println("Error performing encryption", err)
			os.Exit(1)
		}
		bar.Add(20)
		prompt.Println("Decryption Successful!")
		process.Printf("The following files have been created in the current directory - decrypted data: %v", decDataFileName)
		os.Exit(0)

	} else if err := scanner.Err(); err != nil {
		e.Println("Error reading input:", err)
	}
}
