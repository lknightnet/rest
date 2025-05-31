package main

import (
	"backend-mobAppRest/config"
	"backend-mobAppRest/internal"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	internal.Run(cfg)
}
