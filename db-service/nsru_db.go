package dbservice

import (
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	metering "github.com/Ankr-network/dccn-common/metering"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type NSRUDBService interface {
	Insert() error
	GetNamespaceResourceUsage()(*[]*metering.NamespaceResourceMetering, error)
    GetNamespaceResourceUsageWithTimeSpan(start int64, end int64)(*[]*metering.NamespaceResourceMetering, error)
	GetNamespaceResourceUsageWithTimeSpanNamespace(namespace string ,start int64, end int64) (*[]*metering.NamespaceResourceMetering, error)
	Close()
}





// UserDB implements DBService
type NSRUDB struct {
	database *mgo.Database
	metering *mgo.Collection
}

// New returns DBService.
func NewNSRUDB() (*NSRUDB, error) {

		database := micro2.GetDB("dccn_nsrc")
	metering := database.C("dccn_nsrc_metering")
	return &NSRUDB{
		database: database,
		metering: metering,
	}, nil
}

func (p *NSRUDB) Close() {
	//p.Close()
}




func (p *NSRUDB)GetNamespaceResourceUsage()(*[]*metering.NamespaceResourceMetering, error){
	var list []*metering.NamespaceResourceMetering
	if err :=  p.metering.Find(bson.M{}).All(&list); err != nil {
		return nil, err
	}
	return &list, nil

}


func (p *NSRUDB)GetNamespaceResourceUsageWithTimeSpan(start int64, end int64)(*[]*metering.NamespaceResourceMetering, error){
	var list []*metering.NamespaceResourceMetering
	if err :=  p.metering.Find(bson.M{ "timestamp" : 	bson.M{
		"$gte": start,
		"$lt": end,
	}}).All(&list); err != nil {
		return nil, err
	}
	return &list, nil

}

func (p *NSRUDB) GetNamespaceResourceUsageWithTimeSpanNamespace(namespace string ,start int64, end int64) (*[]*metering.NamespaceResourceMetering, error) {
	var list []*metering.NamespaceResourceMetering
	if err :=  p.metering.Find(bson.M{"name":namespace, "timestamp" : 	bson.M{
		"$gte": start,
		"$lt": end,
	}}).All(&list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (p *NSRUDB)Insert() error {
       return p.metering.Insert(bson.M{"ttest":"test"})
}
