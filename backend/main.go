package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"io"
	"log"
	"net/http"
	"teirxserver/src/cfg"
	"teirxserver/src/core"
	"teirxserver/src/dbapi"
	"teirxserver/src/txlog"
)

func movieGet() {
	err := cfg.LoadAppConfig("teirxcfg.json")
	if err != nil {
		log.Panic(err)
		return
	}

	appConfig := cfg.GetAppConfig()

	url2 := "http://www.omdbapi.com/?i=tt3896198&apikey=2b90379"
	url := fmt.Sprintf("http://www.omdbapi.com/?type='movie&s='Blade runner'&apikey=%s", appConfig.Keys.Omdb)

	fmt.Println(url)
	fmt.Println(url2)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

	// Pretty print JSON
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "    ")
	if err != nil {
		return
	}

	fmt.Println(string(prettyJSON.Bytes()))
}

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

    core.PrepareEnvironment()

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
