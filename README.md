# minenotyours

A small command-line tool for encrypting and decrypting files with a password.
Encryption uses **AES-256-GCM** for the data and **Argon2id** to turn your
password into a key. The installed command is called **`mine`**.

> There are two planned front-ends to the same crypto core:
>
> - **`mine` (CLI)** — the pure command-line version. **This is what ships today.**
> - **TUI** — an interactive terminal UI. *Planned for later; not implemented yet.*

---

## How it works

When you encrypt a file, `mine` overwrites it **in place** with the following layout:

```
[ salt (16 bytes) ][ nonce (12 bytes) ][ ciphertext + GCM auth tag ]
```

- A random 16-byte **salt** is generated per encryption and stored at the front of
  the file. Decryption reads it back to re-derive the same key from your password.
- **Argon2id** derives a 32-byte key (AES-256) from your password + salt using:
  `memory = 64 MiB`, `iterations = 3`, `parallelism = 2`.
- **AES-256-GCM** encrypts and authenticates the data. On decryption, a wrong
  password (or a tampered file) fails the authentication check with
  `cipher: message authentication failed`, and the file is left untouched.

---

## Requirements

- [Go](https://go.dev/dl/) 1.26+ installed.
- Your Go bin directory (`go env GOPATH`\bin, e.g. `C:\Users\<you>\go\bin`) must be
  on your `PATH`. `go install` puts the `mine` binary there.

---

## Install

The project is a multi-module repo that uses local `replace` directives, so install
it from a local clone (not `go install ...@latest`):

```sh
git clone https://github.com/jansuria/minenotyours.git
cd minenotyours/mine
go install .
```

This builds and installs `mine` (`mine.exe` on Windows) into your Go bin directory.
Confirm it's available:

```sh
mine
```

You should see a usage error rather than "command not found". If you get
"command not found", your Go bin directory isn't on your `PATH`.

---

## Usage

```sh
mine encrypt -file <path>
mine decrypt -file <path>
```

`-file` is the only flag, and it's required. The **password is entered
interactively** — `mine` prompts for it and hides your typing, so the password
never appears on screen, in your shell history, or in the process list.

- On **encrypt**, you're asked to enter the password twice (to catch typos). If the
  two entries don't match, nothing is encrypted.
- On **decrypt**, you're asked once. A wrong password fails safely and leaves the
  file unchanged.

> The password prompt requires a real terminal. You can't pipe a password in
> (e.g. `echo pw | mine ...`) — that's intentional, to keep secrets off the
> command line.

### Paths are relative to your current directory

`-file` is resolved against the directory you run `mine` from, so you can work with
files anywhere:

```sh
cd ~/Documents/secrets
mine encrypt -file taxes.txt   # encrypts ~/Documents/secrets/taxes.txt
mine decrypt -file taxes.txt   # restores it
```

Absolute paths work too:

```sh
mine encrypt -file "C:\Users\me\notes.txt"
```

### Example round-trip

```sh
echo hello > secret.txt
mine encrypt -file secret.txt   # prompts for password (twice); secret.txt is now ciphertext
mine decrypt -file secret.txt   # prompts for password; secret.txt is "hello" again
```

---

## Exit codes

| Code | Meaning                                                             |
|------|--------------------------------------------------------------------|
| `0`  | Success                                                            |
| `1`  | Operation failed (e.g. wrong password, file not found)            |
| `2`  | Misuse (unknown command, missing `-file`)                         |

---

## ⚠️ Important notes

- **Encryption is in place and destructive.** The original file is overwritten with
  its encrypted form. Keep a backup of anything irreplaceable until you've confirmed
  you can decrypt it.
- **There is no password recovery.** If you forget the password, the data is gone.
- The password is read **interactively and never echoed**, so it stays out of your
  shell history and the process list. Still, treat this tool as a learning/personal
  utility rather than a hardened secret manager.

---

## Project layout

This is a multi-module Go workspace:

| Module                | Path        | Responsibility                                   |
|-----------------------|-------------|--------------------------------------------------|
| `minenotyours/mine`   | `mine/`     | CLI entry point: arg parsing, flags, exit codes  |
| `minenotyours/fileio` | `fileio/`   | Thin wrappers that set parameters and call crypto|
| `minenotyours/mycrypto`| `mycrypto/`| AES-GCM encryption/decryption + Argon2 hashing   |

`mine` depends on `fileio`, which depends on `mycrypto`, wired together with
relative `replace` directives in each `go.mod`.

---

## Roadmap

- [x] CLI: `encrypt` / `decrypt` with `-file`
- [x] Interactive password prompt (hidden input, confirm-on-encrypt)
- [ ] TUI version
- [ ] Optional `-out` flag to write to a new file instead of in place

---

## License

[MIT](LICENSE)
