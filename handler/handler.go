package handler

import (
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"github.com/Ankr-network/dccn-common/util"
	"github.com/Ankr-network/dccn-fees/db-service"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"strings"
	"time"
)






type Handler struct {
	db dbservice.DBService
}

func NewHandler(db dbservice.DBService) *Handler {
	handler := &Handler{
		db: db,
	}
	return handler
}

func (p *Handler) ClusterDashBoard(
	ctx context.Context, req *dcmgr.DashBoardRequest) (*dcmgr.DashBoardResponse, error) {

	//	log.Printf("context %+v", ctx)
	uid2 := util.GetUserID(ctx)
	log.Printf("uid2 %s \n", uid2)
	//
	uid := req.Uid

	rsp := dcmgr.DashBoardResponse{}
	rsp.TotalIncome = int32(p.db.GetTotalIncome(uid))

	log.Printf("total income %d \n", rsp.TotalIncome)

	//todo
	rsp.CurrentUsage = &dcmgr.Usage{}
	rsp.CurrentUsage.CpuTotal = 10000
	rsp.CurrentUsage.CpuUsed = 5000
	rsp.CurrentUsage.MemoryTotal = 10000
	rsp.CurrentUsage.MemoryUsed = 50000
	rsp.CurrentUsage.StorageTotal = 100000
	rsp.CurrentUsage.StorageUsed = 70000


	now := time.Now()
	currentYear, currentMonth, currentDay := now.Date()
	currentLocation := now.UTC().Location()

	firstOfMonth := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, currentLocation)
	startFirstDayOf7Days := firstOfMonth.AddDate(0, 0, -7)
	endLastDayOf7Days := firstOfMonth.AddDate(0, 0, 0)
	endLastDayOf7Days = endLastDayOf7Days.Add(-time.Second)



	start := startFirstDayOf7Days.Unix()
	end := endLastDayOf7Days.Unix()

	//layOut := "2006-01-02  12:23:45"
	log.Printf("start %s  %d end  %s %d \n", startFirstDayOf7Days, start, endLastDayOf7Days, end)

	rsp.Week = make([]*dcmgr.Income, 0)
	if list, err := p.db.GetDailyFees(uid,start,end); err != nil {
		log.Println(err.Error())
		log.Println("DataCenterList failure")
		return nil, err
	} else {
		log.Printf("week records  count: %d \n",len(*list))
		for _, record :=range *list {
            income := dcmgr.Income{}
            income.Usage = &dcmgr.Usage{}
			income.Income = record.Fees
			income.Date = record.CreateDate

			income.Usage.CpuTotal = record.Usage.CpuTotal
			income.Usage.CpuUsed = record.Usage.CpuUsed
			income.Usage.MemoryTotal = record.Usage.MemoryTotal
			income.Usage.MemoryUsed = record.Usage.MemoryUsed
			income.Usage.StorageTotal = record.Usage.StorageTotal
			income.Usage.StorageUsed = record.Usage.StorageUsed
			rsp.Week = append(rsp.Week, &income)
		}

	}

   // monthly
	startFirstDayOf30days := firstOfMonth.AddDate(0, 0, -30)
	endLastDayOf30Days := firstOfMonth.AddDate(0, 0, 0)
	endLastDayOf30Days = endLastDayOf30Days.Add(-time.Second)



	start = startFirstDayOf30days.Unix()
	end = endLastDayOf30Days.Unix()

	//layOut := "2006-01-02  12:23:45"
	log.Printf("start %s  %d end  %s %d \n", startFirstDayOf7Days, start, endLastDayOf7Days, end)

	rsp.Month = make([]*dcmgr.Income, 0)
	if list, err := p.db.GetDailyFees(uid,start,end); err != nil {
		log.Println(err.Error())
		log.Println("DataCenterList failure")
		return nil, err
	} else {
		log.Printf("monthly record  %d \n", len(*list))

		for _, record :=range *list {
			income := dcmgr.Income{}
			income.Usage = &dcmgr.Usage{}
			income.Income = record.Fees
			income.Date = record.CreateDate

			income.Usage.CpuTotal = record.Usage.CpuTotal
			income.Usage.CpuUsed = record.Usage.CpuUsed
			income.Usage.MemoryTotal = record.Usage.MemoryTotal
			income.Usage.MemoryUsed = record.Usage.MemoryUsed
			income.Usage.StorageTotal = record.Usage.StorageTotal
			income.Usage.StorageUsed = record.Usage.StorageUsed
			rsp.Month = append(rsp.Month, &income)
		}

	}


	// := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)

	startFirstDayOf1year := firstOfMonth.AddDate(-1, 0, 0)
	endLastDayOf1year := firstOfMonth.AddDate(0, 0, 0)
	endLastDayOf1year = endLastDayOf1year.Add(-time.Second)

	start = startFirstDayOf1year.Unix()
	end = endLastDayOf1year.Unix()


	//layOut := "2006-01-02  12:23:45"
	log.Printf("start %s  %d end  %s %d \n", startFirstDayOf1year, start, endLastDayOf1year, end)

	rsp.Year = make([]*dcmgr.Income, 0)
	if list, err := p.db.GetMonthlyFees(uid,start,end); err != nil {
		log.Println(err.Error())
		log.Println("DataCenterList failure")
		return nil, err
	} else {
		log.Printf("year record  %d \n", len(*list))

		for _, record :=range *list {
			income := dcmgr.Income{}
			income.Usage = &dcmgr.Usage{}
			income.Income = record.Fees
			income.Date = record.CreateDate

			income.Usage.CpuTotal = record.Usage.CpuTotal
			income.Usage.CpuUsed = record.Usage.CpuUsed
			income.Usage.MemoryTotal = record.Usage.MemoryTotal
			income.Usage.MemoryUsed = record.Usage.MemoryUsed
			income.Usage.StorageTotal = record.Usage.StorageTotal
			income.Usage.StorageUsed = record.Usage.StorageUsed
			rsp.Year = append(rsp.Year, &income)
		}

	}

	return &rsp, nil
}


func (p *Handler) UserHistoryFeesList(ctx context.Context, req *dcmgr.HistoryFeesRequest) (*dcmgr.HistoryFeesResponse, error){
	log.Printf("UserHistoryFeesList  %+v \n", req)
	rsp := &dcmgr.HistoryFeesResponse{}
	rsp.Records = make([]*dcmgr.MonthRecord, 0)

	now := time.Now()
	currentLocation := now.UTC().Location()

	startArray := strings.Split(req.Start, "-")
	endArray := strings.Split(req.End, "-")

	if len(startArray) != 3  || len(endArray) != 3 {
		log.Printf("parse start or end error %s %s \n", req.Start, req.End)
		return rsp, nil
	}

	sYear , _ := strconv.Atoi(startArray[0])
	sMonth , _ := strconv.Atoi(startArray[1])
	eYear , _ := strconv.Atoi(endArray[0])
	eMonth , _ := strconv.Atoi(endArray[1])

	firstOfMonth := time.Date(sYear,  time.Month(sMonth) , 01, 0, 0, 0, 0, currentLocation)
	lastOfMonth := time.Date(eYear,  time.Month(eMonth) , 01, 0, 0, 0, 0, currentLocation)
	startTimeStamp := firstOfMonth.Unix()
	endTimeStamp := lastOfMonth.Unix()

	log.Printf("start year %d start month %d end year %d end month  %d  \n", sYear, sMonth, eYear, eMonth )




	records, error := p.db.GetMonthClearingWithTimeSpan(req.Uid, startTimeStamp, endTimeStamp)

	if error != nil {
		log.Printf("UserHistoryFeesList error %s \n", error.Error())
		return rsp, nil
	}

	for _ , record := range *records {
		ns := dcmgr.MonthRecord{}
		ns.Amount = record.Total
		ns.Invoice = record.ID
		ns.Method = "ERC20"
		ns.PaymentDate = strconv.FormatInt(record.PaidDate.Seconds, 10 )

		log.Printf("record %+v \n", ns)
	    rsp.Records = append(rsp.Records, &ns)
	}
	return rsp, nil
}

func (p *Handler) MonthFeesDetail(ctx context.Context, req *dcmgr.FeesDetailRequest) (*dcmgr.FeesDetailResponse, error){
	 rsp := &dcmgr.FeesDetailResponse{}

	now := time.Now()
	currentLocation := now.UTC().Location()


	s := strings.Split(req.Month, "-")

	if len(s) != 3 {
		return rsp, nil
	}

	year , _ := strconv.Atoi(s[0])
	month , _ := strconv.Atoi(s[1])

	firstOfMonth := time.Date(year,  time.Month(month) , 01, 0, 0, 0, 0, currentLocation)
	firstOfMonthTimeStamp := firstOfMonth.Unix()

	log.Printf("year %d  month %d timestamp %d  \n", year, month , firstOfMonthTimeStamp)

	log.Printf("uid %s %d \n", req.Uid, firstOfMonthTimeStamp)
	record, error := p.db.GetMonthlyClearing(req.Uid, firstOfMonthTimeStamp)



	if error == nil {
		rsp.Start = strconv.FormatInt(record.Start.Seconds, 10)
		rsp.End = strconv.FormatInt(record.End.Seconds, 10)
		rsp.Account = record.UID
		rsp.Attn = record.Name
		rsp.Credits = record.Credit
		rsp.InvoiceNumber = record.ID
		rsp.Tax = record.Tax
		rsp.Total = record.Total
		rsp.Charges = record.Charge
		rsp.IssueDate = strconv.FormatInt(record.CreateDate.Seconds, 10)
		rsp.NsFees = make([]*dcmgr.NamespaceFees, 0)
		for k, v := range record.Namespace {
			ns := dcmgr.NamespaceFees{Name:k, Charge:v}
			rsp.NsFees = append(rsp.NsFees , &ns)
		}
	}else{
		log.Printf("MonthFeesDetail error %s \n", error.Error())
	}

	return rsp, nil
}
