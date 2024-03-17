package sound

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aethiopicuschan/cubism-go/sound"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error {
	return nil
}

var initialized = false

type DefaultSound struct {
	streamer beep.StreamSeekCloser
	format   beep.Format
	ctrl     *beep.Ctrl
}

func LoadSound(fp string) (s sound.Sound, err error) {
	ds := &DefaultSound{}
	buf, err := os.ReadFile(fp)
	if err != nil {
		return
	}
	return ds, ds.Decode(fp, buf)
}

func (s *DefaultSound) Decode(fp string, buf []byte) (err error) {
	if s.ctrl != nil {
		return
	}
	f, err := detectFormat(fp)
	if err != nil {
		return
	}
	switch f {
	case "wav":
		s.streamer, s.format, err = wav.Decode(bytes.NewReader(buf))
	case "mp3":
		s.streamer, s.format, err = mp3.Decode(nopCloser{bytes.NewReader(buf)})
	default:
		err = fmt.Errorf("unsupported format: %s", f)
		return
	}
	if err != nil {
		return
	}
	s.ctrl = &beep.Ctrl{Streamer: s.streamer}
	if !initialized {
		err = speaker.Init(s.format.SampleRate, s.format.SampleRate.N(time.Second/10))
	}
	return
}

func (s *DefaultSound) Play() {
	s.streamer.Seek(0)
	s.ctrl.Paused = false
	speaker.Play(s.ctrl)
}

func (s *DefaultSound) Close() {
	s.ctrl.Paused = true
	s.streamer.Seek(0)
}

func detectFormat(fp string) (f string, err error) {
	ext := filepath.Ext(fp)
	switch ext {
	case ".wav", ".wave":
		f = "wav"
	case ".mp3":
		f = "mp3"
	default:
		err = fmt.Errorf("unsupported format: %s", ext)
	}
	return
}
