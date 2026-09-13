// Package auth signs administrators in and out, and keeps the record of what
// they did.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"xermess/internal/model"
	"xermess/internal/store"
)

// SessionLifetime is how long a session lasts before the administrator has to
// sign in again.
const SessionLifetime = 12 * time.Hour

// ErrInvalidCredentials is returned for a username that does not exist, a
// wrong password, and an account that may not sign in. They are one error on
// purpose: telling them apart tells an attacker which usernames are real.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrNoSession is returned when a request carries no usable session.
var ErrNoSession = errors.New("no active session")

// Service signs administrators in and out. Every query it makes goes through
// the store, so this file is about what signing in means rather than about
// how the rows are fetched.
type Service struct {
	store *store.Store
}

// New returns a Service backed by the given store.
func New(st *store.Store) *Service {
	return &Service{store: st}
}

// Request describes where a call came from, which is recorded on the session
// and in the log.
type Request struct {
	IP        string
	UserAgent string
}

// Login checks the credentials and starts a session. The token it returns is
// the only copy: the database keeps a hash of it, so a leaked database cannot
// be used to sign in.
func (s *Service) Login(ctx context.Context, username, password string, req Request) (string, *model.AdminUser, error) {
	admin, err := s.store.AdminByUsername(ctx, username)

	switch {
	case errors.Is(err, store.ErrNotFound):
		// Still hash something, so a missing username and a wrong password
		// take the same time to answer.
		_ = bcrypt.CompareHashAndPassword([]byte("$2a$10$"+hex.EncodeToString(make([]byte, 26))), []byte(password))
		s.record(ctx, nil, username, "admin.login_failed", req, "unknown username")
		return "", nil, ErrInvalidCredentials
	case err != nil:
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		s.record(ctx, &admin.ID, admin.Username, "admin.login_failed", req, "wrong password")
		return "", nil, ErrInvalidCredentials
	}

	if !admin.CanSignIn(time.Now()) {
		s.record(ctx, &admin.ID, admin.Username, "admin.login_blocked", req, string(admin.Status))
		return "", nil, ErrInvalidCredentials
	}

	token, err := newToken()
	if err != nil {
		return "", nil, err
	}

	session := model.AdminUserSession{
		AdminUserID: admin.ID,
		TokenHash:   hashToken(token),
		ExpiresAt:   time.Now().Add(SessionLifetime),
		// There is no second factor yet. When there is, this starts false and
		// is set once the factor is cleared.
		MFAPassed: true,
		UserAgent: req.UserAgent,
		IP:        req.IP,
	}
	if err := s.store.CreateSession(ctx, &session); err != nil {
		return "", nil, err
	}

	if err := s.store.MarkAdminSignedIn(ctx, admin, time.Now(), req.IP); err != nil {
		return "", nil, err
	}

	s.record(ctx, &admin.ID, admin.Username, "admin.login", req, "")

	return token, admin, nil
}

// Authenticate returns the administrator the token belongs to. It is what the
// middleware calls on every request.
func (s *Service) Authenticate(ctx context.Context, token string) (*model.AdminUser, error) {
	if token == "" {
		return nil, ErrNoSession
	}

	session, err := s.store.SessionByTokenHash(ctx, hashToken(token))

	switch {
	case errors.Is(err, store.ErrNotFound):
		return nil, ErrNoSession
	case err != nil:
		return nil, err
	}

	if !session.IsActive(time.Now()) {
		return nil, ErrNoSession
	}

	admin, err := s.store.AdminByID(ctx, session.AdminUserID)
	if err != nil {
		return nil, ErrNoSession
	}

	if !admin.CanSignIn(time.Now()) {
		return nil, ErrNoSession
	}

	return admin, nil
}

// Logout revokes the session the token belongs to. Signing out twice is not
// an error.
func (s *Service) Logout(ctx context.Context, token string, req Request) error {
	if token == "" {
		return nil
	}

	session, err := s.store.SessionByTokenHash(ctx, hashToken(token))
	switch {
	case errors.Is(err, store.ErrNotFound):
		return nil
	case err != nil:
		return err
	}

	if err := s.store.RevokeSession(ctx, session, time.Now()); err != nil {
		return err
	}

	if admin, err := s.store.AdminByID(ctx, session.AdminUserID); err == nil {
		s.record(ctx, &admin.ID, admin.Username, "admin.logout", req, "")
	}

	return nil
}

// record writes one line to the activity log. A failure to write the log must
// not fail the request that caused it, so the error is swallowed on purpose;
// the caller has already done the thing being recorded.
func (s *Service) record(ctx context.Context, adminID *uuid.UUID, actor, action string, req Request, note string) {
	entry := model.AuditLog{
		AdminUserID: adminID,
		ActorEmail:  actor,
		Action:      action,
		TargetType:  "admin_user",
		IP:          req.IP,
		UserAgent:   req.UserAgent,
	}
	if note != "" {
		entry.Metadata = map[string]any{"reason": note}
	}
	if adminID != nil {
		entry.TargetID = adminID.String()
	}

	_ = s.store.WriteAudit(ctx, &entry)
}

// newToken returns a session token: 256 bits of randomness, hex encoded.
func newToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}
	return hex.EncodeToString(b[:]), nil
}

// hashToken is what gets stored. SHA-256 is right here, unlike for passwords:
// the token is long and random, so it cannot be guessed by brute force.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
