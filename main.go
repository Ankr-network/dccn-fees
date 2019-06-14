package main

import (
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"github.com/Ankr-network/dccn-fees/db-service"
	"github.com/Ankr-network/dccn-fees/handler"
	"log"
)

var (
	db  dbservice.DBService
	err error
)

func main() {
	Init()

	if db, err = dbservice.New(); err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	startHandler()
}

// Init starts handler to listen.
func Init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	config := micro2.LoadConfigFromEnv()
	config.Show()
}

func startHandler() {


	service := micro2.NewService()

	handler := handler.NewHandler(db)
	dcmgr.RegisterFeesServiceServer(service.GetServer(), handler)
	service.Start()
}
