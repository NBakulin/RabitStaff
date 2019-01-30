package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func getIoReader() io.Reader {
	reader := bytes.NewReader([]byte{0, 1, 2, 3, 4})
	return reader
}

func encryptRsa(message []byte) []byte {
	privateKey := readPrivateKey()
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &privateKey.PublicKey, message, []byte("orders"))

	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	return ciphertext
}

func dencryptRsa(ciphertext []byte) []byte {
	privateKey := readPrivateKey()
	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, &privateKey, ciphertext, []byte("orders"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return plainText
}
