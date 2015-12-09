package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
)

// note: names here are obfuscated to make it a bit harder to
// reverse engineer the code and extract decryption code
// not a big challenge

var (
	// must be 32 bytes
	poem = []byte("lost in the wilderness is my dog")
	ivy  = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
)

func init() {
	n := len(poem)
	if n != 32 {
		LogFatalf("len(poem)=%d and should be 32\n", n)
	}
	for i := 0; i < n; i++ {
		poem[i] = poem[i] + 1
	}
	n = len(ivy)
	if n != aes.BlockSize {
		LogFatalf("len(ivy)=%d and should be %d\n", aes.BlockSize)
	}
}

// for ad-hoc testing
func testPoem(plaintext string) {
	fmt.Printf("%s\n", plaintext)
	ciphertext := encrypt(plaintext)
	fmt.Printf("%0x\n", ciphertext)
	result, err := decrypt(ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", result)
}

func encrypt(s string) string {
	if s == "" {
		return ""
	}
	plaintext := []byte(s)
	block, err := aes.NewCipher(poem)
	if err != nil {
		LogFatalf("aes.NewCipher() failed with '%s'\n", err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	copy(iv, ivy)
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return base64.URLEncoding.EncodeToString(ciphertext)
}

func decrypt(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	ciphertext, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("failed to decode '%s'", s)
	}

	block, err := aes.NewCipher(poem)
	if err != nil {
		LogFatalf("aes.NewCipher() failed with '%s'\n", err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("text too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
