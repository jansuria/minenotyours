package main

import (
	"flag"
	"fmt"
	"minenotyours/fileio"
	"os"

	"golang.org/x/term"
)

func ReadNewPassword() (string, error) {
	fmt.Print("Enter Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	fmt.Print("Confirm Password: ")
	confirmPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	if string(password) != string(confirmPassword) {
		return "", fmt.Errorf("passwords do not match")
	}

	return string(password), nil

}

func ReadPassword() (string, error) {
	fmt.Print("Enter Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func main() {
	args := os.Args

	if len(args) >= 2 {
		switch args[1] {
		case "encrypt":
			encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
			encryptFile := encryptCmd.String("file", "", "path to file")
			encryptCmd.Parse(args[2:])
			if *encryptFile == "" {
				fmt.Println("Error: -file cannot be empty")
				encryptCmd.Usage()
				os.Exit(2)
			}
			encryptPassword, err := ReadNewPassword()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			if err := fileio.CallEncryption(encryptPassword, *encryptFile); err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
		case "decrypt":
			decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
			decryptFile := decryptCmd.String("file", "", "path to file")
			decryptCmd.Parse(args[2:])
			if *decryptFile == "" {
				fmt.Println("Error: -file cannot be empty")
				decryptCmd.Usage()
				os.Exit(2)
			}
			decryptPassword, err := ReadPassword()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			if err := fileio.CallDecryption(decryptPassword, *decryptFile); err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
		default:
			fmt.Println("Not a valid command")
			os.Exit(2)
		}
	} else {
		fmt.Println("Not valid")
		os.Exit(2)
	}
}
