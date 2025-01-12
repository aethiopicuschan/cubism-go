package model

// struct for *.model3.json files
type ModelJson struct {
	Version        int `json:"Version"`
	FileReferences struct {
		Moc         string   `json:"Moc"`
		Textures    []string `json:"Textures"`
		Physics     string   `json:"Physics"`
		Pose        string   `json:"Pose"`
		DisplayInfo string   `json:"DisplayInfo"`
		Expressions []struct {
			Name string `json:"Name"`
			File string `json:"File"`
		} `json:"Expressions"`
		Motions map[string][]struct {
			File        string  `json:"File"`
			FadeInTime  float64 `json:"FadeInTime"`
			FadeOutTime float64 `json:"FadeOutTime"`
			Sound       string  `json:"Sound"`
			MotionSync  string  `json:"MotionSync"`
		} `json:"Motions"`
		UserData string `json:"UserData"`
	} `json:"FileReferences"`
	Groups   []Group   `json:"Groups"`
	HitAreas []HitArea `json:"HitAreas"`
}

type Group struct {
	Target string   `json:"Target"`
	Name   string   `json:"Name"`
	Ids    []string `json:"Ids"`
}

type HitArea struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}
