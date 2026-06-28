module minenotyours/mine

go 1.26.4

require minenotyours/fileio v0.0.0-00010101000000-000000000000

require (
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/term v0.44.0 // indirect
	minenotyours/mycrypto v0.0.0-00010101000000-000000000000 // indirect
)

replace minenotyours/fileio => ../fileio

replace minenotyours/mycrypto => ../mycrypto
