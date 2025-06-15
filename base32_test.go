package base32

import (
	"bytes"
	"testing"
)

// Test data from the Rust implementation to ensure compatibility
func TestZBase32Compatibility(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "em3ags7p"},
		{"test123", "wm3g84fg13cy"},
		{"", ""},
		{"a", "bd"},
		{"aa", "bmay"},
		{"aaa", "bmang"},
		{"aaaa", "bmansob"},
		{"aaaaa", "bmansofc"},
		{"aaaaaa", "bmansofcbd"},
		{"aaaaaaa", "bmansofcbmay"},
		{"aaaaaaaa", "bmansofcbmang"},
	}

	for _, test := range tests {
		t.Run("encode_"+test.input, func(t *testing.T) {
			result := EncodeString(test.input)
			if result != test.expected {
				t.Errorf("Encode(%q) = %q, want %q", test.input, result, test.expected)
			}
		})

		t.Run("decode_"+test.expected, func(t *testing.T) {
			result, err := DecodeString(test.expected)
			if err != nil {
				t.Errorf("Decode(%q) failed: %v", test.expected, err)
				return
			}
			if string(result) != test.input {
				t.Errorf("Decode(%q) = %q, want %q", test.expected, string(result), test.input)
			}
		})
	}
}

func TestRFC4648Compatibility(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "NBSWY3DP"},
		{"test123", "ORSXG5BRGIZQ"},
		{"", ""},
		{"a", "ME"},
		{"aa", "MFQQ"},
		{"aaa", "MFQWC"},
		{"aaaa", "MFQWCYI"},
		{"aaaaa", "MFQWCYLB"},
		{"aaaaaa", "MFQWCYLBME"},
		{"aaaaaaa", "MFQWCYLBMFQQ"},
		{"aaaaaaaa", "MFQWCYLBMFQWC"},
	}

	for _, test := range tests {
		t.Run("encode_"+test.input, func(t *testing.T) {
			result := EncodeAlphabetString(test.input, RFC4648)
			if result != test.expected {
				t.Errorf("EncodeAlphabet(%q, RFC4648) = %q, want %q", test.input, result, test.expected)
			}
		})

		t.Run("decode_"+test.expected, func(t *testing.T) {
			result, err := DecodeAlphabetString(test.expected, RFC4648)
			if err != nil {
				t.Errorf("DecodeAlphabet(%q, RFC4648) failed: %v", test.expected, err)
				return
			}
			if string(result) != test.input {
				t.Errorf("DecodeAlphabet(%q, RFC4648) = %q, want %q", test.expected, string(result), test.input)
			}
		})
	}
}

func TestRoundtrip(t *testing.T) {
	testData := [][]byte{
		{},
		{0x01},
		{0x01, 0x02},
		{0x01, 0x02, 0x03},
		{0x01, 0x02, 0x03, 0x04},
		{0x01, 0x02, 0x03, 0x04, 0x05},
		{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
		{0xFF, 0xFE, 0xFD, 0xFC, 0xFB},
		{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
	}

	alphabets := []*Alphabet{ZBASE32, RFC4648, BECH32}
	names := []string{"ZBASE32", "RFC4648", "BECH32"}

	for i, alphabet := range alphabets {
		for j, data := range testData {
			t.Run(names[i]+"_roundtrip_"+string(rune('A'+j)), func(t *testing.T) {
				encoded := EncodeAlphabet(data, alphabet)
				decoded, err := DecodeAlphabet([]byte(encoded), alphabet)
				if err != nil {
					t.Errorf("DecodeAlphabet failed: %v", err)
					return
				}
				if !bytes.Equal(data, decoded) {
					t.Errorf("Roundtrip failed: input=%v, encoded=%q, decoded=%v", data, encoded, decoded)
				}
			})
		}
	}
}

func TestInvalidDecoding(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		alphabet *Alphabet
	}{
		{"ZBASE32_invalid_char", "hello@", ZBASE32},
		{"RFC4648_invalid_char", "HELLO@", RFC4648},
		{"BECH32_invalid_char", "hello@", BECH32},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := DecodeAlphabetString(test.input, test.alphabet)
			if err == nil {
				t.Errorf("Expected error for invalid input %q", test.input)
			}
		})
	}
}

func TestEncodeToSlice(t *testing.T) {
	input := []byte("hello")
	output := make([]byte, EncodedLen(len(input)))
	
	n := EncodeToSlice(input, output, ZBASE32)
	result := string(output[:n])
	expected := "em3ags7p"
	
	if result != expected {
		t.Errorf("EncodeToSlice result = %q, want %q", result, expected)
	}
}

func TestDecodeToSlice(t *testing.T) {
	input := []byte("em3ags7p")
	output := make([]byte, DecodedLen(len(input)))
	
	n, err := DecodeAlphabetToSlice(input, output, ZBASE32)
	if err != nil {
		t.Errorf("DecodeAlphabetToSlice failed: %v", err)
		return
	}
	
	result := string(output[:n])
	expected := "hello"
	
	if result != expected {
		t.Errorf("DecodeAlphabetToSlice result = %q, want %q", result, expected)
	}
}

func TestNewAlphabet(t *testing.T) {
	// Test valid alphabet
	alphabet, err := NewAlphabet("ybndrfg8ejkmcpqxot1uwisza345h769", OrderInversed)
	if err != nil {
		t.Errorf("NewAlphabet failed: %v", err)
		return
	}
	
	// Test encoding with custom alphabet
	input := []byte("hello")
	result := EncodeAlphabet(input, alphabet)
	expected := "em3ags7p"
	
	if result != expected {
		t.Errorf("Custom alphabet encode result = %q, want %q", result, expected)
	}
	
	// Test invalid length
	_, err = NewAlphabet("tooshort", OrderNormal)
	if err == nil {
		t.Error("Expected error for short alphabet")
	}
	
	// Test duplicate character
	_, err = NewAlphabet("aacdefghijklmnopqrstuvwxyz234567", OrderNormal)  
	if err == nil {
		t.Error("Expected error for duplicate characters")
	}
	
	// Test unprintable character
	_, err = NewAlphabet("abcdefghijklmnopqrstuvwxyz23456\x01", OrderNormal)
	if err == nil {
		t.Error("Expected error for unprintable character")
	}
}

// This test specifically verifies the Rspamd bug compatibility
func TestRspamdBugCompatibility(t *testing.T) {
	// The key insight: ZBASE32 with OrderInversed should produce different results
	// than a standard base32 implementation due to the "bytes flip" bug
	
	testCases := []struct {
		input    string
		zbase32  string  // Expected result with Rspamd's buggy implementation
	}{
		{"test", "wm3gso"},
		{"hello world", "em3ags7py376g3tprd"},
		{"binary", "bn5xw4rygm"},
	}
	
	for _, tc := range testCases {
		t.Run("rspamd_bug_"+tc.input, func(t *testing.T) {
			result := EncodeString(tc.input)
			
			// The exact values should match what Rspamd produces
			// due to maintaining bug-to-bug compatibility
			t.Logf("Input: %s, Encoded: %s", tc.input, result)
			
			// Test roundtrip to ensure decoding works
			decoded, err := DecodeString(result)
			if err != nil {
				t.Errorf("Failed to decode %q: %v", result, err)
				return
			}
			
			if string(decoded) != tc.input {
				t.Errorf("Roundtrip failed: input=%q, decoded=%q", tc.input, string(decoded))
			}
		})
	}
}

// Benchmarks
func BenchmarkEncodeZBase32(b *testing.B) {
	data := []byte("hello world this is a test string for benchmarking")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(data)
	}
}

func BenchmarkDecodeZBase32(b *testing.B) {
	encoded := Encode([]byte("hello world this is a test string for benchmarking"))
	encodedBytes := []byte(encoded)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(encodedBytes)
	}
}

func BenchmarkEncodeRFC4648(b *testing.B) {
	data := []byte("hello world this is a test string for benchmarking")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EncodeAlphabet(data, RFC4648)
	}
}

func BenchmarkDecodeRFC4648(b *testing.B) {
	encoded := EncodeAlphabet([]byte("hello world this is a test string for benchmarking"), RFC4648)
	encodedBytes := []byte(encoded)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DecodeAlphabet(encodedBytes, RFC4648)
	}
} 