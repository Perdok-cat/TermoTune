package player

import(
	
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


