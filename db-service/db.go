package dbservice

import (
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"errors"
	"github.com/Ankr-network/dccn-common/protos"
	"gopkg.in/mgo.v2"
	"github.com/golang/protobuf/ptypes/timestamp"
	"gopkg.in/mgo.v2/bson"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	"log"
)


type DataCenterRecord struct {
	DcId              string
	ClusterName       string
	GeoLocation       *common_proto.GeoLocation
	DcStatus          common_proto.DCStatus
	DcAttributes      *common_proto.DataCenterAttributes
	DcHeartbeatReport *common_proto.DCHeartbeatReport
	UserId            string
	Clientcert        string
}
type UserRecord struct {
	ID               string             `bson:"id"`
	Email            string             `bson:"email"`
	Name             string             `bson:"name"`

}

type DBService interface {
	GetUser(id string) (*UserRecord, error)
	GetCluser(id string) (*DataCenterRecord, error)
	InsertDailyFees(record *DailyFeesRecord) error
	InsertMonthlyFees(record *MonthlyFeesRecord) error
	InsertMonthlyClearing(record *MonthlyClearing) error
	GetDailyFees(uid string, start_time int64, end_time int64)(*[]*DailyFeesRecord, error)
	GetDailyFeesForProvider(uid string, start_time int64, end_time int64)(*[]*DailyFeesRecord, error)
	GetMonthlyFees(uid string, start_time int64, end_time int64)(*[]*MonthlyFeesRecord, error)
	GetMonthlyFeesForProvider(uid string, start_time int64, end_time int64)(*[]*MonthlyFeesRecord, error)
	GetNamespace(namespaceId string) (*NamespaceRecord, error)
	GetDailyFeesWithUidAndDate(uid string, date int64, namespace string) (*DailyFeesRecord, error)
	GetMonthlyClearing(uid string, date int64) (*MonthlyClearing, error)
	GetMonthlyClearingForUser(uid string, date int64) (*MonthlyClearing, error)  // only fees > 0
	GetDailyFeesWithTimeSpan(start int64, end int64)(*[]*DailyFeesRecord, error)
	GetDailyFeesWithTimeSpanForProvider(start int64, end int64)(*[]*DailyFeesRecord, error)
	GetDailyFeesWithTimeSpanAndUid(uid string, start int64, end int64)(*[]*DailyFeesRecord, error)
	GetMonthFeesWithTimeSpan(start int64, end int64)(*[]*MonthlyFeesRecord, error)
	GetMonthFeesWithTimeSpanForProvider(start int64, end int64)(*[]*MonthlyFeesRecord, error)
	GetMonthClearingWithTimeSpanForUser(uid string, start int64, end int64)(*[]*MonthlyClearing, error)
	GetMonthClearingWithTimeSpanForProvider(uid string, start int64, end int64)(*[]*MonthlyClearing, error) // fees > 0
	GetTotalIncomeForProvider(uid string) int
	GetClearingRecord(id string) (*MonthlyClearing, error)
	GetClearingRecordForUser(id string) (*MonthlyClearing, error) //fees > 0 for user
	ConvertClearingStatus(status ClearingStatus) string
	GetClusterByUserID(uid string) (*DataCenterRecord, error)
	Close()
}

type UserType int

const (
	ClusterProvider    UserType = 0
	ClusterUser        UserType = 1
)


type ClearingStatus int

const (
	UnPaid    ClearingStatus = 0
	Paid      ClearingStatus = 1

)

type Usage struct {
	CpuTotal             int32
	CpuUsed              int32
	MemoryTotal          int32
	MemoryUsed           int32
	StorageTotal         int32
	StorageUsed          int32
}

type DailyFeesRecord struct {
	UID       string
	ClusterId  string
	Namespace  string
    Fees       int32    // cent of dollar
	UserType    UserType
	Date        int64
	Usage        Usage
	CreateDate  *timestamp.Timestamp
	Count       int32

}

type MonthlyFeesRecord struct {
	UID       string
	Namespace  string
	Fees       int32    // cent of dollar
	UserType    UserType
	Month        int64
	Usage        Usage
	CreateDate  *timestamp.Timestamp
	Count       int32
}


type MonthlyClearing struct {
	ID        string
	UID       string
	Name      string
	Namespace  map[string]int32
	Total   int32    // cent of dollar
	Charge  int32
	Credit   int32
	Tax      int32
	UserType    UserType
	Month        int64
	CreateDate  *timestamp.Timestamp
	PaidDate  *timestamp.Timestamp
	Start  *timestamp.Timestamp
	End  *timestamp.Timestamp
	Status      ClearingStatus
}

// UserDB implements DBService
type DB struct {
	daily *mgo.Collection
	monthly *mgo.Collection
	clearing *mgo.Collection
	namespace *mgo.Collection
	cluser    *mgo.Collection
	user    *mgo.Collection
}

type NamespaceRecord struct {
	ID                   string // short hash of uid+name+cluster_id
	Name                 string
	NameUpdating         string
	UID                  string
	ClusterID            string //id of cluster
	ClusterName          string //name of cluster
	LastModifiedDate     *timestamp.Timestamp
	CreationDate         *timestamp.Timestamp
	CpuLimit             uint32
	CpuLimitUpdating     uint32
	MemLimit             uint32
	MemLimitUpdating     uint32
	StorageLimit         uint32
	StorageLimitUpdating uint32
	Status               common_proto.NamespaceStatus
	Event                common_proto.NamespaceEvent
}

// New returns DBService.
func New() (*DB, error) {
	config := micro2.GetConfig()
	config.DatabaseName = "dccn"
	dailyCollection := micro2.GetCollection("dailyfees")
	monthlyCollection := micro2.GetCollection("monthlyfees")
	clearingCollection := micro2.GetCollection("monthlyclearing")
	namespace := micro2.GetCollection("namespace")
	cluser := micro2.GetCollection("datacenter")
	userCollection  := micro2.GetCollection("user")
	return &DB{
		daily: dailyCollection,
		monthly:monthlyCollection,
		clearing:clearingCollection,
		namespace: namespace,
		cluser:cluser,
		user:userCollection,
	}, nil
}

func (p *DB) Close() {
	//p.Close()
}


// Create creates a new data center item if it not exists
func (p *DB) InsertDailyFees(record *DailyFeesRecord) error {

	if len(record.UID) == 0 {
		log.Printf("InsertDailyFees error, uid does not exist.")
		return nil
	}

	_, error := p.GetDailyFeesWithUidAndDate(record.UID, record.Date, record.Namespace)
	if error == nil { //exit old record
		log.Printf("----> update record %+v \n", record)
		return p.daily.Update(bson.M{"uid": record.UID, "date": record.Date, "namespace": record.Namespace},
			bson.M{"$set": record})

	}else{
		log.Printf("----> insert record %+v \n", record)
		return p.daily.Insert(record)

	}
}

// Create creates a new data center item if it not exists
func (p *DB) InsertMonthlyFees(record *MonthlyFeesRecord) error {
	_, error := p.GetMonthlyFeesUidAndDate(record.UID, record.Month, record.Namespace)


	if error == nil { //exit old record
		log.Printf("----> update monthly record %+v \n", record)
		return p.monthly.Update(bson.M{"uid": record.UID, "month": record.Month, "namespace": record.Namespace},
			bson.M{"$set": record})

	}else{
		log.Printf("----> insert monthly record %+v \n", record)
		return p.monthly.Insert(record)

	}

}

// Create creates a new data center item if it not exists
func (p *DB) InsertMonthlyFeesForProvider(record *MonthlyFeesRecord) error {
	_, error := p.GetMonthlyFeesUidAndDate(record.UID, record.Month, record.Namespace)


	if error == nil { //exit old record
		log.Printf("----> update monthly record %+v \n", record)
		return p.monthly.Update(bson.M{"uid": record.UID, "month": record.Month, "usertype" : ClusterUser, "namespace": record.Namespace},
			bson.M{"$set": record})

	}else{
		log.Printf("----> insert monthly record %+v \n", record)
		return p.monthly.Insert(record)

	}

}



func (p *DB)GetMonthlyClearing(uid string, date int64) (*MonthlyClearing, error){
	var record MonthlyClearing
	if err :=  p.clearing.Find(bson.M{"uid": uid, "month": date}).One(&record); err != nil {
		return nil, err
	}
	return &record, nil

}

func (p *DB)GetMonthlyClearingForUser(uid string, date int64) (*MonthlyClearing, error){
	var record MonthlyClearing
	if err :=  p.clearing.Find(bson.M{"uid": uid, "month": date, "usertype" : ClusterUser}).One(&record); err != nil {
		return nil, err
	}
	return &record, nil
}




func (p *DB) InsertMonthlyClearing(record *MonthlyClearing) error {
	r, error := p.GetMonthlyClearing(record.UID, record.Month)

	if error == nil && r.Status == Paid  {
		return errors.New("can not change paid record")
	}


	if error == nil { //exit old record
		log.Printf("find record, can not insert clearing twice  %+v \n", record)
		//return p.clearing.Update(bson.M{"uid": record.UID, "month": record.Month, "namespace": record.Namespace, "total": record.Total},
		//	bson.M{"$set": record})

	}else{
		log.Printf("----> insert monthly clearing %+v \n", record)
		return p.clearing.Insert(record)

	}

	return nil
}

func (p *DB) GetDailyFees(uid string, start_time int64, end_time int64) (*[]*DailyFeesRecord, error) {
	log.Printf("GetDailyFees uid %s start %d end %d \n", uid, start_time, end_time)
	var list []*DailyFeesRecord
	if err :=  p.daily.Find(bson.M{"uid":uid,"date" : 	bson.M{
		"$gte": start_time,
		"$lt": end_time,
	}}).All(&list); err != nil {
		return nil, err
	}

	return &list, nil
}


func (p *DB) GetDailyFeesForProvider(uid string, start_time int64, end_time int64) (*[]*DailyFeesRecord, error) {
	log.Printf("GetDailyFees uid %s start %d end %d \n", uid, start_time, end_time)
	var list []*DailyFeesRecord
	if err :=  p.daily.Find(bson.M{"uid":uid, "usertype" : ClusterProvider, "date" : 	bson.M{
		"$gte": start_time,
		"$lt": end_time,
	}}).All(&list); err != nil {
		return nil, err
	}

	return &list, nil
}

func (p *DB) GetMonthlyFees(uid string, start_time int64, end_time int64) (*[]*MonthlyFeesRecord, error) {
	var list []*MonthlyFeesRecord

	if err :=  p.monthly.Find(bson.M{ "uid":uid, "month" : 	bson.M{
		"$gte": start_time,
		"$lt": end_time,
	}}).All(&list); err != nil {
		return nil, err
	}

	return &list, nil
}


func (p *DB) GetMonthlyFeesForProvider(uid string, start_time int64, end_time int64) (*[]*MonthlyFeesRecord, error) {
	var list []*MonthlyFeesRecord

	if err :=  p.monthly.Find(bson.M{ "uid":uid,
		"usertype" : ClusterProvider,
		"month" : 	bson.M{
		"$gte": start_time,
		"$lt": end_time,
	}}).All(&list); err != nil {
		return nil, err
	}

	return &list, nil
}



func (p *DB) GetMonthlyFeesUidAndDate(uid string, date int64, namespace string) (*MonthlyFeesRecord, error) {
	var record MonthlyFeesRecord
	if err :=  p.monthly.Find(bson.M{"uid": uid, "month": date, "namespace": namespace}).One(&record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (p *DB) GetDailyFeesWithUidAndDate(uid string, month int64, namespace string) (*DailyFeesRecord, error) {
	var record DailyFeesRecord
	if err :=  p.daily.Find(bson.M{"uid": uid, "date": month, "namespace":namespace}).One(&record); err != nil {
		return nil, err
	}

	return &record, nil
}


func (p *DB)GetDailyFeesWithTimeSpan(start int64, end int64)(*[]*DailyFeesRecord, error){
	var list []*DailyFeesRecord
	if err :=  p.daily.Find(bson.M{
		"usertype" : ClusterUser ,
		"date" : 	bson.M{
		"$gte": start,
		"$lt": end,
	}}).All(&list); err != nil {
		return nil, err
	}
	return &list, nil

}

func (p *DB)GetDailyFeesWithTimeSpanForProvider(start int64, end int64)(*[]*DailyFeesRecord, error){
	var list []*DailyFeesRecord
	if err :=  p.daily.Find(bson.M{
		"usertype" : ClusterProvider ,
		"date" : 	bson.M{
			"$gte": start,
			"$lt": end,
		}}).All(&list); err != nil {
		return nil, err
	}
	return &list, nil

}


// this functon only for user spending which fees > 0

func (p *DB)GetDailyFeesWithTimeSpanAndUid(uid string, start int64, end int64)(*[]*DailyFeesRecord, error){
	var list []*DailyFeesRecord
	if err :=  p.daily.Find(bson.M{ "uid":uid,
		"usertype" : ClusterUser,
		"fees":bson.M{"$gte" : 0},
		"date" : 	bson.M{
		"$gte": start,
		"$lt": end,
	}}).All(&list); err != nil {
		return nil, err
	}
	return &list, nil
}




func (p *DB) GetNamespace(namespaceId string) (*NamespaceRecord, error) {
	var record NamespaceRecord

	if err :=  p.namespace.Find(bson.M{"id": namespaceId}).One(&record); err != nil {
		return &record, err
	}

	return &record , nil
}


// Get gets user item by id.
func (p *DB) GetCluser(id string) (*DataCenterRecord, error) {
	var center DataCenterRecord
	log.Printf("get datacetner %s \n", id)
	err := p.cluser.Find(bson.M{"dcid": id}).One(&center)
	return &center, err
}

func (p *DB) GetMonthFeesWithTimeSpan(start int64, end int64)(*[]*MonthlyFeesRecord, error){

	var list []*MonthlyFeesRecord
	if err :=  p.monthly.Find(bson.M{
		"usertype" : ClusterUser ,
		"month" : 	bson.M{
		"$gte": start,
		"$lt": end,
	}}).All(&list); err != nil {
		return nil, err
	}
	return &list, nil
}


func (p *DB) GetMonthFeesWithTimeSpanForProvider(start int64, end int64)(*[]*MonthlyFeesRecord, error){

	var list []*MonthlyFeesRecord
	if err :=  p.monthly.Find(bson.M{
		"usertype" : ClusterProvider ,
		"month" : 	bson.M{
			"$gte": start,
			"$lt": end,
		}}).All(&list); err != nil {
		return nil, err
	}
	return &list, nil
}


func (p *DB)GetTotalIncomeForProvider(uid string) int {
	pipe := p.daily.Pipe([]bson.M{bson.M{"$match": bson.M{"uid": uid, "usertype" : ClusterProvider}},bson.M{"$group": bson.M{"_id": "$uid",
		 "total": bson.M{"$sum": "$fees"}}}})
	resp := []bson.M{}
	iter := pipe.Iter()
	iter.All(&resp)

	if len(resp) > 0 {
		return  int(resp[0]["total"].(int))
	}
	 return 0
}

// Get gets user item by email.
func (p *DB) GetUser(id string) (*UserRecord, error) {
	var user UserRecord
	err := p.user.Find(bson.M{"id": id}).One(&user)
	if err != nil {
		return nil, errors.New(ankr_default.DbError+err.Error())
	}
	return &user, nil
}

func (p *DB) GetMonthClearingWithTimeSpanForUser(uid string, start int64, end int64)(*[]*MonthlyClearing, error){
	log.Printf("GetMonthClearingWithTimeSpan  uid %s start %d end %d \n", uid, start, end)
	var list []*MonthlyClearing
	if err :=  p.clearing.Find(bson.M{ "uid" : uid , "usertype" : ClusterUser,  "month" : 	bson.M{
		"$gte": start,
		"$lt": end,
	}}).All(&list); err != nil {
		return nil, err
	}
	return &list, nil

}


func (p *DB) GetMonthClearingWithTimeSpanForProvider(uid string, start int64, end int64)(*[]*MonthlyClearing, error){
	log.Printf("GetMonthClearingWithTimeSpan  uid %s start %d end %d \n", uid, start, end)
	var list []*MonthlyClearing
	if err :=  p.clearing.Find(bson.M{ "uid" : uid , "usertype" : ClusterProvider,  "month" : 	bson.M{
		"$gte": start,
		"$lt": end,
	}}).All(&list); err != nil {
		return nil, err
	}
	return &list, nil

}


func (p *DB) GetClearingRecord(id string) (*MonthlyClearing, error){
	var record MonthlyClearing
	log.Printf("GetClearingRecord from %s \n", id)
	err := p.clearing.Find(bson.M{"id": id}).One(&record)
	if err != nil {
		return nil, errors.New(ankr_default.DbError+err.Error())
	}

	log.Printf("GetClearingRecord record %+v \n", record)
	return &record, nil

}


func (p *DB) GetClearingRecordForUser(id string) (*MonthlyClearing, error){
	var record MonthlyClearing
	log.Printf("GetClearingRecord from %s \n", id)
	err := p.clearing.Find(bson.M{"id": id, "usertype" : ClusterUser}).One(&record)
	if err != nil {
		return nil, errors.New(ankr_default.DbError+err.Error())
	}

	log.Printf("GetClearingRecord record %+v \n", record)
	return &record, nil

}

func  (p *DB)ConvertClearingStatus(status ClearingStatus) string{
	if status == 1 {
		return "Paid"
	}

	return "Unpaid"

}


func (p *DB) GetClusterByUserID(uid string) (*DataCenterRecord, error) {
	var center DataCenterRecord
	log.Printf("GetClusterByUserID uid %s \n", uid)
	err := p.cluser.Find(bson.M{"userid": uid}).One(&center)
	if err != nil {
		return nil, errors.New(ankr_default.DbError+err.Error())
	}
	log.Printf("cluster %+v \n", center)
	return &center, err
}




