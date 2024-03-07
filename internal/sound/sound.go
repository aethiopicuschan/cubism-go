package sound

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"time"

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

func Play(f string, buf []byte) (err error) {
	var streamer beep.StreamSeekCloser
	var format beep.Format
	switch f {
	case "wav":
		streamer, format, err = wav.Decode(bytes.NewReader(buf))
	case "mp3":
		streamer, format, err = mp3.Decode(nopCloser{bytes.NewReader(buf)})
	default:
		err = fmt.Errorf("unsupported format: %s", f)
		return
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
	return
}

func DetectFormat(fp string) (f string, err error) {
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
