package main

import (
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lexruntimeservice"
)

type application struct {
	config   *Config
	errorLog *log.Logger
	infoLog  *log.Logger
	lex      *lexruntimeservice.LexRuntimeService
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
		wg:       &sync.WaitGroup{},
	}

	// initialise lex using environment
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	svc := lexruntimeservice.New(sess)
	infoLog.Println("testing", svc.ServiceID, svc.ServiceName, svc.APIVersion)
	_, err = svc.PostText(&lexruntimeservice.PostTextInput{
		BotAlias:  &app.config.Lex.Alias,
		BotName:   &app.config.Lex.BotName,
		InputText: aws.String("hello"),
		UserId:    aws.String("dummy"),
	})
	if err != nil {
		errorLog.Fatalf("error communicating with lex: %s", err)
	}
	infoLog.Println("successfully connected to lex")
	app.lex = svc

	err = app.serve()
	if err != nil {
		app.errorLog.Fatal(err)
	}
}
