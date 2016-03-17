package vmail
import (
	"fmt"
	"./vproto"
	"strings"
	"gopkg.in/mgo.v2/bson"
)

type TransportChannels struct {
	VServer chan vproto.VMessage
}

type MessagePlane struct {
	Channels		TransportChannels
	Mongo 			*MongoStore
	StorageDomain	string
}

func getDomain(address string) string{
	return address[strings.Index(address, "@") + 1:]
}

func allRecepients(message *vproto.VMessage) []string {
	return append(message.GetReceivers(), message.GetHiddenReceivers()...)
}

func newMail(mongo *MongoStore, username string) []vproto.VMessage {
	return mongo.msgsFromDB(bson.M{}, username+"_rx")
}

func (this *MessagePlane) Start(){
	for {
		var msg vproto.VMessage
		select {
		case msg = <- this.Channels.VServer:
			fmt.Println("vServer Message")
		}
		if getDomain(msg.GetSender()) == this.StorageDomain {

		}
		this.Mongo.toDB(msg)
	}
}