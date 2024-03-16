package model

// *.exp3.json用の構造体
type ExpJson struct {
	Name       string
	Type       string `json:"Type"`
	Parameters []struct {
		Id    string  `json:"Id"`
		Value float64 `json:"Value"`
		Blend string  `json:"Blend"`
	} `json:"Parameters"`
}
