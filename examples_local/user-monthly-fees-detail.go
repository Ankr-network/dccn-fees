package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var addr = "localhost:50051"


func main() {

	log.SetFlags(log.LstdFlags | log.Llongfile)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Println(err.Error())
		}
	}(conn)

	dcClient := dcmgr.NewFeesServiceClient(conn)


	md := metadata.New(map[string]string{
		"token": "",
	})

	//log.Printf("get access_token after login %s \n", access_token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	tokenContext, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// var userTasks []*common_proto.Task
	if rsp, err := dcClient.MonthFeesDetail(tokenContext, &dcmgr.FeesDetailRequest{Uid:"1111", Month:""}); err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Printf("MonthFeesDetail %+v \n", rsp)


	}

}
