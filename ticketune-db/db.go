package ticketune_db

import (
	"database/sql"
	"log"

	"github.com/amatsagu/tempest"
	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "ticketune-db.sqlite3"

// DB wraps the sql.DB for ticketune
type DB struct {
	db *sql.DB // The underlying database connection
}

var TicketuneDB *DB // Global instance of the DB

func InitDB() (err error) {
	if TicketuneDB != nil && TicketuneDB.db != nil && TicketuneDB.db.Ping() == nil {
		return
	}

	TicketuneDB, err = openDB()
	if err != nil {
		return
	}
	return TicketuneDB.db.Ping()
}

func GetDB() *DB {
	if TicketuneDB == nil || TicketuneDB.db == nil {
		if err := InitDB(); err != nil {
			log.Fatalln("failed to initialize the database", err)
		}
	}
	return TicketuneDB
}

// Open (or or create) the ticketune database and ensure the support_tickets table exists
func openDB() (*DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	createTable := `CREATE TABLE IF NOT EXISTS support_tickets (
	       user_id TEXT PRIMARY KEY,
	       thread_id TEXT NOT NULL,
	       created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
       );`
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}
	ticketuneDB := &DB{db: db}

	return ticketuneDB, nil
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
func (d *DB) GetUserThread(userID tempest.Snowflake) (threadID tempest.Snowflake, err error) {
	row := d.db.QueryRow(
		`SELECT thread_id FROM support_tickets WHERE user_id = ?`,
		userID,
	)
	err = row.Scan(&threadID)
	return
}

// GetThreadUser returns the user ID associated with a thread ID.
func (d *DB) GetThreadUser(threadID tempest.Snowflake) (userID tempest.Snowflake, err error) {
	row := d.db.QueryRow(
		`SELECT user_id FROM support_tickets WHERE thread_id = ?`,
		threadID,
	)
	err = row.Scan(&userID)
	return
}

// DeleteUserThread removes a user's thread record.
func (d *DB) DeleteUserThread(userID string) error {
	_, err := d.db.Exec(`DELETE FROM support_tickets WHERE user_id = ?`, userID)
	return err
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.db.Close()
}

// Cleanup threads older than 1 month.
// Unused in favor of explicit thread closing, but kept for potential future use.
func (d *DB) CleanupOldThreads() error {
	_, err := d.db.Exec(`DELETE FROM support_tickets WHERE created_at < datetime('now', '-1 month')`)
	return err
}

// Delete a thread record by thread ID, and return the user ID that was associated with it.
func (d *DB) CloseThread(threadId tempest.Snowflake) (userID tempest.Snowflake, err error) {
	row := d.db.QueryRow(`DELETE FROM support_tickets WHERE thread_id = ? RETURNING user_id`, threadId)
	err = row.Scan(&userID)
	return
}
