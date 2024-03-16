package model

// *.userdata3.json用の構造体
type UserDataJson struct {
	Version int `json:"Version"`
	Meta    struct {
		UserDataCount     int `json:"UserDataCount"`
		TotalUserDataSize int `json:"TotalUserDataSize"`
	} `json:"Meta"`
	UserData []struct {
		Target string `json:"Target"`
		Id     string `json:"Id"`
		Value  string `json:"Value"`
	} `json:"UserData"`
}
