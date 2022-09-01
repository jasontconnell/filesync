package conf

import "github.com/jasontconnell/conf"

type Config struct {
	Role     string   `json:"role"`
	Path     string   `json:"path"`
	Clients  []string `json:"clients"`
	Bind     string   `json:"bind"`
	Schedule string   `json:"schedule"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	err := conf.LoadConfig(filename, &config)
	return config, err
}
