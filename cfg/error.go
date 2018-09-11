package cfg

import "errors"

var (
	ErrEnoughArgs     = errors.New("must specify a config file")
	ErrNotToml        = errors.New("config file must be a toml file")
	ErrFileOpen       = errors.New("failed to open file")
	ErrFileCreate     = errors.New("failed to create file")
	ErrTomlParse      = errors.New("failed to parsing toml")
	ErrPrintVersion   = errors.New("print version then exit")
	ErrExampleCreated = errors.New("a default config file has been created")
)
