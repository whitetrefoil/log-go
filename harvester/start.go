package harvester

import (
	"fmt"
	"os"
	"os/signal"
	"whitetrefoil.com/log-go/hs"
	loggerPkg "whitetrefoil.com/log-go/logger"
)

var logger *loggerPkg.Logger

func Start(config *config) {
	logger = GetLogger(config.Name)

	sysSig := make(chan os.Signal, 1)
	signal.Notify(sysSig, os.Interrupt)
	signal.Notify(sysSig, os.Kill)

	hb, err := NewHeartbeat(config)
	if err != nil {
		logger.Logln("The heart failed to beat :( due to: ", err)
		os.Exit(-1)
		return
	}

	handler := NewHandler(config.Files)
	addr := fmt.Sprintf("%s:%d", config.Api.Host, config.Api.Port)
	server := hs.NewServer(addr, handler, logger)
	server.Start()

	go func() {
		sig := <-sysSig
		logger.Logf("Received signal \"%v\"...", sig)
		hb.RIP()
		server.Stop()
	}()

	<-hb.Ended
	<-server.Ended
	logger.Logln("Everything has been shutdown, nice and clean, bye~")
}
