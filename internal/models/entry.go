// Package models defines the core data structures for GopherVault.
//
// 📚 Go concept: A "package" is Go's way of grouping related code.
// The package name matches the directory name (models).
package models

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Entry represents a single piece of knowledge saved in the vault.
//
// 📚 Go concept: A "struct" is a collection of named fields — similar to
// a class in other languages, but without inheritance.
//
// The backtick annotations (e.g., `json:"id"`) are called "struct tags".
// They tell encoding/json what key name to use when reading/writing JSON.
type Entry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewEntry creates a new Entry with a generated ID and timestamps.
//
// 📚 Go concept: Functions that start with a capital letter are "exported"
// (public). Lowercase = unexported (private to the package).
func NewEntry(title, body, category string, tags []string) *Entry {
	// 📚 Go concept: "now" is a short variable declaration using :=
	// Go infers the type from the right side — no need to write "var now time.Time"
	now := time.Now().UTC()

	return &Entry{
		ID:        generateID(),
		Title:     title,
		Body:      body,
		Category:  category,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// HasTag reports whether the entry has the given tag (case-insensitive).
//
// 📚 Go concept: This is a "method" — a function attached to a type.
// "(e *Entry)" is the "receiver". It means: call this on an Entry value.
// The * means we receive a pointer (no copy is made), which is more efficient.
func (e *Entry) HasTag(tag string) bool {
	tag = strings.ToLower(tag)
	// 📚 Go concept: "range" iterates over slices and maps.
	// "_" discards the index (we only want the value "t").
	for _, t := range e.Tags {
		if strings.ToLower(t) == tag {
			return true
		}
	}
	return false
}

// MatchesSearch reports whether the entry contains the query string
// in its title, body, or tags.
func (e *Entry) MatchesSearch(query string) bool {
	query = strings.ToLower(query)
	if strings.Contains(strings.ToLower(e.Title), query) {
		return true
	}
	if strings.Contains(strings.ToLower(e.Body), query) {
		return true
	}
	for _, t := range e.Tags {
		if strings.Contains(strings.ToLower(t), query) {
			return true
		}
	}
	return false
}

// String returns a short one-line summary of the entry.
//
// 📚 Go concept: Implementing the "Stringer" interface. Any type with a
// String() string method can be printed nicely with fmt.Println().
func (e *Entry) String() string {
	tags := strings.Join(e.Tags, ", ")
	if tags == "" {
		tags = "none"
	}
	return fmt.Sprintf("[%s] %s (tags: %s)", e.ID[:6], e.Title, tags)
}

// generateID returns a short random hex ID (8 chars).
//
// 📚 Go concept: lowercase function = unexported. Only code inside
// the "models" package can call this. It's a pure implementation detail.
func generateID() string {
	// 📚 Go concept: fmt.Sprintf formats a string, like printf in C.
	// %08x means: print as hex, padded to 8 chars with leading zeros.
	return fmt.Sprintf("%08x", rand.Uint32())
}
