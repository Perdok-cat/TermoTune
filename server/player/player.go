package player

import (
	"log"
	"os"
	"sync"
	"time"
	"strconv"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/perdokcat/TermoTune/config"
	"github.com/perdokcat/TermoTune/logger"
	"github.com/perdokcat/TermoTune/shared"
	"go.uber.org/zap"
	
)

var Instance *Player

var once sync.Once


type lmeta struct {
	_lcurrentPos time.Duration
	_lcurrentDur time.Duration
}

type Player struct {
	Queue       *MusicQueue
	playerState shared.PState
	done        chan struct{}
	initialised bool
	Director    *Director
	Tasks       map[string]shared.Task
	Vol         uint8
	_lmeta      lmeta
	mu          sync.Mutex
}

func NewPlayer() *Player {
	if _, err := os.Stat(config.GetConfig().TermoTunePath); os.IsNotExist(err) {
		logger.LogInfor("TermoTune file is not exists, i will create it")
		err := os.Mkdir(config.GetConfig().TermoTunePath, 0o755)
		if err != nil {
			log.Fatal(err)
		}
	}
	director, err := NewDeFaultDirector()
	if err != nil {
		log.Fatal(err)
	}
	return &Player {
		Queue: NewMusicQueue(),
		playerState: shared.Stopped,
		done: make(chan struct{}),
		initialised: false,
		Director: director,
		Vol: 100,
		Tasks: make(map[string]shared.Task),

	}
}

func (p *Player) _setMeta(m lmeta) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p._lmeta = m
}

func (p * Player) _getMeta() lmeta {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p._lmeta
}

func(p *Player) GetPlayer() *Player{
	once.Do(func() {
		Instance = NewPlayer()
	})
	return Instance
}

func (p*Player) Play() error {
	if(p.Queue.IsEmpty()) {
		return logger.NewTermoTuneError("queue is empty")
	}
	music := p.Queue.GetCurrentMusic()
	if music == nil {
		return logger.NewTermoTuneError("Failed to get current Music")
	}
	if !p.initialised {
		err := speaker.Init(
			music.Format.SampleRate,
			music.Format.SampleRate.N(time.Second/10),
		)
		if err != nil {
			return err
		}
	}else {
		if !p.isSpeakerLocked() {
			speaker.Clear()
		}else {
			speaker.Unlock()
			speaker.Clear()
		}
	}
	p.setPlayerSate(shared.Playing)

	go func() {
		done := make(chan struct{})
		music.SetVolume(p.Vol)
		speaker.Play(
			beep.Seq(
				music.Volume,
				beep.Callback(func() {done <- struct{}{}}),
			),
		)
		<-done
		p.Next()
	}()
	return nil
}

func (p *Player) getPlayerState() shared.PState {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.playerState
}

func (p *Player) setPlayerSate(plState shared.PState) {
	p.mu.Lock()
	p.playerState = plState
	p.mu.Unlock()
}

func (p *Player) Next() error {
	if p.Queue.IsEmpty() {
		return logger.NewTermoTuneError("queue is empty")
	}
	if p.getPlayerState() == shared.Stopped {
		return logger.NewTermoTuneError("player is stopped")
	}
	if p.getPlayerState() == shared.Paused {
		p.Resume()
	}
	currentMusic := p.Queue.GetCurrentMusic()

	if currentMusic == nil {
		return logger.NewTermoTuneError("failed to get current music")
	}
	
	if err := currentMusic.SetPositionD(0); err != nil {
		return logger.NewTermoTuneError("failed to set possition ")
	}

	p.Queue.QueueNext()
	err := p.Play()
	if  err != nil {
		return logger.NewTermoTuneError("failed to play next music")
	}
	return nil
}

func (p * Player) Prev() error {
	if p.Queue.IsEmpty() {
		return logger.NewTermoTuneError("Queue is empty")
	}
	if p.playerState == shared.Stopped {
		return logger.NewTermoTuneError("Player is Stopped")
	}
	if p.playerState == shared.Paused {
		p.Resume()
	}
	currentMusic := p.Queue.GetCurrentMusic()
	if currentMusic == nil {
		return logger.NewTermoTuneError("failed to get current music")
	}
	if err := currentMusic.SetPositionD(0); err != nil {
		return logger.NewTermoTuneError("failed to set Postition duration")
	}
	p.Queue.QueuePrev()
	err := p.Play()
	if err != nil {
		return err
	}
	return nil 

}

func (p *Player) Stop() error {
	clear(p.Tasks)

	state := p.getPlayerState()

	if state == shared.Stopped {
		return nil
	}
	
	if state == shared.Paused {
		p.Resume()
	}
	speaker.Clear()
	p.Queue.Clear()
	p.setPlayerSate(shared.Stopped)
	return nil	
}

func (p *Player) Pause() error {
	state := p.getPlayerState()
	if state == shared.Paused || state == shared.Stopped {
		return nil
	}
	p._setMeta(
		lmeta{
			_lcurrentDur: p.GetCurrMusicDuration(),
			_lcurrentPos: p.GetCurrMusicPosition(),
		},
	)
	p.setPlayerSate(shared.Paused)
	speaker.Lock()
	return nil
}

func (p * Player) Resume() error {
	state := p.getPlayerState()
	if state == shared.Playing || state == shared.Stopped {
		return nil
	}
	p.setPlayerSate(shared.Playing)
	speaker.Unlock()
	return nil
}

func (p *Player) Seek(d time.Duration) error {
	state := p.getPlayerState()
	if state == shared.Stopped {
		return logger.NewTermoTuneError("player is not runing")
	}
	if p.Queue.IsEmpty() {
		return logger.NewTermoTuneError("queue is empty")
	}
	currentMusic := p.Queue.GetCurrentMusic()
	if currentMusic == nil {
		return logger.NewTermoTuneError("failed to get current music")
	}
	if state == shared.Paused {
		err := p.Resume()
		if err != nil {
			return err
		}
		defer p.Pause()
	}
	if err := currentMusic.Seek(d); err != nil {
		return logger.NewTermoTuneError("failed to seek")
	}
	return nil
}

func (p *Player) Volume(
	vp uint8,
) error {
	if p.getPlayerState() == shared.Stopped {
		return nil
	}
	p.Vol = vp
	currentMusic := p.Queue.GetCurrentMusic()
	if currentMusic == nil {
		return logger.NewTermoTuneError("Current music is null")
	}
	p.concernSpeakerLock(
		func() {
			currentMusic.SetVolume(
				vp,
			)
		},
	)
	return nil
}

func (p *Player) Remove(music shared.IntOrString) error {
	if p.Queue.IsEmpty() {
		return logger.NewTermoTuneError("queue is empty")
	}
	if p.Queue.Size() == 1 {
		p.Stop()
		return nil
	} else {
		var m *Music
		if music.IsInt {
			musicIndex := music.IntVal
			logger.LogInfo(
				"Removing music by index",
				zap.String("index", strconv.Itoa(musicIndex)),
			)
			m = p.Queue.GetMusicByIndex(
				musicIndex,
			)
		} else {
			musicName := music.StrVal
			m = p.Queue.GetMusicByName(
				musicName,
			)
			logger.LogInfo(
				"Removing music by name",
				zap.String("Name", musicName),
			)
		}
		if m == nil {
			return logger.NewTermoTuneError("music is nil")
		}

		if m.Name == p.Queue.GetCurrentMusic().Name {
			p.Queue.Remove(m)
			return p.Next()
		}

		p.Queue.Remove(
			m,
		)
	}
	return nil
}

// play list method

func (p *Player) CreatePlayList(plname string) error {
	_, err := p.Director.Db.GetPlaylist(
		plname,
	)

	if err == nil {
		return logger.NewTermoTuneError("play list already exists")
	}
	err = p.Director.Db.AddPlaylist(
		plname,
	)
	if err != nil {
		return logger.NewTermoTuneError("failed to create playlist")
	}
	return nil
}

func (p *Player) RemovePlayList(plname string) error {
	pl, err := p.Director.Db.GetPlaylist(
		plname,
	)
	if err != nil {
		return logger.NewTermoTuneError("Play List is not Exist")
	}
	err = p.Director.Db.RemovePlaylist(
		pl.Name,
	)
	if err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

func (p *Player) GetPlayListsNames() ([]string, error) {
	lists, err := p.Director.Db.GetPlaylists()
	if err != nil {
		return nil, logger.NewTermoTuneError("Failed to get playlistnames")
	}

	var names []string
	for _, list := range lists {
		names = append(names, list.Name)
	}
	return names, nil
}

func (p *Player) RemoveMusicFromPlayList(plname string, music shared.IntOrString) error {
	pl, err := p.Director.Db.GetPlaylist(
		plname,
	)
	if err != nil {
		return logger.NewTermoTuneError("PlayList is not exist")
	}
	var m *Music
	ms, err := p.Director.Db.GetMusicsFromPlaylist(
		pl.Name,
	)
	if music.IsInt {
		index := music.IntVal
		if index < 0 || index >= len(ms) {
			return logger.NewTermoTuneError("index out of range")
		}
		err = p.Director.Db.RemoveMusicFromPlaylist(
			pl.Name,
			ms[index].Name,
		)
	} else {
		name := music.StrVal
		err = p.Director.Db.RemoveMusicFromPlaylist(
			pl.Name,
			name,
		)
	}

	if err != nil {
		return logger.NewTermoTuneError("failed to get song from playlist")
	}

	// check if the exists in the queue and remove it
	for _, music := range p.Queue.queue {
		if hashData(music.Data) == hashData(m.Data) {
			p.Queue.Remove(
				&music,
			)
		}
	}
	return nil
}

func (p *Player) GetPlayListMusicNames(plname string) ([]string, error) {
	pl, err := p.Director.Db.GetPlaylist(
		plname,
	)
	if err != nil {
		return nil, logger.NewTermoTuneError("failed to get playlist")
	}
	songs, err := p.Director.Db.GetMusicsFromPlaylist(
		pl.Name,
	)
	if err != nil {
		return nil, logger.NewTermoTuneError("failed to get song from playlist")
	}
	var names []string
	logger.LogInfo(
		"Playlist name :",
		pl.Name,
	)
	for _, song := range songs {
		logger.LogInfo(
			"Music name :",
			song.Name,
		)
		names = append(names, song.Name)
	}
	return names, nil
}

