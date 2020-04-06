// GNU GPL v3 License
// Copyright (c) 2020 github.com:go-trellis

package conf

type Config struct {
	Servers []*Server `yaml:"server,omitempty" json:"server,omitempty"`
}

type Server struct {
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
}
