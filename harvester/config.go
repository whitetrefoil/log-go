package harvester

import (
	"fmt"
	"time"
	"whitetrefoil.com/log-go/cfg"
)

type heartbeatConfig struct {
	Timeout *cfg.Duration `toml:"timeout"`
	Freq    *cfg.Duration `toml:"freq"`
}

type apiConfig struct {
	Host string `toml:"host"`
	Port uint16 `toml:"port"`
}

type config struct {
	Name      string             `toml:"name"`
	ServerUrl string             `toml:"central_url"`
	Api       *apiConfig         `toml:"api"`
	Heartbeat *heartbeatConfig   `toml:"heartbeat"`
	Files     *map[string]string `toml:"files"`
}

var DefaultConfig *config = &config{
	Name:      "HARVESTER",
	ServerUrl: "http://localhost:48752",
	Api: &apiConfig{
		Host: "0.0.0.0",
		Port: 48751,
	},
	Heartbeat: &heartbeatConfig{
		Timeout: &cfg.Duration{5 * time.Second},
		Freq:    &cfg.Duration{10 * time.Second},
	},
	Files: &map[string]string{
		"mail": "/var/log/mail.log",
	},
}

func GetConfig(argv []string) (*config, error) {
	def := DefaultConfig
	api := *def.Api
	hb := *def.Heartbeat
	logs := *def.Files
	def.Api = &api
	def.Heartbeat = &hb
	def.Files = &logs

	err := cfg.Get(argv, def)
	if err != nil {
		return nil, err
	}
	if def.Api.Host == "0.0.0.0" {
		fmt.Println("[WARN] There's a wired problem with central server when connecting to harvester via 0.0.0.0, will instead use localhost...")
		def.Api.Host = "localhost"
	}
	return def, nil
}
