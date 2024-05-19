package service

import (
	"fmt"

	"github.com/go-crypt/crypt"
	"github.com/go-crypt/crypt/algorithm/argon2"
)

type UsersService struct {
	hasher  *argon2.Hasher
	decoder *crypt.Decoder
}

func NewUsersService() (*UsersService, error) {
	hasher, err := argon2.New(argon2.WithProfileRFC9106LowMemory(), argon2.WithVariantID())
	if err != nil {
		return nil, fmt.Errorf("creating hasher: %w", err)
	}

	decoder := crypt.NewDecoder()
	err = argon2.RegisterDecoderArgon2id(decoder)
	if err != nil {
		return nil, fmt.Errorf("registering argon2id decoder: %w", err)
	}

	return &UsersService{
		hasher:  hasher,
		decoder: decoder,
	}, nil
}

func (s *UsersService) EncodePasswordDigest(password string) (string, error) {
	digest, err := s.hasher.Hash(password)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}

	return digest.Encode(), nil
}

func (s *UsersService) MatchPassword(passwordDigest, passwordToCheck string) (bool, error) {
	digest, err := s.decoder.Decode(passwordDigest)
	if err != nil {
		return false, fmt.Errorf("decoding password: %w", err)
	}

	return digest.Match(passwordToCheck), nil
}
