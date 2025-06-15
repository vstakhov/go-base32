package main

import (
	"fmt"
	"log"

	"github.com/vstakhov/base32"
)

func main() {
	fmt.Println("Go Base32 Library - Rspamd Compatible")
	fmt.Println("=====================================")

	// Example data
	examples := []string{
		"hello",
		"test123",
		"hello world",
		"binary data example",
	}

	fmt.Println("\n1. ZBASE32 Encoding (Rspamd Compatible):")
	fmt.Println("----------------------------------------")
	for _, example := range examples {
		encoded := base32.EncodeString(example)
		decoded, err := base32.DecodeString(encoded)
		if err != nil {
			log.Printf("Error decoding %q: %v", encoded, err)
			continue
		}
		fmt.Printf("Input:   %q\n", example)
		fmt.Printf("Encoded: %q\n", encoded)
		fmt.Printf("Decoded: %q\n", string(decoded))
		fmt.Printf("Match:   %t\n\n", example == string(decoded))
	}

	fmt.Println("2. RFC 4648 Standard Base32:")
	fmt.Println("-----------------------------")
	for _, example := range examples {
		encoded := base32.EncodeAlphabetString(example, base32.RFC4648)
		decoded, err := base32.DecodeAlphabetString(encoded, base32.RFC4648)
		if err != nil {
			log.Printf("Error decoding %q: %v", encoded, err)
			continue
		}
		fmt.Printf("Input:   %q\n", example)
		fmt.Printf("Encoded: %q\n", encoded)
		fmt.Printf("Decoded: %q\n", string(decoded))
		fmt.Printf("Match:   %t\n\n", example == string(decoded))
	}

	fmt.Println("3. BECH32 Encoding:")
	fmt.Println("-------------------")
	for _, example := range examples {
		encoded := base32.EncodeAlphabetString(example, base32.BECH32)
		decoded, err := base32.DecodeAlphabetString(encoded, base32.BECH32)
		if err != nil {
			log.Printf("Error decoding %q: %v", encoded, err)
			continue
		}
		fmt.Printf("Input:   %q\n", example)
		fmt.Printf("Encoded: %q\n", encoded)
		fmt.Printf("Decoded: %q\n", string(decoded))
		fmt.Printf("Match:   %t\n\n", example == string(decoded))
	}

	fmt.Println("4. Custom Alphabet Example:")
	fmt.Println("----------------------------")
	customAlphabet, err := base32.NewAlphabet("0123456789ABCDEFGHIJKLMNOPQRSTUV", base32.OrderNormal)
	if err != nil {
		log.Printf("Error creating custom alphabet: %v", err)
		return
	}

	example := "custom example"
	encoded := base32.EncodeAlphabet([]byte(example), customAlphabet)
	decoded, err := base32.DecodeAlphabet([]byte(encoded), customAlphabet)
	if err != nil {
		log.Printf("Error decoding %q: %v", encoded, err)
		return
	}
	fmt.Printf("Input:   %q\n", example)
	fmt.Printf("Encoded: %q\n", encoded)
	fmt.Printf("Decoded: %q\n", string(decoded))
	fmt.Printf("Match:   %t\n\n", example == string(decoded))

	fmt.Println("5. Demonstrating Rspamd Bug Compatibility:")
	fmt.Println("-------------------------------------------")
	// This shows that our ZBASE32 implementation maintains the same "bug"
	// as Rspamd's implementation for compatibility
	rspamdExamples := []string{"test", "hello", "data"}
	for _, example := range rspamdExamples {
		zbase32 := base32.EncodeString(example)
		rfc4648 := base32.EncodeAlphabetString(example, base32.RFC4648)
		fmt.Printf("Input:        %q\n", example)
		fmt.Printf("ZBASE32:      %q (Rspamd compatible with reversed octets)\n", zbase32)
		fmt.Printf("RFC 4648:     %q (Standard base32)\n", rfc4648)
		fmt.Printf("Different:    %t\n\n", zbase32 != rfc4648)
	}
} 