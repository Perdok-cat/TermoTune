package db

import (
	"github.com/doug-martin/goqu/v9"
	"fmt"
	"database/sql"
	"crypto/md5"
	"encoding/hex"
)

type Music struct {
	Name string 
	Source string 
	Key string 
	Data []byte
	Hash string 
}

func (d *DB) InitMusic() error {
	_, err := d.db.Exec(
		`CREATE TABLE IF NOT EXISTS music (
      name TEXT UNIQUE,
      source TEXT,
      key TEXT,
      data BLOB,
      hash TEXT UNIQUE NOT NULL,
      PRIMARY KEY (source, key)
    )`,
	)
	return err
}


func (d *DB) GetMusic(source string, key string) (*Music, error) {
	var music Music 
	found, err := d.goquDB.From("music").
		Where(goqu.Ex{"source": source, "key": key}).
		ScanStruct(&music)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &music, nil
}


func (d *DB) UpdateMusic(music *Music) error {
    _, err := d.goquDB.Insert("music").
        Rows(music).
        OnConflict(
            goqu.DoUpdate("source, key", 
                goqu.Record{
                    "name": music.Name, 
                    "data": music.Data, 
                    "hash": music.Hash,
                },
            ),
        ).Executor().Exec()

    return err
}

func (d *DB) InsertUniqueMusicName(music *Music) error {
	
    var newName string

	ds := d.goquDB.From("music").
		Select("name").
		Where(goqu.I("name").Like(fmt.Sprintf("%s%%", music.Name))).
		Order(goqu.I("name").Desc()).
		Limit(1)

	found, err := ds.ScanVal(&newName)
	if err != nil && err == sql.ErrNoRows {
		return err
	}	

	if !found {
		newName = fmt.Sprintf("%s_1", music.Name)
	} else {
		var suffix int 
		_, err = fmt.Sscanf(newName, "%s_%d", &newName, &suffix)
		if err != nil {	
			return err
		}
		suffix++
		newName = fmt.Sprintf("%s_%d", music.Name, suffix)
	}

	_, err = d.goquDB.Insert("music").
		Rows(music).
		OnConflict(
			goqu.DoUpdate("source, key",
				goqu.Record{
					"name": music.Name,
					"data": music.Data,
					"hash": music.Hash,
				},
			),
		).Executor().Exec()
	
	return err
}

func (d *DB) AddMusic (music *Music) error {
	var count int 
	found, err := d.goquDB.From("music").
		Select("name").
		Where(goqu.I("name").Like(fmt.Sprintf("%s%%", music.Name))).
		ScanVal(&count)
	
	if found && count > 0 {
		d.InsertUniqueMusicName(music)
	}

	insertMusic := d.goquDB.Insert("music").
		Rows(
			goqu.Record{
				"name":   music.Name,
				"source": music.Source,
				"key":    music.Key,
				"data":   music.Data,
				"hash":   music.Hash,
			},
		)

	_, err = insertMusic.Executor().Exec()
	return err
}

func (d *DB) GetMusicByName(name string) (Music, error) {
	var music Music
    found, err := d.goquDB.From("music").
        Select("name", "source", "key", "data", "hash").
        Where(goqu.I("name").Eq(name)).
        ScanStruct(&music)
    if err != nil {
        return Music{}, err
    }
    if !found {
        return Music{}, sql.ErrNoRows
    }
    return music, nil
}

func (d *DB) GetMusicByHash(hash string) (Music, error) {
	var music Music
	found, err := d.goquDB.From("music").
		Select("name", "source", "key", "data", "hash").
		Where(goqu.I("hash").Eq(hash)).
		ScanStruct(&music)
	if err != nil {
		return Music{}, err
	}
	if !found {
		return Music{}, sql.ErrNoRows
	}
	return music, nil
}


func hash(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (m Music) GetHash() string {
	return m.Hash
}

func (d *DB) GetMusicByHashPrefix(hash_p string) (Music, error) {
	var music Music 

	query, args, err := d.goquDB.From("music").
		Select("name", "source", "data", "key", "hash").
		Where(goqu.L("SUBSTRING(hash, 1 , ?)", shared.HashPrefixLength).Eq(goqu.L("SUBSTRING(?,1,?)", hash_p, shared.HashPrefixLength))).
		ToSQL()

	if err != nil {
		return music, err
	}

	err = d.db.QueryRow(query, args...).Scan(
		&music.Name,
		&music.Source,
		&music.Key,
		&music.Data,
		&music.Hash,
	)
	return music, err
}

func(d * DB)FilterMusic(query string)([]Music, error) {
	sql, args, err := goqu.From("music").
	Select("name", "source", "key", "data", "hash").
	Where(goqu.C("name").Like(fmt.Sprintf("%%%s%%", query))).
	ToSQL()

	if err != nil {
		return nil, err
	}

	rows, err := d.db.Query(sql, args...)
	if err != nil {
		return nil , err
	}

	defer rows.Close()

	var musics []Music
	for rows.Next() {
		var music Music
		err := rows.Scan(
			&music.Name,
            &music.Source,
            &music.Key,
            &music.Data,
            &music.Hash,
		)
		if err != nil {
			return nil, err
		}
		musics = append(musics, music)
	}
	return musics, nil
}

func (d * DB) CleanCache() error  {
	_, err := d.goquDB.Delete("music").
		Where(goqu.C("name").NotIn(goqu.From("music_playlist").Select("music_name"))).
		Executor().Exec()
	
	return err
}

func (d *DB) GetCachedMusics() ([]Music, error) {
    var musics []Music
    
    err := d.goquDB.From("music").
        Select("name", "source", "key", "data", "hash").
        Where(goqu.C("hash").NotIn(
            goqu.From("music").
                InnerJoin(goqu.T("music_playlist"), goqu.On(goqu.I("music.name").Eq(goqu.I("music_playlist.music_name")))).
                Select("hash"),
        )).
        ScanStructs(&musics)
    
    return musics, err
}