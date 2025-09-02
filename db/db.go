/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 * SPDX-FileContributor: patapancakes
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package db

import (
	"database/sql"

	"github.com/amatsagu/tempest"
	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "ticketune-db.sqlite3"

// DB wraps the sql.DB for ticketune
type DB struct {
	db *sql.DB // The underlying database connection
}

var TicketuneDB *DB // Global instance of the DB

func Init() error {
	var err error
	TicketuneDB, err = open()
	if err != nil {
		return err
	}

	err = TicketuneDB.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func Get() *DB {
	return TicketuneDB
}

// Open (or or create) the ticketune database and ensure the support_tickets table exists
func open() (*DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS support_tickets (
	       user_id TEXT PRIMARY KEY,
	       thread_id TEXT NOT NULL,
	       created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
       );`)
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

// SetUserThread stores or updates a user's thread info.
func (d *DB) SetUserThread(userID tempest.Snowflake, threadID tempest.Snowflake) error {
	_, err := d.db.Exec(
		`INSERT OR REPLACE INTO support_tickets (user_id, thread_id) VALUES (?, ?)`,
		userID,
		threadID,
	)

	return err
}

// GetUserThread returns the thread ID and creation time for a user, or empty string and zero time if not found.
func (d *DB) GetUserThread(userID tempest.Snowflake) (tempest.Snowflake, error) {
	row := d.db.QueryRow(
		`SELECT thread_id FROM support_tickets WHERE user_id = ?`,
		userID,
	)

	var threadID tempest.Snowflake
	err := row.Scan(&threadID)
	if err != nil {
		return tempest.Snowflake(0), err
	}

	return threadID, nil
}

// GetThreadUser returns the user ID associated with a thread ID.
func (d *DB) GetThreadUser(threadID tempest.Snowflake) (tempest.Snowflake, error) {
	row := d.db.QueryRow(
		`SELECT user_id FROM support_tickets WHERE thread_id = ?`,
		threadID,
	)

	var userID tempest.Snowflake
	err := row.Scan(&userID)
	if err != nil {
		return tempest.Snowflake(0), err
	}

	return userID, nil
}

// DeleteUserThread removes a user's thread record.
func (d *DB) DeleteUserThread(userID string) error {
	_, err := d.db.Exec(`DELETE FROM support_tickets WHERE user_id = ?`, userID)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.db.Close()
}

// Cleanup threads older than 1 month.
// Unused in favor of explicit thread closing, but kept for potential future use.
func (d *DB) CleanupOldThreads() error {
	_, err := d.db.Exec(`DELETE FROM support_tickets WHERE created_at < datetime('now', '-1 month')`)
	if err != nil {
		return err
	}

	return nil
}

// Delete a thread record by thread ID, and return the user ID that was associated with it.
func (d *DB) CloseThread(threadId tempest.Snowflake) (tempest.Snowflake, error) {
	row := d.db.QueryRow(`DELETE FROM support_tickets WHERE thread_id = ? RETURNING user_id`, threadId)

	var userID tempest.Snowflake
	err := row.Scan(&userID)
	if err != nil {
		return tempest.Snowflake(0), err
	}

	return userID, nil
}
