package player

import (
	"os"
	"path/filepath"

	"github.com/perdokcat/TermoTune/logger"
	"github.com/perdokcat/TermoTune/server/player/db"
	"github.com/perdokcat/TermoTune/shared"
)

type callback func (m db.Music) error


func (p *Player) AddMusicFromHash(
	hash string,
	how callback,
) error {
	m, err := p.Director.DB.GetMusicByHashPrefix(hash)
	if err != nil {
		return err
	}
	if how != nil {
		err = how(m)
			if err != nil {
				return err
			}
	}
	return nil
}


func (p * Player) AddMusicFromFile(
	file_path string, 
	how callback,
) error {
	logger.LogInfor(
		"Reading File",
	)
	data , err := os.ReadFile(file_path)
	if err != nil  {
		return err
	}
	
	var m db.Music
	m,err = p.Director.DB.GetMusicByHash(hashData(data))

	if err != nil {
		music ,err := NewMusic(
			filepath.Base(file_path),
			data,
		)
		if err != nil {
			logger.LogInfor("Failed to create new Music", err)
			return err
		}

		err = p.Director.DB.AddMusic(
			&db.Music{
				Name: music.Name,
				Data: music.Data,
				Source: "local",
				Key: file_path,
			},
		)
		if err != nil {
			logger.LogInfor("Failed to add music", err)
			return err
		}
		m , err = p.Director.GetMusicByHash(hashData(data))
		if err != nil {
			return err
		}else {
			logger.LogWarn("file found in db" , file_path)
		}
	}
	if how != nil {
		err = how(m)
		if err != nil {
			return err
			}
	}
	return nil
}

func (p * Player) AddMusicFromDir(dir_path string , how callback) error {
	dir , err := os.Open(dir_path)
	if err != nil {
		return err
	}
	entries, err := dir.ReadDir(0) 
	if err != nil{
		return nil
	}

	for _ , entry := range entries {
		if !entry.IsDir() {
			music_path := filepath.Join(dir_path, entry.Name())
			err := p.AddMusicFromFile(music_path, how)
			if err != nil  {
				logger.LogWarn("skipping music", entry.Name() + "with error" , err)
			}
			continue
		}
	}
	return nil
}

// the unique is the unique id of the music in the engine it can be url or id
func(p * Player) AddMusicFromOnline()  {
	return err	
}