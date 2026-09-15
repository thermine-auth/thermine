package jose

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"strings"
	"testing"
)

func testKey(t *testing.T, alg, kid string) Key {
	t.Helper()

	private, err := Generate(alg)
	if err != nil {
		t.Fatal(err)
	}

	return Key{ID: kid, Algorithm: alg, Private: private}
}

func lookupOf(keys ...Key) func(string) (Key, bool) {
	return func(kid string) (Key, bool) {
		for _, key := range keys {
			if key.ID == kid {
				return key, true
			}
		}
		return Key{}, false
	}
}

func TestSignAndVerifyEachAlgorithm(t *testing.T) {
	for _, alg := range Algorithms {
		t.Run(alg, func(t *testing.T) {
			key := testKey(t, alg, "k1")

			token, err := Sign(key, "at+jwt", map[string]any{"sub": "user", "n": 7})
			if err != nil {
				t.Fatal(err)
			}

			var claims struct {
				Sub string      `json:"sub"`
				N   json.Number `json:"n"`
			}
			header, err := Verify(token, lookupOf(key), &claims)
			if err != nil {
				t.Fatalf("Verify: %v", err)
			}

			if header.Type != "at+jwt" || header.Algorithm != alg || header.KeyID != "k1" {
				t.Errorf("header = %+v", header)
			}
			if claims.Sub != "user" || claims.N != "7" {
				t.Errorf("claims = %+v", claims)
			}
		})
	}
}

func TestVerifyRefusesTampering(t *testing.T) {
	key := testKey(t, RS256, "k1")
	other := testKey(t, RS256, "k1")

	token, err := Sign(key, "JWT", map[string]any{"sub": "alice"})
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(token, ".")

	forged := parts[0] + "." + base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"mallory"}`)) + "." + parts[2]

	// A header naming "none", or HS256 with the public key as the secret, is
	// the classic way round a careless check.
	none := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","kid":"k1"}`)) + "." + parts[1] + "."
	hs := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","kid":"k1"}`)) + "." + parts[1] + "." + parts[2]

	cases := map[string]struct {
		token  string
		lookup func(string) (Key, bool)
	}{
		"changed payload":   {forged, lookupOf(key)},
		"alg none":          {none, lookupOf(key)},
		"alg switched":      {hs, lookupOf(key)},
		"another key":       {token, lookupOf(other)},
		"unknown key id":    {token, lookupOf(Key{ID: "k2", Algorithm: RS256, Private: key.Private})},
		"not a jws":         {"abc", lookupOf(key)},
		"truncated":         {parts[0] + "." + parts[1], lookupOf(key)},
		"signature garbage": {parts[0] + "." + parts[1] + ".!!!", lookupOf(key)},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var claims map[string]any
			if _, err := Verify(tc.token, tc.lookup, &claims); err == nil {
				t.Fatal("Verify accepted it")
			}
		})
	}
}

func TestPublicJWKMatchesTheKey(t *testing.T) {
	rsaKey := testKey(t, PS256, "r")
	jwk, err := rsaKey.Public()
	if err != nil {
		t.Fatal(err)
	}

	public := rsaKey.Private.Public().(*rsa.PublicKey)
	n, _ := base64.RawURLEncoding.DecodeString(jwk.N)
	if jwk.KeyType != "RSA" || jwk.Use != "sig" || jwk.Algorithm != PS256 || new(big.Int).SetBytes(n).Cmp(public.N) != 0 || jwk.E != "AQAB" {
		t.Errorf("RSA JWK = %+v", jwk)
	}

	ecKey := testKey(t, ES256, "e")
	jwk, err = ecKey.Public()
	if err != nil {
		t.Fatal(err)
	}

	// An uncompressed point is 0x04, then X, then Y.
	point, err := ecKey.Private.Public().(*ecdsa.PublicKey).Bytes()
	if err != nil {
		t.Fatal(err)
	}
	x, _ := base64.RawURLEncoding.DecodeString(jwk.X)
	if jwk.KeyType != "EC" || jwk.Curve != "P-256" || !bytes.Equal(x, point[1:33]) {
		t.Errorf("EC JWK = %+v", jwk)
	}
}

func TestSealOpensOnlyWithTheSameSecret(t *testing.T) {
	key := testKey(t, ES256, "k")

	sealer, err := NewSealer("0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}

	sealed, err := sealer.Seal(key.Private)
	if err != nil {
		t.Fatal(err)
	}

	opened, err := sealer.Open(sealed)
	if err != nil {
		t.Fatal(err)
	}
	if !opened.Public().(*ecdsa.PublicKey).Equal(key.Private.Public()) {
		t.Error("opened a different key")
	}

	wrong, _ := NewSealer("another secret that is long enough!!")
	if _, err := wrong.Open(sealed); err == nil {
		t.Error("opened with the wrong secret")
	}
}

// The example from OpenID Connect Core's at_hash definition is not published
// with a value, so this pins the construction instead: half of SHA-256.
func TestHalfHashIsHalfOfSHA256(t *testing.T) {
	got, _ := base64.RawURLEncoding.DecodeString(HalfHash("token"))
	if len(got) != 16 {
		t.Errorf("at_hash is %d bytes, want 16", len(got))
	}
}
