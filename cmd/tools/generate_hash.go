// Package main provides a command-line tool for generating bcrypt password hashes
// with configurable cost parameters for the ChatLogger API application.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Define command-line flags
	var password string
	var cost int

	flag.StringVar(&password, "password", "", "Password to hash")
	flag.IntVar(&cost, "cost", 10, "Bcrypt cost (4-31, higher is more secure but slower)")
	flag.Parse()

	// If no password provided via flag, check for it as a positional argument
	if password == "" && len(flag.Args()) > 0 {
		password = flag.Args()[0]
	}

	// If still no password, prompt or exit
	if password == "" {
		fmt.Println("Usage: generate_hash [-password=<password>] [-cost=<cost>] [password]")
		fmt.Println("  or provide password via stdin")
		fmt.Print("Enter password to hash: ")

		_, err := fmt.Scanln(&password)
		if err != nil {
			log.Fatalf("Error reading password: %v", err)
		}

		if password == "" {
			os.Exit(1)
		}
	}

	// Validate cost parameter
	if cost < 4 || cost > 31 {
		log.Fatalf("Invalid cost parameter: %d (must be between 4-31)", cost)
	}

	// Generate hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Fatalf("Error generating hash: %v", err)
	}

	// Print the hash
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash: %s\n", hashedPassword)
}
