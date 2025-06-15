// Package base32 provides base32 encoding/decoding functionality compatible with Rspamd
// 
// This package maintains bug-to-bug compatibility with Rspamd's base32 implementation,
// including the reversed octets order bug in ZBASE32 encoding.
package base32

import (
	"fmt"
)

const AlphabetSize = 32

// EncodeOrder represents the bit ordering during encoding
type EncodeOrder int

const (
	OrderNormal EncodeOrder = iota
	OrderInversed
)

// Alphabet defines the 32 characters used for Base32 encoding
type Alphabet struct {
	encodeSymbols [AlphabetSize]byte
	decodeBytes   [256]byte
	encodeOrder   EncodeOrder
}

// NewAlphabet creates a new alphabet from a string with specified encode order
func NewAlphabet(alphabet string, order EncodeOrder) (*Alphabet, error) {
	if len(alphabet) != AlphabetSize {
		return nil, fmt.Errorf("invalid length - must be %d bytes", AlphabetSize)
	}

	var symbols [AlphabetSize]byte
	var decodeBytes [256]byte
	duplicates := make(map[byte]bool)

	// Initialize decode table with invalid values
	for i := range decodeBytes {
		decodeBytes[i] = 0xff
	}

	for i, b := range []byte(alphabet) {
		// Check for printable characters
		if b < 32 || b > 126 {
			return nil, fmt.Errorf("unprintable byte: %#04x", b)
		}

		// Check for duplicates
		if duplicates[b] {
			return nil, fmt.Errorf("duplicated byte: %#04x", b)
		}
		duplicates[b] = true

		symbols[i] = b
		decodeBytes[b] = byte(i)
	}

	return &Alphabet{
		encodeSymbols: symbols,
		decodeBytes:   decodeBytes,
		encodeOrder:   order,
	}, nil
}

// Helper function to initialize decode tables
func initDecodeTable(alphabet string) [256]byte {
	var decodeBytes [256]byte
	// Initialize all entries to invalid (0xff)
	for i := range decodeBytes {
		decodeBytes[i] = 0xff
	}
	// Set valid entries
	for i, b := range []byte(alphabet) {
		decodeBytes[b] = byte(i)
	}
	return decodeBytes
}

// Predefined alphabets

// ZBASE32 alphabet with reversed order for Rspamd compatibility
// This maintains the bug-to-bug compatibility with Rspamd's implementation
var ZBASE32 = &Alphabet{
	encodeSymbols: [32]byte{'y', 'b', 'n', 'd', 'r', 'f', 'g', '8', 'e', 'j', 'k', 'm', 'c', 'p', 'q', 'x', 'o', 't', '1', 'u', 'w', 'i', 's', 'z', 'a', '3', '4', '5', 'h', '7', '6', '9'},
	decodeBytes:   initDecodeTable("ybndrfg8ejkmcpqxot1uwisza345h769"),
	encodeOrder:   OrderInversed,
}

// RFC4648 alphabet with normal order
var RFC4648 = &Alphabet{
	encodeSymbols: [32]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '2', '3', '4', '5', '6', '7'},
	decodeBytes:   initDecodeTable("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"),
	encodeOrder:   OrderNormal,
}

// BECH32 alphabet with normal order
var BECH32 = &Alphabet{
	encodeSymbols: [32]byte{'q', 'p', 'z', 'r', 'y', '9', 'x', '8', 'g', 'f', '2', 't', 'v', 'd', 'w', '0', 's', '3', 'j', 'n', '5', '4', 'k', 'h', 'c', 'e', '6', 'm', 'u', 'a', '7', 'l'},
	decodeBytes:   initDecodeTable("qpzry9x8gf2tvdw0s3jn54khce6mua7l"),
	encodeOrder:   OrderNormal,
}

// DecodeError represents an error during decoding
type DecodeError struct {
	Msg    string
	Offset int
	Byte   byte
}

func (e DecodeError) Error() string {
	if e.Offset >= 0 {
		return fmt.Sprintf("%s at offset %d, byte %#02x", e.Msg, e.Offset, e.Byte)
	}
	return e.Msg
}

// EncodedLen returns the length of encoding for the given input length
func EncodedLen(bytesLen int) int {
	minBytes := bytesLen / 5
	rem := bytesLen % 5
	return minBytes*8 + rem*2 + 1
}

// EncodeToSlice encodes input using the specified alphabet into the output buffer
// Returns the number of bytes written to the output buffer
func EncodeToSlice(input []byte, output []byte, alphabet *Alphabet) int {
	encodeTable := alphabet.encodeSymbols
	remain := int32(-1)
	o := 0

	if alphabet.encodeOrder == OrderInversed {
		// Rspamd compatible encoding with reversed bit order
		for i := 0; i < len(input); i++ {
			switch i % 5 {
			case 0:
				// 8 bits of input and 3 to remain
				x := int32(input[i])
				output[o] = encodeTable[x&0x1F]
				o++
				remain = x >> 5
			case 1:
				// 11 bits of input, 1 to remain
				inp := int32(input[i])
				x := remain | (inp << 3)
				output[o] = encodeTable[x&0x1F]
				o++
				output[o] = encodeTable[(x>>5)&0x1F]
				o++
				remain = x >> 10
			case 2:
				// 9 bits of input, 4 to remain
				inp := int32(input[i])
				x := remain | (inp << 1)
				output[o] = encodeTable[x&0x1F]
				o++
				remain = x >> 5
			case 3:
				// 12 bits of input, 2 to remain
				inp := int32(input[i])
				x := remain | (inp << 4)
				output[o] = encodeTable[x&0x1F]
				o++
				output[o] = encodeTable[(x>>5)&0x1F]
				o++
				remain = (x >> 10) & 0x3
			case 4:
				// 10 bits of output, nothing to remain
				inp := int32(input[i])
				x := remain | (inp << 2)
				output[o] = encodeTable[x&0x1F]
				o++
				output[o] = encodeTable[(x>>5)&0x1F]
				o++
				remain = -1
			}
		}
	} else {
		// Standard encoding
		for i := 0; i < len(input); i++ {
			switch i % 5 {
			case 0:
				// 8 bits of input and 3 to remain
				inp := int32(input[i])
				x := inp >> 3
				output[o] = encodeTable[x&0x1F]
				o++
				remain = (inp & 7) << 2
			case 1:
				// 11 bits of input, 1 to remain
				inp := int32(input[i])
				x := (remain << 6) | inp
				output[o] = encodeTable[(x>>6)&0x1F]
				o++
				output[o] = encodeTable[(x>>1)&0x1F]
				o++
				remain = (x & 0x1) << 4
			case 2:
				// 9 bits of input, 4 to remain
				inp := int32(input[i])
				x := (remain << 4) | inp
				output[o] = encodeTable[(x>>4)&0x1F]
				o++
				remain = (x & 15) << 1
			case 3:
				// 12 bits of input, 2 to remain
				inp := int32(input[i])
				x := (remain << 7) | inp
				output[o] = encodeTable[(x>>7)&0x1F]
				o++
				output[o] = encodeTable[(x>>2)&0x1F]
				o++
				remain = (x & 3) << 3
			case 4:
				// 10 bits of output, nothing to remain
				inp := int32(input[i])
				x := (remain << 5) | inp
				output[o] = encodeTable[(x>>5)&0x1F]
				o++
				output[o] = encodeTable[x&0x1F]
				o++
				remain = -1
			}
		}
	}

	if remain >= 0 {
		output[o] = encodeTable[remain&0x1F]
		o++
	}

	return o
}

// EncodeAlphabet encodes input using the specified alphabet
func EncodeAlphabet(input []byte, alphabet *Alphabet) string {
	encodedSize := EncodedLen(len(input))
	buf := make([]byte, encodedSize)
	encLen := EncodeToSlice(input, buf, alphabet)
	return string(buf[:encLen])
}

// Encode encodes input using the default ZBASE32 alphabet (Rspamd compatible)
func Encode(input []byte) string {
	return EncodeAlphabet(input, ZBASE32)
}

// DecodedLen returns the maximum decoded length for the given encoded length
func DecodedLen(bytesLen int) int {
	fullChunks := bytesLen / 8
	remainder := bytesLen % 8
	return fullChunks*5 + remainder
}

// DecodeAlphabet decodes input using the specified alphabet
func DecodeAlphabet(input []byte, alphabet *Alphabet) ([]byte, error) {
	buffer := make([]byte, DecodedLen(len(input)))
	actualLen, err := DecodeAlphabetToSlice(input, buffer, alphabet)
	if err != nil {
		return nil, err
	}
	return buffer[:actualLen], nil
}

// DecodeAlphabetToSlice decodes input using the specified alphabet into the provided buffer
// Returns the number of bytes written and any error
func DecodeAlphabetToSlice(input []byte, buffer []byte, alphabet *Alphabet) (int, error) {
	processedBits := 0
	acc := uint32(0)
	o := 0

	if alphabet.encodeOrder == OrderInversed {
		// Rspamd compatible decoding with reversed bit order
		for i, c := range input {
			if processedBits >= 8 {
				// Emit from left to right
				processedBits -= 8
				buffer[o] = byte(acc & 0xFF)
				o++
				acc = acc >> 8
			}

			decoded := alphabet.decodeBytes[c]
			if decoded == 0xff {
				return 0, DecodeError{
					Msg:    "invalid byte",
					Offset: i,
					Byte:   c,
				}
			}

			acc = (uint32(decoded) << processedBits) | acc
			processedBits += 5
		}

		if processedBits > 0 {
			buffer[o] = byte(acc & 0xFF)
			o++
		}
	} else {
		// Standard decoding
		for i, c := range input {
			decoded := alphabet.decodeBytes[c]
			if decoded == 0xff {
				return 0, DecodeError{
					Msg:    "invalid byte",
					Offset: i,
					Byte:   c,
				}
			}

			acc = (acc << 5) | uint32(decoded)
			processedBits += 5

			if processedBits >= 8 {
				processedBits -= 8
				// Emit from right to left
				buffer[o] = byte((acc >> processedBits) & 0xFF)
				o++
				acc = acc & ((1 << processedBits) - 1)
			}
		}
	}

	return o, nil
}

// Decode decodes input using the default ZBASE32 alphabet (Rspamd compatible)
func Decode(input []byte) ([]byte, error) {
	return DecodeAlphabet(input, ZBASE32)
}

// String-based convenience functions

// EncodeString encodes a string using the default ZBASE32 alphabet
func EncodeString(input string) string {
	return Encode([]byte(input))
}

// DecodeString decodes a string using the default ZBASE32 alphabet
func DecodeString(input string) ([]byte, error) {
	return Decode([]byte(input))
}

// EncodeAlphabetString encodes a string using the specified alphabet
func EncodeAlphabetString(input string, alphabet *Alphabet) string {
	return EncodeAlphabet([]byte(input), alphabet)
}

// DecodeAlphabetString decodes a string using the specified alphabet
func DecodeAlphabetString(input string, alphabet *Alphabet) ([]byte, error) {
	return DecodeAlphabet([]byte(input), alphabet)
} 