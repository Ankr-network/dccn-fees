package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"github.com/Ankr-network/dccn-common/util"
	"github.com/Ankr-network/dccn-fees/db-service"
	"github.com/golang/protobuf/ptypes/timestamp"
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

// this api for cluster provider dashboard
func (p *Handler) ClusterDashBoard(ctx context.Context, req *dcmgr.DashBoardRequest) (*dcmgr.DashBoardResponse, error) {
	uid := util.GetUserID(ctx)
	rsp := &dcmgr.DashBoardResponse{}
	rsp.CurrentUsage = &dcmgr.Usage{}
	rsp.TotalIncome = -int32(p.db.GetTotalIncomeForProvider(uid))
	log.Printf("total income %d  for user [%s] \n", rsp.TotalIncome, uid)
	cluster, err := p.db.GetClusterByUserID(uid)
	if err != nil {
		log.Printf("user does not register cluster \n")
		return rsp, nil
	}

	metrics := Metrics{}
	if err := json.Unmarshal([]byte(cluster.DcHeartbeatReport.Metrics), &metrics); err != nil {
		log.Printf("error metrics: metrics is empty  ")
		return rsp, nil
	} else {
		rsp.CurrentUsage.CpuTotal = int32(metrics.TotalCPU)
		rsp.CurrentUsage.CpuUsed = int32(metrics.UsedCPU)
		rsp.CurrentUsage.MemoryTotal = int32(metrics.TotalMemory)
		rsp.CurrentUsage.MemoryUsed = int32(metrics.UsedMemory)
		rsp.CurrentUsage.StorageTotal = int32(metrics.TotalStorage)
		rsp.CurrentUsage.StorageUsed = int32(metrics.UsedStorage)
	}

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
	if list, err := p.db.GetDailyFeesForProvider(uid, start, end); err != nil {
		log.Println(err.Error())
		log.Println("DataCenterList failure")
		return nil, err
	} else {
		log.Printf("week records  count: %d \n", len(*list))
		for _, record := range *list {
			log.Printf("print daily record %+v \n", record)
			income := dcmgr.Income{}
			income.Usage = &dcmgr.Usage{}
			income.Income = -record.Fees
			income.Date = &timestamp.Timestamp{Seconds: record.Date}

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
	if list, err := p.db.GetDailyFeesForProvider(uid, start, end); err != nil {
		log.Println(err.Error())
		log.Println("DataCenterList failure")
		return nil, err
	} else {
		log.Printf("monthly record  %d \n", len(*list))

		for _, record := range *list {
			income := dcmgr.Income{}
			income.Usage = &dcmgr.Usage{}
			income.Income = -record.Fees
			income.Date = &timestamp.Timestamp{Seconds: record.Date}

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
	if list, err := p.db.GetMonthlyFeesForProvider(uid, start, end); err != nil {
		log.Println(err.Error())
		log.Println("DataCenterList failure")
		return nil, err
	} else {
		log.Printf("year record  %d \n", len(*list))

		for _, record := range *list {
			income := dcmgr.Income{}
			income.Usage = &dcmgr.Usage{}
			income.Income = -record.Fees

			income.Date = &timestamp.Timestamp{Seconds: record.Month}

			income.Usage.CpuTotal = record.Usage.CpuTotal
			income.Usage.CpuUsed = record.Usage.CpuUsed
			income.Usage.MemoryTotal = record.Usage.MemoryTotal
			income.Usage.MemoryUsed = record.Usage.MemoryUsed
			income.Usage.StorageTotal = record.Usage.StorageTotal
			income.Usage.StorageUsed = record.Usage.StorageUsed
			rsp.Year = append(rsp.Year, &income)
		}

	}

	return rsp, nil
}

func (p *Handler) UserHistoryFeesList(ctx context.Context, req *dcmgr.HistoryFeesRequest) (*dcmgr.HistoryFeesResponse, error) {
	log.Printf("UserHistoryFeesList  %+v \n", req)
	// access control
	uid := util.GetUserID(ctx)
	res := newResource(req.TeamId, "historyFees")
	if err := checkAccess(ctx, uid, res, "list"); err != nil {
		return nil, handleError(err)
	}

	rsp := &dcmgr.HistoryFeesResponse{}
	rsp.Records = make([]*dcmgr.MonthRecord, 0)

	now := time.Now()
	currentLocation := now.UTC().Location()

	startArray := strings.Split(req.Start, "-")
	endArray := strings.Split(req.End, "-")

	if len(startArray) != 3 || len(endArray) != 3 {
		log.Printf("parse start or end error %s %s \n", req.Start, req.End)
		return rsp, nil
	}

	sYear, _ := strconv.Atoi(startArray[0])
	sMonth, _ := strconv.Atoi(startArray[1])
	eYear, _ := strconv.Atoi(endArray[0])
	eMonth, _ := strconv.Atoi(endArray[1])

	firstOfMonth := time.Date(sYear, time.Month(sMonth), 01, 0, 0, 0, 0, currentLocation)
	lastOfMonth := time.Date(eYear, time.Month(eMonth), 01, 0, 0, 0, 0, currentLocation)
	startTimeStamp := firstOfMonth.Unix()
	endTimeStamp := lastOfMonth.Unix()

	log.Printf("start year %d start month %d end year %d end month  %d  \n", sYear, sMonth, eYear, eMonth)

	records, err := p.db.GetMonthClearingWithTimeSpanForUser(req.TeamId, startTimeStamp, endTimeStamp)
	if err != nil {
		log.Printf("UserHistoryFeesList error %s \n", err.Error())
		return rsp, nil
	}

	for _, record := range *records {
		ns := dcmgr.MonthRecord{}
		ns.Amount = record.Total
		ns.Invoice = record.ID
		ns.Method = "ERC20"
		ns.PaymentDate = strconv.FormatInt(record.PaidDate.Seconds, 10)
		ns.PaymentStatus = p.db.ConvertClearingStatus(record.Status)

		log.Printf("record %+v \n", ns)
		rsp.Records = append(rsp.Records, &ns)
	}
	return rsp, nil
}

func (p *Handler) MonthFeesDetail(ctx context.Context, req *dcmgr.FeesDetailRequest) (*dcmgr.FeesDetailResponse, error) {
	// access control
	uid := util.GetUserID(ctx)
	res := newResource(req.TeamId, "monthFeesDetail")

	if err := checkAccess(ctx, uid, res, "get"); err != nil {
		return nil, handleError(err)
	}

	rsp := &dcmgr.FeesDetailResponse{}

	now := time.Now()
	currentLocation := now.UTC().Location()

	if len(req.TeamId) == 0 {
		log.Printf("error for MonthFeesDetail, teamId is empty")
		return rsp, nil
	}

	if len(req.Month) == 0 {
		return p.CacluateCurrentMonthFees(req.TeamId)
	}

	currentYear, currentMonth, _ := now.Date()

	s := strings.Split(req.Month, "-")

	if len(s) != 3 {
		return rsp, nil
	}

	year, _ := strconv.Atoi(s[0])
	month, _ := strconv.Atoi(s[1])

	firstOfCurrentMonth := time.Date(currentYear, currentMonth, 01, 0, 0, 0, 0, currentLocation)
	firstOfMonth := time.Date(year, time.Month(month), 01, 0, 0, 0, 0, currentLocation)
	firstOfMonthTimeStamp := firstOfMonth.Unix()

	if firstOfCurrentMonth == firstOfMonth {
		return p.CacluateCurrentMonthFees(req.TeamId)
	}

	log.Printf("year %d  month %d timestamp %d  \n", year, month, firstOfMonthTimeStamp)

	log.Printf("uid %s %d \n", req.TeamId, firstOfMonthTimeStamp)
	record, error := p.db.GetMonthlyClearingForUser(req.TeamId, firstOfMonthTimeStamp)

	if error == nil {
		rsp.Start = strconv.FormatInt(record.Start.Seconds, 10)
		rsp.End = strconv.FormatInt(record.End.Seconds, 10)
		rsp.Account = record.TeamID
		rsp.Attn = record.Name
		rsp.Credits = record.Credit
		rsp.InvoiceNumber = record.ID
		rsp.Tax = record.Tax
		rsp.Total = record.Total
		rsp.Charges = record.Charge
		rsp.IssueDate = strconv.FormatInt(record.CreateDate.Seconds, 10)
		rsp.NsFees = make([]*dcmgr.NamespaceFees, 0)
		for k, v := range record.Namespace {
			ns := dcmgr.NamespaceFees{Name: k, Charge: v}
			rsp.NsFees = append(rsp.NsFees, &ns)
		}
	} else {
		log.Printf("MonthFeesDetail error %s \n", error.Error())
	}

	return rsp, nil
}

func (p *Handler) CacluateCurrentMonthFees(uid string) (*dcmgr.FeesDetailResponse, error) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.UTC().Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	startLastMonth := firstOfMonth.AddDate(0, 0, 0)
	endlastMonth := firstOfMonth.AddDate(0, 1, 0)

	endlastMonth = endlastMonth.Add(-time.Second)

	start := startLastMonth.Unix()
	end := endlastMonth.Unix()
	list, _ := p.db.GetDailyFeesWithTimeSpanAndUid(uid, start, end)

	records := make(map[string]*dbservice.MonthlyFeesRecord, 0)

	log.Printf("total record of daily fees %d  from %d to %d  \n", len(records), start, end)

	for _, record := range *list {

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
			r.TeamID = record.TeamID
			r.Fees = record.Fees
			r.Count = 1

			now := time.Now().Unix()
			r.CreateDate = &timestamp.Timestamp{Seconds: now}

			records[namespace] = &r
		}
	}

	r := dbservice.MonthlyClearing{}
	user, _ := p.db.GetUser(uid)

	r.Namespace = make(map[string]int32)
	r.UserType = 0
	r.Month = start
	r.CreateDate = &timestamp.Timestamp{Seconds: now.Unix()}
	r.Start = &timestamp.Timestamp{Seconds: start}
	r.End = &timestamp.Timestamp{Seconds: end}
	r.PaidDate = &timestamp.Timestamp{Seconds: 0}
	r.TeamID = uid
	if user != nil {
		r.Name = user.Name
	}
	r.Charge = 0
	r.Status = dbservice.UnPaid

	count := 0
	for _, v := range records {
		v.Usage.CpuUsed += v.Usage.CpuUsed / v.Count
		v.Usage.MemoryUsed += v.Usage.MemoryUsed / v.Count
		v.Usage.StorageUsed += v.Usage.StorageUsed / v.Count

		if count == 0 {
			r.Namespace[v.Namespace] = v.Fees
			r.UserType = v.UserType
			r.Charge = v.Fees
			//r.Usage = record.Usage
		} else {

			r.Namespace[v.Namespace] = v.Fees
			r.Charge += v.Fees
		}
		count++
	}

	r.Total = r.Charge

	rsp := &dcmgr.FeesDetailResponse{}

	rsp.Start = strconv.FormatInt(r.Start.Seconds, 10)
	rsp.End = strconv.FormatInt(r.End.Seconds, 10)
	rsp.Account = r.TeamID
	rsp.Attn = r.Name
	rsp.Credits = r.Credit
	rsp.InvoiceNumber = r.ID
	rsp.Tax = r.Tax
	rsp.Total = r.Total
	rsp.Charges = r.Charge
	rsp.IssueDate = strconv.FormatInt(r.CreateDate.Seconds, 10)
	rsp.NsFees = make([]*dcmgr.NamespaceFees, 0)
	for k, v := range r.Namespace {
		ns := dcmgr.NamespaceFees{Name: k, Charge: v}
		rsp.NsFees = append(rsp.NsFees, &ns)
	}

	return rsp, nil
}

func (p *Handler) InvoiceDetail(ctx context.Context, req *dcmgr.InvoiceDetailRequest) (*dcmgr.FeesDetailResponse, error) {
	// access control
	uid := util.GetUserID(ctx)
	res := newResource(req.TeamId, fmt.Sprintf("invoices/%s", req.InvoiceId))
	if err := checkAccess(ctx, uid, res, "get"); err != nil {
		return nil, handleError(err)
	}

	rsp := &dcmgr.FeesDetailResponse{}

	invoice_id := req.InvoiceId
	log.Printf("InvoiceDetail for invoiceid  %s \n", invoice_id)
	record, error := p.db.GetClearingRecordForUser(invoice_id)

	if error != nil {
		log.Printf("InvoiceDetail error %s \n", error.Error())
		return rsp, nil
	}

	if record.TeamID != req.TeamId {
		log.Printf("InvoiceDetail error teamid[%s] !=  record's teamid [%s] \n", req.TeamId, record.TeamID)
		return rsp, nil
	}

	if error == nil && record.TeamID == req.TeamId {
		rsp.Start = strconv.FormatInt(record.Start.Seconds, 10)
		rsp.End = strconv.FormatInt(record.End.Seconds, 10)
		rsp.Account = record.TeamID
		rsp.Attn = record.Name
		rsp.Credits = record.Credit
		rsp.InvoiceNumber = record.ID
		rsp.Tax = record.Tax
		rsp.Total = record.Total
		rsp.Charges = record.Charge
		rsp.IssueDate = strconv.FormatInt(record.CreateDate.Seconds, 10)
		rsp.NsFees = make([]*dcmgr.NamespaceFees, 0)
		for k, v := range record.Namespace {
			ns := dcmgr.NamespaceFees{Name: k, Charge: v}
			rsp.NsFees = append(rsp.NsFees, &ns)
		}
	} else {

	}

	return rsp, nil
}
