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
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTUwMTUzMTgsImp0aSI6IjQ4NTQ5YjQxLWUzNjYtNGIxMi05NTc3LTU0M2Y5NTE5Y2JlZiIsImlzcyI6ImFua3IubmV0d29yayJ9.A0p3KyxIKZHAZb_buPgadKj3d40Rlw_hSpsFBrNLjuw",
	})

	//log.Printf("get access_token after login %s \n", access_token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	tokenContext, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// var userTasks []*common_proto.Task
	if rsp, err := dcClient.ClusterDashBoard(tokenContext, &dcmgr.DashBoardRequest{TeamId: "e76dddd7-b370-4748-8b1f-4cffefb82a78"}); err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Printf("total income %d  days of one week %d \n", rsp.TotalIncome, len(rsp.Week))
		for i := 0; i < len(rsp.Week); i++ {
			d := rsp.Week[i]
			tm := time.Unix(d.Date.Seconds, 0)
			layOut := "2006-01-02"
			fmt.Printf("income  %d  cputotal %d  cpuused %d memey total %d  memey used %d data %s \n", d.Income, d.Usage.CpuTotal, d.Usage.CpuUsed, d.Usage.MemoryTotal, d.Usage.MemoryUsed, tm.Format(layOut))
		}

		fmt.Printf("total income %d  days of one month %d \n", rsp.TotalIncome, len(rsp.Month))

		for i := 0; i < len(rsp.Month); i++ {
			d := rsp.Month[i]
			tm := time.Unix(d.Date.Seconds, 0)
			layOut := "2006-01-02"
			fmt.Printf("income  %d  cputotal %d  cpuused %d memey total %d  memey used %d data %s \n", d.Income, d.Usage.CpuTotal, d.Usage.CpuUsed, d.Usage.MemoryTotal, d.Usage.MemoryUsed, tm.Format(layOut))
		}

		fmt.Printf("total income %d  days of one year %d \n", rsp.TotalIncome, len(rsp.Year))
		for i := 0; i < len(rsp.Year); i++ {
			d := rsp.Year[i]
			tm := time.Unix(d.Date.Seconds, 0)
			layOut := "2006-01-02"
			fmt.Printf("income  %d  cputotal %d  cpuused %d memey total %d  memey used %d data %s \n", d.Income, d.Usage.CpuTotal, d.Usage.CpuUsed, d.Usage.MemoryTotal, d.Usage.MemoryUsed, tm.Format(layOut))
		}

	}

}
