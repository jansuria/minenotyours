package fileio

import (
	"minenotyours/mycrypto"
)

func CallDecryption(password string) error {

	path := `C:\Users\Jan-PC\Documents\Codes\Go\minenotyours\mine\test.txt`

	parameters := &mycrypto.ArgonParameters{
		Memory:      64 * 1024,
		Iteration:   3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
	return mycrypto.DecryptFile(path, password, *parameters)
}
