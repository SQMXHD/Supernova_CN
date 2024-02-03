package Encryptors

import (
	"Supernova/Converters"
	"Supernova/Output"
	"Supernova/Maldev/Crypto"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
)

// Rc4Context represents the state of the RC4 encryption algorithm.
type Rc4Context struct {
	i uint32
	j uint32
	s [256]uint8
}

const (
	// chars defines the set of characters used to generate a random key and IV.
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+{}[]"

	// keySize specifies the size (in bytes) of the encryption key.
	keySize = 32

	// ivSize specifies the size (in bytes) of the initialization vector (IV).
	ivSize = 16
)

// PKCS7Padding function
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// AESEncryption function
func AESEncryption(key []byte, iv []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Apply PKCS7 padding to ensure plaintext length is a multiple of the block size
	paddedData := PKCS7Padding(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(paddedData))

	// Create a new CBC mode encrypter
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

	return ciphertext, nil
}

// Chacha20Encryption function
func Chacha20Encryption(plaintext []byte, key []byte) ([]byte, error) {
    
	ciphertext, err := crypto.Chacha20Encrypt(plaintext, key) // Encrypt
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

// GenerateRandomBytes function
func GenerateRandomBytes(length int) []byte {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("[!] Failed to generate a random key.")
	}
	return randomBytes
}

// GenerateRandomPassphrase function
func GenerateRandomPassphrase(length int) string {

	charSetLength := big.NewInt(int64(len(chars)))
	passphrase := make([]byte, length)

	for i := range passphrase {
		randomIndex, err := rand.Int(rand.Reader, charSetLength)
		if err != nil {
			fmt.Println("Error generating random number:", err)
			return ""
		}
		passphrase[i] = chars[randomIndex.Int64()]
	}

	return string(passphrase)
}

// XOREncryption function performs XOR encryption on input shellcode using a multi xor key.
func XOREncryption(shellcode []byte, key []byte) []byte {
	encrypted := make([]byte, len(shellcode))
	keyLen := len(key)

	for i := 0; i < len(shellcode); i++ {
		encrypted[i] = shellcode[i] ^ key[i%keyLen]
	}

	return encrypted
}

// RC4Encryption function implements the RC4 encryption algorithm
func RC4Encryption(data []byte, key []byte) []byte {
	var s [256]byte

	// Initialize the S array with values from 0 to 255
	for i := 0; i < 256; i++ {
		s[i] = byte(i)
	}

	j := 0
	// KSA (Key Scheduling Algorithm) - Initial permutation of S array based on the key
	for i := 0; i < 256; i++ {
		j = (j + int(s[i]) + int(key[i%len(key)])) % 256
		s[i], s[j] = s[j], s[i]
	}

	encrypted := make([]byte, len(data))
	i, j := 0, 0
	// PRGA (Pseudo-Random Generation Algorithm) - Generate encrypted output
	for k := 0; k < len(data); k++ {
		i = (i + 1) % 256
		j = (j + int(s[i])) % 256
		s[i], s[j] = s[j], s[i]
		// XOR encrypted byte with generated pseudo-random byte from S array
		encrypted[k] = data[k] ^ s[(int(s[i])+int(s[j]))%256]
	}

	return encrypted
}

// CaesarEncryption function implements the Caesar encryption algorithm
func CaesarEncryption(shellcode []byte, shift int) []byte {
	encrypted := make([]byte, len(shellcode))
	for i, char := range shellcode {
		// Apply Caesar cipher encryption
		encrypted[i] = byte((int(char) + shift) % 256)
		//encryptedChar := char + byte(shift)
		//encrypted[i] = encryptedChar
	}
	return encrypted
}

// DetectEncryption function
func DetectEncryption(cipher string, shellcode string, key int) (string, int, []byte, string, []byte) {
	// Set logger for errors
	logger := log.New(os.Stderr, "[!] ", 0)

	// Set cipher to lower
	cipher = strings.ToLower(cipher)

	// Convert shellcode to bytes
	shellcodeInBytes := []byte(shellcode)

	// Set key size
	shift := key

	switch cipher {
	case "chacha20":
		// Call function named GenerateRandomBytes
		chacha20Key := GenerateRandomBytes(32)

		// Print generated Chacha2 key
		fmt.Printf("[+] Generated Chacha20 key: ")

		// Call function named PrintKeyDetails
		Output.PrintKeyDetails(chacha20Key)

		// Call function named Chacha20Encryption
		encryptedShellcode, err := Chacha20Encryption(shellcodeInBytes, chacha20Key)
		if err != nil {
			panic(err)
		}

		// Call function named FormatShellcode
		shellcodeFormatted := Converters.FormatShellcode(encryptedShellcode)

		return shellcodeFormatted, len(encryptedShellcode), chacha20Key, "", nil
	case "b64chacha20":
		// Call function named GenerateRandomBytes
		chacha20Key := GenerateRandomBytes(32)

		// Print generated Chacha2 key
		fmt.Printf("[+] Generated Chacha20 key: ")

		// Call function named PrintKeyDetails
		Output.PrintKeyDetails(chacha20Key)

		// Call function named Chacha20Encryption
		encryptedShellcode, err := Chacha20Encryption(shellcodeInBytes, chacha20Key)
		if err != nil {
			panic(err)
		}

		// Convert encryptedShellcode to Base64
		encryptedShellcodeBase64 := base64.StdEncoding.EncodeToString(encryptedShellcode)

		// Convert Base64 string to []byte
		encryptedShellcodeBytes := []byte(encryptedShellcodeBase64)
		
		// Print length changed notification
		fmt.Printf("[+] New Payload size: %d bytes\n\n", len(encryptedShellcodeBytes))

		// Call function named FormatShellcode
		shellcodeFormatted := Converters.FormatShellcode(encryptedShellcodeBytes)

		return shellcodeFormatted, len(encryptedShellcode), chacha20Key, "", nil
	case "xor":
		// Call function named GenerateRandomBytes
		xorKey := GenerateRandomBytes(shift)

		// Print generated XOR key
		fmt.Printf("[+] Generated XOR key: ")

		// Call function named PrintKeyDetails
		Output.PrintKeyDetails(xorKey)

		// Call function named XOREncryption
		encryptedShellcode := XOREncryption(shellcodeInBytes, xorKey)

		// Call function named FormatShellcode
		shellcodeFormatted := Converters.FormatShellcode(encryptedShellcode)

		return shellcodeFormatted, len(encryptedShellcode), xorKey, "", nil
	case "b64xor":
		// Call function named GenerateRandomBytes
		xorKey := GenerateRandomBytes(shift)

		// Print generated XOR key
		fmt.Printf("[+] Generated XOR key: ")

		// Call function named PrintKeyDetails
		Output.PrintKeyDetails(xorKey)

		// Call function named XOREncryption
		encryptedShellcode := XOREncryption(shellcodeInBytes, xorKey)

		// Convert encryptedShellcode to Base64
		encryptedShellcodeBase64 := base64.StdEncoding.EncodeToString(encryptedShellcode)

		// Convert Base64 string to []byte
		encryptedShellcodeBytes := []byte(encryptedShellcodeBase64)

		// Print length changed notification
		fmt.Printf("[+] New Payload size: %d bytes\n\n", len(encryptedShellcodeBytes))

		// Call function named FormatShellcode
		shellcodeFormatted := Converters.FormatShellcode(encryptedShellcodeBytes)

		return shellcodeFormatted, len(encryptedShellcode), xorKey, "", nil
	case "rot":
		// Print selected shift key
		fmt.Printf("[+] Selected Shift key: %d\n\n", shift)

		// Call function named CaesarEncryption
		encryptedShellcode := CaesarEncryption(shellcodeInBytes, shift)

		// Call function named FormatShellcode
		shellcodeFormatted := Converters.FormatShellcode(encryptedShellcode)

		return shellcodeFormatted, len(encryptedShellcode), nil, "", nil
	case "aes":
		// Generate a random 32-byte key and a random 16-byte IV
		key := GenerateRandomBytes(keySize)
		iv := GenerateRandomBytes(ivSize)

		// Print generated key
		fmt.Printf("[+] Generated key (32-byte): ")

		// Call function named PrintKeyDetails
		Output.PrintKeyDetails(key)

		// Print generated key
		fmt.Printf("[+] Generated IV (16-byte): ")

		// Call function named PrintKeyDetails
		Output.PrintKeyDetails(iv)

		// Print AES-256-CBC notification
		fmt.Printf("[+] Using AES-256-CBC encryption\n\n")

		// Encrypt the shellcode using AES-256-CBC
		encryptedShellcode, err := AESEncryption(key, iv, shellcodeInBytes)
		if err != nil {
			panic(err)
		}

		// Print length changed notification
		fmt.Printf("[+] New Payload size: %d bytes\n\n", len(encryptedShellcode))

		// Call function named FormatShellcode
		shellcodeFormatted := Converters.FormatShellcode(encryptedShellcode)

		return shellcodeFormatted, len(encryptedShellcode), key, "", iv
	case "b64aes":
		// Generate a random 32-byte key and a random 16-byte IV
		key := GenerateRandomBytes(keySize)
		iv := GenerateRandomBytes(ivSize)

		// Print generated key
		fmt.Printf("[+] Generated key (32-byte): ")

		// Call function named PrintKeyDetails
		Output.PrintKeyDetails(key)

		// Print generated key
		fmt.Printf("[+] Generated IV (16-byte): ")

		// Call function named PrintKeyDetails
		Output.PrintKeyDetails(iv)

		// Print AES-256-CBC notification
		fmt.Printf("[+] Using AES-256-CBC encryption\n\n")

		// Encrypt the shellcode using AES-256-CBC
		encryptedShellcode, err := AESEncryption(key, iv, shellcodeInBytes)
		if err != nil {
			panic(err)
		}

		// Print length changed notification
		fmt.Printf("[+] New Payload size: %d bytes\n\n", len(encryptedShellcode))

		// Convert encryptedShellcode to Base64
		encryptedShellcodeBase64 := base64.StdEncoding.EncodeToString(encryptedShellcode)

		// Convert Base64 string to []byte
		encryptedShellcodeBytes := []byte(encryptedShellcodeBase64)

		// Print length changed notification
		fmt.Printf("[+] New Payload size: %d bytes\n\n", len(encryptedShellcodeBytes))

		// Call function named FormatShellcode
		shellcodeFormatted := Converters.FormatShellcode(encryptedShellcodeBytes)

		return shellcodeFormatted, len(encryptedShellcode), key, "", iv
	case "rc4":
		// Call function named GenerateRandomPassphrase
		randomPassphrase := GenerateRandomPassphrase(key)

		// Convert passphrase to bytes
		rc4Key := []byte(randomPassphrase)

		// Print generated passphrase
		fmt.Printf("[+] Generated passphrase: %s\n\n", randomPassphrase)

		// Call function named RC4Encryption
		encryptedShellcode := RC4Encryption(shellcodeInBytes, rc4Key)

		// Call function named FormatShellcode
		shellcodeFormatted := Converters.FormatShellcode(encryptedShellcode)

		return shellcodeFormatted, len(encryptedShellcode), rc4Key, randomPassphrase, nil
	case "b64rc4":
		// Call function named GenerateRandomPassphrase
		randomPassphrase := GenerateRandomPassphrase(key)

		// Convert passphrase to bytes
		rc4Key := []byte(randomPassphrase)

		// Print generated passphrase
		fmt.Printf("[+] Generated passphrase: %s\n\n", randomPassphrase)

		// Call function named RC4Encryption
		encryptedShellcode := RC4Encryption(shellcodeInBytes, rc4Key)

		// Convert encryptedShellcode to Base64
		encryptedShellcodeBase64 := base64.StdEncoding.EncodeToString(encryptedShellcode)
		
		// Convert Base64 string to []byte
		encryptedShellcodeBytes := []byte(encryptedShellcodeBase64)

		// Print length changed notification
		fmt.Printf("[+] New Payload size: %d bytes\n\n", len(encryptedShellcodeBytes))

		// Call function named FormatShellcode
		shellcodeFormatted := Converters.FormatShellcode(encryptedShellcodeBytes)

		return shellcodeFormatted, len(encryptedShellcode), rc4Key, randomPassphrase, nil
	default:
		logger.Fatal("Unsupported encryption cipher")
		return "", 0, nil, "", nil
	}
}
