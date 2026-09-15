// Package jose signs and checks the JSON Web Tokens this server issues, and
// publishes the keys it signs them with.
//
// It is small on purpose. The server only ever signs with its own asymmetric
// keys — RS256, PS256 and ES256 — and only ever checks tokens it signed
// itself, so a general JOSE library would bring algorithms, encryption modes
// and header options that nothing here may accept. What is not written here
// cannot be tricked into being used: there is no "none", no HMAC, and a
// token's own header never chooses the key.
package jose

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// The algorithms a key can sign with.
const (
	RS256 = "RS256"
	PS256 = "PS256"
	ES256 = "ES256"
)

// Algorithms lists every algorithm, RS256 first: it is what every OpenID
// Connect client has to support, so ID tokens use it.
var Algorithms = []string{RS256, PS256, ES256}

// ErrInvalid is returned for a token that is malformed, signed by a key this
// server does not know, or whose signature does not match. The reasons are
// one error: none of them is the caller's to act on differently.
var ErrInvalid = errors.New("invalid token")

var b64 = base64.RawURLEncoding

// Key is a private key and what it is published as.
type Key struct {
	ID        string
	Algorithm string
	Private   crypto.Signer
}

// Header is the part of a JWS header this package writes and reads.
type Header struct {
	Algorithm string `json:"alg"`
	KeyID     string `json:"kid,omitempty"`
	Type      string `json:"typ,omitempty"`
}

// Generate makes a new private key for an algorithm: a 2048-bit RSA key for
// RS256 and PS256, a P-256 key for ES256.
func Generate(algorithm string) (crypto.Signer, error) {
	switch algorithm {
	case RS256, PS256:
		return rsa.GenerateKey(rand.Reader, 2048)
	case ES256:
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	default:
		return nil, fmt.Errorf("jose: unsupported algorithm %q", algorithm)
	}
}

// Sign returns a compact JWS of the claims, signed with the key. `typ` is the
// header's type: "at+jwt" for an access token (RFC 9068), "JWT" for an ID
// token.
func Sign(key Key, typ string, claims any) (string, error) {
	header, err := json.Marshal(Header{Algorithm: key.Algorithm, KeyID: key.ID, Type: typ})
	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	input := b64.EncodeToString(header) + "." + b64.EncodeToString(payload)

	signature, err := sign(key, []byte(input))
	if err != nil {
		return "", err
	}

	return input + "." + b64.EncodeToString(signature), nil
}

func sign(key Key, input []byte) ([]byte, error) {
	digest := sha256.Sum256(input)

	switch key.Algorithm {
	case RS256:
		private, ok := key.Private.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("jose: %s needs an RSA key", key.Algorithm)
		}
		return rsa.SignPKCS1v15(rand.Reader, private, crypto.SHA256, digest[:])
	case PS256:
		private, ok := key.Private.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("jose: %s needs an RSA key", key.Algorithm)
		}
		return rsa.SignPSS(rand.Reader, private, crypto.SHA256, digest[:], &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
	case ES256:
		private, ok := key.Private.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("jose: %s needs an EC key", key.Algorithm)
		}
		r, s, err := ecdsa.Sign(rand.Reader, private, digest[:])
		if err != nil {
			return nil, err
		}
		// JWS wants r and s as two fixed 32-byte numbers side by side, not
		// the ASN.1 structure Go's SignASN1 would make.
		out := make([]byte, 64)
		r.FillBytes(out[:32])
		s.FillBytes(out[32:])
		return out, nil
	default:
		return nil, fmt.Errorf("jose: unsupported algorithm %q", key.Algorithm)
	}
}

// Verify checks a compact JWS against the key `lookup` returns for its key id,
// and decodes its claims into `claims`. The key decides the algorithm: a token
// whose header names another one is refused, so a header can never switch the
// check to something weaker.
//
// Verify checks the signature only. What the claims have to say — issuer,
// expiry, audience — is the caller's to check, since it differs by token.
func Verify(token string, lookup func(kid string) (Key, bool), claims any) (Header, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Header{}, ErrInvalid
	}

	rawHeader, err := b64.DecodeString(parts[0])
	if err != nil {
		return Header{}, ErrInvalid
	}

	var header Header
	if err := json.Unmarshal(rawHeader, &header); err != nil {
		return Header{}, ErrInvalid
	}

	key, ok := lookup(header.KeyID)
	if !ok || key.Algorithm != header.Algorithm {
		return Header{}, ErrInvalid
	}

	signature, err := b64.DecodeString(parts[2])
	if err != nil {
		return Header{}, ErrInvalid
	}

	if !verify(key, []byte(parts[0]+"."+parts[1]), signature) {
		return Header{}, ErrInvalid
	}

	payload, err := b64.DecodeString(parts[1])
	if err != nil {
		return Header{}, ErrInvalid
	}

	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.UseNumber()
	if err := decoder.Decode(claims); err != nil {
		return Header{}, ErrInvalid
	}

	return header, nil
}

func verify(key Key, input, signature []byte) bool {
	digest := sha256.Sum256(input)

	switch public := key.Private.Public().(type) {
	case *rsa.PublicKey:
		switch key.Algorithm {
		case RS256:
			return rsa.VerifyPKCS1v15(public, crypto.SHA256, digest[:], signature) == nil
		case PS256:
			return rsa.VerifyPSS(public, crypto.SHA256, digest[:], signature, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash}) == nil
		}
	case *ecdsa.PublicKey:
		if key.Algorithm != ES256 || len(signature) != 64 {
			return false
		}
		r := new(big.Int).SetBytes(signature[:32])
		s := new(big.Int).SetBytes(signature[32:])
		return ecdsa.Verify(public, digest[:], r, s)
	}

	return false
}

// JWK is a public key as a JSON Web Key (RFC 7517).
type JWK struct {
	KeyType   string `json:"kty"`
	Use       string `json:"use"`
	Algorithm string `json:"alg"`
	KeyID     string `json:"kid"`

	// RSA
	N string `json:"n,omitempty"`
	E string `json:"e,omitempty"`

	// EC
	Curve string `json:"crv,omitempty"`
	X     string `json:"x,omitempty"`
	Y     string `json:"y,omitempty"`
}

// Public returns the key's public half as a JWK, which is what the JWKS
// endpoint publishes.
func (k Key) Public() (JWK, error) {
	jwk := JWK{Use: "sig", Algorithm: k.Algorithm, KeyID: k.ID}

	switch public := k.Private.Public().(type) {
	case *rsa.PublicKey:
		jwk.KeyType = "RSA"
		jwk.N = b64.EncodeToString(public.N.Bytes())
		jwk.E = b64.EncodeToString(big.NewInt(int64(public.E)).Bytes())
	case *ecdsa.PublicKey:
		point, err := public.Bytes()
		if err != nil {
			return JWK{}, err
		}
		// Uncompressed point: 0x04, then X and Y, 32 bytes each on P-256.
		jwk.KeyType = "EC"
		jwk.Curve = "P-256"
		jwk.X = b64.EncodeToString(point[1:33])
		jwk.Y = b64.EncodeToString(point[33:65])
	default:
		return JWK{}, fmt.Errorf("jose: unsupported key type %T", public)
	}

	return jwk, nil
}

// HalfHash is the at_hash of OpenID Connect Core 3.1.3.6: the left half of the
// SHA-256 of a value, base64url encoded. All three algorithms here use
// SHA-256, so it is the same for each.
func HalfHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return b64.EncodeToString(sum[:len(sum)/2])
}
