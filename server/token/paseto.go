package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto      paseto.V2
	symetricKey []byte
}

func NewPasetoMaker(symetricKey string) (Maker, error) {
	if len(symetricKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("Invalid symetric key size should be greater or equal to %d", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:      *paseto.NewV2(),
		symetricKey: []byte(symetricKey),
	}
	return maker, nil
}

func (p *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}
	return p.paseto.Encrypt(p.symetricKey, payload, nil)
}

func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := p.paseto.Decrypt(token, p.symetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if err = payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil
}
