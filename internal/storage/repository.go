// Package storage defines how GopherVault persists data.
//
// 📚 Go concept: Interfaces in Go provide "implicit implementation".
// A type doesn't need to declare that it implements an interface (there is no
// "implements" keyword like in Java or TypeScript).
// If a struct implements all the methods declared in this interface,
// it automatically satisfies the interface!
package storage

import "github.com/Uvaisbugh/gophervault/internal/models"

// Repository defines the contract for storing and retrieving entries.
// By coding to this interface instead of a concrete file or database,
// our service layer doesn't care whether data is stored in JSON, SQLite,
// or in-memory for testing.
type Repository interface {
	// Save adds a new entry or updates an existing one.
	Save(entry *models.Entry) error

	// FindByID looks up an entry by its full ID or ID prefix.
	FindByID(id string) (*models.Entry, error)

	// FindAll returns all entries stored in the repository.
	FindAll() ([]*models.Entry, error)

	// Delete removes an entry by ID.
	Delete(id string) error
}
