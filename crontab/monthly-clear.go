package main

import (
	"github.com/Ankr-network/dccn-fees/db-service"
	"github.com/golang/protobuf/ptypes/timestamp"
	"log"
	"os"
	"strings"
	"time"
	"github.com/google/uuid"
)

var (
	dccn_fees  dbservice.DBService
	err error
)

func getAvailableID() string{
	for {
		clearID := uuid.New().String()
		a := strings.Split(clearID, "-")
		ID := a[0]
		_, error := dccn_fees.GetClearingRecord(ID)

		log.Printf("create new ID %s for clear \n ", ID)

		if error != nil { // no found , ok
			return ID
		}
	}

}

func main() {

	dccn_fees, _ = dbservice.New()



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


	list, _  := dccn_fees.GetMonthFeesWithTimeSpan(start, end)

	records := make(map[string]*dbservice.MonthlyClearing, 0)

	log.Printf("total record of daily fees %d  from %d to %d  \n", len(records), start, end)

	for _, record :=range *list {

        uid := record.UID


		log.Printf("process uid fees %s ", uid)

        if len(uid) == 0 {
        	log.Printf("uid is empty \n")
        	continue
		}

		value, ok := records[uid]

		if ok {
			value.Namespace[record.Namespace] = record.Fees
			value.Charge += record.Fees

		} else {
			now := time.Now()
            //log.Printf("process uid %s \n", record.UID)
			user, _ := dccn_fees.GetUser(record.UID)

			r := dbservice.MonthlyClearing{}
			//r.Usage = record.Usage
			r.Namespace = make(map[string]int32)
			r.Namespace[record.Namespace] = record.Fees
			r.UserType = record.UserType
			r.Month = start
			r.CreateDate =  &timestamp.Timestamp{Seconds: now.Unix()}
			r.Start =   &timestamp.Timestamp{Seconds: start}
			r.End =  &timestamp.Timestamp{Seconds: end}
			r.PaidDate = &timestamp.Timestamp{Seconds: 0}
			r.UID = record.UID
			r.Name = user.Name
			r.Charge = record.Fees
			r.Status = dbservice.UnPaid

			records[uid] = &r
		}
	}


    log.Printf("total monthly clearing %d \n", len(records))

	for _ , v := range records {

		v.Total = v.Charge

		v.ID = getAvailableID()
		dccn_fees.InsertMonthlyClearing(v)
	}


}

