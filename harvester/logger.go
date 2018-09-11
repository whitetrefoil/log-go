package harvester

import loggerPkg "whitetrefoil.com/log-go/logger"

func GetLogger(tag string) *loggerPkg.Logger {
	return loggerPkg.NewLogger(tag)
}
