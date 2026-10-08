package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2Config struct {
	Hash       []byte
	Salt       []byte
	TimeCost   uint32
	MemoryCost uint32
	Thread     uint8
	KeyLen     uint32
}

func GenerateHashPassword(password string) (string, error) {
	cf := &Argon2Config{
		TimeCost:   3,
		MemoryCost: 64 * 1024,
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
		cf.MemoryCost,
		cf.Thread, // this correspond with core count
		cf.KeyLen, // password hash byte length
	)
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		cf.MemoryCost,
		cf.TimeCost,
		cf.Thread,
		base64.RawStdEncoding.EncodeToString(cf.Salt),
		base64.RawStdEncoding.EncodeToString(cf.Hash),
	)
	return encodedHash, nil
}

// generates a random bits of provided size
func generateSalt(size uint32) ([]byte, error) {
	salt := make([]byte, size)
	if _, randomizeErr := rand.Read(salt); randomizeErr != nil {
		return nil, fmt.Errorf("salt generation error: %w", randomizeErr)
	}
	return salt, nil
}

func parseArgon2Hash(hashedPassword string) (*Argon2Config, error) {
	// "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
	config := &Argon2Config{}
	components := strings.Split(hashedPassword, "$")
	if len(components) != 6 {
		return nil, errors.New("invalid hash format structure")
	}

	if !strings.HasPrefix(components[1], "argon2id") {
		return nil, errors.New("unsupported algorithm variant")
	}

	var version int
	fmt.Sscanf(components[2], "v=%d", &version)

	fmt.Sscanf(components[3], "m=%d,t=%d,p=%d", &config.MemoryCost, &config.TimeCost, &config.Thread)

	salt, err := base64.RawStdEncoding.DecodeString(components[4])
	if err != nil {
		return nil, fmt.Errorf("salt decoding failed: %w", err)
	}
	config.Salt = salt

	hash, err := base64.RawStdEncoding.DecodeString(components[5])
	if err != nil {
		return nil, fmt.Errorf("hash decoding failed: %w", err)
	}
	config.Hash = hash
	config.KeyLen = uint32(len(hash))

	return config, nil
}

func VerifyPassword(storedHash, providedPassword string) (bool, error) {
	config, err := parseArgon2Hash(storedHash)
	if err != nil {
		return false, fmt.Errorf("hash parsing failed: %w", err)
	}

	computedHash := argon2.IDKey(
		[]byte(providedPassword),
		config.Salt,
		config.TimeCost,
		config.MemoryCost,
		config.Thread,
		config.KeyLen,
	)

	return subtle.ConstantTimeCompare(config.Hash, computedHash) == 1, nil
}
