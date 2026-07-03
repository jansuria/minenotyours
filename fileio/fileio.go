package fileio

import (
	"minenotyours/mycrypto"
)

	var parameters = &mycrypto.ArgonParameters{
		Memory:      64 * 1024,
		Iteration:   3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

func CallDecryption(password string, path string) error {

	return mycrypto.DecryptFile(path, password, *parameters)
}

func CallEncryption(password string, path string) error {
	
	return mycrypto.EncryptFile(path, password, *parameters)
}
