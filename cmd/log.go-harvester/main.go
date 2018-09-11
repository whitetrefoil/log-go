package main

import (
	"os"
	"whitetrefoil.com/log-go/harvester"
)

func main() {
	cfg, err := harvester.GetConfig(os.Args)
	if err != nil {
		os.Exit(-1)
		return
	}
	harvester.Start(cfg)
}
