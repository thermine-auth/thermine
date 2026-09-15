package mfa

// The answers of these endpoints are the auth package's own types —
// auth.MFAStatus and auth.Enrolment — and lists of recovery codes, wrapped in
// one named field each so a client reads them the same way as every other
// answer.
