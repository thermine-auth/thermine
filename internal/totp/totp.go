// Package totp is time-based one-time passwords (RFC 6238): the six-digit
// codes an authenticator app shows, and checking them.
//
// The parameters are the ones every authenticator app assumes when an
// otpauth:// URI does not say otherwise — SHA-1, six digits, thirty-second
// steps — so they are fixed here rather than configurable: an app that reads
// a different value from the URI is an app that shows the wrong codes.
package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Period is the length of one step, and Digits the length of a code.
const (
	Period = 30 * time.Second
	Digits = 6
)

// Skew is how many steps before and after the current one are accepted, so a
// phone whose clock is a little off, or a code typed as it changes, still
// works.
const Skew = 1

var encoding = base32.StdEncoding.WithPadding(base32.NoPadding)

// NewSecret returns a new secret: 160 random bits, the size RFC 4226
// recommends for SHA-1, base32 encoded as authenticator apps expect.
func NewSecret() (string, error) {
	var b [20]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("totp: generate secret: %w", err)
	}

	return encoding.EncodeToString(b[:]), nil
}

// Step is the number of the thirty-second step a moment falls in.
func Step(at time.Time) int64 {
	return at.Unix() / int64(Period/time.Second)
}

// Code is the code for one step.
func Code(secret string, step int64) (string, error) {
	key, err := encoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", fmt.Errorf("totp: secret is not base32: %w", err)
	}

	var counter [8]byte
	binary.BigEndian.PutUint64(counter[:], uint64(step))

	mac := hmac.New(sha1.New, key)
	mac.Write(counter[:])
	sum := mac.Sum(nil)

	// Dynamic truncation, RFC 4226 section 5.3.
	offset := sum[len(sum)-1] & 0x0f
	value := binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7fffffff

	return fmt.Sprintf("%0*d", Digits, value%1_000_000), nil
}

// Verify checks a code against the steps around `at`, and returns the step it
// matched. A step at or before `after` is refused even when the code is right:
// that code has been used already, and accepting it again would let whoever
// saw it over someone's shoulder use it too.
func Verify(secret, code string, at time.Time, after int64) (int64, bool) {
	code = strings.ReplaceAll(strings.TrimSpace(code), " ", "")
	if len(code) != Digits {
		return 0, false
	}

	now := Step(at)
	for step := now - Skew; step <= now+Skew; step++ {
		if step <= after {
			continue
		}

		want, err := Code(secret, step)
		if err != nil {
			return 0, false
		}

		if subtle.ConstantTimeCompare([]byte(want), []byte(code)) == 1 {
			return step, true
		}
	}

	return 0, false
}

// URI is the otpauth:// URI an authenticator app reads from a QR code. The
// label names the account within the issuer, so one app can hold several.
func URI(issuer, account, secret string) string {
	label := url.PathEscape(issuer) + ":" + url.PathEscape(account)

	query := url.Values{}
	query.Set("secret", secret)
	query.Set("issuer", issuer)
	query.Set("algorithm", "SHA1")
	query.Set("digits", fmt.Sprint(Digits))
	query.Set("period", fmt.Sprint(int(Period/time.Second)))

	return "otpauth://totp/" + label + "?" + query.Encode()
}
