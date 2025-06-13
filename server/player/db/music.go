package db

import (
	"github.com/doug-martin/goqu/v9"
)

type Music struct {
	ID     int64      `db:"id"`
	Title  string     `db:"title"`
	Artist string     `db:"artist"`
	Album  string     `db:"album"`
	Genre  string     `db:"genre"`
	Year   int       `db:"year"`
}

type MusicDB struct {
	db *goqu.Database
}

func NewMusicDB(db *goqu.Database) *MusicDB {
	return &MusicDB{
		db: db,
	}
}
