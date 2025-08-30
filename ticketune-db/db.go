package ticketune_db

import (
	"database/sql"
	"log"
	"time"

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

	// Start cleanup goroutine: every 2 days, delete threads older than 1 month
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(48 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ticketuneDB.CleanupOldThreads()
			case <-stop:
				return
			}
		}
	}()
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

// DeleteUserThread removes a user's thread record.
func (d *DB) DeleteUserThread(userID string) error {
	_, err := d.db.Exec(`DELETE FROM support_tickets WHERE user_id = ?`, userID)
	return err
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.db.Close()
}

// DeleteThreadsOlderThan deletes threads older than 1 month
// Called automatically every 2 days by a goroutine
func (d *DB) CleanupOldThreads() error {
	_, err := d.db.Exec(`DELETE FROM support_tickets WHERE created_at < datetime('now', '-1 month')`)
	return err
}

// Delete a thread record by thread ID
func (d *DB) CloseThread(threadId tempest.Snowflake) error {
	_, err := d.db.Exec(`DELETE FROM support_tickets WHERE thread_id = ?`, threadId)
	return err
}
