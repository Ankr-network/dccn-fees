package main

import (
	"encoding/json"
	"fmt"
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-fees/db-service"
	"github.com/golang/protobuf/ptypes/timestamp"
	"log"
	"os"
	"time"
	//	"github.com/micro/go-plugins/broker/rabbitmq"
)

var (
	nccn_db  dbservice.DBService
	nsru_db  *dbservice.NSRUDB
	err_nsru error
)

func InitDB() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	micro2.LoadConfigFromEnv()
	//config := micro2.LoadConfigFromEnv()
	//config.Show()
}


type UsageRecord struct {
	namespace            string
	CpuTotal             int32
	CpuUsed              int32
	MemoryTotal          int32
	MemoryUsed           int32
	StorageTotal         int32
	StorageUsed          int32
	Count                int32
	Start                int64
}

func CalculateFees(usage *UsageRecord) int32{
	reserveFees := float64(usage.CpuTotal) * 0.01  + float64(usage.MemoryTotal) * 0.01 + float64(usage.StorageTotal)  * 0.01
	usedFees := float64(usage.CpuTotal) * 0.02  + float64(usage.MemoryUsed) * 0.02 + float64(usage.StorageUsed)  * 0.02
	return int32(reserveFees + usedFees)
}


func CalcuateFeesAndSaveToDataBase(usage *UsageRecord) *dbservice.DailyFeesRecord{
	log.Printf("CalcuateFeesAndSaveToDataBase  ---->  %+v \n", usage)
	ns , _ :=nccn_db.GetNamespace(usage.namespace)


	ns.UID = "1111"
	ns.ClusterID = "b4668d6b-3257-47a5-a044-9c002be7dccd"


	record := dbservice.DailyFeesRecord{}
	record.UID = ns.UID
	record.ClusterId = ns.ClusterID
	record.UserType = dbservice.ClusterUser

	record.Fees = CalculateFees(usage)
	now := time.Now().Unix()
	record.CreateDate = &timestamp.Timestamp{Seconds: now}
	record.Date = usage.Start
	record.UserType = dbservice.ClusterUser
	record.Namespace = usage.namespace
	record.Usage.CpuTotal = usage.CpuTotal
	record.Usage.CpuUsed = usage.CpuUsed
	record.Usage.MemoryTotal = usage.MemoryTotal
	record.Usage.MemoryUsed = usage.MemoryUsed
	record.Usage.StorageTotal = usage.StorageTotal
	record.Usage.StorageUsed = usage.StorageUsed
	nccn_db.InsertDailyFees(&record)
	return &record
	//

}


type Metrics struct {
	TotalCPU     int64
	UsedCPU      int64
	TotalMemory  int64
	UsedMemory   int64
	TotalStorage int64
	UsedStorage  int64

	ImageCount    int64
	EndPointCount int64
	NetworkIO     int64 // No data
}

func main() {
	InitDB()

	if nsru_db, err_nsru = dbservice.NewNSRUDB(); err_nsru != nil {
		log.Fatal(err_nsru.Error())
	}
	defer nsru_db.Close()


	nccn_db, _ = dbservice.New()



	layOut := "2006-01-02"


	var start int64
	var end int64

	if len(os.Args) == 1 {
		// today
		diff := 24 * time.Hour
		yesterday := time.Now().Add(-diff)

		yesterdayStart := yesterday.Format(layOut)

		log.Printf(">>>> processing day %s \n", yesterdayStart)

		dateStamp, _ := time.Parse(layOut, yesterdayStart)

		//processing yesterday

		start = dateStamp.Unix()
		end = start + 86400



	}else{
		processingDay := os.Args[1]
		dateStamp, error := time.Parse(layOut, processingDay)


		log.Printf(">>>> processing day %s \n", processingDay)

		if error != nil {
			log.Printf("input processing day format error, expected : %s", layOut)
			return
		}

		start = dateStamp.Unix()
		end = start + 86400

	}



	//list, error  := db2.GetNamespaceResourceUsage()

	list, error  := nsru_db.GetNamespaceResourceUsageWithTimeSpan(start, end)

	records := make(map[string]*UsageRecord, 0)

	if error != nil {
		log.Printf("error %s \n", error.Error())
	}else{


		for _, record :=range *list {
			//log.Printf("souce ----> record %+v", record)

			namespace := record.Name
			value, ok := records[namespace]
			if ok {
				//fmt.Println("value: ", value)
				//value.CpuTotal = int32(record.CPU)
				//value.MemoryTotal = int32(record.Mem)
				//value.StorageTotal = int32(record.Disk)
				value.CpuUsed += int32(record.CPUUsed)
				value.MemoryUsed += int32(record.MemUsed)
				value.StorageUsed += int32(record.DiskUsed)
				value.Count += 1
			} else {

				//fmt.Println("key not found")

				r := UsageRecord{}
				r.namespace = namespace
				r.CpuTotal = int32(record.CPU)
				r.MemoryTotal = int32(record.Mem)
				r.StorageTotal = int32(record.Disk)
				r.CpuUsed = int32(record.CPUUsed)
				r.MemoryUsed = int32(record.MemUsed)
				r.StorageUsed = int32(record.DiskUsed)
				r.Count = 1
				r.Start = start

				records[namespace] = &r

			}

			//log.Printf("%+v \n", record)
		}


		for _, v := range records {
			fmt.Printf("%+v \n", v)
		}



		// for each namesapce calculate fees

		clusersRecords := make(map[string]*dbservice.DailyFeesRecord, 0)

		for _, v := range records {
			v.CpuUsed = v.CpuUsed/v.Count
			v.MemoryUsed = v.MemoryUsed/v.Count
			v.StorageUsed = v.StorageUsed/v.Count
			fmt.Printf("user insert ##### %+v \n", v)
			namespaceFees := CalcuateFeesAndSaveToDataBase(v)




			// calculate cluster total  usage and fees by user suage

			value, ok := clusersRecords[namespaceFees.ClusterId]

			if ok {

				value.Usage.CpuUsed += namespaceFees.Usage.CpuUsed
				value.Usage.MemoryUsed += namespaceFees.Usage.MemoryUsed
				value.Usage.StorageUsed += namespaceFees.Usage.StorageUsed
				value.Fees += namespaceFees.Fees
	//			value.Count ++

			}else{

				r := dbservice.DailyFeesRecord{}
				usage := dbservice.Usage{}
				usage.CpuUsed = namespaceFees.Usage.CpuUsed
				usage.MemoryUsed = namespaceFees.Usage.MemoryUsed
				usage.StorageUsed = namespaceFees.Usage.StorageUsed
				r.Usage = usage
				r.Fees = namespaceFees.Fees
				r.ClusterId = namespaceFees.ClusterId
				r.Namespace = namespaceFees.Namespace
		//		r.Count = 0



				now := time.Now().Unix()
				r.CreateDate = &timestamp.Timestamp{Seconds: now}
				r.Date = namespaceFees.Date
				r.UserType = dbservice.ClusterProvider
				r.CreateDate = namespaceFees.CreateDate

				clusersRecords[r.ClusterId] = &r
			}

		}


		// insert cluster provider

		for _ , v := range clusersRecords {
			cluster, _ := nccn_db.GetCluser(v.ClusterId)

			metrics := Metrics{}

			//v.Usage.CpuUsed = v.Usage.CpuUsed/v.count
			//v.Usage.MemoryUsed = v.Usage.MemoryUsed/v.count
			//v.Usage.StorageUsed = v.Usage.StorageUsed/v.count

			v.Fees = - v.Fees   // cluster fees is negitive

			if cluster == nil || cluster.DcHeartbeatReport == nil {
				log.Printf("can not find the cluster record")

			}else{
				json.Unmarshal([]byte(cluster.DcHeartbeatReport.Metrics), &metrics)
				v.Usage.CpuTotal = int32(metrics.TotalCPU)
				v.Usage.MemoryTotal = int32(metrics.TotalMemory)
				v.Usage.StorageTotal = int32(metrics.TotalStorage)
				v.UID = cluster.UserId

			}




			nccn_db.InsertDailyFees(v)
		}


	}


}




