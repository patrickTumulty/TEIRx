package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type Argon2Alg struct {
	Id    int
	Label string
}

const (
	ARGON2I = iota
	ARGON2ID
)

type Argon2Params struct {
	MemoryKiB   uint32
	Iterations  uint32
	Parallelism uint8
	SaltLen     uint32
	KeyLen      uint32
}

func DefaultArgon2Params() *Argon2Params {
	return &Argon2Params{
		MemoryKiB:   64 * 1024,
		Iterations:  4,
		Parallelism: 4,
		SaltLen:     16,
		KeyLen:      32,
	}
}

func EncodeArgon2Hash(str string, p *Argon2Params) (string, error) {
	salt, err := generateRandomBytes(p.SaltLen)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error generating random salt: %s", err))
	}

	hash := argon2.IDKey([]byte(str), salt, p.Iterations, p.MemoryKiB, p.Parallelism, p.KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.MemoryKiB, p.Iterations, p.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func DecodeArgon2Hash(encodedHash string) (p *Argon2Params, salt []byte, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("Invalid hash string")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}

	p = &Argon2Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.MemoryKiB, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLen = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLen = uint32(len(hash))

	return p, salt, hash, nil
}
