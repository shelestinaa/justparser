package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/shelestinaa/justparser/internal/app/apiserver"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs\apiserver.toml", "path to config file")
}

func main() {

	flag.Parse()

	config := new(apiserver.Config)
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	s := apiserver.New(config)
	err := s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
