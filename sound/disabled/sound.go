package disabled

import "github.com/aethiopicuschan/cubism-go/sound"

type Sound struct {
	fp string
}

func LoadSound(fp string) (s sound.Sound, err error) {
	ds := &Sound{
		fp: fp,
	}
	return ds, nil
}

func (s *Sound) Play() (err error) {
	return
}

func (s *Sound) Close() {}
