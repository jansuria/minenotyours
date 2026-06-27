package main

import (
	"flag"
	"fmt"
	"minenotyours/fileio"
	"os"
)

func main() {

	// file = encryptCmd.String("file", "", "path to file")

	args := os.Args

	if len(args) >= 2 {
		switch args[1] {
		case "encrypt":
			encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
			encryptPassword := encryptCmd.String("password", "", "password")
			encryptCmd.Parse(args[2:])
			if *encryptPassword == "" {
				fmt.Println("Error: -password cannot be empty")
				encryptCmd.Usage()
				os.Exit(2)
			}
			if err := fileio.CallEncryption(*encryptPassword); err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
		case "decrypt":
			decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
			decryptPassword := decryptCmd.String("password", "", "password")
			decryptCmd.Parse(args[2:])
			if *decryptPassword == "" {
				fmt.Println("Error: -password cannot be empty")
				decryptCmd.Usage()
				os.Exit(2)
			}
			if err := fileio.CallDecryption(*decryptPassword); err != nil {
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
