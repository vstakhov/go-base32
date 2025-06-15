# Go Base32 Library - Rspamd Compatible

This is a Go translation of the Rust zbase32 library that provides **bug-to-bug compatibility** with [Rspamd](https://rspamd.com)'s base32 implementation.

## Key Features

- **Rspamd Compatibility**: Maintains the exact same "bytes flip" bug as Rspamd's C implementation
- **Multiple Alphabets**: Supports ZBASE32, RFC 4648, and BECH32 alphabets
- **High Performance**: Optimized encoding/decoding with minimal allocations
- **Type Safety**: Proper error handling and type-safe APIs
- **Zero Dependencies**: Pure Go implementation with no external dependencies

## Bug-to-Bug Compatibility

This implementation maintains compatibility with Rspamd's base32 encoding, which includes:

1. **Reversed Octets Order**: ZBASE32 encodes data in reversed octets order due to an initial bug in Rspamd
2. **No Padding**: RFC 4648 encoding doesn't include padding (following Rspamd's behavior)

This ensures that data encoded with Rspamd can be decoded with this library and vice versa.

## Installation

```bash
go get github.com/vstakhov/go-base32
```

## Usage

### Basic ZBASE32 Encoding (Rspamd Compatible)

```go
package main

import (
    "fmt"
    "github.com/vstakhov/go-base32"
)

func main() {
    // Encode a string
    encoded := base32.EncodeString("hello world")
    fmt.Println("Encoded:", encoded) // Output: em3ags7py376g3tprd
    
    // Decode back
    decoded, err := base32.DecodeString(encoded)
    if err != nil {
        panic(err)
    }
    fmt.Println("Decoded:", string(decoded)) // Output: hello world
}
```

### Using Different Alphabets

```go
// RFC 4648 Standard Base32
encoded := base32.EncodeAlphabetString("hello", base32.RFC4648)
fmt.Println("RFC 4648:", encoded) // Output: NBSWY3DP

// BECH32 Alphabet
encoded = base32.EncodeAlphabetString("hello", base32.BECH32)
fmt.Println("BECH32:", encoded) // Output: dpjkcmr0

// Custom Alphabet
alphabet, err := base32.NewAlphabet("0123456789ABCDEFGHIJKLMNOPQRSTUV", base32.OrderNormal)
if err != nil {
    panic(err)
}
encoded = base32.EncodeAlphabet([]byte("hello"), alphabet)
fmt.Println("Custom:", encoded)
```

### Working with Byte Slices

```go
// Encode to a pre-allocated slice
input := []byte("hello")
output := make([]byte, base32.EncodedLen(len(input)))
n := base32.EncodeToSlice(input, output, base32.ZBASE32)
fmt.Println("Encoded:", string(output[:n]))

// Decode to a pre-allocated slice
encoded := []byte("em3ags7p")
decoded := make([]byte, base32.DecodedLen(len(encoded)))
n, err := base32.DecodeAlphabetToSlice(encoded, decoded, base32.ZBASE32)
if err != nil {
    panic(err)
}
fmt.Println("Decoded:", string(decoded[:n]))
```

## API Reference

### Core Functions

- `Encode(input []byte) string` - Encode using ZBASE32 (Rspamd compatible)
- `Decode(input []byte) ([]byte, error)` - Decode using ZBASE32
- `EncodeString(input string) string` - String convenience wrapper for Encode
- `DecodeString(input string) ([]byte, error)` - String convenience wrapper for Decode

### Alphabet-Specific Functions

- `EncodeAlphabet(input []byte, alphabet *Alphabet) string`
- `DecodeAlphabet(input []byte, alphabet *Alphabet) ([]byte, error)`
- `EncodeAlphabetString(input string, alphabet *Alphabet) string`
- `DecodeAlphabetString(input string, alphabet *Alphabet) ([]byte, error)`

### Buffer Management Functions

- `EncodeToSlice(input []byte, output []byte, alphabet *Alphabet) int`
- `DecodeAlphabetToSlice(input []byte, buffer []byte, alphabet *Alphabet) (int, error)`
- `EncodedLen(bytesLen int) int` - Calculate encoded length
- `DecodedLen(bytesLen int) int` - Calculate maximum decoded length

### Alphabet Management

- `NewAlphabet(alphabet string, order EncodeOrder) (*Alphabet, error)`

### Predefined Alphabets

- `ZBASE32` - ZBase32 alphabet with reversed order (Rspamd compatible)
- `RFC4648` - RFC 4648 standard base32 alphabet
- `BECH32` - Bitcoin BECH32 alphabet

## Performance

The library is optimized for performance:

```
BenchmarkEncodeZBase32-16       10721727     95.51 ns/op
BenchmarkDecodeZBase32-16       10683105    109.1 ns/op
BenchmarkEncodeRFC4648-16       12482202     95.10 ns/op
BenchmarkDecodeRFC4648-16       12807644     88.49 ns/op
```

## Compatibility Testing

The implementation includes comprehensive tests that verify compatibility with the original Rust implementation:

```bash
go test -v
```

Key test cases include:
- Basic encoding/decoding roundtrips
- Compatibility with known Rust test vectors
- Error handling for invalid inputs
- All three alphabet types (ZBASE32, RFC4648, BECH32)
- Custom alphabet creation and validation

## The Rspamd Bug

The "bytes flip" bug in Rspamd's base32 implementation causes the encoded data to have reversed octets order compared to standard base32 implementations. This library maintains this exact behavior when using the ZBASE32 alphabet with `OrderInversed` encoding order.

**Example showing the difference:**

```go
input := "test"
zbase32 := base32.EncodeString(input)        // "wm3g84b" (Rspamd compatible)
rfc4648 := base32.EncodeAlphabetString(input, base32.RFC4648) // "ORSXG5A" (Standard)
fmt.Printf("Different: %t\n", zbase32 != rfc4648) // true
```

## License

This project is licensed under Apache 2.0 license, same as the original Rust implementation.

## Acknowledgments

This Go implementation is a faithful translation of the original Rust zbase32 library, maintaining full compatibility with Rspamd's base32 encoding behavior. 
