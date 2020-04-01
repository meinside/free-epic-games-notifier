package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3" // sqlite3

	"github.com/meinside/free-epic-games-notifier/extractor"
)

// Database is a sqlite database for caching free games
type Database struct {
	db *sql.DB
	sync.RWMutex
}

// Open opens a local database
func Open(filepath string) (*Database, error) {
	var db *sql.DB
	var err error

	if db, err = sql.Open("sqlite3", filepath); err != nil {
		return nil, fmt.Errorf("failed to open database: %s", err)
	}

	_sqlite := &Database{
		db: db,
	}

	// caches table and index
	if _, err = db.Exec(`create table if not exists caches(
				id integer primary key autoincrement,
				title text not null,
				store_url text not null,
				image_url text default null,
				time integer default (strftime('%s', 'now'))
			)`); err != nil {
		return nil, fmt.Errorf("failed to create caches table: %s", err)
	}
	if _, err = db.Exec(`create index if not exists idx_caches1 on caches(
				title
			)`); err != nil {
		return nil, fmt.Errorf("failed to create index idx_caches1: %s", err)
	}

	return _sqlite, nil
}

// Close closes database
func (d *Database) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// CacheGame caches a game to the database
func (d *Database) CacheGame(game extractor.FreeGame) (err error) {
	d.Lock()
	defer d.Unlock()

	var stmt *sql.Stmt
	if stmt, err = d.db.Prepare(`insert into caches(title, store_url, image_url) values(?, ?, ?)`); err != nil {
		log.Printf("* failed to prepare a statement: %s", err)
	} else {
		defer stmt.Close()
		if _, err = stmt.Exec(game.Title, game.StoreURL, game.ImageURL); err != nil {
			log.Printf("* failed to cache free game into local database: %s", err)
		}
	}

	return err
}

// IsCachedGame checks if a game with given title is already cached or not
func (d *Database) IsCachedGame(title string) (cached bool, err error) {
	d.RLock()
	defer d.RUnlock()

	var stmt *sql.Stmt
	if stmt, err = d.db.Prepare(`select count(id) from caches where title = ?`); err != nil {
		log.Printf("* failed to prepare a statement: %s", err)
	} else {
		defer stmt.Close()

		var rows *sql.Rows
		if rows, err = stmt.Query(title); err != nil {
			log.Printf("* failed to count caches from local database: %s", err)
		} else {
			defer rows.Close()

			var cnt int
			if rows.Next() {
				rows.Scan(&cnt)

				if cnt > 0 {
					cached = true
				}
			}
		}
	}

	return cached, err
}
