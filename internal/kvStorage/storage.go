// Package kvStorage provides a simple key-value storage implementation using SQLite.
// It is used for persistent storage of application data, including HomeKit pairing information
// and the deCONZ API key. This package implements the storage interface required by the
// HomeKit Accessory Protocol (HAP) library.
package kvStorage

import (
	"database/sql"
	"errors"

	// Import SQLite driver
	_ "github.com/glebarez/go-sqlite"
)

// Storage represents a key-value storage backed by SQLite.
// It provides methods for storing, retrieving, and deleting binary data
// associated with string keys.
type Storage struct {
	// conn is the database connection to the SQLite database
	conn *sql.DB
}

// New creates a new Storage instance with the specified database file.
// If the database file doesn't exist, it will be created.
// If the kv_store table doesn't exist, it will be created.
//
// Parameters:
//   - path: The path to the SQLite database file
//
// Returns:
//   - *Storage: A pointer to the initialized Storage
//   - error: An error if the database could not be opened or the table could not be created
func New(path string) (*Storage, error) {
	// Open the SQLite database
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Create the kv_store table if it doesn't exist
	// The table has two columns: key (TEXT, primary key) and value (BLOB)
	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS kv_store (key TEXT PRIMARY KEY, value BLOB);"); err != nil {
		return nil, err
	}

	return &Storage{db}, nil
}

// Set stores a value for the given key.
// If the key already exists, its value will be updated.
//
// Parameters:
//   - key: The key to store the value under
//   - value: The binary data to store
//
// Returns:
//   - error: An error if the value could not be stored
func (s *Storage) Set(key string, value []byte) error {
	// Insert a new row or update an existing one if the key already exists
	_, err := s.conn.Exec(`INSERT INTO kv_store(key, value) VALUES(?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value;`, key, value)
	return err
}

// Get retrieves the value for the given key.
// If the key doesn't exist, it returns nil without an error.
//
// Parameters:
//   - key: The key to retrieve the value for
//
// Returns:
//   - []byte: The stored binary data, or nil if the key doesn't exist
//   - error: An error if the value could not be retrieved
func (s *Storage) Get(key string) ([]byte, error) {
	var val []byte
	// Query the value for the given key
	err := s.conn.QueryRow(`SELECT value FROM kv_store WHERE key = ?;`, key).Scan(&val)

	// If the key doesn't exist, return nil without an error
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return val, err
}

// Delete removes the value for the given key.
// If the key doesn't exist, this is a no-op.
//
// Parameters:
//   - key: The key to delete the value for
//
// Returns:
//   - error: An error if the value could not be deleted
func (s *Storage) Delete(key string) error {
	// Delete the row for the given key
	_, err := s.conn.Exec(`DELETE FROM kv_store WHERE key = ?;`, key)
	return err
}

// KeysWithSuffix returns a list of keys that end with the given suffix.
// This is used by the HAP library to find all keys related to a specific accessory.
//
// Parameters:
//   - suffix: The suffix to search for
//
// Returns:
//   - []string: A slice of keys that end with the given suffix
//   - error: An error if the keys could not be retrieved
func (s *Storage) KeysWithSuffix(suffix string) ([]string, error) {
	// Create a LIKE parameter for the SQL query
	// The % is a wildcard that matches any characters before the suffix
	likeParam := "%" + suffix

	// Query all keys that match the LIKE parameter
	rows, err := s.conn.Query(`SELECT key FROM kv_store WHERE key LIKE ?;`, likeParam)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect all matching keys into a slice
	var keys []string
	for rows.Next() {
		var key string
		if err = rows.Scan(&key); err == nil {
			keys = append(keys, key)
		}
	}

	return keys, nil
}
