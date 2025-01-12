package model

// struct for *.physics3.json files
type PhysicsJson struct {
	Version int `json:"Version"`
	Meta    struct {
		PhysicsSettingCount int `json:"PhysicsSettingCount"`
		TotalInputCount     int `json:"TotalInputCount"`
		TotalOutputCount    int `json:"TotalOutputCount"`
		VertexCount         int `json:"VertexCount"`
		EffectiveForces     struct {
			Gravity struct {
				X float64 `json:"X"`
				Y float64 `json:"Y"`
			} `json:"Gravity"`
			Wind struct {
				X float64 `json:"X"`
				Y float64 `json:"Y"`
			} `json:"Wind"`
		} `json:"EffectiveForces"`
		PhysicsDictionary []struct {
			Id   string `json:"Id"`
			Name string `json:"Name"`
		} `json:"PhysicsDictionary"`
	} `json:"Meta"`
	PhysicsSettings []struct {
		Id    string `json:"Id"`
		Input []struct {
			Source struct {
				Target string `json:"Target"`
				Id     string `json:"Id"`
			} `json:"Source"`
			Weight  float64 `json:"Weight"`
			Type    string  `json:"Type"`
			Reflect bool    `json:"Reflect"`
		} `json:"Input"`
		Output []struct {
			Destination struct {
				Target string `json:"Target"`
				Id     string `json:"Id"`
			} `json:"Destination"`
			VertexIndex int     `json:"VertexIndex"`
			Scale       float64 `json:"Scale"`
			Weight      float64 `json:"Weight"`
			Type        string  `json:"Type"`
			Reflect     bool    `json:"Reflect"`
		} `json:"Output"`
		Vertices []struct {
			Position struct {
				X float64 `json:"X"`
				Y float64 `json:"Y"`
			} `json:"Position"`
			Mobility     float64 `json:"Mobility"`
			Delay        float64 `json:"Delay"`
			Acceleration float64 `json:"Acceleration"`
			Radius       float64 `json:"Radius"`
		} `json:"Vertices"`
		Normalization struct {
			Position struct {
				Minimum float64 `json:"Minimum"`
				Default float64 `json:"Default"`
				Maximum float64 `json:"Maximum"`
			} `json:"Position"`
			Angle struct {
				Minimum float64 `json:"Minimum"`
				Default float64 `json:"Default"`
				Maximum float64 `json:"Maximum"`
			} `json:"Angle"`
		} `json:"Normalization"`
	} `json:"PhysicsSettings"`
}
