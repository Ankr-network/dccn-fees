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
	if rsp, err := dcClient.UserHistoryFeesList(tokenContext, &dcmgr.HistoryFeesRequest{Uid:"admin9880", Start:"2019-03-01",  End: "2019-06-01"}); err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Printf("UserHistoryFeesList %+v \n", rsp)
		//for i := 0; i < len(rsp.Week); i++ {
		//	d := rsp.Week[i]
		//	tm := time.Unix(d.Date.Seconds, 0)
		//	layOut := "2006-01-02"
		//	fmt.Printf("income  %d  cputotal %d  cpuused %d memey total %d  memey used %d data %s \n", d.Income, d.Usage.CpuTotal,d.Usage.CpuUsed, d.Usage.MemoryTotal, d.Usage.MemoryUsed, tm.Format(layOut))
		//}



	}

}
