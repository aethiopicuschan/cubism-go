package model

// struct for *.cdi3.json files
type CdiJson struct {
	Version    int `json:"Version"`
	Parameters []struct {
		Id      string `json:"Id"`
		GroupId string `json:"GroupId"`
		Name    string `json:"Name"`
	} `json:"Parameters"`
	ParameterGroups []struct {
		Id      string `json:"Id"`
		GroupId string `json:"GroupId"`
		Name    string `json:"Name"`
	} `json:"ParameterGroups"`
	Parts []struct {
		Id   string `json:"Id"`
		Name string `json:"Name"`
	} `json:"Parts"`
}
