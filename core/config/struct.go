package config

var c struct {
	DB struct {
		User string `json:"user"`
		Pass string `json:"pass"`
		Host string `json:"host"`
		Port uint16 `json:"port"`
		Name string `json:"name"`
	} `json:"db"`

	IsDebug bool `json:"is_debug"`
}
