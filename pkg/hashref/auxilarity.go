package hashref

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

// GetHashTypeAndValue identifies if the provided input is a
// text, file or sha256 hashsum
func GetHashTypeAndValue(input string) (HashType, string) {
	dat, err := os.ReadFile(input)
	if err != nil {
		return analyzeText(input)
	}
	log.Println("Input is a File!")
	return File, CalculateHash(dat)
}

// analyzeText identifies if the provided input is a sha256 hashsum
// or just a normal text
func analyzeText(input string) (HashType, string) {
	_, err := hex.DecodeString(input)
	if len(input) != 64 || err != nil {
		log.Println("Input is a Text!")
		return Text, CalculateHash([]byte(input))
	}
	log.Println("Input is a Hash!")
	return Hash, input
}

// CalculateHash calculates the sha256 hashsum to a provided byte slice
func CalculateHash(raw []byte) string {
	log.Printf("Calculate hash from %v bytes\n", len(raw))
	h := sha256.New()
	h.Write(raw)
	return fmt.Sprintf("%x", h.Sum(nil))
}
