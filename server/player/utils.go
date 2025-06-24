package player

import(
	"bytes"
	"github.com/gopxl/beep"
	"io"
	"github.com/gopxl/beep/mp3"
	"crypto/md5"
	"encoding/hex"
	"github.com/perdokcat/TermoTune/shared"
	"github.com/gopxl/beep/speaker"
)

	


type customReadCloser struct {
	io.Reader
	io.Seeker
}
func (crc *customReadCloser) Close() error {
	return nil
}

func MusicDecode(data []byte) (beep.StreamSeekCloser, beep.Format, error) {
	reader := bytes.NewReader(data)
	readerCloser := &customReadCloser{Reader: reader, Seeker: reader}
	return mp3.Decode(readerCloser)
}

func hashData(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

func (p *Player) isSpeakerLocked() bool {
	return p.getPlayerState() == shared.Paused
}

func (p *Player) concernSpeakerLock(callback func()) {
	if p.isSpeakerLocked() {
		speaker.Unlock()
		callback()
		speaker.Lock()
	} else {
		callback()
	}
}