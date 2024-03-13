package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const keyFile = "key.txt"

func generateKey() ([]byte, error) {
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(keyFile, key, 0644)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func readKey() ([]byte, error) {
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func encryptImage(inputPath string, outputPath string, key []byte) error {
	// Open the input image file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Create the output encrypted image file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Generate a new AES cipher block using the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Generate a nonce with the correct length (block size of the cipher)
	nonce := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// Write the nonce to the beginning of the output file
	if _, err := outputFile.Write(nonce); err != nil {
		return err
	}

	// Create a writer that encrypts data written to it
	stream := cipher.NewCTR(block, nonce)
	writer := &cipher.StreamWriter{S: stream, W: outputFile}

	// Copy the plaintext image data to the encrypted image file
	if _, err := io.Copy(writer, inputFile); err != nil {
		return err
	}

	return nil
}

func decryptImage(inputPath string, outputPath string, key []byte) error {
	// Open the input encrypted image file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Read nonce from the beginning of the file
	nonce := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inputFile, nonce); err != nil {
		return err
	}

	// Create AES cipher block using the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Create a new stream for decryption
	stream := cipher.NewCTR(block, nonce)

	// Create a reader that decrypts data read from it
	reader := &cipher.StreamReader{S: stream, R: inputFile}

	// Create the output decrypted image file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Copy the encrypted image data to the output decrypted image file
	if _, err := io.Copy(outputFile, reader); err != nil {
		return err
	}

	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Check if key file exists, if not, generate a new key
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		fmt.Println("Key file not found, generating a new key...")
		key, err := generateKey()
		if err != nil {
			panic(err)
		}
		fmt.Println("New key generated and saved to key.txt.")
		fmt.Println("Please store this key securely.")
		fmt.Printf("Key: %x\n", key)
	} else {
		fmt.Println("Using existing key from key.txt")
	}

	key, err := readKey()
	if err != nil {
		panic(err)
	}

	fmt.Print("Enter 'encrypt' or 'decrypt' mode: ")
	mode, _ := reader.ReadString('\n')
	mode = mode[:len(mode)-1]

	fmt.Print("Enter the path to the input image file: ")
	inputPath, _ := reader.ReadString('\n')
	inputPath = inputPath[:len(inputPath)-1]

	fmt.Print("Enter the path to save the output image file: ")
	outputPath, _ := reader.ReadString('\n')
	outputPath = outputPath[:len(outputPath)-1]

	switch mode {
	case "encrypt":
		if err := encryptImage(inputPath, outputPath, key); err != nil {
			panic(err)
		}
		fmt.Println("Image encrypted successfully.")
	case "decrypt":
		if err := decryptImage(inputPath, outputPath, key); err != nil {
			panic(err)
		}
		fmt.Println("Image decrypted successfully.")
	default:
		fmt.Println("Invalid mode. Please enter 'encrypt' or 'decrypt'.")
	}
}

