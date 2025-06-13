package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/sven-seyfert/apiprobe/internal/logger"
)

// Obfuscate encodes a string with Base64 and random characters
// to produce a token.
func Obfuscate(data string) string {
	core := base64.StdEncoding.EncodeToString([]byte(data))
	core = strings.ReplaceAll(core, "=", "-")

	return fmt.Sprintf("ey%s.%s%s%s.%s", chars(9), chars(2), core, chars(6), chars(24)) //nolint:mnd
}

// Deobfuscate decodes a token back into its original plaintext.
func Deobfuscate(data string) string {
	if data == "" {
		return ""
	}

	core := data[14 : len(data)-31]
	core = strings.ReplaceAll(core, "-", "=")

	byteString, err := base64.StdEncoding.DecodeString(core)
	if err != nil {
		logger.Warnf("Decryption failed: %s", err)

		return ""
	}

	return string(byteString)
}

// chars generates a random alphanumeric string of the specified length.
func chars(length int) string {
	const alphaNum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	byteString := make([]byte, length)

	for idx := range byteString {
		nBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphaNum))))

		byteString[idx] = alphaNum[nBig.Int64()]
	}

	return string(byteString)
}

// HexHash returns a cryptographically secure random hex string of length 10.
// It reads the needed random bytes, encodes them, and truncates to 10 chars.
func HexHash() (string, error) {
	const charCount = 10

	byteLen := (charCount + 1) / 2 //nolint:mnd
	buf := make([]byte, byteLen)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	hexStr := hex.EncodeToString(buf)

	return hexStr[:charCount], nil
}
