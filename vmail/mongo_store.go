package vmail

import (
	"gopkg.in/mgo.v2"
	"./vmail_proto"
	"gopkg.in/mgo.v2/bson"
)

type MongoStore struct {
	session 		*mgo.Session
	database		*mgo.Database
	storageDomain	string
}

type VMailMail struct {
	vmail_proto.VMessage
}

func (this *MongoStore) msgsFromDB(query bson.M, collection string) []vmail_proto.VMessage{
	var results []vmail_proto.VMessage
	c := this.database.C(collection)
	c.Find(query).All(&results)
	return results
}

func (this *MongoStore) toDB(message vmail_proto.VMessage){
	if getDomain(message.GetSender()) == this.storageDomain {
		cname := message.GetSender()+"_tx"
		collection := this.database.C(cname)
		collection.Insert(&message)
	}
	receivers := allRecepients(&message)
	for _, receiver := range receivers {
		if getDomain(receiver) == this.storageDomain {
			cname := receiver+"_rx"
			collection := this.database.C(cname)
			collection.Insert(&message)
		}
	}
}

func (this *MongoStore) Init(storageDomain string, mongoServers string){
	this.storageDomain = storageDomain
	var err error
	this.session, err = mgo.Dial(mongoServers)
	if err != nil {
		panic(err)
	}
	this.database = this.session.DB("vmail")
}