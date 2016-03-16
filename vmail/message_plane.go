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

func getDomain(address string) string{
	return address[strings.Index(address, "@") + 1:]
}

func allRecepients(message *vproto.VMessage) []string {
	return append(message.GetReceivers(), message.GetHiddenReceivers()...)
}

func newMail(mongo *MongoStore, username string) []vproto.VMessage {
	return mongo.msgsFromDB(bson.M{}, username+"_rx")
}

func MessagePlane(channels TransportChannels, mongo *MongoStore) {
	for {
		var msg vproto.VMessage
		select {
		case msg = <- channels.VServer:
			fmt.Println("vServer Message")
		}
		mongo.toDB(msg)
	}
}