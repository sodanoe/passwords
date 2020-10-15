package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/recoilme/slowpoke"
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
	cypher := "cypher"
	dbfile := "db/main.db"

	defer slowpoke.CloseAll()

	cypherFile, err := os.OpenFile(cypher, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer cypherFile.Close()

	data := []byte("our super secret text")

	// Инициализация среза для хранения данных в памяти

	allData := make([]string, 0)

	kkey, err := ioutil.ReadFile("cypher")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		nummero := strconv.Itoa(i)
		dbkey := []byte(nummero)
		res, _ := slowpoke.Get(dbfile, dbkey)
		if string(res) != "" {
			decdat, _ := Decrypt(kkey, res)
			ressult := string(decdat)
			allData = append(allData, ressult)
		}
	}

	fmt.Println(allData)

	var i int
	for i != 4 {
		fmt.Println("Меню: ")
		fmt.Println("1. Вывод данных")
		fmt.Println("2. Пополнить таблицу")
		fmt.Println("3. Вывести какой-то хуйни")
		fmt.Println("4. Выход")

		fmt.Scan(&i)
		switch i {
		case 1:
			for i := 0; i < 10; i++ {
				nummero := strconv.Itoa(i)
				dbkey := []byte(nummero)
				//get from database
				res, _ := slowpoke.Get(dbfile, dbkey)
				//result
				if string(res) != "" {
					decdat, _ := Decrypt(kkey, res)
					ressult := string(decdat)
					fmt.Print(string(dbkey), " is ")
					fmt.Println(ressult)
				}
			}
		case 2:
			{
				var name string
				var password string
				for {
					fmt.Println("Please input username: ")
					fmt.Scan(&name)
					if name == "stop" {
						break
					}

					fmt.Println("Please input password: ")
					fmt.Scan(&password)

					result := name + ":" + password
					allData = append(allData, result)
				}
				fmt.Println(allData)
			}
		case 3:
			{
				kkey, err := ioutil.ReadFile("cypher")
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
		case 4:
			if err != nil {
				log.Fatal(err)
			}

			for key, value := range allData {
				//Database store
				keyvalue := strconv.Itoa(key)
				currentData := []byte(value)
				chipheredData, _ := Encrypt(kkey, currentData)
				slowpoke.Set(dbfile, []byte(keyvalue), chipheredData) // instead of chipheredData
				fmt.Println("||", value, "||")
			}
			break
		}
	}
	// key, err := GenerateKey()
	// Write key to cypher file
	// err = ioutil.WriteFile(cypher, key, 0644)
	// if err != nil {
	//     return
	// }
}
