package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Aragon2Config struct {
	Hash       []byte
	Salt       []byte
	TimeCost   uint32
	MemoryCont uint32
	Thread     uint8
	KeyLen     uint32
}

func GenerateHashPassword(password string) (string, error) {
	cf := &Aragon2Config{
		TimeCost:   2,
		MemoryCont: 64 * 1024,
		Thread:     4,
		KeyLen:     32,
	}
	salt, saltGenError := generateSalt(16)
	if saltGenError != nil {
		return "", fmt.Errorf("Password hashing failed: %w", saltGenError)
	}
	cf.Salt = salt
	cf.Hash = argon2.IDKey(
		[]byte(password),
		cf.Salt, // randomly generated 128 bit
		cf.TimeCost,
		// The memory usage.
		// Higher values improve security but increase resource usage.
		cf.MemoryCont,
		cf.Thread, // this correspond with core count
		cf.KeyLen, // password hash byte length
	)
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		cf.MemoryCont,
		cf.TimeCost,
		cf.Thread,
		base64.RawStdEncoding.EncodeToString(cf.Salt),
		base64.RawStdEncoding.EncodeToString(cf.Hash),
	)
	return encodedHash, nil
}

func ComparePasswordHash(hashedPassword, password string) error {
	/// aaaaaaa why
	cf, parseErr := parseArgon2Hash(hashedPassword)
	if parseErr != nil {
		return parseErr
	}
	ProvidedPassHash := argon2.IDKey(
		[]byte(password),
		cf.Salt, // randomly generated 128 bit
		cf.TimeCost,
		// The memory usage.
		// Higher values improve security but increase resource usage.
		cf.MemoryCont,
		cf.Thread, // this correspond with core count
		cf.KeyLen, // password hash byte length
	)
	if match := subtle.ConstantTimeCompare(cf.Hash, ProvidedPassHash); match != 1 {
		return fmt.Errorf("password did not match the hashed pass")
	}
	return nil
}

// generates a random bits of provided size
func generateSalt(size uint32) ([]byte, error) {
	salt := make([]byte, size)
	if _, randomizeErr := rand.Read(salt); randomizeErr != nil {
		return nil, fmt.Errorf("salt generation error: %w", randomizeErr)
	}
	return salt, nil
}

func parseArgon2Hash(hashedPassword string) (*Aragon2Config, error) {
	// "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
	config := &Aragon2Config{}
	par := strings.Split(hashedPassword, "$")
	if len(par) != 6 {
		return nil, fmt.Errorf("invalid hash format")

	}
	// format check here

	fmt.Println()

	return config, nil
}
