package model

// struct for *.pose3.json files
type PoseJson struct {
	Type       string  `json:"Type"`
	FadeInTime float64 `json:"FadeInTime"`
	Groups     [][]struct {
		Id   string   `json:"Id"`
		Link []string `json:"Link"`
	} `json:"Groups"`
}
