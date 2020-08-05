package main

import (
	"Users/sergeynogin/DEV/golang/passwords/menu"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	// "os"
)

//Encrypt - шифрование
func Encrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

//Decrypt – расшифровка конечно же
func Decrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

//GenerateKey – создание ключа шифрования
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func main() {

	//File openings

	// kkey := []byte("le'gj[rjv")

	// cypher := "cypher.txt"
	// shit := "shit.txt"

	// cypherFile, err := os.OpenFile(cypher, os.O_RDWR, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer cypherFile.Close()

	// shitFile, err := os.OpenFile(shit, os.O_RDWR, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer shitFile.Close()

	//

	data := []byte("our super secret text")

	// key, err := GenerateKey()

	//Write key to cypher.txt file

	// err = ioutil.WriteFile(cypher, key, 0644)
	// if err != nil {
	//     return
	// }
	//
	//

	menu.PrintShit()

	kkey, err := ioutil.ReadFile("cypher.txt")
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	ciphertext, err := Encrypt(kkey, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ciphertext: %s\n", hex.EncodeToString(ciphertext))
	plaintext, err := Decrypt(kkey, ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("plaintext: %s\n", plaintext)
}
