package player

import (
	"sync"
)

type MusicQueue struct {
	queue   []Music
	current int
	mu      *sync.Mutex
}

func NewMusicQueue() *MusicQueue {
	return &MusicQueue{
		queue:   make([]Music, 0), // using silde like a queue
		current: 0,
		mu:      &sync.Mutex{},
	}
}

func (q *MusicQueue) GetTitles() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	titles := make([]string, 0)
	for _, music := range q.queue {
		titles = append(titles, music.Name)
	}
	return titles
}

func (q *MusicQueue) GetCurrentIndex() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.current
}

func (q *MusicQueue) SetCurrentIndex(current_index int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.current = current_index
}

func (q *MusicQueue) GetMusicByIndex(index int) *Music {
	if index < 0 && index > len(q.queue) {
		return nil
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	return &q.queue[index]
}

func (q *MusicQueue) GetMusicByName(music_name string) *Music {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, music := range q.queue {
		if music.Name == music_name {
			return &music
		}
	}
	return nil
}

func (q *MusicQueue) IsEmpty() bool {
	if len(q.queue) == 0 {
		return true
	}
	return false
}

func (q *MusicQueue) GetCurrentMusic() *Music {
	q.mu.Lock()
	q.mu.Unlock()
	if q.IsEmpty() {
		return nil
	}

	return q.GetMusicByIndex(q.GetCurrentIndex())
}

func (q *MusicQueue) EnQueue(_music Music) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, music := range q.queue {
		if hashData(music.Data) == hashData(_music.Data) {
			return
		}
	}
	q.queue = append(q.queue, _music)
}

func (q *MusicQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.queue)
}

func (q *MusicQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, music := range q.queue {
		music.Streamer().Close()
	}
	q.queue = make([]Music, 0)
	q.SetCurrentIndex(0)
}

func (q *MusicQueue) Remove(_music *Music) {
	index := -1
	for i, music := range q.queue {
		if hashData(music.Data) == hashData(_music.Data) {
			index = i
			break
		}
	}
	if index < 0 {
		return
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue[index].Streamer().Close()
	q.queue = append(q.queue[:index], q.queue[index+1:]...)
}

func (q *MusicQueue) QueueNext() {
	index := q.GetCurrentIndex() + 1
	if index > q.Size()-1 {
		index = 0
	}
	q.SetCurrentIndex(index)
}

func (q * MusicQueue) QueuePrev() {
	index := q.GetCurrentIndex() - 1
	
	if index < 0 {
		index = q.Size() - 1
	}

	q.SetCurrentIndex(index)

}
