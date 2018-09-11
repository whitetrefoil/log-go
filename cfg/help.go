package cfg

import (
	"fmt"
	"whitetrefoil.com/log-go"
)

var helpMessage = fmt.Sprintf("\n"+
	"log.go version %s\n"+
	"\n"+
	"Usage:\n"+
	"  To start a central server:\n"+
	"    $ log.io-server <server_config_file.toml>\n"+
	"  To start a harvester:\n"+
	"    $ log.io-harvester <harvester_config_file.toml>\n"+
	"  To print version number in plain text:\n"+
	"    $ log.io-server version\n"+
	"    or\n"+
	"    $ log.io-harvester version\n"+
	"\n"+
	"Config File:\n"+
	"  Config files above must have extname \".toml\" and in TOML format.\n"+
	"  If you don't have a config file, give a non-existed file name,\n"+
	"  we will create one there will all default options for you.\n",
	log_go.Version,
)

func printHelpMessage() {
	fmt.Print(helpMessage)
}
