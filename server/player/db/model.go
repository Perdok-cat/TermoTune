package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"github.com/doug-martin/goqu/v9"
)


type DB struct {
	db *sql.DB
	goquDB *goqu.Database
	path string
}

func NewDb(path string) (*DB, error) {
    dir := filepath.Dir(path)
    err := os.MkdirAll(dir, 0o755)
    if err != nil {
        return nil, err
    }
	sqlDB, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }
    goquDB := goqu.New("sqlite3", sqlDB)
    return &DB{
		db:      sqlDB,
		goquDB:    goquDB,
        path:    path,
    }, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func LoadDb(path string) (*DB, error) {
	db, err := NewDb(path)
	if err != nil {
		return nil, err
	}
	err = db.InitMusic()
	if err != nil {
		return nil, err
	}

	err = db.InitPlaylist()
	if err != nil {
		return nil, err
	}

	err = db.InitMusicPlaylist()
	if err != nil {
		return nil, err
	}

	return db, nil
}

