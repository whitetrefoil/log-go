package main

import (
	"os"
	"whitetrefoil.com/log-go/server"
)

func main() {
	cfg, err := server.GetConfig(os.Args)
	if err != nil {
		os.Exit(-1)
		return
	}
	server.Start(cfg)
}
