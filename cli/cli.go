package main

import (
	"minenotyours/fileio"
	"os"
)

func main() {

	args := os.Args

	switch args[1] {
	case "encrypt":
		fileio.CallEncryption(args[2])
	case "decrypt":
		fileio.CallDecryption(args[2])
	}
}
