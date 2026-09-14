// Package query reads list parameters out of a URL's query string: pages and
// yes-or-no filters, held to something sensible so a handler never has to.
package query

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// MaxOffset caps how far into a list a request may page.
const MaxOffset = 1_000_000

// Int reads a whole number, falling back to a default for one that is missing,
// not a number or negative, and holding it to the cap.
func Int(c *gin.Context, name string, fallback, max int) int {
	raw := c.Query(name)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return fallback
	}

	return min(value, max)
}

// Bool reads a yes-or-no filter: "true" or "false", or nil for anything else,
// which means not filtering at all.
func Bool(c *gin.Context, name string) *bool {
	switch c.Query(name) {
	case "true":
		value := true
		return &value
	case "false":
		value := false
		return &value
	default:
		return nil
	}
}
