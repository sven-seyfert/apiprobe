package db

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"

	"github.com/sven-seyfert/apiprobe/internal/logger"
)

// Init opens or creates the SQLite database file at './db/store.db',
// ensures that the 'secrets' table exists and returns the active
// connection to the caller.
func Init() (*sqlite.Conn, error) {
	// Create database.
	conn, err := sqlite.OpenConn("./db/store.db", sqlite.OpenReadWrite, sqlite.OpenCreate)
	if err != nil {
		logger.Errorf("Failed to open database. Error: %v", err)

		return nil, err
	}

	// Create table if it does not exist.
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS secrets (
			hash   TEXT PRIMARY KEY,
			secret TEXT NOT NULL
		);`

	err = sqlitex.ExecuteTransient(conn, createTableSQL, nil)
	if err != nil {
		logger.Errorf("Failed to create database. Error: %v", err)

		return nil, err
	}

	return conn, nil
}

// InsertSeedData checks if the 'secrets' table is empty;
// if so, reads './db/seed.csv', constructs a bulk-insert SQL statement
// and populates the table. Returns an error if any operation fails.
func InsertSeedData(conn *sqlite.Conn) error {
	// Check if the table is empty.
	count, err := GetTableEntryCount(conn)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	// Insert data (bulk insert).
	SQLValues, err := readSeedData()
	if err != nil {
		return err
	}

	bulkInsertSQL := "INSERT INTO secrets(hash, secret) VALUES" + SQLValues

	err = sqlitex.ExecuteTransient(conn, bulkInsertSQL, nil)
	if err != nil {
		logger.Errorf("Failed to insert data. Error: %v", err)

		return err
	}

	return nil
}

// GetTableEntryCount returns the total number of rows in the 'secrets'
// table by executing 'SELECT COUNT(*)'.
func GetTableEntryCount(conn *sqlite.Conn) (int, error) {
	var count int

	countSQL := "SELECT COUNT(*) FROM secrets"
	err := sqlitex.ExecuteTransient(conn, countSQL, &sqlitex.ExecOptions{
		Args:  nil,
		Named: nil,
		ResultFunc: func(stmt *sqlite.Stmt) error {
			count = stmt.ColumnInt(0)

			return nil
		},
	})
	if err != nil {
		logger.Errorf("Failed to query table count. Error: %v", err)

		return 0, err
	}

	return count, nil
}

// readSeedData reads './db/seed.csv', each line containing 'hash,secret'
// and returns a string suitable for a SQL VALUES clause for a bulk insert.
func readSeedData() (string, error) {
	file, err := os.Open("./db/seed.csv")
	if err != nil {
		logger.Errorf("Failure opening file. Error: %v", err)

		return "", err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var values []string

	for {
		record, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}

		if readErr != nil {
			logger.Errorf("Failure reading file. Error: %v", readErr)

			return "", readErr
		}

		hash := record[0]
		secret := record[1]
		values = append(values, fmt.Sprintf("('%s', '%s')", hash, secret))
	}

	joinedValues := strings.Join(values, ",\n    ")
	bulkInsertSQL := fmt.Sprintf("\n    %s;", joinedValues)

	return bulkInsertSQL, nil
}

// InsertSecret stores a new (hash, secret) pair into the 'secrets' table
// using parameterized SQL to avoid injection. Returns an error if insertion fails.
func InsertSecret(conn *sqlite.Conn, hash string, secret string) error {
	stmt, _, err := conn.PrepareTransient("INSERT INTO secrets(hash, secret) VALUES (?, ?)")
	if err != nil {
		logger.Errorf("Failed to prepare insert statement. Error: %v", err)

		return err
	}

	defer func() {
		if err = stmt.Finalize(); err != nil {
			logger.Errorf("Failed to finalize statement. Error: %v", err)
		}
	}()

	stmt.BindText(1, hash)
	stmt.BindText(2, secret) //nolint:mnd

	if _, err = stmt.Step(); err != nil {
		logger.Errorf("Failed to execute insert statement. Error: %v", err)

		return err
	}

	return nil
}

// SelectHash queries the 'secrets' table for the given hash
// and returns its stored secret. Returns an empty string if no row is found
// or an error on failure.
func SelectHash(conn *sqlite.Conn, hash string) (string, error) {
	stmt, _, err := conn.PrepareTransient("SELECT secret FROM secrets WHERE hash = ?")
	if err != nil {
		logger.Errorf("Failed to prepare select statement. Error: %v", err)

		return "", err
	}

	defer func() {
		if err = stmt.Finalize(); err != nil {
			logger.Errorf("Failed to finalize statement. Error: %v", err)
		}
	}()

	stmt.BindText(1, hash)

	hasRow, err := stmt.Step()
	if err != nil {
		logger.Errorf("Failed to execute select statement. Error: %v", err)

		return "", err
	}

	if !hasRow {
		return "", nil
	}

	secret := stmt.ColumnText(0)

	return secret, nil
}
