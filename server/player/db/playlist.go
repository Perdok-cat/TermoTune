package db

import(
	"github.com/doug-martin/goqu/v9"
	"fmt"
	"strings"
)


type PlayList struct {
	name string
}


func (d *DB) InitPlaylist() error {
	_, err := d.db.Exec(
		`CREATE TABLE IF NOT EXISTS playlist (
      name TEXT PRIMARY KEY
    )`,
	)
	return err
}

// relation ship
func (d *DB) InitMusicPlaylist() error {
	_, err := d.db.Exec(
		`CREATE TABLE IF NOT EXISTS music_playlist (
      music_name TEXT,
      playlist_name TEXT,
      PRIMARY KEY (music_name, playlist_name),
      FOREIGN KEY (music_name) REFERENCES music (name),
      FOREIGN KEY (playlist_name) REFERENCES playlist (name)
    )`,
	)
	return err
}

func (d * DB) GetPlaylist(name string) (PlayList, error){
	var playlist PlayList

	_, err := d.goquDB.From("playlist").
		Select("name").
		Where(goqu.C("name").Eq(name)).
		ScanStruct(&playlist)

	return playlist, err
}


func (d * DB) GetPlayLists()([]PlayList, error){
	var playlists []PlayList

	err := d.goquDB.From("playlist").
	Select("name").
	ScanStructs(&playlists)

	return playlists, err
}

func(d * DB) AddPlayList(playlist_name string) error {
	_, err := d.goquDB.Insert("playlist").
	Rows(goqu.Record{"name" : playlist_name}).
	OnConflict(goqu.DoNothing()).
	Executor().Exec()


	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed"){
		return fmt.Errorf(
			"Playlist %s already exists",playlist_name,
		)
	}
	return nil
}

func (d * DB) RemovePlaylist(playlist_name string) error {
	_ , err := d.goquDB.Delete("playlist").
	Where(goqu.C("name").Eq(playlist_name)).
	Executor().Exec()

	if err != nil {
		return err 
	}

	_ , err = d.goquDB.Delete("music_playlist").
	Where(goqu.C("playlist_name").Eq(playlist_name)).
	Executor().Exec()

	if err != nil {
		return err
	}
	return nil
}

func (d * DB) AddMusicToPlaylist(music_name, playlist_name string) error {
	_ , err := d.goquDB.Insert("music_playlist").
	Rows(goqu.Record{"music_name" : music_name , "playlist_name": playlist_name}).
	Executor().Exec()

	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return fmt.Errorf(
			"Music %s already in playlist %s",
			music_name,
			playlist_name,
		)
	}

	return nil
}

func (d * DB) RemoveMusicFromPlaylist(
	music_name string,
	playlist_name string,

) error{
	_ , err := d.goquDB.Delete("music_playlist").
	Where(goqu.C("music_name").Eq(music_name), goqu.C("playlist_name").Eq(playlist_name)).
	Executor().Exec()

	if err != nil {
		return err 
	}
	return nil 

}

func (d * DB) GetMusicFromPlaylist(
	playlist_name string,
) ([]Music, error) {
	var musics []Music

	_, err := d.goquDB.From("music_playlist").
	Select().
	Where(goqu.C("playlist_name").Eq(playlist_name)).
	ScanStruct(&musics)

	if err != nil {
		return nil,err 
	}

	return musics, nil
}


