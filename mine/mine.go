package main

import (
	"flag"
	"fmt"
	"minenotyours/fileio"
	"os"
)

func main() {
	args := os.Args

	if len(args) >= 2 {
		switch args[1] {
		case "encrypt":
			encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
			encryptPassword := encryptCmd.String("password", "", "file password")
			encryptFile := encryptCmd.String("file", "", "path to file")
			encryptCmd.Parse(args[2:])
			if *encryptPassword == "" {
				fmt.Println("Error: -password cannot be empty")
				encryptCmd.Usage()
				os.Exit(2)
			}
			if *encryptFile == "" {
				fmt.Println("Error: -file cannot be empty")
				encryptCmd.Usage()
				os.Exit(2)
			}
			if err := fileio.CallEncryption(*encryptPassword, *encryptFile); err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
		case "decrypt":
			decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
			decryptPassword := decryptCmd.String("password", "", "file password")
			decryptFile := decryptCmd.String("file", "", "path to file")
			decryptCmd.Parse(args[2:])
			if *decryptPassword == "" {
				fmt.Println("Error: -password cannot be empty")
				decryptCmd.Usage()
				os.Exit(2)
			}
			if *decryptFile == "" {
				fmt.Println("Error: -file cannot be empty")
				decryptCmd.Usage()
				os.Exit(2)
			}
			if err := fileio.CallDecryption(*decryptPassword, *decryptFile); err != nil {
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
