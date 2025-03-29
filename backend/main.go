package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"teirxserver/src/cfg"
	"teirxserver/src/core"
	"teirxserver/src/dbapi"
	"teirxserver/src/txlog"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := cfg.LoadAppConfig("teirxcfg.json")
	if err != nil {
		log.Panic(err)
		return
	}

	appConfig := cfg.GetAppConfig()

	logLevel := txlog.Str2LogLevel(appConfig.Logging.Level)
	err = txlog.InitLogging(logLevel, appConfig.Logging.Filepath)
	if err != nil {
		log.Panic(err)
		return
	}

	txlog.TxLogInfo("**** Teirx Server Starting")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	signalRcvChan := make(chan bool, 1)
	go func() {
		sig := <-signalChan
		txlog.TxLogInfo("Recieved Signal: %s", sig)
		signalRcvChan <- true
	}()

	srv := core.NewHTTPServer()
	srv.Start()

	<-signalRcvChan

	txlog.TxLogInfo("**** Teirx Server Closing")

	srv.Stop()

	dbapi.GetDBConnection().Close()
}
