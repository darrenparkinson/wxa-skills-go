package main

import (
	"log"
	"os"
	"sync"
)

type application struct {
	config   *Config
	errorLog *log.Logger
	infoLog  *log.Logger
	models   models
	wg       *sync.WaitGroup
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := &application{
		config:   cfg,
		errorLog: errorLog,
		infoLog:  infoLog,
		models:   newModels(),
		wg:       &sync.WaitGroup{},
	}
	err = app.serve()
	if err != nil {
		app.errorLog.Fatal(err)
	}

}
