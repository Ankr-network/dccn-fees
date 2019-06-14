package main

import (
	"github.com/Ankr-network/dccn-fees/db-service"
	"github.com/golang/protobuf/ptypes/timestamp"
	"log"
	"os"
	"time"
)

var (
	dccn  dbservice.DBService
	dccn_error error
)



func main() {
	//InitDB()



	dccn, _ = dbservice.New()



	layOut := "2006-01-02"


	var start int64
	var end int64

	if len(os.Args) == 1 {
		// today
		now := time.Now()
		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.UTC().Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		startLastMonth := firstOfMonth.AddDate(0, -1, 0)
		endlastMonth := firstOfMonth.AddDate(0, 0, 0)

		endlastMonth = endlastMonth.Add(-time.Second)


		start = startLastMonth.Unix()
		end = endlastMonth.Unix()

		log.Printf("month from %s to %s ", startLastMonth.String(), endlastMonth.String())


		//processing yesterday
	}else{

		processingDay := os.Args[1]
		d, _ := time.Parse(layOut, processingDay)
		year, month, _ := d.Date()
		now := time.Now()
		currentLocation := now.UTC().Location()


		firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, currentLocation)
		startLastMonth := firstOfMonth.AddDate(0, 0, 0)
		endlastMonth := firstOfMonth.AddDate(0, 1, 0)
		endlastMonth = endlastMonth.Add(-time.Second)

		log.Printf("month from %s to %s ", startLastMonth, endlastMonth)

		start = startLastMonth.Unix()
		end = endlastMonth.Unix()

	}


	list, _  := dccn.GetDailyFeesWithTimeSpan(start, end)

	records := make(map[string]*dbservice.MonthlyFeesRecord, 0)

	log.Printf("total record of daily fees %d  from %d to %d  \n", len(records), start, end)

	for _, record :=range *list {


		namespace := record.Namespace
		if record.UserType == dbservice.ClusterProvider {
			namespace = record.ClusterId
		}

		log.Printf("name space %s ", namespace)

		value, ok := records[namespace]

		if ok {
			value.Usage.CpuUsed += record.Usage.CpuUsed
			value.Usage.MemoryUsed += record.Usage.MemoryUsed
			value.Usage.StorageUsed += record.Usage.StorageUsed
			value.Count += 1
			value.Fees += record.Fees
		} else {

			r := dbservice.MonthlyFeesRecord{}
			r.Usage = record.Usage
			r.Namespace = record.Namespace
			r.UserType = record.UserType
			r.Month = start
			r.UID = record.UID
			r.Fees = record.Fees
			r.Count = 1

			now := time.Now().Unix()
			r.CreateDate = &timestamp.Timestamp{Seconds: now}

			records[namespace] = &r
		}
	}


	for _ , v := range records {
		v.Usage.CpuUsed += v.Usage.CpuUsed/v.Count
		v.Usage.MemoryUsed += v.Usage.MemoryUsed/v.Count
		v.Usage.StorageUsed += v.Usage.StorageUsed/v.Count

	   dccn.InsertMonthlyFees(v)
	}


}