package model

// struct for *.exp3.json files
type ExpJson struct {
	Name       string
	Type       string `json:"Type"`
	Parameters []struct {
		Id    string  `json:"Id"`
		Value float64 `json:"Value"`
		Blend string  `json:"Blend"`
	} `json:"Parameters"`
}
