package server

import (
	"time"
	"whitetrefoil.com/log-go/cfg"
)

type heartbeatConfig struct {
	Timeout   *cfg.Duration `toml:"timeout"`
	CheckFreq *cfg.Duration `toml:"check_freq"`
}

type wwwConfig struct {
	Host string `toml:"host"`
	Port uint16 `toml:"port"`
}

type apiConfig struct {
	Host string `toml:"host"`
	Port uint16 `toml:"port"`
}

type config struct {
	Name      string           `toml:"name"`
	Api       *apiConfig       `toml:"api"`
	Www       *wwwConfig       `toml:"www"`
	Heartbeat *heartbeatConfig `toml:"heartbeat"`
}

var DefaultConfig *config = &config{
	Name: "SERVER",
	Api: &apiConfig{
		Host: "localhost",
		Port: 48752,
	},
	Www: &wwwConfig{
		Host: "0.0.0.0",
		Port: 48753,
	},
	Heartbeat: &heartbeatConfig{
		Timeout:   &cfg.Duration{20 * time.Second},
		CheckFreq: &cfg.Duration{10 * time.Second},
	},
}

func GetConfig(argv []string) (*config, error) {
	def := DefaultConfig
	www := *def.Www
	api := *def.Api
	hb := *def.Heartbeat
	def.Www = &www
	def.Api = &api
	def.Heartbeat = &hb

	err := cfg.Get(argv, def)
	if err != nil {
		return nil, err
	}
	return def, nil
}
