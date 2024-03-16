package model

// *.pose3.json用の構造体
type PoseJson struct {
	Type       string  `json:"Type"`
	FadeInTime float64 `json:"FadeInTime"`
	Groups     [][]struct {
		Id   string   `json:"Id"`
		Link []string `json:"Link"`
	} `json:"Groups"`
}
