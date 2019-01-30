package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"os"
)

func main1() {

	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	var publickey *rsa.PublicKey
	publickey = &privatekey.PublicKey

	// save private and public key separately
	privatekeyfile, err := os.Create("private.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	privatekeyencoder := gob.NewEncoder(privatekeyfile)
	privatekeyencoder.Encode(privatekey)
	privatekeyfile.Close()

	publickeyfile, err := os.Create("public.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	publickeyencoder := gob.NewEncoder(publickeyfile)
	publickeyencoder.Encode(publickey)
	publickeyfile.Close()

	// save PEM file
	pemfile, err := os.Create("private.pem")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// http://golang.org/pkg/encoding/pem/#Block
	var pemkey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}

	err = pem.Encode(pemfile, pemkey)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pemfile.Close()
}

func readPrivateKey() rsa.PrivateKey {
	privateKey := rsa.PrivateKey{}
	file, _ := os.Open("private.key")
	defer file.Close()

	decoder := gob.NewDecoder(file)
	decoder.Decode(&privateKey)
	return privateKey
}