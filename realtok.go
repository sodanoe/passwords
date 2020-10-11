package main

import (
	"bufio"
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
	wad := "data"
	dbfile := "db/main.db"

	defer slowpoke.CloseAll()

	cypherFile, err := os.OpenFile(cypher, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer cypherFile.Close()

	dataFile, err := os.OpenFile(wad, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer dataFile.Close()

	data := []byte("our super secret text")

	// Инициализация среза для хранения данных в памяти

	allData := make([]string, 1)

	kkey, err := ioutil.ReadFile("cypher")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(dataFile)
	for scanner.Scan() {
		crap := []byte(scanner.Text())
		// --- Туть ошибка слайса slice bounds out of range [:12] with capacity 8
		realData, _ := Decrypt(kkey, crap)
		strrealData := string(realData)
		if strrealData != "" {
			allData = append(allData, strrealData)
		}
		//fmt.Println(allData)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

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
			//fmt.Println(allData)
			// plaintext, err := Decrypt(kkey, ciphertext)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// fmt.Printf("plaintext: %s\n", plaintext)
			for i := 0; i < 10; i++ {
				nummero := strconv.Itoa(i)
				dbkey := []byte(nummero)
				//get from database
				res, _ := slowpoke.Get(dbfile, dbkey)
				//result
				if string(res) != "" {
					decdat, _ := Decrypt(kkey, res)
					ressult := string(decdat)
					fmt.Println(ressult)
					allData = append(allData, ressult)
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
			//var _ = os.Remove(dbfile)
			err = ioutil.WriteFile(wad, []byte(""), 0644)
			//kkey, err := ioutil.ReadFile("cypher")
			if err != nil {
				log.Fatal(err)
			}

			for i := 1; i < len(allData); i++ {
				//Database store
				fmt.Println("im writing")
				keyvalue := strconv.Itoa(i)
				currentData := []byte(allData[i])
				chipheredData, _ := Encrypt(kkey, currentData)
				slowpoke.Set(dbfile, []byte(keyvalue), chipheredData)
				//
				// currentData := []byte(allData[i])
				// chipheredData, _ := Encrypt(kkey, currentData)
				// f, err := os.OpenFile(wad, os.O_APPEND|os.O_WRONLY, 0600)
				// if err != nil {
				// 	panic(err)
				// }
				// defer f.Close()
				// if _, err = f.Write(chipheredData); err != nil {
				// 	panic(err)
				// }
				// if _, err = f.WriteString("\n"); err != nil {
				// 	panic(err)
				// }
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
