package vmail

import (
	"gopkg.in/mgo.v2"
	"./vproto"
	"gopkg.in/mgo.v2/bson"
)

type MongoStore struct {
	session 		*mgo.Session
	database		*mgo.Database
	storageDomain	string
}

func (this *MongoStore) msgsFromDB(query bson.M, collection string) []vproto.VMessage{
	var results []vproto.VMessage
	c := this.database.C(collection)
	c.Find(query).All(&results)
	return results
}

func (this *MongoStore) toDB(message vproto.VMessage){
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

func (this *MongoStore) Close(){
	this.session.Close()
}