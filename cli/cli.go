package cli

import (
	"bufio"
	"os"

	pal "github.com/abusomani/go-palette/palette"
)

var p = pal.New()

func PrintWelcomeMessage() {
	p.SetOptions(pal.WithBackground(pal.Color(pal.BrightGreen)), pal.WithForeground(pal.Black))
	// Create a new scanner to read from stdin
	scanner := bufio.NewScanner(os.Stdin)

	// Prompt the user for input
	p.Println("Enter a file name below and press Enter: ")

	// Read user input
	if scanner.Scan() {
		userInput := scanner.Text()
		p.Println("You entered:", userInput)
	} else if err := scanner.Err(); err != nil {
		p.Println("Error reading input:", err)
	}
}
